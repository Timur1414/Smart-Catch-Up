package validators

import (
	"net/mail"
	"slices"
)

func validatePassword(passwordStr string) error {
	password := []rune(passwordStr)
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 72 {
		return ErrPasswordTooLong
	}
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasInvalid := false
	for i := range len(password) {
		switch {
		case 'A' <= password[i] && password[i] <= 'Z':
			hasUpper = true
		case 'a' <= password[i] && password[i] <= 'z':
			hasLower = true
		case '0' <= password[i] && password[i] <= '9':
			hasDigit = true
		default:
			hasInvalid = true
		}
	}
	if !hasUpper {
		return ErrPasswordDoNotContainsUppercase
	}
	if !hasLower {
		return ErrPasswordDoNotContainsLowercase
	}
	if !hasDigit {
		return ErrPasswordDoNotContainsDigits
	}
	if hasInvalid {
		return ErrPasswordContainsInvalidCharacters
	}
	return nil
}

func validateConfirmPassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordsNotEqual
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) == 0 || len(email) >= 255 {
		return ErrEmailInvalidSize
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}
	if addr.Address != email {
		return ErrEmailInvalid
	}
	return nil
}

func validateMinValue(value, min int) error {
	if value < min {
		return ErrLessThenMin
	}
	return nil
}

func validateMaxValue(value, max int) error {
	if value > max {
		return ErrGreaterThanMax
	}
	return nil
}

func validateNotEmpty(value string) error {
	if value == "" {
		return ErrEmpty
	}
	return nil
}

func validateContains[T comparable](value T, allowedValues []T) error {
	if !slices.Contains(allowedValues, value) {
		return ErrNotAllowed
	}
	return nil
}

func validateAllContains[T comparable](values, allowedValues []T) error {
	for _, value := range values {
		err := validateContains(value, allowedValues)
		if err != nil {
			return ErrNotAllowed
		}
	}
	return nil
}
