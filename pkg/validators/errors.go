package validators

import "errors"

var (
	ErrPasswordTooShort                  = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong                   = errors.New("password must be less than 72 characters")
	ErrPasswordDoNotContainsUppercase    = errors.New("password must contain uppercase letters")
	ErrPasswordDoNotContainsLowercase    = errors.New("password must contain lowercase letters")
	ErrPasswordDoNotContainsDigits       = errors.New("password must contain digits")
	ErrPasswordContainsInvalidCharacters = errors.New("password contain invalid characters")
	ErrPasswordsNotEqual                 = errors.New("passwords not equal")
)

var (
	ErrEmailInvalidSize = errors.New("email address must be between 0 and 255 characters")
	ErrEmailInvalid     = errors.New("email address must be a valid email address")
)

var (
	ErrLessThenMin    = errors.New("less than minimum value")
	ErrGreaterThanMax = errors.New("greater than maximum value")
)

var (
	ErrEmpty = errors.New("value must not be empty")
)

var (
	ErrNotAllowed = errors.New("value not allowed")
)
