package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/config"
	authDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/auth/delivery/grpc"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/cache"
	commonConfig "github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/Timur1414/Smart-Catch-Up/pkg/crypto"
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

	err = logger.InitLogger(cfg.Common.Debug, cfg.Auth.LogPrefix)
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

	postgresCtx, postgresCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer postgresCancel()
	dbPool, err := database.New(postgresCtx, cfg.Auth.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}
	defer dbPool.Close()
	log.Info("Connected to Postgres")

	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer redisCancel()
	redisClient, err := cache.New(redisCtx, cfg.Auth.Redis)
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
	jwtRepo := repository.NewRefreshTokenRedis(redisClient, cfg)
	twoFactorRepo := repository.NewTwoFactorTokenRedis(redisClient, cfg)
	log.Info("Repository initialized")
	encryptor, err := crypto.NewEncryptor([]byte(cfg.Auth.Encryptor.Secret))
	if err != nil {
		log.Fatal("Failed to initialize encryptor", zap.Error(err))
		return
	}
	userUseCase := usecase.NewUser(userRepo, encryptor, cfg)
	jwtUseCase, err := usecase.NewJwt(jwtRepo, twoFactorRepo, cfg)
	if err != nil {
		log.Fatal("Failed to initialize jwt", zap.Error(err))
	}
	log.Info("UseCase initialized")
	authGrpcServer := authDelivery.NewAuthServer(userUseCase, jwtUseCase)
	log.Info("Auth gRPC server initialized")

	grpcPort := cfg.Auth.Port
	var listenConfig net.ListenConfig
	lisCtx, lisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer lisCancel()
	lis, err := listenConfig.Listen(lisCtx, "tcp", ":"+grpcPort)
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
