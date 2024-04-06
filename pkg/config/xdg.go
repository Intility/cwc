package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	serviceName = "cwc" // The name of our application
)

// helper function to get the XDG config path.
func XdgConfigPath() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		// XDG_CONFIG_HOME was not set, use the default "~/.config"
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting user home directory: %w", err)
		}

		xdgConfigHome = filepath.Join(homeDir, ".config")
	}

	configDir := filepath.Join(xdgConfigHome, serviceName) // use serviceName to create a subdirectory for our application

	// Ensure that the config directory exists
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return configDir, nil
}
