package credentials

import (
	"errors"
)

const (
	// StaticProviderName the name of this provider
	StaticProviderName = "StaticProvider"
)

var (

	// ErrorStaticAPIKeyNotFound is returned when the API Key can not be found in the configuration file
	ErrorStaticAPIKeyNotFound = errors.New("api_key not found in configuration file")

	// ErrorStaticAPISecretNotFound is returned when the API Secret can not be found in the configuration file
	ErrorStaticAPISecretNotFound = errors.New("api_secret not found in configuration file")
)

// StaticProvider provides credentials that are allocated at startup only.
type StaticProvider struct {
	APIKey    string
	APISecret string
}

// NewStaticCredentials create a new credential using static API information
func NewStaticCredentials(key string, secret string) *Credentials {
	return NewCredentials(&StaticProvider{
		APIKey:    key,
		APISecret: secret,
	})
}

// Retrieve Provider.Retrieve()
func (p *StaticProvider) Retrieve() (Value, error) {

	if p.APIKey == "" {
		return Value{ProviderName: StaticProviderName}, ErrorStaticAPIKeyNotFound
	}

	if p.APISecret == "" {
		return Value{ProviderName: StaticProviderName}, ErrorStaticAPISecretNotFound
	}

	return Value{
		ProviderName: StaticProviderName,
		APIKey:       p.APIKey,
		Secret:       p.APISecret,
	}, nil
}

// IsExpired Provider.IsExpired()
func (p *StaticProvider) IsExpired() bool {
	return false
}
