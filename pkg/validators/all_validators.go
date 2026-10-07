package validators

import (
	"errors"
	"net/mail"
	"slices"
	"strconv"
)

func validatePassword(passwordStr string) error {
	password := []rune(passwordStr)
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return errors.New("password must be less than 72 characters")
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
		return errors.New("password must contain uppercase letters")
	}
	if !hasLower {
		return errors.New("password must contain lowercase letters")
	}
	if !hasDigit {
		return errors.New("password must contain digits")
	}
	if hasInvalid {
		return errors.New("password must contain invalid characters")
	}
	return nil
}

func validateConfirmPassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("password must match confirm password")
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) == 0 || len(email) >= 255 {
		return errors.New("email address must be between 0 and 255 characters")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email address must be a valid email address")
	}
	if addr.Address != email {
		return errors.New("email address must be a valid email address")
	}
	return nil
}

func validateMinValue(value, min int) error {
	if value < min {
		return errors.New("value must be greater than or equal to " + strconv.Itoa(min))
	}
	return nil
}

func validateMaxValue(value, max int) error {
	if value > max {
		return errors.New("value must be less than or equal to " + strconv.Itoa(max))
	}
	return nil
}

func validateNotEmpty(value string) error {
	if value == "" {
		return errors.New("value must not be empty")
	}
	return nil
}

func validateContains[T comparable](value T, allowedValues []T) error {
	if !slices.Contains(allowedValues, value) {
		return errors.New("value not allowed")
	}
	return nil
}

func validateAllContains[T comparable](values, allowedValues []T) error {
	for _, value := range values {
		err := validateContains(value, allowedValues)
		if err != nil {
			return err
		}
	}
	return nil
}
