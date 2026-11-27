package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/config"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/aggregator/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/cache"
	commonConfig "github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
	"github.com/Timur1414/Smart-Catch-Up/pkg/llm"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
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
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../../configs/config.yaml"
	}
	cfg, err := commonConfig.Load[config.Config](configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = logger.InitLogger(cfg.Common.Debug, cfg.Aggregator.LogPrefix)
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
	dbPool, err := database.New(postgresCtx, cfg.Aggregator.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}
	defer dbPool.Close()
	log.Info("Connected to Postgres")

	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer redisCancel()
	redisClient, err := cache.New(redisCtx, cfg.Aggregator.Redis)
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

	webClient := &http.Client{}
	if cfg.Aggregator.Proxy.On {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		proxyUrl := &url.URL{
			Scheme: cfg.Aggregator.Proxy.Scheme,
			User:   url.UserPassword(cfg.Aggregator.Proxy.Login, cfg.Aggregator.Proxy.Password),
			Host:   cfg.Aggregator.Proxy.Host,
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
		webClient.Transport = transport
	}
	openRouterClient, err := llm.NewClient(cfg.Aggregator.Llm, webClient)
	if err != nil {
		log.Fatal("Failed to initialize open router client", zap.Error(err))
		return
	}

	blockRepo := repository.NewBlockPostgres(dbPool)
	notificationRepo := repository.NewNotificationPostgres(dbPool)
	digestRepo := repository.NewDigestPostgres(dbPool)
	userRepo := repository.NewUserPostgres(dbPool)
	settingsRepo := repository.NewSettingsPostgres(dbPool)
	cacheRepo := repository.NewCacheRedis(redisClient, cfg)
	log.Info("Repository initialized")
	blockUsecase := usecase.NewBlock(blockRepo, notificationRepo, digestRepo, settingsRepo, openRouterClient, cacheRepo)
	log.Info("Usecase initialized")

	batchSize := cfg.Aggregator.BatchSize
	workerCount := cfg.Aggregator.WorkerCount
	interval := cfg.Aggregator.Interval
	aggregationCtx := context.Background()

	for {
		maxUserID, err := userRepo.GetMaxId(aggregationCtx)
		if err != nil {
			log.Fatal("Failed to get max user id", zap.Error(err))
			return
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
