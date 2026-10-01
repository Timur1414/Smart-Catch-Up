package grpc

import (
	"context"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
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
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	userId := 1
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePairByUser(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
	user, err := obj.userUsecase.GetByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePairByUser(ctx, user.Id)
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
	accessToken, refreshToken, err := obj.tokenUsecase.CreatePairByRefreshToken(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return &authpb.RefreshResponse{
		UserId:       1,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (obj *AuthServer) Logout(ctx context.Context, request *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	log := logger.GetLogger()
	if request.GetUserId() == 0 || request.GetAccessToken() == "" || request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	err := obj.tokenUsecase.DeleteRefreshToken(ctx, request.GetRefreshToken())
	if err != nil {
		return &authpb.LogoutResponse{
			Success: false,
		}, err
	}
	log.Info("Logout", zap.Int64("user_id", request.GetUserId()))
	return &authpb.LogoutResponse{
		Success: true,
	}, nil
}
