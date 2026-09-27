package http

import "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/usecase"

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(usecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

type SettingsHandler struct {
	usecase usecase.SettingsUseCase
}

func NewSettingsHandler(usecase usecase.SettingsUseCase) *SettingsHandler {
	return &SettingsHandler{usecase: usecase}
}
