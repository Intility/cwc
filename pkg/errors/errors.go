package errors

import (
	"errors"
	"fmt"
	"strings"
)

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

type SaveConfigError struct {
	Message string
	Err     error
}

func (e *SaveConfigError) Error() string {
	return e.Message + ": " + e.Err.Error()
}

// ConfigValidationError collects all configuration validation errors.
type ConfigValidationError struct {
	Errors []string
}

func (e ConfigValidationError) Error() string {
	return fmt.Sprintf("config validation failed: %s", strings.Join(e.Errors, ", "))
}

// AsConfigValidationError attempts to convert an error to a *ConfigValidationError and returns it with a boolean indicating success.
func AsConfigValidationError(err error) (*ConfigValidationError, bool) {
	var validationErr *ConfigValidationError
	if err != nil {
		ok := errors.As(err, &validationErr)
		return validationErr, ok
	}
	return nil, false
}
