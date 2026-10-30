package domain

import "errors"

var (
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrInvalidTotpCode      = errors.New("invalid TOTP code")
	ErrTotpAlreadyExists    = errors.New("TOTP already exists")
)

var (
	ErrSecretTooShort   = errors.New("secret too short")
	ErrVersionIsEmpty   = errors.New("version is empty")
	ErrInvalidTokenType = errors.New("invalid token type")
)

var (
	ErrNothingInTable        = errors.New("no rows in result set")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrDuplicatedData        = errors.New("duplicated data")
	ErrSettingsAlreadyExists = errors.New("settings already exists")
	ErrDigestAlreadyExists   = errors.New("digest already exists")

	ErrNotFound = errors.New("not found")
)
