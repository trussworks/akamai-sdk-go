package credentials

import "testing"

type stubProvider struct {
	creds   AuthValue
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (AuthValue, error) {
	s.expired = false
	s.creds.ProviderName = "stubProvider"
	return s.creds, s.err
}

func (s *stubProvider) IsExpired() bool {
	return s.expired
}

func TestCredentialsGet(t *testing.T) {
	c := NewCredentials(&stubProvider{
		creds: AuthValue{
			ClientSecret: "client_secret",
			ClientToken:  "client_token",
			AccessToken:  "access_token",
			Host:         "host",
		},
		expired: true,
	})

	creds, err := c.Get()
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

func TestCredentialsExpire(t *testing.T) {
	stub := &stubProvider{}
	c := NewCredentials(stub)

	stub.expired = false
	if !c.IsExpired() {
		t.Errorf("Expected to start out expired")
	}
	c.Expire()
	if !c.IsExpired() {
		t.Errorf("Expected to be expired")
	}

	c.forceRefresh = false
	if c.IsExpired() {
		t.Errorf("Expected not to be expired")
	}

	stub.expired = true
	if !c.IsExpired() {
		t.Errorf("Expected to be expired")
	}
}

func TestCredentialsGetWithProviderName(t *testing.T) {
	stub := &stubProvider{}

	c := NewCredentials(stub)

	creds, err := c.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e, a := creds.ProviderName, "stubProvider"; e != a {
		t.Errorf("Expected provider name to match, %v got %v", e, a)
	}
}
