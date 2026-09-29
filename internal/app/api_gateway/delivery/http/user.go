package http

import (
	"net/http"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
	"github.com/Timur1414/Smart-Catch-Up/pkg/context_helper"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(usecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (obj *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Get profile request")
	user := domain.User{
		Id:        1,
		IsStaff:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     "abc@example.com",
		Password:  "abc",
		VkId:      1,
	}
	settings := domain.Settings{
		Id:        1,
		UserId:    1,
		Interval:  10,
		Important: nil,
		AvatarUrl: "/base.png",
		FirstName: "John",
		LastName:  "Doe",
		UpdatedAt: time.Now(),
	}
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	response := NewUserResponse(requestId, user, settings)
	web_helpers.WriteResponseJSON(w, response.Code, response)
}

func (obj *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Update profile request")
}

func (obj *UserHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Update profile avatar")
}

func (obj *UserHandler) IsStaff(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerWithRequestId(r.Context())
	log.Info("Is staff request")
	user := domain.User{
		Id:        1,
		IsStaff:   true,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Email:     "",
		Password:  "",
		VkId:      0,
	}
	requestId := context_helper.GetRequestIdFromContext(r.Context())
	request := NewUserRoleResponse(requestId, user)
	web_helpers.WriteResponseJSON(w, request.Code, request)
}
