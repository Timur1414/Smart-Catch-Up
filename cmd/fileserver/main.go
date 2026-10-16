package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
	fileServerDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file:", err)
			return
		}
	}
	DEBUG := os.Getenv("DEBUG") == "true"

	err := logger.InitLogger(DEBUG)
	if err != nil {
		fmt.Println("Error initializing logger: ", err)
		return
	}
	defer func() {
		err = logger.Close()
		if err != nil {
			fmt.Println("Error closing logger: ", err)
		}
	}()
	log := logger.GetLogger()

	fileRepo, err := repository.NewFileLocalStore()
	if err != nil {
		log.Fatal("failed to create file repository", zap.Error(err))
		return
	}
	log.Info("Repository initialized")
	fileUsecase := usecase.NewFileUseCaseImplementation(fileRepo)
	log.Info("UseCase initialized")
	fileGrpcServer := fileServerDelivery.NewFileServer(fileUsecase)
	log.Info("FileServer gRPC server initialized")
	grpcPort := os.Getenv("FILESERVER_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50053"
	}
	var config net.ListenConfig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lis, err := config.Listen(ctx, "tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.String("port", grpcPort), zap.Error(err))
	}
	server := grpc.NewServer()

	fileserverpb.RegisterFileServerServer(server, fileGrpcServer)
	log.Info("FileServer gRPC server started", zap.String("port", grpcPort))
	if err = server.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
