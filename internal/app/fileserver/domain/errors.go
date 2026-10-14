package domain

import "errors"

var (
	ErrEmptyName = errors.New("storage: empty filename")
)
