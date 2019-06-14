package credentials

import (
	"fmt"
)

// ChainProvider provides chaining of multiple providers so we can gracefully
// go through the provider list until we find one that works.
type ChainProvider struct {
	Providers []Provider
	current   Provider
}

// NewChainCredentials creates a new credential using the chaining provider
func NewChainCredentials(providers []Provider) *Credentials {
	return NewCredentials(&ChainProvider{
		Providers: append([]Provider{}, providers...),
	})
}

// Retrieve - see Provider.Retrieve()
func (c *ChainProvider) Retrieve() (Value, error) {
	var errs []error

	for _, p := range c.Providers {
		credentials, err := p.Retrieve()
		if err == nil {
			c.current = p
			return credentials, nil
		}
		errs = append(errs, err)
	}

	c.current = nil

	return Value{}, fmt.Errorf("No valid providers in chain %s", errs)
}

// IsExpired - see Provider.IsExpired()
func (c *ChainProvider) IsExpired() bool {
	if c.current != nil {
		return c.current.IsExpired()
	}
	return true
}
