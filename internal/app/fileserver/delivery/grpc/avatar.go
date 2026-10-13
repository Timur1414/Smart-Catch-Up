package grpc

import (
	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/usecase"
)

type FileServer struct {
	fileserverpb.UnimplementedFileServerServer
	usecase usecase.SettingsUseCase
}

func NewFileServer(usecase usecase.SettingsUseCase) *FileServer {
	return &FileServer{usecase: usecase}
}
