package config

type APIKeyStorage interface {
	GetAPIKey() (string, error)
	SetAPIKey(apiKey string) error
	ClearAPIKey() error
}
