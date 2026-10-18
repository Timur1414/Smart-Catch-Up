package csrf

import "errors"

var (
	ErrCsrfSecret           = errors.New("CSRF secret is empty")
	ErrInvalidCsrf          = errors.New("CSRF token invalid")
	ErrInvalidCsrfSignature = errors.New("CSRF token signature invalid")
)
