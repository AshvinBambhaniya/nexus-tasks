package utils

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// ValidateMessage is the default message for validation errors.
const ValidateMessage = "fields are invalid."

// ValidateEmail validates an email address.
func ValidateEmail(email string) (bool, error) {
	return regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
}

// ValidatorErrorString returns a formatted error string for validation errors.
func ValidatorErrorString(err error) string {
	var msg string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msg += strings.ToLower(err.Field()) + ","
		}
		msg = strings.TrimSuffix(msg, ",")
		msg = fmt.Sprintf("%s %s", msg, ValidateMessage)
		return msg
	}
	return msg
}
