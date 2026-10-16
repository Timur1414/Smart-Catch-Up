package domain

import "errors"

var (
	ErrEmptyName         = errors.New("empty filename")
	ErrFailedToCreateDir = errors.New("failed to create dir")
	ErrInvalidExtension  = errors.New("invalid extension")
)
