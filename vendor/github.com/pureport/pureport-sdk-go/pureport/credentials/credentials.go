// Package credentials provides ...
package credentials

import (
	"sync"
	"time"
)

// Value is the gathered credentials information
type Value struct {

	// Pureport API Key
	APIKey string

	// Pureport Secret Access Key
	Secret string

	// Pureport Session Token
	SessionToken string

	// Pureport Refresh Token
	RefreshToken string

	// Provider used to get credentials
	ProviderName string
}

// The Provider interface defines a component who is capable of requesting
// credentials for a session and managing credential expiration.
//
// The Provider instance does not need to handle locking of resources. This
// will be managed by the Credentials object.
type Provider interface {
	Retrieve() (Value, error)
	IsExpired() bool
}

// Expiry stores the expiration data use by credentials providers
type Expiry struct {

	// The expiration time of the credential
	expiration time.Time

	// Function to use to get the current time
	// This can be changed by unit tests.
	CurrentTime func() time.Time
}

// SetExpiration set the time of expiration that is used by Provider.IsExpired()
//
// Specifying a value for window will cause the expiration to occur prior to the
// actual expiration time. This will allow a new token to be requested before the
// existing token has reach is expiration time.
func (e *Expiry) SetExpiration(expiration time.Time, window time.Duration) {
	e.expiration = expiration
	if window > 0 {
		e.expiration = e.expiration.Add(-window)
	}
}

// IsExpired return if the current Expiry time has be reached
func (e *Expiry) IsExpired() bool {
	currentTime := e.CurrentTime
	if currentTime == nil {
		currentTime = time.Now
	}
	return e.expiration.Before(currentTime())
}

// ExpiresAt returns the expiration time
func (e *Expiry) ExpiresAt() time.Time {
	return e.expiration
}

// The Credentials object provides threadsafe access to credentials
// provided by the Pureport API. It takes care of caching credentials and requesting
// new access tokens when they expire.
type Credentials struct {
	credentials  Value
	forceRefresh bool
	m            sync.RWMutex
	provider     Provider
}

// NewCredentials creates new credentials information for the specified provider
func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider:     provider,
		forceRefresh: true,
	}
}

func (c *Credentials) printCredentials() {

	printPwd := ""
	runeCount := 0

	for _, char := range c.credentials.Secret {
		runeCount++

		if runeCount < 5 {
			printPwd += string(char)
		} else {
			printPwd += "*"
		}
	}

	log.Debugf("Found Credentials: key=%s, secret=%s", c.credentials.APIKey, printPwd)
}

// Get the current value of the credentials
func (c *Credentials) Get() (Value, error) {
	c.m.RLock()

	// Check to see if the credentials have expired
	if !c.isExpired() {
		credentials := c.credentials
		c.m.RUnlock()

		c.printCredentials()
		return credentials, nil
	}
	c.m.RUnlock()

	// We need to request new credentials so get the R/W lock
	c.m.Lock()
	defer c.m.Unlock()

	if c.isExpired() {
		credentials, err := c.provider.Retrieve()
		if err != nil {
			return Value{}, err
		}

		c.credentials = credentials
		c.forceRefresh = false
	}

	c.printCredentials()

	return c.credentials, nil
}

// Expire force refresh of the credentials even if they haven't expired
func (c *Credentials) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
}

// IsExpired whether the current credentials are expired
func (c *Credentials) IsExpired() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.isExpired()
}

// Private function to check if the session token has expired
func (c *Credentials) isExpired() bool {
	return c.forceRefresh || c.provider.IsExpired()
}
