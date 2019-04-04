package credentials

import "errors"

// StaticProviderName provides a name of Static provider
const StaticProviderName = "StaticProvider"

var (
	// ErrStaticCredentialsEmpty is emitted when static credentials are empty.
	ErrStaticCredentialsEmpty = errors.New("EmptyStaticCreds: static credentials are empty")
)

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	AuthValue
}

// NewStaticCredentials returns a pointer to a new Credentials object
// wrapping a static credentials value provider.
func NewStaticCredentials(cs, ct, at, ah string) *Credentials {
	return NewCredentials(&StaticProvider{AuthValue: AuthValue{
		ClientSecret: cs,
		ClientToken:  ct,
		AccessToken:  at,
		Host:         ah,
	}})
}

// NewStaticCredentialsFromCreds returns a pointer to a new Credentials object
// wrapping the static credentials value provide. Same as NewStaticCredentials
// but takes the creds AuthValue instead of individual fields
func NewStaticCredentialsFromCreds(creds AuthValue) *Credentials {
	return NewCredentials(&StaticProvider{AuthValue: creds})
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s *StaticProvider) Retrieve() (AuthValue, error) {
	if s.ClientSecret == "" || s.ClientToken == "" || s.AccessToken == "" || s.Host == "" {
		return AuthValue{ProviderName: StaticProviderName}, ErrStaticCredentialsEmpty
	}

	if len(s.AuthValue.ProviderName) == 0 {
		s.AuthValue.ProviderName = StaticProviderName
	}
	return s.AuthValue, nil
}

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s *StaticProvider) IsExpired() bool {
	return false
}
