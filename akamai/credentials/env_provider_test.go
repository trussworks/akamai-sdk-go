package credentials

import (
	"os"
	"testing"
)

func TestEnvProviderRetrieve(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_CLIENT_SECRET", "client_secret")
	os.Setenv("AKAMAI_CLIENT_TOKEN", "client_token")
	os.Setenv("AKAMAI_ACCESS_TOKEN", "access_token")
	os.Setenv("AKAMAI_HOST", "host")

	e := EnvProvider{}
	creds, err := e.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "client_secret", creds.ClientSecret; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "client_token", creds.ClientToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "access_token", creds.AccessToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "host", creds.Host; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEnvProviderNoClientToken(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_CLIENT_SECRET", "secret")

	e := EnvProvider{}
	_, err := e.Retrieve()
	if e, a := ErrClientTokenNotFoundEnv, err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEnvProviderNoClientSecret(t *testing.T) {
	os.Clearenv()

	e := EnvProvider{}
	_, err := e.Retrieve()
	if e, a := ErrClientSecretNotFoundEnv, err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEnvProviderNoAccessToken(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_CLIENT_SECRET", "secret")
	os.Setenv("AKAMAI_CLIENT_TOKEN", "token")

	e := EnvProvider{}
	_, err := e.Retrieve()
	if e, a := ErrAccessTokenNotFoundEnv, err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEnvProviderNoHost(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_CLIENT_SECRET", "secret")
	os.Setenv("AKAMAI_CLIENT_TOKEN", "token")
	os.Setenv("AKAMAI_ACCESS_TOKEN", "access")

	e := EnvProvider{}
	_, err := e.Retrieve()
	if e, a := ErrAkamaiHostNotFoundEnv, err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEnvProviderIsExpired(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_CLIENT_SECRET", "secret")
	os.Setenv("AKAMAI_CLIENT_TOKEN", "token")
	os.Setenv("AKAMAI_ACCESS_TOKEN", "access")
	os.Setenv("AKAMAI_HOST", "host")

	e := EnvProvider{}

	if !e.IsExpired() {
		t.Errorf("Expect creds to be expired before retrieve.")
	}

	_, err := e.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e.IsExpired() {
		t.Errorf("Expect creds to not be expired after retrieve.")
	}
}
