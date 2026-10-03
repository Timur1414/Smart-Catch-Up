package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	generatorDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/generator/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/generator/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
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
	log.Info("Repository initialized")
	notificationUsecase := usecase.NewNotification(notificationRepo)
	log.Info("UseCase initialized")
	generatorServer := generatorDelivery.NewGeneratorServer(notificationUsecase)
	log.Info("Generator gRPC server initialized")

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}
	var config net.ListenConfig
	lis, err := config.Listen(ctx, "tcp", ":"+grpcPort)
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
