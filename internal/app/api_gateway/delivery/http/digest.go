package http

import (
	"net/http"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/config"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
	"go.uber.org/zap"
)

type DigestHandler struct {
	usecase      usecase.DigestUseCase
	blockUsecase usecase.BlockUseCase
}

const importantType = "important"

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
	cachedDigest, err := obj.usecase.GetFromCache(r.Context(), user)
	if err == nil {
		importantBlocks := make([]domain.BlockPart, 0)
		baseBlocks := make([]domain.BlockPart, 0)
		for _, part := range cachedDigest.Categories {
			blockPart := domain.BlockPart{
				ClusterType: part.ClusterType,
				Text:        part.Text,
				Start:       part.Start,
				End:         part.End,
				Rank:        part.Rank,
			}
			if part.ClusterType == importantType {
				importantBlocks = append(importantBlocks, blockPart)
			} else {
				baseBlocks = append(baseBlocks, blockPart)
			}
		}
		response := NewDigestResponse(requestId, digest, importantBlocks, baseBlocks)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	log.Warn("Cache miss")
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
	baseBlocks := make([]domain.BlockPart, 0)
	for _, part := range block.Parts {
		if part.ClusterType == "important" {
			importantBlocks = append(importantBlocks, part)
		} else {
			baseBlocks = append(baseBlocks, part)
		}
	}
	response := NewDigestResponse(requestId, digest, importantBlocks, baseBlocks)
	web_helpers.WriteResponseJSON(w, response.Code, response)
	err = obj.usecase.UpdateCache(r.Context(), user, block)
	if err != nil {
		log.Warn("failed to update cache", zap.Error(err))
	}
}

type DigestExtensionHandler struct {
	usecase      usecase.DigestUseCase
	blockUsecase usecase.BlockUseCase
	cfg          config.ApiConfig
}

func NewDigestExtensionHandler(usecase usecase.DigestUseCase, blockUsecase usecase.BlockUseCase, cfg config.ApiConfig) *DigestExtensionHandler {
	return &DigestExtensionHandler{
		usecase:      usecase,
		blockUsecase: blockUsecase,
		cfg:          cfg,
	}
}

func (obj *DigestExtensionHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get digest request")
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	user := domain.User{Id: obj.cfg.ExtensionUserId}
	digest, err := obj.usecase.GetByUser(r.Context(), user)
	if err != nil {
		response := web_helpers.NewServerErrorResponse(requestId)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	cachedDigest, err := obj.usecase.GetFromCache(r.Context(), user)
	if err == nil {
		importantBlocks := make([]domain.BlockPart, 0)
		baseBlocks := make([]domain.BlockPart, 0)
		for _, part := range cachedDigest.Categories {
			blockPart := domain.BlockPart{
				ClusterType: part.ClusterType,
				Text:        part.Text,
				Start:       part.Start,
				End:         part.End,
				Rank:        part.Rank,
			}
			if part.ClusterType == importantType {
				importantBlocks = append(importantBlocks, blockPart)
			} else {
				baseBlocks = append(baseBlocks, blockPart)
			}
		}
		response := NewDigestResponse(requestId, digest, importantBlocks, baseBlocks)
		web_helpers.WriteResponseJSON(w, response.Code, response)
		return
	}
	log.Warn("Cache miss")
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
	baseBlocks := make([]domain.BlockPart, 0)
	for _, part := range block.Parts {
		if part.ClusterType == "important" {
			importantBlocks = append(importantBlocks, part)
		} else {
			baseBlocks = append(baseBlocks, part)
		}
	}
	response := NewDigestResponse(requestId, digest, importantBlocks, baseBlocks)
	web_helpers.WriteResponseJSON(w, response.Code, response)
	err = obj.usecase.UpdateCache(r.Context(), user, block)
	if err != nil {
		log.Warn("failed to update cache", zap.Error(err))
	}
}
