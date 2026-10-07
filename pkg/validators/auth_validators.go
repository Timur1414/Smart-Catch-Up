package validators

import "github.com/Timur1414/Smart-Catch-Up/internal/web_helpers"

func ValidateRegisterUser(email, password, confirmPassword string) []web_helpers.ValidationError {
	var errors []web_helpers.ValidationError
	err := validateEmail(email)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "email", Message: err.Error()})
	}
	err = validatePassword(password)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "password", Message: err.Error()})
	}
	err = validateConfirmPassword(password, confirmPassword)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "confirm_password", Message: err.Error()})
	}
	return errors
}
