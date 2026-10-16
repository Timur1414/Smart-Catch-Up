package grpc

import (
	"bytes"
	"context"

	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileServer struct {
	fileserverpb.UnimplementedFileServerServer
	usecase usecase.FileUseCase
}

func NewFileServer(usecase usecase.FileUseCase) *FileServer {
	return &FileServer{usecase: usecase}
}

func (obj *FileServer) SaveFile(ctx context.Context, request *fileserverpb.SaveFileRequest) (*fileserverpb.SaveFileResponse, error) {
	if request.GetData() == nil || request.GetExtension() == "" {
		return nil, status.Error(codes.InvalidArgument, "data is empty")
	}
	reader := bytes.NewReader(request.GetData())
	filename, err := obj.usecase.StoreFile(ctx, reader, request.GetExtension())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &fileserverpb.SaveFileResponse{
		Filename: filename,
	}, nil
}
