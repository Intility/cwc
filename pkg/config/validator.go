package config

import "github.com/intility/cwc/pkg/errors"

func DefaultValidator(cfg *Config) error {
	var validationErrors []string

	if cfg.APIKey() == "" {
		validationErrors = append(validationErrors, "apiKey must be provided and not be empty")
	}

	if cfg.Endpoint == "" {
		validationErrors = append(validationErrors, "endpoint must be provided and not be empty")
	}

	if cfg.ModelDeployment == "" {
		validationErrors = append(validationErrors, "modelDeployment must be provided and not be empty")
	}

	if len(validationErrors) > 0 {
		return &errors.ConfigValidationError{Errors: validationErrors}
	}

	return nil
}
