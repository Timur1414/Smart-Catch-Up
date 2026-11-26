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
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/config"
	delivery "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/http"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/middleware"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/secure/csrf"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/repository"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/cache"
	commonConfig "github.com/Timur1414/Smart-Catch-Up/pkg/config"
	"github.com/Timur1414/Smart-Catch-Up/pkg/database"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbPool, err := database.New(ctx, cfg.Api.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}
	defer dbPool.Close()
	log.Info("Connected to Postgres")

	redisClient, err := cache.New(ctx, cfg.Api.Redis)
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

	authConn, err := grpc.NewClient(
		cfg.Api.Targets.AuthTarget,
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
		cfg.Api.Targets.GeneratorTarget,
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
		cfg.Api.Targets.FileserverTarget,
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
	blockRepo := repository.NewBlockPostgres(dbPool)
	cacheRepo := repository.NewCacheRedis(redisClient, cfg)
	log.Info("Repository initialized")
	userUseCase := usecase.NewUser(userRepo)
	settingsUseCase := usecase.NewSettings(settingsRepo)
	notificationUseCase := usecase.NewNotification(notificationRepo, cfg)
	digestUseCase := usecase.NewDigest(digestRepo, cacheRepo)
	blockUseCase := usecase.NewBlock(blockRepo)
	log.Info("UseCase initialized")
	userHandler := delivery.NewUserHandler(userUseCase, settingsUseCase, fileserverClient, notificationUseCase, cfg)
	digestHandler := delivery.NewDigestHandler(digestUseCase, blockUseCase)
	digestExtensionHandler := delivery.NewDigestExtensionHandler(digestUseCase, blockUseCase, cfg)
	adminHandler := delivery.NewAdminHandler(generatorClient, userUseCase, notificationUseCase)
	authHandler := delivery.NewAuthHandler(authClient, cfg)
	notificationsHandler := delivery.NewNotificationHandler(notificationUseCase)
	log.Info("Handler initialized")

	csrfServer, err := csrf.NewCsrf(cfg)
	if err != nil {
		log.Fatal("Failed to initialize CSRF", zap.Error(err))
		return
	}

	webMux := http.NewServeMux()
	webMux.HandleFunc("POST /auth/login", authHandler.Login)
	webMux.HandleFunc("POST /auth/login/2fa", authHandler.Login2FA)
	webMux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	webMux.HandleFunc("POST /auth/register", authHandler.Register)
	webMux.HandleFunc("POST /auth/logout", authHandler.Logout)
	webMux.HandleFunc("GET /api/profile", userHandler.GetProfile)
	webMux.HandleFunc("POST /api/profile", userHandler.UpdateProfile)
	webMux.HandleFunc("POST /api/profile/avatar", userHandler.UpdateProfileAvatar)
	webMux.HandleFunc("POST /api/profile/2fa/setup", authHandler.Setup2FA)
	webMux.HandleFunc("POST /api/profile/2fa/enable", authHandler.Enable2FA)
	webMux.HandleFunc("POST /api/profile/2fa/disable", authHandler.Disable2FA)
	webMux.HandleFunc("GET /api/is_staff", userHandler.IsStaff)
	webMux.HandleFunc("GET /api/digest", digestHandler.Get)
	webMux.HandleFunc("GET /api/notifications", notificationsHandler.Get)
	webMux.HandleFunc("GET /api/notifications/all", notificationsHandler.GetAll)
	webMux.HandleFunc("POST /api/admin/generate", adminHandler.Generate1)
	webMux.HandleFunc("POST /api/admin/generate_n", adminHandler.GenerateN)
	webMux.HandleFunc("GET /api/users_ids", userHandler.GetAllowedIds)
	webMux.HandleFunc("GET /api/notification_types", notificationsHandler.GetAllTypes)
	webHandler := middleware.CSRFMiddleware(webMux, csrfServer, cfg)
	webHandler = middleware.AuthMiddleware(webHandler, authClient, cfg)
	webHandler = middleware.CORSMiddleware(webHandler, cfg)
	webHandler = middleware.PanicMiddleware(webHandler)
	webHandler = middleware.AccessLogMiddleware(webHandler)

	extMux := http.NewServeMux()
	extMux.HandleFunc("GET /api/ext/digest", digestExtensionHandler.Get)
	extHandler := middleware.CORSMiddleware(extMux, cfg)

	rootMux := http.NewServeMux()
	rootMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	rootMux.Handle("/api/", webHandler)
	rootMux.Handle("/auth/", webHandler)
	rootMux.Handle("/api/ext/", extHandler)
	rootHandler := middleware.PanicMiddleware(rootMux)
	rootHandler = middleware.AccessLogMiddleware(rootHandler)

	addr := ":" + cfg.Api.Port
	server := http.Server{
		Addr:         addr,
		Handler:      rootHandler,
		ReadTimeout:  cfg.Api.ReadTimeout,
		WriteTimeout: cfg.Api.WriteTimeout,
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
