package config

import (
	"fmt"
	"os"
)

type OSFileManager struct{}

func (y *OSFileManager) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}

func (y *OSFileManager) Write(path string, content []byte, perm os.FileMode) error {
	err := os.WriteFile(path, content, perm)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
