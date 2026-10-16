package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"

	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
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
	fileserver      fileserverpb.FileServerClient
}

const maxAvatarSize = 5 * 1024 * 1024

func NewUserHandler(userUsecase usecase.UserUseCase, settings usecase.SettingsUseCase, fileserverClient fileserverpb.FileServerClient) *UserHandler {
	return &UserHandler{
		userUsecase:     userUsecase,
		settingsUsecase: settings,
		fileserver:      fileserverClient,
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
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize)
	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		log.Warn("Failed to read avatar from multipart form", zap.Error(err))
		response := web_helpers.NewValidationErrorResponse(requestId, []web_helpers.ValidationError{
			web_helpers.NewValidationError("avatar", "Файл аватара не найден или превышает допустимый размер (5 МБ)"),
		})
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Error("failed to close file", zap.Error(err))
		}
	}()
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		response := web_helpers.NewValidationErrorResponse(requestId, []web_helpers.ValidationError{
			web_helpers.NewValidationError("avatar", "Файл должен иметь расширение (например, .jpg, .png)"),
		})
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error("Failed to read file bytes", zap.Error(err))
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	settings, err := obj.settingsUsecase.GetByUser(r.Context(), domain.User{Id: userId})
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	fileserverResponse, err := obj.fileserver.SaveFile(r.Context(), &fileserverpb.SaveFileRequest{
		Data:      fileBytes,
		Extension: ext,
	})
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	err = obj.settingsUsecase.UpdateAvatar(r.Context(), settings.Id, fileserverResponse.GetFilename())
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
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
