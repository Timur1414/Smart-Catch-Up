package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/validators"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
	"go.uber.org/zap"
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
	var request UpdateProfileRequest
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
	user := domain.User{Id: userId}
	settings, err := obj.settingsUsecase.GetByUser(r.Context(), user)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	validationErrors := validators.ValidateUpdateProfile(request.Email, request.FirstName, request.LastName)
	if len(validationErrors) > 0 {
		response := web_helpers.NewValidationErrorResponse(requestId, validationErrors)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	user.Email = request.Email
	settings.FirstName = request.FirstName
	settings.LastName = request.LastName
	err = obj.userUsecase.Update(r.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicatedData) {
			response := web_helpers.NewValidationErrorResponse(requestId, []web_helpers.ValidationError{})
			response.Code = http.StatusConflict
			response.Message = err.Error()
			web_helpers.WriteResponseJSON(w, response.Code, response)
			return
		}
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	err = obj.settingsUsecase.Update(r.Context(), settings)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewUpdateProfileResponse(requestId, http.StatusOK, "Ok", user, settings)
	web_helpers.WriteResponseJSON(w, response.Code, response)
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
