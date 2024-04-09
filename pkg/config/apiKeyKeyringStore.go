package config

import (
	"fmt"
	"os/user"

	"github.com/zalando/go-keyring"
)

type UsernameRetriever func() (*user.User, error)

type APIKeyKeyringStore struct {
	serviceName       string
	usernameRetriever UsernameRetriever
}

func NewAPIKeyKeyringStore(serviceName string, usernameRetriever UsernameRetriever) *APIKeyKeyringStore {
	return &APIKeyKeyringStore{
		serviceName:       serviceName,
		usernameRetriever: usernameRetriever,
	}
}

func (k *APIKeyKeyringStore) GetAPIKey() (string, error) {
	usr, err := k.usernameRetriever()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}

	apiKey, err := keyring.Get(k.serviceName, usr.Username)
	if err != nil {
		return "", fmt.Errorf("error getting API key from keyring: %w", err)
	}

	return apiKey, nil
}

func (k *APIKeyKeyringStore) SetAPIKey(apiKey string) error {
	usr, err := k.usernameRetriever()
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	username := usr.Username
	err = keyring.Set(k.serviceName, username, apiKey)

	if err != nil {
		return fmt.Errorf("error storing API key in keyring: %w", err)
	}

	return nil
}

func (k *APIKeyKeyringStore) ClearAPIKey() error {
	usr, err := k.usernameRetriever()
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
