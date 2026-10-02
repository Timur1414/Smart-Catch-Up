package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	apiDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/http"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (obj *responseWriter) WriteHeader(code int) {
	obj.StatusCode = code
	obj.ResponseWriter.WriteHeader(code)
}

func (obj *responseWriter) Write(b []byte) (int, error) {
	if obj.StatusCode == 0 {
		obj.StatusCode = http.StatusOK
	}
	return obj.ResponseWriter.Write(b)
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		defer func() {
			if err := recover(); err != nil {
				requestId, ok := r.Context().Value("request_id").(string)
				if !ok {
					log.Error("[panic middleware] panic",
						zap.Any("err", err))
				} else {
					log.Error("[panic middleware] panic",
						zap.String("request_id", requestId),
						zap.Any("err", err))
				}
				response := web_helpers.NewServerErrorResponse(requestId)
				web_helpers.WriteResponseJSON(w, response.Code, response)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetAccessLogger()
		requestId := uuid.New().String()
		ctx := context.WithValue(r.Context(), context_helper.ContextKeyRequestId, requestId)
		log.Info("Request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("request_id", requestId))
		wr := &responseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		timeStart := time.Now()
		next.ServeHTTP(wr, r.WithContext(ctx))
		duration := time.Since(timeStart)
		log.Info("Response",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("request_id", requestId),
			zap.Int("status_code", wr.StatusCode),
			zap.String("duration", duration.String()))
	})
}

func AuthMiddleware(next http.Handler, authServer authpb.AuthClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		path := r.URL.Path
		if (strings.HasPrefix(path, "/auth/") && path != "/auth/logout") || path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		requestId := context_helper.GetRequestIdFromContext(r.Context())
		cookie, err := apiDelivery.GetAccessCookie(r.Context(), r)
		if err != nil {
			log.Error("[auth middleware] failed to get access cookie", zap.Error(err))
			response := web_helpers.NewUnauthorizedResponse(requestId)
			web_helpers.WriteResponseJSON(w, response.Code, response)
			return
		}
		authResponse, err := authServer.IsAuth(r.Context(), &authpb.IsAuthRequest{AccessToken: cookie.Value})
		if err != nil || !authResponse.GetIsAuth() {
			response := web_helpers.NewUnauthorizedResponse(requestId)
			web_helpers.WriteResponseJSON(w, response.Code, response)
			return
		}
		ctx := context.WithValue(r.Context(), context_helper.ContextKeyUser, authResponse.GetUserId())
		log.Info("[auth middleware] auth success]", zap.Int64("user_id", authResponse.GetUserId()), zap.String("path", path))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
