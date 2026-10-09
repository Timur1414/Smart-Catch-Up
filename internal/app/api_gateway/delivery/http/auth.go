package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/validators"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AccessTokenName            = "access_token"
	RefreshTokenName           = "refresh_token"
	AccessTokenExpirationTime  = time.Minute * 5
	RefreshTokenExpirationTime = time.Hour * 24 * 7
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
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		log.Error("failed to decode request", zap.Error(err))
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Login(r.Context(), &authpb.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	if authResponse.GetMfaRequired() {
		response := NewLoginMFAResponse(requestId, authResponse.GetTempToken())
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	WriteAuthCookies(w, authResponse.GetAccessToken(), authResponse.GetRefreshToken())
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	log.Info("Refresh request")
	refreshToken, err := GetRefreshCookie(r.Context(), r)
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Refresh(r.Context(), &authpb.RefreshRequest{
		RefreshToken: refreshToken.Value,
	})
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	WriteAuthCookies(w, authResponse.GetAccessToken(), authResponse.GetRefreshToken())
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Register request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		log.Error("failed to decode request", zap.Error(err))
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	errors := validators.ValidateRegisterUser(request.Email, request.Password, request.ConfirmPassword)
	if len(errors) > 0 {
		log.Warn("validation errors", zap.Any("errors", errors))
		response := web_helpers.NewValidationErrorResponse(requestId, errors)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Register(r.Context(), &authpb.RegisterRequest{
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	WriteAuthCookies(w, authResponse.GetAccessToken(), authResponse.GetRefreshToken())
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Logout request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	accessToken, err := GetAccessCookie(r.Context(), r)
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	refreshToken, err := GetRefreshCookie(r.Context(), r)
	if err != nil {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Logout(r.Context(), &authpb.LogoutRequest{
		UserId:       int64(userId),
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	})
	if err != nil || !authResponse.GetSuccess() {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	ClearAuthCookies(w)
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Login2FA(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Login2FA request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	var request Login2FARequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Login2FA(r.Context(), &authpb.Login2FARequest{
		TempToken: request.TempToken,
		Code:      request.Code,
	})
	if err != nil {
		if status.Code(err) == codes.ResourceExhausted {
			response := web_helpers.NewTooManyRequestsResponse(requestId)
			web_helpers.WriteResponseJSON(w, response.Code, response)
			return
		}
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	WriteAuthCookies(w, authResponse.GetAccessToken(), authResponse.GetRefreshToken())
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Setup2FA request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	authResponse, err := obj.authServer.Setup2FA(r.Context(), &authpb.Setup2FARequest{
		UserId: int64(userId),
	})
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewSetup2FAResponse(requestId, authResponse.GetSecret(), authResponse.GetOtpauthUrl())
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Enable2FA(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Enable2FA request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	var request Enable2FARequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Enable2FA(r.Context(), &authpb.Enable2FARequest{
		UserId: int64(userId),
		Code:   request.Code,
	})
	if err != nil || !authResponse.GetSuccess() {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewEnable2FAResponse(requestId, authResponse.GetBackupCodes())
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AuthHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Disable2FA request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	var request Disable2FARequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	authResponse, err := obj.authServer.Disable2FA(r.Context(), &authpb.Disable2FARequest{
		UserId: int64(userId),
		Code:   request.Code,
	})
	if err != nil || !authResponse.GetSuccess() {
		response := web_helpers.NewUnauthorizedResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func GetAccessCookie(ctx context.Context, r *http.Request) (*http.Cookie, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	cookie, err := r.Cookie(AccessTokenName)
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
		Name:     AccessTokenName,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(AccessTokenExpirationTime).UTC(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    refreshToken,
		Path:     "/auth/",
		Expires:  time.Now().Add(RefreshTokenExpirationTime).UTC(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().AddDate(0, -1, 0).UTC(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/auth/",
		Expires:  time.Now().AddDate(0, -1, 0).UTC(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
