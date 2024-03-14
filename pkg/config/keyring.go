package config

import (
	"fmt"
	"os/user"

	"github.com/zalando/go-keyring"
)

func getAPIKeyFromKeyring() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}

	apiKey, err := keyring.Get(serviceName, usr.Username)
	if err != nil {
		return "", fmt.Errorf("error getting API key from keyring: %w", err)
	}

	return apiKey, nil
}

func storeAPIKeyInKeyring(apiKey string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	username := usr.Username
	err = keyring.Set(serviceName, username, apiKey)

	if err != nil {
		return fmt.Errorf("error storing API key in keyring: %w", err)
	}

	return nil
}

func clearAPIKeyInKeyring() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	username := usr.Username
	err = keyring.Delete(serviceName, username)

	if err != nil {
		return fmt.Errorf("error deleting API key from keyring: %w", err)
	}

	return nil
}
