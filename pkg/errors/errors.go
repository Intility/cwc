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
	return "config validation failed: " + strings.Join(e.Errors, ", ")
}

// AsConfigValidationError attempts to convert an error to a
// *ConfigValidationError and returns it with a boolean indicating success.
func AsConfigValidationError(err error) (*ConfigValidationError, bool) {
	var validationErr *ConfigValidationError
	if err != nil {
		ok := errors.As(err, &validationErr)
		return validationErr, ok
	}

	return nil, false
}

func IsConfigValidationError(err error) bool {
	var validationErr ConfigValidationError
	return errors.As(err, &validationErr)
}

// FileNotExistError is an error type for when a file does not exist.
type FileNotExistError struct {
	FileName string
}

func (e FileNotExistError) Error() string {
	return fmt.Sprintf("file %s does not exist", e.FileName)
}

func IsFileNotExistError(err error) bool {
	var fileDoesNotExistError FileNotExistError
	return errors.As(err, &fileDoesNotExistError)
}

type GitNotInstalledError struct {
	Message string
}

func (e GitNotInstalledError) Error() string {
	return e.Message
}

func IsGitNotInstalledError(err error) bool {
	var gitNotInstalledError GitNotInstalledError
	return errors.As(err, &gitNotInstalledError)
}

type NotAGitRepositoryError struct {
	Message string
}

func (e NotAGitRepositoryError) Error() string {
	return e.Message
}

func IsNotAGitRepositoryError(err error) bool {
	var notAGitRepositoryError NotAGitRepositoryError
	return errors.As(err, &notAGitRepositoryError)
}

type NoPromptProvidedError struct {
	Message string
}

func (e NoPromptProvidedError) Error() string {
	return e.Message
}

type TemplateNotFoundError struct {
	TemplateName string
}

func (e TemplateNotFoundError) Error() string {
	return "template not found: " + e.TemplateName
}

func IsTemplateNotFoundError(err error) bool {
	var templateNotFoundError TemplateNotFoundError
	return errors.As(err, &templateNotFoundError)
}

type SuppressedError struct{}

func (e SuppressedError) Error() string {
	return "error suppressed"
}

type ArgParseError struct {
	Message string
}

func (e ArgParseError) Error() string {
	return e.Message
}
