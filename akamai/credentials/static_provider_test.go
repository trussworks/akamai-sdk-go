package credentials

import "testing"

func TestStaticProviderGet(t *testing.T) {
	s := StaticProvider{
		AuthValue: AuthValue{
			ClientSecret: "client_secret",
			ClientToken:  "client_token",
			AccessToken:  "access_token",
			Host:         "host",
		},
	}

	creds, err := s.Retrieve()
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

func TestStaticProviderIsExpired(t *testing.T) {
	s := StaticProvider{
		AuthValue: AuthValue{
			ClientSecret: "client_secret",
			ClientToken:  "client_token",
			AccessToken:  "access_token",
			Host:         "host",
		},
	}

	if s.IsExpired() {
		t.Errorf("Expect static credentials to never expire")
	}
}
