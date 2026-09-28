package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	authDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/auth/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/cache"
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
		"postgres://%s:%s@%s:%s/%s",
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

	redisAddr := fmt.Sprintf("redis://:%s@%s:%s/0", os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient, err := cache.New(ctx, redisAddr)
	if err != nil {
		log.Fatal("Failed to connect to redis", zap.Error(err))
	}
	defer func() {
		err = redisClient.Close()
		if err != nil {
			log.Fatal("Failed to close redis connection", zap.Error(err))
		}
	}()
	log.Info("Connected to Redis")

	userRepo := repository.NewUserPostgres(dbPool)
	jwtRepo := repository.NewRefreshTokenRedis(redisClient)
	log.Info("Repository initialized")
	userUseCase := usecase.NewUser(userRepo)
	jwtUseCase := usecase.NewRefreshToken(jwtRepo)
	log.Info("UseCase initialized")
	authGrpcServer := authDelivery.NewAuthServer(userUseCase, jwtUseCase)
	log.Info("Auth gRPC server initialized")

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	var config net.ListenConfig
	lis, err := config.Listen(ctx, "tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.String("port", grpcPort), zap.Error(err))
	}
	server := grpc.NewServer()

	authpb.RegisterAuthServer(server, authGrpcServer)
	log.Info("Auth gRPC server started", zap.String("port", grpcPort))
	if err = server.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
