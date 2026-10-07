package domain

import "errors"

var (
	ErrFailedToHashPassword = errors.New("failed to hash password")
)

var (
	ErrSecretTooShort    = errors.New("secret too short")
	ErrVersionIsEmpty    = errors.New("version is empty")
	ErrFailedToSignToken = errors.New("failed to sign token")
)

var (
	ErrNothingInTable        = errors.New("no rows in result set")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrDuplicatedData        = errors.New("duplicated data")
	ErrSettingsAlreadyExists = errors.New("settings already exists")

	ErrNotFound = errors.New("not found")
)
