package grpc

import (
	"context"
	"errors"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/usecase"
)

type GeneratorServer struct {
	generatorpb.UnimplementedGeneratorServer
	usecase usecase.NotificationUseCase
}

func NewGeneratorServer(usecase usecase.NotificationUseCase) *GeneratorServer {
	return &GeneratorServer{usecase: usecase}
}

func (obj *GeneratorServer) Generate1(ctx context.Context, request *generatorpb.Generate1Request) (*generatorpb.Generate1Response, error) {
	if request.GetNotificationType() == "" || request.GetText() == "" || request.GetUserId() == 0 {
		return nil, errors.New("invalid parameters")
	}
	_, err := obj.usecase.Create(ctx, request.GetNotificationType(), request.GetText(), int(request.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &generatorpb.Generate1Response{
		Status: true,
	}, nil
}

func (obj *GeneratorServer) GenerateN(ctx context.Context, request *generatorpb.GenerateNRequest) (*generatorpb.GenerateNResponse, error) {
	if request.GetNotificationTypes() == nil || request.GetNumber() == 0 || request.GetUserIds() == nil {
		return nil, errors.New("invalid parameters")
	}
	userIds := make([]int, len(request.GetUserIds()))
	for i := 0; i < len(request.GetUserIds()); i++ {
		userIds[i] = int(request.GetUserIds()[i])
	}
	err := obj.usecase.BulkCreate(ctx, request.GetNotificationTypes(), int(request.GetNumber()), userIds)
	if err != nil {
		return nil, err
	}
	return &generatorpb.GenerateNResponse{
		Status: true,
	}, nil
}
