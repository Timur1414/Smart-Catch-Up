package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	config "github.com/Timur1414/Smart-Catch-Up/internal/app/generator/config"
	generatorDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/generator/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/usecase"
	commonConfig "github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
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

	DEBUG := cfg.Common.Debug

	err = logger.InitLogger(DEBUG)
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

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbPool, err := database.New(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}
	defer dbPool.Close()
	log.Info("Connected to Postgres")

	notificationRepo := repository.NewNotificationPostgres(dbPool)
	userRepo := repository.NewUserPostgres(dbPool)
	settingsRepo := repository.NewSettingsPostgres(dbPool)
	log.Info("Repository initialized")
	notificationUsecase := usecase.NewNotification(notificationRepo, userRepo, settingsRepo)
	log.Info("UseCase initialized")
	generatorServer := generatorDelivery.NewGeneratorServer(notificationUsecase)
	log.Info("Generator gRPC server initialized")

	grpcPort := cfg.Generator.Port
	var listenConfig net.ListenConfig
	lis, err := listenConfig.Listen(ctx, "tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.String("port", grpcPort), zap.Error(err))
	}
	server := grpc.NewServer()

	generatorpb.RegisterGeneratorServer(server, generatorServer)
	log.Info("Generator gRPC server started", zap.String("port", grpcPort))
	if err = server.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
