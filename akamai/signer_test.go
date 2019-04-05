package akamai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/mojotalantikite/akamai-sdk-go/akamai/credentials"
	"github.com/stretchr/testify/assert"
)

var (
	testFile = "../testdata/testdata.json"

	akamaiTestHost         = "https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/"
	akamaiTestAccessToken  = "akab-access-token-xxx-xxxxxxxxxxxxxxxx"
	akamaiTestClientToken  = "akab-client-token-xxx-xxxxxxxxxxxxxxxx"
	akamaiTestClientSecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx="
	nonce                  = "nonce-xx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	timestamp              = "20140321T19:34:21+0000"
)

type JSONTests struct {
	Tests []Test `json:"tests"`
}

type Test struct {
	Name    string `json:"testName"`
	Request struct {
		Method  string              `json:"method"`
		Path    string              `json:"path"`
		Headers []map[string]string `json:"headers"`
		Data    string              `json:"data"`
	} `json:"request"`
	ExpectedAuthorization string `json:"expectedAuthorization"`
}

func TestBuildTime(t *testing.T) {
	var ctx signingCtx
	ctx.buildTime()
	time := ctx.formattedTime
	expected := regexp.MustCompile(`^\d{4}[0-1][0-9][0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\+0000$`)
	if assert.Regexp(t, expected, time, "Fail: Regex do not match") {
		t.Log("Pass: Regex matches")
	}
}

func TestBuildNonce(t *testing.T) {
	var ctx signingCtx
	ctx.buildNonce()
	first := ctx.nonce
	for i := 0; i < 100; i++ {
		ctx.buildNonce()
		second := ctx.nonce
		assert.NotEqual(t, first, second, "Fail: Nonce matches")
	}
}

func TestCreateAuthHeader(t *testing.T) {
	var edgegrid JSONTests
	byt, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Test file not found, err %s", err)
	}

	err = json.Unmarshal(byt, &edgegrid)
	if err != nil {
		t.Fatalf("JSON is not parsable, err %s", err)
	}

	url, err := url.Parse(akamaiTestHost)
	if err != nil {
		t.Fatalf("URL is not parsable, err %s", err)
	}

	creds := credentials.NewStaticCredentialsFromCreds(credentials.AuthValue{
		ClientSecret: akamaiTestClientSecret,
		ClientToken:  akamaiTestClientToken,
		AccessToken:  akamaiTestAccessToken,
		Host:         akamaiTestHost,
	})

	signer := NewSigner(creds)

	for _, edge := range edgegrid.Tests {
		url.Path = edge.Request.Path
		req, _ := http.NewRequest(
			edge.Request.Method,
			url.String(),
			bytes.NewBuffer([]byte(edge.Request.Data)),
		)

		for _, header := range edge.Request.Headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		signer.Timestamp = timestamp
		signer.Nonce = nonce

		signer.Sign(req, bytes.NewReader([]byte(edge.Request.Data)))

		t.Errorf("Expected: %s, received %s", edge.ExpectedAuthorization, req.Header.Get("Authorization"))

		/*
			if assert.Equal(t, edge.ExpectedAuthorization, actual, fmt.Sprintf("Fail: %s", edge.Name)) {
				t.Logf("Pass: %s\n", edge.Name)
				t.Logf("Expected: %s - Actual %s", edge.ExpectedAuthorization, actual)
			}
		*/
	}

}
