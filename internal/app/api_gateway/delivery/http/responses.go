package http

import (
	"net/http"
	"time"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"
)

type LoginSuccessResponse web_helpers.UuidResponse

func NewLoginSuccessResponse(requestId string) LoginSuccessResponse {
	return LoginSuccessResponse{
		SimpleResponse: web_helpers.SimpleResponse{
			Code:    http.StatusOK,
			Message: "Ok",
		},
		RequestId: requestId,
	}
}

type LoginErrorResponse web_helpers.ValidationErrorResponse

func NewLoginErrorResponse(requestId string, code int, message string, errors []web_helpers.ValidationError) LoginErrorResponse {
	return LoginErrorResponse{
		UuidResponse: web_helpers.NewUuidResponse(code, message, requestId),
		Errors:       errors,
	}
}

type RegisterSuccessResponse web_helpers.UuidResponse

func NewRegisterSuccessResponse(requestId string) RegisterSuccessResponse {
	return RegisterSuccessResponse{
		SimpleResponse: web_helpers.SimpleResponse{
			Code:    http.StatusOK,
			Message: "Ok",
		},
		RequestId: requestId,
	}
}

type RegisterErrorResponse web_helpers.ValidationErrorResponse

func NewRegisterErrorResponse(requestId string, code int, message string, errors []web_helpers.ValidationError) RegisterErrorResponse {
	return RegisterErrorResponse{
		UuidResponse: web_helpers.NewUuidResponse(code, message, requestId),
		Errors:       errors,
	}
}

type LogoutSuccessResponse web_helpers.UuidResponse

func NewLogoutSuccessResponse(requestId string) LogoutSuccessResponse {
	return LogoutSuccessResponse{
		SimpleResponse: web_helpers.SimpleResponse{
			Code:    http.StatusOK,
			Message: "Ok",
		},
		RequestId: requestId,
	}
}

type UserResponse struct {
	web_helpers.UuidResponse
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	AvatarUrl string    `json:"avatar_url"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsStaff   bool      `json:"is_staff"`
}

func NewUserResponse(requestId string, user domain.User, settings domain.Settings) UserResponse {
	return UserResponse{
		UuidResponse: web_helpers.NewUuidResponse(http.StatusOK, "Ok", requestId),
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		AvatarUrl:    settings.AvatarUrl,
		FirstName:    settings.FirstName,
		LastName:     settings.LastName,
		IsStaff:      user.IsStaff,
	}
}

type UpdateProfileResponse UserResponse

func NewUpdateProfileResponse(requestId string, code int, message string, user domain.User, settings domain.Settings) UpdateProfileResponse {
	return UpdateProfileResponse{
		UuidResponse: web_helpers.NewUuidResponse(code, message, requestId),
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		AvatarUrl:    settings.AvatarUrl,
		FirstName:    settings.FirstName,
		LastName:     settings.LastName,
		IsStaff:      user.IsStaff,
	}
}

type UserRoleResponse struct {
	web_helpers.UuidResponse
	IsStaff bool `json:"is_staff"`
}

func NewUserRoleResponse(requestId string, user domain.User) UserRoleResponse {
	return UserRoleResponse{
		UuidResponse: web_helpers.NewUuidResponse(http.StatusOK, "Ok", requestId),
		IsStaff:      user.IsStaff,
	}
}

type ImportantNotificationResponse struct {
	Notification string    `json:"notification"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewImportantNotificationResponse(notification domain.BlockPart) ImportantNotificationResponse {
	return ImportantNotificationResponse{
		Notification: notification.Text,
		Timestamp:    notification.Start,
	}
}

type NotificationCategoryResponse struct {
	CategoryType string `json:"category_type"`
	Text         string `json:"text"`
}

func NewNotificationCategoryResponse(notificationCategory domain.BlockPart) NotificationCategoryResponse {
	return NotificationCategoryResponse{
		CategoryType: notificationCategory.ClusterType,
		Text:         notificationCategory.Text,
	}
}

type DigestResponse struct {
	web_helpers.UuidResponse
	CreatedAt  time.Time                       `json:"created_at"`
	Important  []ImportantNotificationResponse `json:"important"`
	Categories []NotificationCategoryResponse  `json:"categories"`
}

func NewDigestResponse(requestId string, digest domain.Digest, importantBlocks []domain.BlockPart, categoriesBlocks []domain.BlockPart) DigestResponse {
	importantResponses := make([]ImportantNotificationResponse, len(importantBlocks))
	for i, blockPart := range importantBlocks {
		importantResponses[i] = NewImportantNotificationResponse(blockPart)
	}
	categoriesResponses := make([]NotificationCategoryResponse, len(categoriesBlocks))
	for i, blockPart := range categoriesBlocks {
		categoriesResponses[i] = NewNotificationCategoryResponse(blockPart)
	}
	return DigestResponse{
		UuidResponse: web_helpers.NewUuidResponse(http.StatusOK, "Ok", requestId),
		CreatedAt:    digest.CreatedAt,
		Important:    importantResponses,
		Categories:   categoriesResponses,
	}
}
