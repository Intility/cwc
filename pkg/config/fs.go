package config

import (
	"fmt"
	"os"
)

type OSFileManager struct{}

func (y *OSFileManager) Read(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %w", err)
		}

		return nil, fmt.Errorf("error stating file: %w", err)
	}

	// even though the path is a variable, it is safe to assume that the file exists
	// in a safe location as we are using XDG Base Directory Specification
	data, err := os.ReadFile(path) // #nosec
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
