package http

import (
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
)

type UserHandler struct {
	userUsecase     usecase.UserUseCase
	settingsUsecase usecase.SettingsUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase, settings usecase.SettingsUseCase) *UserHandler {
	return &UserHandler{
		userUsecase:     userUsecase,
		settingsUsecase: settings,
	}
}

func (obj *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get profile request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	user, err := obj.userUsecase.GetById(r.Context(), userId)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	settings, err := obj.settingsUsecase.GetByUser(r.Context(), user)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewUserResponse(requestId, user, settings)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Update profile request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
}

func (obj *UserHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Update profile avatar")
}

func (obj *UserHandler) IsStaff(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Is staff request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	user, err := obj.userUsecase.GetById(r.Context(), userId)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	request := NewUserRoleResponse(requestId, user)
	web_helpers.WriteResponseJSON(w, request.Code, request)
}

func (obj *UserHandler) GetAllowedIds(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get allowed user ids request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	users, settings, err := obj.userUsecase.GetAllShortUsers(r.Context())
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewAllowedUserIdsResponse(requestId, users, settings)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
