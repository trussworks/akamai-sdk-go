package credentials

import (
	"errors"
	"sync"
	"time"
)

// AuthValue is a struct holding the Akamai credentials for each field
type AuthValue struct {
	// Akamai client_secret
	ClientSecret string

	// Akamai client_token
	ClientToken string

	// Akamai access_token
	AccessToken string

	// Akamai host
	Host string

	// Provider used to get credentials
	ProviderName string
}

// Provider is an interface for a component that will provide a CredentialValue
// This can be used to read from an environment or config file, or any other
// method that returns the authentication credentials as an AuthValue.
type Provider interface {
	Retrieve() (AuthValue, error)

	// IsExpired returns if the credentials are no longer valid, and need
	// to be retrieved.
	IsExpired() bool
}

// Credentials provides concurrency safe retrieval of Akamai credentials.
type Credentials struct {
	creds        AuthValue
	forceRefresh bool

	m sync.RWMutex

	provider Provider
}

// NewCredentials returns a pointer to a new Credentials with the provider set.
func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider:     provider,
		forceRefresh: true,
	}
}

// Get returns the credentials AuthValue or error in the case of failure.
//
// Will return the cached credentials Value if it has not expired. If the
// credentials Value has expired the Provider's Retrieve() will be called
// to refresh the credentials.
func (c *Credentials) Get() (AuthValue, error) {
	// Check the cached credentials first with just the read lock.
	c.m.RLock()
	if !c.isExpired() {
		creds := c.creds
		c.m.RUnlock()
		return creds, nil
	}
	c.m.RUnlock()

	// Credentials are expired need to retrieve the credentials taking the full
	// lock.
	c.m.Lock()
	defer c.m.Unlock()

	if c.isExpired() {
		creds, err := c.provider.Retrieve()
		if err != nil {
			return AuthValue{}, err
		}
		c.creds = creds
		c.forceRefresh = false
	}

	return c.creds, nil
}

// Expire expires the credentials and forces them to be retrieved on the
// next call to Get().
//
// This will override the Provider's expired state, and force Credentials
// to call the Provider's Retrieve().
func (c *Credentials) Expire() {
	c.m.Lock()
	defer c.m.Unlock()

	c.forceRefresh = true
}

// IsExpired returns if the credentials are no longer valid, and need
// to be retrieved.
//
// If the Credentials were forced to be expired with Expire() this will
// reflect that override.
func (c *Credentials) IsExpired() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.isExpired()
}

// isExpired helper method wrapping the definition of expired credentials.
func (c *Credentials) isExpired() bool {
	return c.forceRefresh || c.provider.IsExpired()
}

// ExpiresAt provides access to the functionality of the Expirer interface of
// the underlying Provider, if it supports that interface.  Otherwise, it returns
// an error.
func (c *Credentials) ExpiresAt() (time.Time, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	expirer, ok := c.provider.(Expirer)
	if !ok {
		return time.Time{}, errors.New("provider does not support ExpiresAt()")
	}
	if c.forceRefresh {
		// set expiration time to the distant past
		return time.Time{}, nil
	}
	return expirer.ExpiresAt(), nil
}

// A Expiry provides shared expiration logic to be used by credentials
// providers to implement expiry functionality.
//
// The best method to use this struct is as an anonymous field within the
// provider's struct.
//
// Example:
//     type AkamaiProvider struct {
//         Expiry
//         ...
//     }
type Expiry struct {
	// The date/time when to expire on
	expiration time.Time

	// If set will be used by IsExpired to determine the current time.
	// Defaults to time.Now if CurrentTime is not set.  Available for testing
	// to be able to mock out the current time.
	CurrentTime func() time.Time
}

// SetExpiration sets the expiration IsExpired will check when called.
//
// If window is greater than 0 the expiration time will be reduced by the
// window value.
//
// Using a window is helpful to trigger credentials to expire sooner than
// the expiration time given to ensure no requests are made with expired
// tokens.
func (e *Expiry) SetExpiration(expiration time.Time, window time.Duration) {
	e.expiration = expiration
	if window > 0 {
		e.expiration = e.expiration.Add(-window)
	}
}

// IsExpired returns if the credentials are expired.
func (e *Expiry) IsExpired() bool {
	curTime := e.CurrentTime
	if curTime == nil {
		curTime = time.Now
	}
	return e.expiration.Before(curTime())
}

// ExpiresAt returns the expiration time of the credential
func (e *Expiry) ExpiresAt() time.Time {
	return e.expiration
}

// An Expirer is an interface that Providers can implement to expose the expiration
// time, if known.  If the Provider cannot accurately provide this info,
// it should not implement this interface.
type Expirer interface {
	// The time at which the credentials are no longer valid
	ExpiresAt() time.Time
}
