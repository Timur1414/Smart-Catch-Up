package grpc

import "github.com/Timur1414/Smart-Catch-Up/internal/app/generator/usecase"

type Generator struct {
	usecase usecase.NotificationUseCase
}

func NewGenerator(usecase usecase.NotificationUseCase) *Generator {
	return &Generator{usecase: usecase}
}
