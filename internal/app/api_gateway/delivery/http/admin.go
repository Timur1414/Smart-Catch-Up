package http

import (
	"net/http"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
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
	generatorResponse, err := obj.generatorServer.Generate1(r.Context(), &generatorpb.Generate1Request{
		NotificationType: "abc",
		Text:             "text",
		UserId:           int64(1),
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
	generatorResponse, err := obj.generatorServer.GenerateN(r.Context(), &generatorpb.GenerateNRequest{
		NotificationTypes: []string{"abc"},
		Number:            int64(2),
		UserIds:           []int64{int64(1), int64(2)},
	})
	if err != nil || !generatorResponse.GetStatus() {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	response := web_helpers.NewOkResponse()
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
