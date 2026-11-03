package http

import (
	"encoding/json"
	"net/http"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/validators"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
	"go.uber.org/zap"
)

type AdminHandler struct {
	generatorServer     generatorpb.GeneratorClient
	userUsecase         usecase.UserUseCase
	notificationUsecase usecase.NotificationUseCase
}

func NewAdminHandler(server generatorpb.GeneratorClient, userUsecase usecase.UserUseCase, notificationUsecase usecase.NotificationUseCase) *AdminHandler {
	return &AdminHandler{
		generatorServer:     server,
		userUsecase:         userUsecase,
		notificationUsecase: notificationUsecase,
	}
}

func (obj *AdminHandler) Generate1(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Generate 1 request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	var request Generate1Request
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		log.Error("failed to parse request", zap.Error(err))
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	allowedTypes, err := obj.notificationUsecase.GetAllTypes(r.Context())
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	exists, err := obj.userUsecase.IsUserIdsExists(r.Context(), []int{request.UserId})
	if err != nil || !exists {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	errors := validators.ValidateGenerate1(request.Text, request.NotificationType, request.UserId, allowedTypes)
	if len(errors) > 0 {
		log.Warn("validation errors", zap.Any("errors", errors))
		response := web_helpers.NewValidationErrorResponse(requestId, errors)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	generatorResponse, err := obj.generatorServer.Generate1(r.Context(), &generatorpb.Generate1Request{
		NotificationType: request.NotificationType,
		Text:             request.Text,
		UserId:           int64(request.UserId),
	})
	if err != nil || !generatorResponse.GetStatus() {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *AdminHandler) GenerateN(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Generate N request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	var request GenerateNRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Error("failed to close request body", zap.Error(err))
		}
	}()
	if err != nil {
		log.Error("failed to parse request", zap.Error(err))
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	userIds := make([]int64, len(request.UserIds))
	for i, userId := range request.UserIds {
		userIds[i] = int64(userId)
	}
	allowedTypes, err := obj.notificationUsecase.GetAllTypes(r.Context())
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	exists, err := obj.userUsecase.IsUserIdsExists(r.Context(), request.UserIds)
	if err != nil || !exists {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	errors := validators.ValidateGenerateN(request.Number, request.NotificationTypes, request.UserIds, allowedTypes)
	if len(errors) > 0 {
		log.Warn("validation errors", zap.Any("errors", errors))
		response := web_helpers.NewValidationErrorResponse(requestId, errors)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	generatorResponse, err := obj.generatorServer.GenerateN(r.Context(), &generatorpb.GenerateNRequest{
		NotificationTypes: request.NotificationTypes,
		Number:            int64(request.Number),
		UserIds:           userIds,
	})
	if err != nil || !generatorResponse.GetStatus() {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
