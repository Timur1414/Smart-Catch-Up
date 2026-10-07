package http

import (
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	web_helpers2 "github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
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
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewNotificationsResponse(requestId, notifications)
	web_helpers2.WriteResponseJSON(w, response.Code, response)
}

func (obj *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get notifications request")
	userId := context_helper.GetUserIdFromContext(r.Context())
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	notifications, err := obj.usecase.GetAllByUser(r.Context(), userId)
	if err != nil {
		response := web_helpers2.NewServerErrorResponse(requestId)
		web_helpers2.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := NewNotificationsResponse(requestId, notifications)
	web_helpers2.WriteResponseJSON(w, response.Code, response)
}
