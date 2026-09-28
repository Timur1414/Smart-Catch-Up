package grpc

import (
	"context"

	authpb "github.com/Timur1414/Smart-Catch-Up/api/proto/auth"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	userUsecase  usecase.UserUseCase
	tokenUsecase usecase.RefreshTokenUseCase
}

func NewAuthServer(userUseCase usecase.UserUseCase, tokenUseCase usecase.RefreshTokenUseCase) *AuthServer {
	return &AuthServer{
		userUsecase:  userUseCase,
		tokenUsecase: tokenUseCase,
	}
}

func (obj AuthServer) Register(ctx context.Context, request *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	return &authpb.RegisterResponse{
		UserId:       1,
		AccessToken:  "token",
		RefreshToken: "token",
	}, nil
}

func (obj AuthServer) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	user, err := obj.userUsecase.GetByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &authpb.LoginResponse{
		UserId:       int64(user.Id),
		AccessToken:  "token",
		RefreshToken: "token",
	}, nil
}

func (obj AuthServer) Refresh(ctx context.Context, request *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	if request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	return &authpb.RefreshResponse{
		UserId:       1,
		AccessToken:  "token",
		RefreshToken: "token",
	}, nil
}

func (obj AuthServer) Logout(ctx context.Context, request *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if request.GetUserId() == 0 || request.GetAccessToken() == "" || request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	return &authpb.LogoutResponse{
		Success: true,
	}, nil
}
