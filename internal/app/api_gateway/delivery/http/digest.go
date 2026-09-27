package http

import "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"

type DigestHandler struct {
	usecase usecase.DigestUseCase
}

func NewDigestHandler(usecase usecase.DigestUseCase) *DigestHandler {
	return &DigestHandler{usecase: usecase}
}
