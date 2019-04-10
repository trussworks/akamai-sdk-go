package credentials

import (
	"os"
	"testing"
)

func TestSharedCredentialsProvider(t *testing.T) {
	os.Clearenv()

	p := SharedCredentialsProvider{Filename: "example_edgerc", Profile: ""}
	creds, err := p.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "clientSecret", creds.ClientSecret; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "clientToken", creds.ClientToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "accessToken", creds.AccessToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "akamaiHost", creds.Host; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestSharedCredentialsProviderIsExpired(t *testing.T) {
	os.Clearenv()

	p := SharedCredentialsProvider{Filename: "example_edgerc", Profile: ""}

	if !p.IsExpired() {
		t.Errorf("Expect creds to be expired before retrieve")
	}

	_, err := p.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if p.IsExpired() {
		t.Errorf("Expect creds to not be expired after retrieve")
	}
}

func TestSharedCredentialsProviderWithAKAMAI_ENVRC_FILE(t *testing.T) {
	os.Clearenv()
	os.Setenv("AKAMAI_ENVRC_FILE", "example_edgerc")
	p := SharedCredentialsProvider{}
	creds, err := p.Retrieve()

	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "clientSecret", creds.ClientSecret; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "clientToken", creds.ClientToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "accessToken", creds.AccessToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "akamaiHost", creds.Host; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestSharedCredentialsProviderWithoutHostFromProfile(t *testing.T) {
	os.Clearenv()

	p := SharedCredentialsProvider{Filename: "example_edgerc", Profile: "no_host"}
	creds, _ := p.Retrieve()

	if v := creds.Host; len(v) != 0 {
		t.Errorf("Expect no host, %v", v)
	}
}
