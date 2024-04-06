package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/errors"
)

type Provider interface {
	GetConfig() (*Config, error)
	NewFromConfigFile() (openai.ClientConfig, error)
	GetConfigDir() (string, error)
	SaveConfig(config *Config) error
	ClearConfig() error
}

type FileManager interface {
	Read(path string) ([]byte, error)
	Write(path string, content []byte, perm os.FileMode) error
}
type FileReader func(path string) ([]byte, error)

type Marshaller interface {
	Unmarshal(in []byte, out interface{}) error
	Marshal(in interface{}) ([]byte, error)
}

type Parser func(in []byte, out interface{}) error

type Validator func(cfg *Config) error

type DefaultProviderOptions struct {
	ConfigPath  string
	FileManager FileManager
	Marshaller  Marshaller
	Validator   Validator
	KeyStore    APIKeyStorage
}

type DefaultProvider struct {
	configPath  string
	fileManager FileManager
	marshaller  Marshaller
	keyStore    APIKeyStorage
	validate    Validator
}

func NewDefaultProvider() *DefaultProvider {
	return NewDefaultProviderWithOptions(DefaultProviderOptions{
		ConfigPath:  "",
		FileManager: &OSFileManager{},
		Marshaller:  &YamlMarshaller{},
		Validator:   DefaultValidator,
		KeyStore:    NewKeyRingAPIKeyStorage("cwc", user.Current),
	})
}

func NewDefaultProviderWithOptions(opts DefaultProviderOptions) *DefaultProvider {
	if opts.ConfigPath == "" {
		path, err := DefaultConfigPath()
		if err != nil {
			path = configFileName
		}

		opts.ConfigPath = path
	}

	if opts.FileManager == nil {
		opts.FileManager = &OSFileManager{}
	}

	if opts.Marshaller == nil {
		opts.Marshaller = &YamlMarshaller{}
	}

	if opts.Validator == nil {
		opts.Validator = DefaultValidator
	}

	if opts.KeyStore == nil {
		opts.KeyStore = NewKeyRingAPIKeyStorage("cwc", user.Current)
	}

	return &DefaultProvider{
		configPath:  opts.ConfigPath,
		fileManager: opts.FileManager,
		marshaller:  opts.Marshaller,
		validate:    opts.Validator,
		keyStore:    opts.KeyStore,
	}
}

func (c *DefaultProvider) GetConfig() (*Config, error) {
	if c.configPath == "" {
		path, err := DefaultConfigPath()
		if err != nil {
			return nil, err
		}

		c.configPath = path
	}

	data, err := c.fileManager.Read(c.configPath)
	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			"config file does not exist",
			"please run `cwc login` to create a new config file.",
		}}
	}

	var cfg Config
	err = c.marshaller.Unmarshal(data, &cfg)

	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			"invalid config file format",
			"please run `cwc login` to create a new config file.",
		}}
	}

	apiKey, err := c.keyStore.GetAPIKey()
	if err != nil {
		return nil, errors.ConfigValidationError{Errors: []string{
			err.Error(),
			"please run `cwc login` to create a new config file.",
		}}
	}

	cfg.SetAPIKey(apiKey)

	return &cfg, nil
}

func (c *DefaultProvider) NewFromConfigFile() (openai.ClientConfig, error) {
	cfg, err := c.GetConfig()
	if err != nil {
		return openai.ClientConfig{}, err
	}

	// validate the configuration
	err = c.validate(cfg)
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

func (c *DefaultProvider) SaveConfig(config *Config) error {
	// validate the configuration
	err := c.validate(config)
	if err != nil {
		return err
	}

	data, err := c.marshaller.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config data: %w", err)
	}

	err = c.keyStore.SetAPIKey(config.APIKey())
	if err != nil {
		return fmt.Errorf("error saving API key in keystore: %w", err)
	}

	if c.configPath == "" {
		c.configPath, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	err = c.fileManager.Write(c.configPath, data, configFilePermissions)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func (c *DefaultProvider) GetConfigDir() (string, error) {
	return filepath.Dir(c.configPath), nil
}

func (c *DefaultProvider) ClearConfig() error {
	err := os.Remove(c.configPath)
	if err != nil {
		return fmt.Errorf("error removing config file: %w", err)
	}

	err = c.keyStore.ClearAPIKey()
	if err != nil {
		return fmt.Errorf("error clearing API key from storage: %w", err)
	}

	return nil
}
