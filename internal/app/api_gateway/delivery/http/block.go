package http

import "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"

type BlockHandler struct {
	usecase usecase.BlockUseCase
}

func NewBlockHandler(usecase usecase.BlockUseCase) *BlockHandler {
	return &BlockHandler{usecase: usecase}
}

type BlockPartHandler struct {
	usecase usecase.BlockPartUseCase
}

func NewBlockPartHandler(usecase usecase.BlockPartUseCase) *BlockPartHandler {
	return &BlockPartHandler{usecase: usecase}
}
