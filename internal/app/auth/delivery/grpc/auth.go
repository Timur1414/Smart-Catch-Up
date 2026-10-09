package grpc

import (
	"context"
	"errors"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/usecase"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	userUsecase  usecase.UserUseCase
	tokenUsecase usecase.JwtUseCase
}

func NewAuthServer(userUseCase usecase.UserUseCase, tokenUseCase usecase.JwtUseCase) *AuthServer {
	return &AuthServer{
		userUsecase:  userUseCase,
		tokenUsecase: tokenUseCase,
	}
}

func (obj *AuthServer) Register(ctx context.Context, request *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" || request.GetConfirmPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	userId, err := obj.userUsecase.Create(ctx, domain.User{Email: request.GetEmail(), Password: request.GetPassword()})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) || errors.Is(err, domain.ErrSettingsAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, domain.ErrFailedToHashPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePair(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &authpb.RegisterResponse{
		UserId:       int64(userId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (obj *AuthServer) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	userRequest := domain.User{Email: request.GetEmail(), Password: request.GetPassword()}
	user, err := obj.userUsecase.IsExists(ctx, userRequest)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if user.TotpEnabled {
		tempToken, err := obj.tokenUsecase.CreateTempToken(ctx, user.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &authpb.LoginResponse{
			UserId:       int64(user.Id),
			AccessToken:  "",
			RefreshToken: "",
			MfaRequired:  true,
			TempToken:    tempToken,
		}, nil
	}
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePair(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &authpb.LoginResponse{
		UserId:       int64(user.Id),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (obj *AuthServer) Refresh(ctx context.Context, request *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	if request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	ok, userId := obj.tokenUsecase.CheckRefreshToken(ctx, request.GetRefreshToken())
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePair(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &authpb.RefreshResponse{
		UserId:       int64(userId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (obj *AuthServer) Logout(ctx context.Context, request *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	if request.GetUserId() == 0 || request.GetAccessToken() == "" || request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	err := obj.tokenUsecase.DeleteRefreshToken(ctx, request.GetRefreshToken())
	if err != nil {
		return &authpb.LogoutResponse{
			Success: false,
		}, status.Error(codes.Unauthenticated, err.Error())
	}
	log.Info("Logout", zap.Int64("user_id", request.GetUserId()))
	return &authpb.LogoutResponse{
		Success: true,
	}, nil
}

func (obj *AuthServer) IsAuth(ctx context.Context, request *authpb.IsAuthRequest) (*authpb.IsAuthResponse, error) {
	if request.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}
	ok, userId := obj.tokenUsecase.CheckAccessToken(ctx, request.GetAccessToken())
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	return &authpb.IsAuthResponse{
		IsAuth: true,
		UserId: int64(userId),
	}, nil
}

func (obj *AuthServer) Login2FA(ctx context.Context, request *authpb.Login2FARequest) (*authpb.Login2FAResponse, error) {
	if request.GetTempToken() == "" || request.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "code and temp token are required")
	}
	ok, userId := obj.tokenUsecase.CheckTempToken(ctx, request.GetTempToken())
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid temp token")
	}
	validCode, err := obj.userUsecase.Verify2FA(ctx, userId, request.GetCode())
	if err != nil || !validCode {
		attempts, err := obj.tokenUsecase.IncrementTempTokenAttempts(ctx, userId)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if attempts >= 3 {
			err = obj.tokenUsecase.DeleteTempToken(ctx, request.GetTempToken())
			return nil, status.Error(codes.ResourceExhausted, "too many failed attempts, please login again")
		}
		return nil, status.Error(codes.Unauthenticated, "invalid code")
	}
	err = obj.tokenUsecase.DeleteTempToken(ctx, request.GetTempToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePair(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create token pair")
	}
	return &authpb.Login2FAResponse{
		UserId:       int64(userId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (obj *AuthServer) Setup2FA(ctx context.Context, request *authpb.Setup2FARequest) (*authpb.Setup2FAResponse, error) {
	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	secret, otpauthUrl, err := obj.userUsecase.Setup2FA(ctx, int(request.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authpb.Setup2FAResponse{
		Secret:     secret,
		OtpauthUrl: otpauthUrl,
	}, nil
}

func (obj *AuthServer) Enable2FA(ctx context.Context, request *authpb.Enable2FARequest) (*authpb.Enable2FAResponse, error) {
	if request.GetUserId() == 0 || request.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id and code are required")
	}
	backupCodes, err := obj.userUsecase.Enable2FA(ctx, int(request.GetUserId()), request.GetCode())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTotpCode) {
			return nil, status.Error(codes.InvalidArgument, "invalid confirmation code")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authpb.Enable2FAResponse{
		Success:     true,
		BackupCodes: backupCodes,
	}, nil
}

func (obj *AuthServer) Disable2FA(ctx context.Context, request *authpb.Disable2FARequest) (*authpb.Disable2FAResponse, error) {
	if request.GetUserId() == 0 || request.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id and code are required")
	}
	err := obj.userUsecase.Disable2FA(ctx, int(request.GetUserId()), request.GetCode())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTotpCode) {
			return nil, status.Error(codes.InvalidArgument, "invalid confirmation code")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authpb.Disable2FAResponse{
		Success: true,
	}, nil
}
