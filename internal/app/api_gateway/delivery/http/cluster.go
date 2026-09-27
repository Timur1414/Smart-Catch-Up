package http

import "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"

type ClusterHandler struct {
	usecase usecase.ClusterUseCase
}

func NewClusterHandler(usecase usecase.ClusterUseCase) *ClusterHandler {
	return &ClusterHandler{usecase: usecase}
}
