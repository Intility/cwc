package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intility/cwc/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

const (
	configFileName        = "cwc.yaml" // The name of the config file we want to save
	configFilePermissions = 0o600      // The permissions we want to set on the config file
	apiVersion            = "2024-02-01"
)

func NewFromConfigFile() (openai.ClientConfig, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return openai.ClientConfig{}, err
	}

	// validate the configuration
	err = ValidateConfig(cfg)
	if err != nil {
		return openai.ClientConfig{}, err
	}

	config := openai.DefaultAzureConfig(cfg.APIKey(), cfg.Endpoint)
	config.APIVersion = apiVersion
	config.AzureModelMapperFunc = func(model string) string {
		return cfg.ModelDeployment
	}

	return config, nil
}

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

// ValidateConfig checks if a Config object has valid data.
func ValidateConfig(cfg *Config) error {
	var validationErrors []string

	if cfg.APIKey() == "" {
		validationErrors = append(validationErrors, "apiKey must be provided and not be empty")
	}

	if cfg.Endpoint == "" {
		validationErrors = append(validationErrors, "endpoint must be provided and not be empty")
	}

	if cfg.ModelDeployment == "" {
		validationErrors = append(validationErrors, "modelDeployment must be provided and not be empty")
	}

	if len(validationErrors) > 0 {
		return &errors.ConfigValidationError{Errors: validationErrors}
	}

	return nil
}

// SaveConfig writes the configuration to disk, and the API key to the keyring.
func SaveConfig(config *Config) error {
	// validate the configuration
	err := ValidateConfig(config)
	if err != nil {
		return err
	}

	configDir, err := xdgConfigPath()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(configDir, configFileName)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config data: %w", err)
	}

	err = storeAPIKeyInKeyring(config.APIKey())

	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, configFilePermissions)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// LoadConfig reads the configuration from disk and loads the API key from the keyring.
func LoadConfig() (*Config, error) {
	// Read data from file or secure store
	configDir, err := xdgConfigPath()
	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			err.Error(),
			"please run `cwc login` to create a new config file.",
		}}
	}

	data, err := os.ReadFile(filepath.Join(configDir, configFileName))
	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			"config file does not exist",
			"please run `cwc login` to create a new config file.",
		}}
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			"invalid config file format",
			"please run `cwc login` to create a new config file.",
		}}
	}

	apiKey, err := getAPIKeyFromKeyring()
	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			err.Error(),
			"please run `cwc login` to create a new config file.",
		}}
	}

	cfg.SetAPIKey(apiKey)

	return &cfg, nil
}

func ClearConfig() error {
	configDir, err := xdgConfigPath()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(configDir, configFileName)

	err = os.Remove(configFilePath)
	if err != nil {
		return fmt.Errorf("error removing config file: %w", err)
	}

	err = clearAPIKeyInKeyring()
	if err != nil {
		return err
	}

	return nil
}

func GetConfigDir() (string, error) {
	return xdgConfigPath()
}
