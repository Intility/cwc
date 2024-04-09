package config

type APIKeyStore interface {
	GetAPIKey() (string, error)
	SetAPIKey(apiKey string) error
	ClearAPIKey() error
}
