package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

	blockRepo := repository.NewBlockPostgres(dbPool)
	notificationRepo := repository.NewNotificationPostgres(dbPool)
	digestRepo := repository.NewDigestPostgres(dbPool)
	log.Info("Repository initialized")
	blockUsecase := usecase.NewBlock(blockRepo, notificationRepo, digestRepo)
	log.Info("Usecase initialized")

	interval := time.Minute

	aggregationCtx := context.Background()
	err = blockUsecase.Aggregate(aggregationCtx)
	if err != nil {
		log.Error("Failed to aggregate", zap.Error(err))
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	log.Info("Ticker started")
	for {
		select {
		case <-aggregationCtx.Done():
			log.Info("Context cancelled, terminating aggregator service")
			return
		case <-ticker.C:
			err = blockUsecase.Aggregate(aggregationCtx)
			if err != nil {
				log.Error("Failed to aggregate", zap.Error(err))
			}
		}
	}
}
