package http

import (
	"encoding/json"
	"net/http"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/config"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	usecase usecase.NotificationUseCase
}

func NewNotificationHandler(usecase usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{usecase: usecase}
}

func (obj *NotificationHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get notifications request")
	userId := context_helper.GetUserIdFromContext(r.Context())
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	notifications, err := obj.usecase.GetByUser(r.Context(), userId)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewNotificationsResponse(requestId, notifications)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get notifications request")
	userId := context_helper.GetUserIdFromContext(r.Context())
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	notifications, err := obj.usecase.GetAllByUser(r.Context(), userId)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewNotificationsResponse(requestId, notifications)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *NotificationHandler) GetAllTypes(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get notification types request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	types, err := obj.usecase.GetAllTypes(r.Context())
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewAllowedNotificationTypesResponse(requestId, types)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

type NotificationExtensionHandler struct {
	generatorServer generatorpb.GeneratorClient
	cfg             *config.Config
}

func NewNotificationExtensionHandler(server generatorpb.GeneratorClient, cfg *config.Config) *NotificationExtensionHandler {
	return &NotificationExtensionHandler{
		generatorServer: server,
		cfg:             cfg,
	}
}

func (obj *NotificationExtensionHandler) Receive(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Receive notifications from extension")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := obj.cfg.Api.ExtensionUserId
	var request ExtensionNotificationsRequest
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
	for _, notification := range request.Notifications {
		generatorResponse, err := obj.generatorServer.Generate1(r.Context(), &generatorpb.Generate1Request{
			NotificationType: "no_type",
			Text:             notification.Payload,
			UserId:           int64(userId),
		})
		if err != nil || !generatorResponse.GetStatus() {
			response := web_helpers.NewServerErrorResponse(requestId)
			web_helpers.WriteResponseJSON(w, response.Code, response)
			return
		}
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
