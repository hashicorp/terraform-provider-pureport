package credentials

import (
	"errors"
	"os"
)

const (
	// EnvironmentProviderName is the name of this provider
	EnvironmentProviderName = "EnvironmentProvider"

	// CredAPIKeyEnvStr - environment variable for the Pureport API Key
	CredAPIKeyEnvStr = "PUREPORT_API_KEY"

	// CredAPISecretEnvStr - environment variable for the Pureport API Secret
	CredAPISecretEnvStr = "PUREPORT_API_SECRET"
)

var (
	// ErrorEnvironmentAPIKeyNotFound for the when the APIKey isn't found in the environment
	ErrorEnvironmentAPIKeyNotFound = errors.New("PUREPORT_API_KEY not found in environment")

	// ErrorEnvironmentAPISecretNotFound for the when the API Secret isn't found in the environment
	ErrorEnvironmentAPISecretNotFound = errors.New("PUREPORT_API_SECRET not found in environment")
)

// An EnvironmentProvider retrieves the base credentials from the current execution
// environment of the running process. Environment credentials never expire.
//
type EnvironmentProvider struct {
	retrieved bool
}

// NewEnvironmentCredentials creates a new Credentials using the EnvironmentProvider
func NewEnvironmentCredentials() *Credentials {
	return NewCredentials(&EnvironmentProvider{})
}

// Retrieve see Provider.Retrieve()
func (e *EnvironmentProvider) Retrieve() (Value, error) {
	e.retrieved = false

	key := os.Getenv(CredAPIKeyEnvStr)
	if key == "" {
		return Value{ProviderName: EnvironmentProviderName}, ErrorEnvironmentAPIKeyNotFound
	}

	secret := os.Getenv(CredAPISecretEnvStr)
	if secret == "" {
		return Value{ProviderName: EnvironmentProviderName}, ErrorEnvironmentAPISecretNotFound
	}

	e.retrieved = true
	return Value{
		APIKey:       key,
		Secret:       secret,
		ProviderName: EnvironmentProviderName,
	}, nil
}

// IsExpired see Provider.IsExpired()
func (e *EnvironmentProvider) IsExpired() bool {
	return !e.retrieved
}
