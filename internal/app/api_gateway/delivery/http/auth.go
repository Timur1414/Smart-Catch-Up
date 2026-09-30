package http

import (
	"context"
	"net/http"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authServer authpb.AuthClient
}

func NewAuthHandler(authServer authpb.AuthClient) *AuthHandler {
	return &AuthHandler{authServer: authServer}
}

func (obj *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Login request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Refresh request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Register request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Logout request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func GetAccessCookie(ctx context.Context, r *http.Request) (*http.Cookie, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	cookie, err := r.Cookie(TokenName)
	if err != nil {
		log.Warn("failed to get token cookie", zap.Error(err))
	}
	return cookie, err
}

func GetRefreshCookie(ctx context.Context, r *http.Request) (*http.Cookie, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	cookie, err := r.Cookie(RefreshTokenName)
	if err != nil {
		log.Warn("failed to get refresh token cookie", zap.Error(err))
	}
	return cookie, err
}

func WriteAuthCookies(w http.ResponseWriter, accessToken string, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(AccessTokenExpirationTime),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    refreshToken,
		Path:     "/auth/",
		Expires:  time.Now().Add(RefreshTokenExpirationTime),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().AddDate(0, -1, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/auth/",
		Expires:  time.Now().AddDate(0, -1, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
