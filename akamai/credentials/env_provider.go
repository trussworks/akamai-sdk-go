// Example of using the environment variable credentials.
//
//     creds := credentials.NewEnvCredentials()
//
//     // Retrieve the credentials value
//     credValue, err := creds.Get()
//     if err != nil {
//         // handle error
//     }
//
package credentials

import (
	"errors"
	"os"
)

const EnvProviderName = "EnvProvider"

var (
	ErrClientSecretNotFoundEnv = errors.New("AKAMAI_CLIENT_SECRET not found in environment")
	ErrClientTokenNotFoundEnv  = errors.New("AKAMAI_CLIENT_TOKEN not found in environment")
	ErrAccessTokenNotFoundEnv  = errors.New("AKAMAI_ACCESS_TOKEN not found in environment")
	ErrAkamaiHostNotFoundEnv   = errors.New("AKAMAI_HOST not found in environment")
)

// An EnvProvider retrieves the credentials from the environment variables
// of the running process.
//
// Environment variables used:
//
// AKAMAI_ACCESS_TOKEN
// AKAMAI_CLIENT_SECRET
// AKAMAI_CLIENT_TOKEN
// AKAMAI_HOST
type EnvProvider struct {
	retrieved bool
}

// NewEnvCredentials returns a pointer to a new Credentials object
// wrapping the environment variable provider.
func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (AuthValue, error) {
	e.retrieved = false

	cs := os.Getenv("AKAMAI_CLIENT_SECRET")
	if cs == "" {
		return AuthValue{ProviderName: EnvProviderName}, ErrClientSecretNotFoundEnv
	}

	ct := os.Getenv("AKAMAI_CLIENT_TOKEN")
	if ct == "" {
		return AuthValue{ProviderName: EnvProviderName}, ErrClientTokenNotFoundEnv
	}

	at := os.Getenv("AKAMAI_ACCESS_TOKEN")
	if at == "" {
		return AuthValue{ProviderName: EnvProviderName}, ErrAccessTokenNotFoundEnv
	}

	ah := os.Getenv("AKAMAI_HOST")
	if ah == "" {
		return AuthValue{ProviderName: EnvProviderName}, ErrAkamaiHostNotFoundEnv

	}

	e.retrieved = true
	return AuthValue{
		ClientSecret: cs,
		ClientToken:  ct,
		AccessToken:  at,
		Host:         ah,
		ProviderName: EnvProviderName,
	}, nil
}
