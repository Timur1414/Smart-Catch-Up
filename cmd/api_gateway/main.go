package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	fileserverpb "github.com/Timur1414/Smart-Catch-Up/api/proto/file_server"
	generatorpb "github.com/Timur1414/Smart-Catch-Up/api/proto/generator"
	delivery "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/http"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/middleware"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/secure/csrf"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
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

	fileserverConn, err := grpc.NewClient(
		"fileserver:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{Timeout: 10 * time.Second}
			return d.DialContext(ctx, "tcp4", addr)
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to fileserver service", zap.Error(err))
	}
	fileserverClient := fileserverpb.NewFileServerClient(fileserverConn)

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
	userHandler := delivery.NewUserHandler(userUseCase, settingsUseCase, fileserverClient, notificationUseCase)
	digestHandler := delivery.NewDigestHandler(digestUseCase, blockUseCase)
	adminHandler := delivery.NewAdminHandler(generatorClient, userUseCase, notificationUseCase)
	authHandler := delivery.NewAuthHandler(authClient)
	notificationsHandler := delivery.NewNotificationHandler(notificationUseCase)
	log.Info("Handler initialized")
	_ = clusterUseCase
	_ = blockPartUseCase

	csrfServer, err := csrf.NewCsrf(os.Getenv("CSRF_SECRET"))
	if err != nil {
		log.Fatal("Failed to initialize CSRF", zap.Error(err))
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/login/2fa", authHandler.Login2FA)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/profile", userHandler.GetProfile)
	mux.HandleFunc("POST /api/profile", userHandler.UpdateProfile)
	mux.HandleFunc("POST /api/profile/avatar", userHandler.UpdateProfileAvatar)
	mux.HandleFunc("POST /api/profile/2fa/setup", authHandler.Setup2FA)
	mux.HandleFunc("POST /api/profile/2fa/enable", authHandler.Enable2FA)
	mux.HandleFunc("POST /api/profile/2fa/disable", authHandler.Disable2FA)
	mux.HandleFunc("GET /api/is_staff", userHandler.IsStaff)
	mux.HandleFunc("GET /api/digest", digestHandler.Get)
	mux.HandleFunc("GET /api/notifications", notificationsHandler.Get)
	mux.HandleFunc("GET /api/notifications/all", notificationsHandler.GetAll)
	mux.HandleFunc("POST /api/admin/generate", adminHandler.Generate1)
	mux.HandleFunc("POST /api/admin/generate_n", adminHandler.GenerateN)
	mux.HandleFunc("GET /api/users_ids", userHandler.GetAllowedIds)
	mux.HandleFunc("GET /api/notification_types", notificationsHandler.GetAllTypes)

	handler := middleware.CSRFMiddleware(mux, csrfServer)
	handler = middleware.AuthMiddleware(handler, authClient)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.AccessLogMiddleware(handler)
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
