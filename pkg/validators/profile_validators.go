package validators

import "github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"

func ValidateUpdateProfile(email, firstName, lastName string) []web_helpers.ValidationError {
	errors := make([]web_helpers.ValidationError, 0)
	err := validateEmail(email)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "email", Message: err.Error()})
	}
	err = validateNotEmpty(firstName)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "first_name", Message: err.Error()})
	}
	err = validateNotEmpty(lastName)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "last_name", Message: err.Error()})
	}
	return errors
}
