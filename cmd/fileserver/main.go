package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/config"
	fileServerDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/fileserver/usecase"
	commonConfig "github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../../configs/config.yaml"
	}
	cfg, err := commonConfig.Load[config.Config](configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = logger.InitLogger(cfg.Common.Debug)
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

	fileRepo, err := repository.NewFileLocalStore(cfg)
	if err != nil {
		log.Fatal("failed to create file repository", zap.Error(err))
		return
	}
	log.Info("Repository initialized")
	fileUsecase := usecase.NewFileUseCaseImplementation(fileRepo)
	log.Info("UseCase initialized")
	fileGrpcServer := fileServerDelivery.NewFileServer(fileUsecase)
	log.Info("FileServer gRPC server initialized")
	grpcPort := cfg.Fileserver.Port
	var listenConfig net.ListenConfig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lis, err := listenConfig.Listen(ctx, "tcp", ":"+grpcPort)
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
