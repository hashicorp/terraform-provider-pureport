package credentials

import ()

// ErrorProvider is a stub credentials provider that will always return an error.
// This provider is returned when the requested provider creation
// failed due to errors.
type ErrorProvider struct {

	// The error to return from Provider.Retrieve()
	Err error

	// The name to set on the Provider.Retrieve()
	ProviderName string
}

// Retrieve see Provider.Retrieve()
func (p ErrorProvider) Retrieve() (Value, error) {
	return Value{ProviderName: p.ProviderName}, p.Err
}

// IsExpired see Provider.IsExpired()
func (p ErrorProvider) IsExpired() bool {
	return false
}
