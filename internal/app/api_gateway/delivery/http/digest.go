package http

import (
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
)

type DigestHandler struct {
	usecase      usecase.DigestUseCase
	blockUsecase usecase.BlockUseCase
}

func NewDigestHandler(usecase usecase.DigestUseCase, blockUsecase usecase.BlockUseCase) *DigestHandler {
	return &DigestHandler{
		usecase:      usecase,
		blockUsecase: blockUsecase,
	}
}

func (obj *DigestHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get digest request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	userId := context_helper.GetUserIdFromContext(r.Context())
	user := domain.User{Id: userId}
	digest, err := obj.usecase.GetByUser(r.Context(), user)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	if digest.BlockId == 0 {
		response := NewDigestResponse(requestId, digest, []domain.BlockPart{}, []domain.BlockPart{})
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	block, err := obj.blockUsecase.GetById(r.Context(), digest.BlockId)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	importantBlocks := make([]domain.BlockPart, 0)
	response := NewDigestResponse(requestId, digest, importantBlocks, block.Parts)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}
