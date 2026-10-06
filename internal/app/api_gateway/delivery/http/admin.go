package http

import (
	"encoding/json"
	"net/http"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
)

type AdminHandler struct {
	generatorServer generatorpb.GeneratorClient
}

func NewAdminHandler(server generatorpb.GeneratorClient) *AdminHandler {
	return &AdminHandler{generatorServer: server}
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
		response := web_helpers.NewServerErrorResponse(requestId)
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
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	userIds := make([]int64, len(request.UserIds))
	for i, userId := range request.UserIds {
		userIds[i] = int64(userId)
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
