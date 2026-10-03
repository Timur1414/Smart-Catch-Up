package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	delivery "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/http"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/internal/middleware"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	err = logger.InitAccessLogger()
	if err != nil {
		log.Error("Error initializing access logger", zap.Error(err))
		return
	}
	defer func() {
		err = logger.AccessClose()
		if err != nil {
			fmt.Println("Error closing access logger: ", err)
		}
	}()

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

	authConn, err := grpc.NewClient(
		"auth:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{Timeout: 10 * time.Second}
			return d.DialContext(ctx, "tcp4", addr)
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	authClient := authpb.NewAuthClient(authConn)

	generatorConn, err := grpc.NewClient(
		"generator:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{Timeout: 10 * time.Second}
			return d.DialContext(ctx, "tcp4", addr)
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to generator service", zap.Error(err))
	}
	generatorClient := generatorpb.NewGeneratorClient(generatorConn)

	userRepo := repository.NewUserPostgres(dbPool)
	settingsRepo := repository.NewSettingsPostgres(dbPool)
	notificationRepo := repository.NewNotificationPostgres(dbPool)
	digestRepo := repository.NewDigestPostgres(dbPool)
	clusterRepo := repository.NewClusterPostgres(dbPool)
	blockRepo := repository.NewBlockPostgres(dbPool)
	blockPartRepo := repository.NewBlockPartPostgres(dbPool)
	log.Info("Repository initialized")
	userUseCase := usecase.NewUser(userRepo)
	settingsUseCase := usecase.NewSettings(settingsRepo)
	notificationUseCase := usecase.NewNotification(notificationRepo)
	digestUseCase := usecase.NewDigest(digestRepo)
	clusterUseCase := usecase.NewCluster(clusterRepo)
	blockUseCase := usecase.NewBlock(blockRepo)
	blockPartUseCase := usecase.NewBlockPart(blockPartRepo)
	log.Info("UseCase initialized")
	userHandler := delivery.NewUserHandler(userUseCase)
	digestHandler := delivery.NewDigestHandler(digestUseCase)
	adminHandler := delivery.NewAdminHandler(generatorClient)
	authHandler := delivery.NewAuthHandler(authClient)
	log.Info("Handler initialized")
	_ = notificationUseCase
	_ = settingsUseCase
	_ = clusterUseCase
	_ = blockUseCase
	_ = blockPartUseCase

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /profile", userHandler.GetProfile)
	mux.HandleFunc("POST /profile", userHandler.UpdateProfile)
	mux.HandleFunc("POST /profile/avatar", userHandler.UpdateProfileAvatar)
	mux.HandleFunc("GET /is_staff", userHandler.IsStaff)
	mux.HandleFunc("GET /digest", digestHandler.Get)
	mux.HandleFunc("POST /admin/generate", adminHandler.Generate1)
	mux.HandleFunc("POST /admin/generate_n", adminHandler.GenerateN)

	handler := middleware.AuthMiddleware(mux, authClient)
	handler = middleware.AccessLogMiddleware(mux)
	handler = middleware.PanicMiddleware(handler)

	addr := ":" + os.Getenv("PORT")
	server := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	log.Info("starting api server",
		zap.String("addr", addr),
	)
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Error starting server", zap.Error(err))
		return
	}
}
