package http

import (
	"net/http"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
)

type DigestHandler struct {
	usecase usecase.DigestUseCase
}

func NewDigestHandler(usecase usecase.DigestUseCase) *DigestHandler {
	return &DigestHandler{usecase: usecase}
}

func (obj DigestHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get digest request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	digest := domain.Digest{
		Id:        1,
		UserId:    1,
		BlockId:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	importantBlocks := make([]domain.BlockPart, 0)
	importantBlocks = append(importantBlocks, domain.BlockPart{
		Id:          2,
		BlockId:     1,
		ClusterType: "a",
		Text:        "text",
		Start:       time.Now(),
		End:         time.Now(),
	})
	categories := make([]domain.BlockPart, 0)
	response := NewDigestResponse(requestId, digest, importantBlocks, categories)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
