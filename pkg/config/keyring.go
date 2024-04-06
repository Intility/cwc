package config

import (
	"fmt"
	"os/user"

	"github.com/zalando/go-keyring"
)

type UsernameRetriever func() (*user.User, error)

type KeyRingAPIKeyStorage struct {
	serviceName       string
	usernameRetriever UsernameRetriever
}

func NewKeyRingAPIKeyStorage(serviceName string, usernameRetriever UsernameRetriever) *KeyRingAPIKeyStorage {
	return &KeyRingAPIKeyStorage{
		serviceName:       serviceName,
		usernameRetriever: usernameRetriever,
	}
}

func (k *KeyRingAPIKeyStorage) GetAPIKey() (string, error) {
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

func (k *KeyRingAPIKeyStorage) SetAPIKey(apiKey string) error {
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

func (k *KeyRingAPIKeyStorage) ClearAPIKey() error {
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
