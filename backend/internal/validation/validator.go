package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Custom validation functions
	validate.RegisterValidation("username", validateUsername)
}

// ValidateStruct validates any struct with validation tags
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return formatValidationError(err)
	}
	return nil
}

// Custom username validation
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username should be alphanumeric with underscores, 3-30 chars
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}

// Format validation errors to be user-friendly
func formatValidationError(err error) error {
	var errors []string

	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors = append(errors, fmt.Sprintf("%s is required", err.Field()))
		case "email":
			errors = append(errors, fmt.Sprintf("%s must be a valid email", err.Field()))
		case "min":
			errors = append(errors, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			errors = append(errors, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
		case "oneof":
			errors = append(errors, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
		default:
			errors = append(errors, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
}
