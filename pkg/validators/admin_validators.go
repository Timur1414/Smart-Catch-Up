package validators

import (
	"github.com/Timur1414/Smart-Catch-Up/pkg/web_helpers"
)

func ValidateGenerate1(text string, notificationType string, userId int, allowedTypes []string, allowedUsers []int) []web_helpers.ValidationError {
	var errors []web_helpers.ValidationError
	err := validateNotEmpty(text)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_1_text_error", Message: err.Error()})
	}
	err = validateContains(notificationType, allowedTypes)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_1_type_error", Message: err.Error()})
	}
	err = validateContains(userId, allowedUsers)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_1_user_error", Message: err.Error()})
	}
	return errors
}

func ValidateGenerateN(number int, notificationTypes []string, userIds []int, allowedTypes []string, allowedUsers []int) []web_helpers.ValidationError {
	var errors []web_helpers.ValidationError
	err := validateMinValue(number, 1)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_n_number_error", Message: err.Error()})
	}
	err = validateAllContains(notificationTypes, allowedTypes)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_n_types_error", Message: err.Error()})
	}
	err = validateAllContains(userIds, allowedUsers)
	if err != nil {
		errors = append(errors, web_helpers.ValidationError{Field: "admin_form_n_users_error", Message: err.Error()})
	}
	return errors
}
