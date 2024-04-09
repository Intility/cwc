package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	configFileName        = "cwc.yaml" // The name of the config file we want to save
	configFilePermissions = 0o600      // The permissions we want to set on the config file
	apiVersion            = "2024-02-01"
)

// SanitizeInput trims whitespaces and newlines from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

type Config struct {
	Endpoint        string `yaml:"endpoint"`
	ModelDeployment string `yaml:"modelDeployment"`
	ExcludeGitDir   bool   `yaml:"excludeGitDir"`
	UseGitignore    bool   `yaml:"useGitignore"`
	// Keep APIKey unexported to avoid accidental exposure
	apiKey string
}

// NewConfig creates a new Config object.
func NewConfig(endpoint, modelDeployment string) *Config {
	return &Config{
		Endpoint:        endpoint,
		ModelDeployment: modelDeployment,
		ExcludeGitDir:   true,
		UseGitignore:    true,
		apiKey:          "",
	}
}

// SetAPIKey sets the confidential field apiKey.
func (c *Config) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// APIKey returns the confidential field apiKey.
func (c *Config) APIKey() string {
	return c.apiKey
}

func GetConfigDir() (string, error) {
	return XdgConfigPath()
}

func DefaultConfigPath() (string, error) {
	cfgPath, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cfgPath, configFileName), nil
}

func IsWSL() bool {
	_, exists := os.LookupEnv("WSL_DISTRO_NAME")
	return exists
}
