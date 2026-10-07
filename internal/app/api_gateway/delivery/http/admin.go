package http

import (
	"context"
	"encoding/json"
	"net/http"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/validators"
	web_helpers2 "github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
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
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	allowedTypes, allowedUsers, err := getAllowedAdminData(r.Context(), obj.notificationUsecase, obj.userUsecase)
	if err != nil {
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	errors := validators.ValidateGenerate1(request.Text, request.NotificationType, request.UserId, allowedTypes, allowedUsers)
	if len(errors) > 0 {
		log.Warn("validation errors", zap.Any("errors", errors))
		response := web_helpers2.NewValidationErrorResponse(requestId, errors)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	generatorResponse, err := obj.generatorServer.Generate1(r.Context(), &generatorpb.Generate1Request{
		NotificationType: request.NotificationType,
		Text:             request.Text,
		UserId:           int64(request.UserId),
	})
	if err != nil || !generatorResponse.GetStatus() {
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers2.NewOkResponse()
	web_helpers2.WriteResponseJSON(w, response.Code, response)
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
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	userIds := make([]int64, len(request.UserIds))
	for i, userId := range request.UserIds {
		userIds[i] = int64(userId)
	}
	allowedTypes, allowedUsers, err := getAllowedAdminData(r.Context(), obj.notificationUsecase, obj.userUsecase)
	if err != nil {
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	errors := validators.ValidateGenerateN(request.Number, request.NotificationTypes, request.UserIds, allowedTypes, allowedUsers)
	if len(errors) > 0 {
		log.Warn("validation errors", zap.Any("errors", errors))
		response := web_helpers2.NewValidationErrorResponse(requestId, errors)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	generatorResponse, err := obj.generatorServer.GenerateN(r.Context(), &generatorpb.GenerateNRequest{
		NotificationTypes: request.NotificationTypes,
		Number:            int64(request.Number),
		UserIds:           userIds,
	})
	if err != nil || !generatorResponse.GetStatus() {
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers2.NewOkResponse()
	web_helpers2.WriteResponseJSON(w, response.Code, response)
}

func getAllowedAdminData(ctx context.Context, notificationUsecase usecase.NotificationUseCase, userUsecase usecase.UserUseCase) ([]string, []int, error) {
	allowedTypes, err := notificationUsecase.GetAllTypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	allowedUsers, err := userUsecase.GetAllIds(ctx)
	if err != nil {
		return nil, nil, err
	}
	return allowedTypes, allowedUsers, nil
}
