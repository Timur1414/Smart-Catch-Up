package http

import "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"

type NotificationHandler struct {
	usecase usecase.NotificationUseCase
}

func NewNotificationHandler(usecase usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{usecase: usecase}
}
