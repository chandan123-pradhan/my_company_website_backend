package utils

import (
	"fmt"
	"strings"

	"blog_project.com/models"
)

// UserValidationError captures validation errors for user input
// 
// This struct is used to store details about validation errors,
// including the field that caused the error and the error message.
type UserValidationError struct {
	Field   string // The name of the field that failed validation
	Message string // The error message associated with the validation failure
}

// ValidateUserInput validates required fields for user registration and login.
// 
// It takes an interface type `user`, which can be either
// `models.RegisterUserModel` or `models.LoginUserModel`, and
// a boolean `isRegistration` to determine if the validation is
// for registration or login.
// 
// The function returns a slice of UserValidationError,
// capturing any validation errors found during the process.
func ValidateUserInput(user interface{}, isRegistration bool) []UserValidationError {
	var errors []UserValidationError

	switch u := user.(type) {
	case models.RegisterUserModel:
		// Validate fields for user registration
		if u.FullName == "" {
			errors = append(errors, UserValidationError{"full_name", "Full name is mandatory"})
		}
		if u.Email == "" {
			errors = append(errors, UserValidationError{"email", "Email is mandatory"})
		}
		if isRegistration {
			if u.Password == "" {
				errors = append(errors, UserValidationError{"password", "Password is mandatory"})
			}
			if u.ProfilePic == "" {
				errors = append(errors, UserValidationError{"profile_pic", "Profile picture is mandatory"})
			}
		}
	case models.LoginUserModel:
		// Validate fields for user login
		if u.Email == "" {
			errors = append(errors, UserValidationError{"email", "Email is mandatory"})
		}
		if u.Password == "" {
			errors = append(errors, UserValidationError{"password", "Password is mandatory"})
		}
	default:
		// Return an error if the user type is invalid
		return []UserValidationError{{"User", "Invalid user type"}}
	}

	return errors
}

// ErrorMessages returns a formatted string of all error messages.
// 
// It takes a slice of UserValidationError and constructs a
// single string that combines all error messages in a readable format.
// 
// Each error message includes the field name and the associated error
// message, separated by a colon.
func ErrorMessages(errors []UserValidationError) string {
	var messages []string
	for _, err := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ") // Join all messages with a comma
}


