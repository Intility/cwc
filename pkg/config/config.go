package config

import (
	"encoding/json"
	"github.com/emilkje/cwc/pkg/errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	configFileName = "cwc.json" // The name of the config file we want to save
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
	config.APIVersion = cfg.ApiVersion
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
	Endpoint        string `json:"endpoint"`
	ApiVersion      string `json:"apiVersion"`
	ModelDeployment string `json:"modelDeployment"`
	// Keep APIKey unexported to avoid accidental exposure
	apiKey string
}

// NewConfig creates a new Config object
func NewConfig(endpoint, apiVersion, modelDeployment string) *Config {
	return &Config{
		Endpoint:        endpoint,
		ApiVersion:      apiVersion,
		ModelDeployment: modelDeployment,
	}
}

// SetAPIKey sets the confidential field apiKey
func (c *Config) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// APIKey returns the confidential field apiKey
func (c *Config) APIKey() string {
	return c.apiKey
}

// ValidateConfig checks if a Config object has valid data.
func ValidateConfig(c *Config) error {
	var validationErrors []string

	if c.APIKey() == "" {
		validationErrors = append(validationErrors, "apiKey must be provided and not be empty")
	}

	if c.Endpoint == "" {
		validationErrors = append(validationErrors, "endpoint must be provided and not be empty")
	}

	if c.ApiVersion == "" {
		validationErrors = append(validationErrors, "apiVersion must be provided and not be empty")
	}

	if c.ModelDeployment == "" {
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

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = storeApiKeyInKeyring(config.APIKey())

	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}

// LoadConfig reads the configuration from disk and loads the API key from the keyring.
func LoadConfig() (*Config, error) {
	// Read data from file or secure store
	configDir, err := xdgConfigPath()
	if err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(configDir, configFileName)
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	apiKey, err := getApiKeyFromKeyring()
	if err != nil {
		return nil, err
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
		return err
	}

	err = clearApiKeyInKeyring()
	if err != nil {
		return err
	}

	return nil
}
