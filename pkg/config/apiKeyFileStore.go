package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type APIKeyFileStore struct {
	Filepath string
}

func NewAPIKeyFileStore(filepath string) *APIKeyFileStore {
	return &APIKeyFileStore{Filepath: filepath}
}

func (fs *APIKeyFileStore) SetAPIKey(key string) error {
	var (
		dirFileMode os.FileMode = 0o700
		keyFileMode os.FileMode = 0o600
	)

	if err := os.MkdirAll(filepath.Dir(fs.Filepath), dirFileMode); err != nil {
		return fmt.Errorf("error creating directories for file storage: %w", err)
	}

	err := os.WriteFile(fs.Filepath, []byte(key), keyFileMode)
	if err != nil {
		return fmt.Errorf("error storing API key in file: %w", err)
	}

	return nil
}

func (fs *APIKeyFileStore) GetAPIKey() (string, error) {
	data, err := os.ReadFile(fs.Filepath)
	if os.IsNotExist(err) {
		return "", nil // Return an empty string if the file does not exist
	} else if err != nil {
		return "", fmt.Errorf("error reading API key from file: %w", err)
	}

	return string(data), nil
}

func (fs *APIKeyFileStore) ClearAPIKey() error {
	err := os.Remove(fs.Filepath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing API key from file: %w", err)
	}

	return nil
}
