package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type RangeJob struct {
	StartId int
	EndId   int
}

func RunAggregation(ctx context.Context, blockUsecase usecase.BlockUseCase, maxUserId, batchSize, workerCount int) {
	log := logger.GetLogger()
	log.Info("Starting aggregation cycle", zap.Int("max_user_id", maxUserId))
	jobs := make(chan RangeJob, workerCount*2)
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					err := blockUsecase.Aggregate(ctx, job.StartId, job.EndId)
					if err != nil {
						log.Error("Worker failed to aggregate",
							zap.Int("start_id", job.StartId),
							zap.Int("end_id", job.EndId),
							zap.Error(err))
					}
				}
			}
		})
	}
	for i := 0; i <= maxUserId; i += batchSize {
		endId := i + batchSize
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return
		case jobs <- RangeJob{StartId: i, EndId: endId}:
		}
	}
	close(jobs)
	wg.Wait()
	log.Info("Aggregation cycle finished successfully")
}

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
	userRepo := repository.NewUserPostgres(dbPool)
	settingsRepo := repository.NewSettingsPostgres(dbPool)
	log.Info("Repository initialized")
	blockUsecase := usecase.NewBlock(blockRepo, notificationRepo, digestRepo, settingsRepo)
	log.Info("Usecase initialized")

	batchSize := 500
	workerCount := 4
	interval := time.Minute
	aggregationCtx := context.Background()

	for {
		maxUserID, err := userRepo.GetMaxId(aggregationCtx)
		if err != nil {
			log.Error("Failed to get max user id", zap.Error(err))
		}
		RunAggregation(aggregationCtx, blockUsecase, maxUserID, batchSize, workerCount)
		select {
		case <-aggregationCtx.Done():
			log.Info("Context cancelled, terminating aggregator service")
			return
		case <-time.After(interval):
		}
	}
}
