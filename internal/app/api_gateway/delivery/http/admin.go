package http

import (
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
)

type AdminHandler struct {
	usecase usecase.NotificationUseCase
}

func NewAdminHandler(usecase usecase.NotificationUseCase) *AdminHandler {
	return &AdminHandler{usecase}
}

func (obj AdminHandler) Generate1(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Generate 1 request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj AdminHandler) GenerateN(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Generate N request")
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
