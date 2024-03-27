package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

type APIKeyStorage interface {
	StoreAPIKey(key string) error
	GetAPIKey() (string, error)
	ClearAPIKey() error
}

type KeyringAPIKeyStorage struct {
	Servicename string
}

type FileAPIKeyStorage struct {
	Filepath string
}

func isWSL() bool {
	if _, exists := os.LookupEnv("WSL_DISTRO_NAME"); exists {
		return true
	}
	return false
}

func (fs *KeyringAPIKeyStorage) StoreAPIKey(key string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	username := usr.Username
	err = keyring.Set(serviceName, username, key)

	if err != nil {
		return fmt.Errorf("error storing API key in keyring: %w", err)
	}

	return nil
}

func (fs *KeyringAPIKeyStorage) ClearAPIKey() error {
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

func (fs *KeyringAPIKeyStorage) GetAPIKey() (string, error) {
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

func (fs *FileAPIKeyStorage) StoreAPIKey(key string) error {

	if err := os.MkdirAll(filepath.Dir(fs.Filepath), 0700); err != nil {
		return fmt.Errorf("error creating directories for file storage: %w", err)
	}

	err := os.WriteFile(fs.Filepath, []byte(key), 0600)
	if err != nil {
		return fmt.Errorf("error storing API key in file: %w", err)
	}

	return nil
}

func (fs *FileAPIKeyStorage) GetAPIKey() (string, error) {

	data, err := os.ReadFile(fs.Filepath)
	if os.IsNotExist(err) {
		return "", nil // Return an empty string if the file does not exist
	} else if err != nil {
		return "", fmt.Errorf("error reading API key from file: %w", err)
	}
	return string(data), nil
}

func (fs *FileAPIKeyStorage) ClearAPIKey() error {
	err := os.Remove(fs.Filepath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing API key from file: %w", err)
	}

	return nil
}

func ResolveAPIKeyStorage() APIKeyStorage {
	if isWSL() {
		usr, err := user.Current()
		if err != nil {
			return nil
		}
		homeDir := usr.HomeDir
		apiKeyStoragePath := homeDir + "/.cwc/apikeystore"
		return &FileAPIKeyStorage{Filepath: apiKeyStoragePath}
	}
	return &KeyringAPIKeyStorage{Servicename: serviceName}
}
