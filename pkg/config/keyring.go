package config

import (
	"github.com/zalando/go-keyring"
	"os/user"
)

func getApiKeyFromKeyring() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	username := usr.Username
	apiKey, err := keyring.Get(serviceName, username)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func storeApiKeyInKeyring(apiKey string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	username := usr.Username
	err = keyring.Set(serviceName, username, apiKey)
	if err != nil {
		return err
	}

	return nil
}

func clearApiKeyInKeyring() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	username := usr.Username
	err = keyring.Delete(serviceName, username)
	if err != nil {
		return err
	}

	return nil
}
