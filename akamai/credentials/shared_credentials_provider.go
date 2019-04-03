// Example of using the shared credentials provider to read from ~/.edgerc
//
//     creds := credentials.NewSharedCredentials()
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
	"path/filepath"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

// SharedCredsProviderName provides a name of SharedCreds provider
const SharedCredsProviderName = "SharedCredentialsProvider"

var (
	ErrSharedCredentialsNotFoundFile    = errors.New(".edgerc file not found")
	ErrSharedCredentialsProfileNotFound = errors.New("could not load profile from .edgerc")
	ErrClientSecretNotFoundFile         = errors.New("client_secret not found in .edgerc")
	ErrClientTokenNotFoundFile          = errors.New("client_token not found in .edgerc")
	ErrAccessTokenNotFoundFile          = errors.New("access_token not found in .edgerc")
	ErrAkamaiHostNotFoundFile           = errors.New("host not found in .edgerc")
)

// SharedCredentialsProvider retrieves credentials from the current user's home
// directory, and keeps track if those credentials are expired.
//
// Documentation on edgerc: https://developer.akamai.com/legacy/introduction/Conf_Client.html
//
// Profile config file: $HOME/.edgerc
//
// Variables:
// client_secret
// client_token
// access_token
// host
type SharedCredentialsProvider struct {
	// Path to the shared credentials file.
	//
	// If empty will look for "AKAMAI_EDGERC_FILE" env variable. If the
	// env value is empty will default to current user's home directory.
	// Linux/OSX: "$HOME/.edgerc"
	Filename string

	// Edgerc file to extract credentials from. If empty
	// will default to environment variable "AKAMAI_EDGERC_PROFILE" or "default" if
	// environment variable is also not set.
	Profile string

	// retrieved states if the credentials have been successfully retrieved.
	retrieved bool
}

// NewSharedCredentials returns a pointer to a new Credentials object
// wrapping the Profile file provider.
func NewSharedCredentials(filename, profile string) *Credentials {
	return NewCredentials(&SharedCredentialsProvider{
		Filename: filename,
		Profile:  profile,
	})
}

// Retrieve reads and extracts the shared credentials from the current
// users home directory.
func (p *SharedCredentialsProvider) Retrieve() (AuthValue, error) {
	p.retrieved = false

	filename, err := p.filename()
	if err != nil {
		return AuthValue{ProviderName: SharedCredsProviderName}, err
	}

	creds, err := loadProfile(filename, p.profile())
	if err != nil {
		return AuthValue{ProviderName: SharedCredsProviderName}, err
	}

	p.retrieved = true
	return creds, nil
}

// IsExpired returns if the shared credentials have expired.
func (p *SharedCredentialsProvider) IsExpired() bool {
	return !p.retrieved
}

// loadProfiles loads from the file pointed to by shared credentials filename for profile.
// The credentials retrieved from the profile will be returned or error. Error will be
// returned if it fails to read from the file, or the data is invalid.
func loadProfile(filename, profile string) (AuthValue, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrSharedCredentialsNotFoundFile
	}

	iniProfile, err := config.GetSection(profile)
	if err != nil {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrSharedCredentialsProfileNotFound
	}

	cs, err := iniProfile.GetKey("client_secret")
	if err != nil || len(cs.String()) == 0 {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrClientSecretNotFoundFile
	}

	ct, err := iniProfile.GetKey("client_token")
	if err != nil || len(ct.String()) == 0 {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrClientTokenNotFoundFile
	}

	at, err := iniProfile.GetKey("access_token")
	if err != nil || len(at.String()) == 0 {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrClientTokenNotFoundFile
	}

	h, err := iniProfile.GetKey("host")
	if err != nil || len(h.String()) == 0 {
		return AuthValue{ProviderName: SharedCredsProviderName}, ErrAkamaiHostNotFoundFile
	}

	return AuthValue{
		ClientSecret: cs.String(),
		ClientToken:  ct.String(),
		AccessToken:  at.String(),
		Host:         h.String(),
		ProviderName: SharedCredsProviderName,
	}, nil
}

// filename returns the filename to use to read Akamai shared credentials.
// We use AKAMAI_ENVRC_FILE as the env variable to store this in.
// If not found will default to ~/.edgerc
//
// Will return an error if the user's home directory path cannot be found.
func (p *SharedCredentialsProvider) filename() (string, error) {
	if len(p.Filename) != 0 {
		return p.Filename, nil
	}

	if p.Filename = os.Getenv("AKAMAI_ENVRC_FILE"); len(p.Filename) != 0 {
		return p.Filename, nil
	}

	// try the default ~/.edgerc location
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.New("could not find user's homedir")
	}

	return filepath.Join(home, ".edgerc"), nil
}

// profile returns the Akamai shared credentials profile.  If empty will read
// environment variable "AKAMAI_PROFILE". If that is not set profile will
// return "default".
func (p *SharedCredentialsProvider) profile() string {
	if p.Profile == "" {
		p.Profile = os.Getenv("AKAMAI_PROFILE")
	}

	if p.Profile == "" {
		p.Profile = "default"
	}

	return p.Profile
}
