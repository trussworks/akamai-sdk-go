package akamai

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/mojotalantikite/terraform-provider-akamai/sdk/credentials"
)

// Signer applies Akamai Edgegrid signing to a given request.
type Signer struct {
	Credentials *credentials.Credentials
}

// NewSigner returns a Signer pointer configured with the credentials
func NewSigner(credentials *credentials.Credentials) *Signer {
	a := &Signer{
		Credentials: credentials,
	}

	return a
}

// Sign signs Akamai requests with the provided body.
func (s *Signer) Sign(req *http.Request, body io.ReadSeeker) (http.Header, error) {
	creds, err := s.Credentials.Get()
	if err != nil {
		return http.Header{}, err

	}

	ctx := &signingCtx{
		Request:    req,
		Body:       body,
		Query:      req.URL.Query(),
		credValues: creds,
	}

	if err := ctx.build(); err != nil {
		return nil, err
	}

	return ctx.SignedHeaderVals, nil
}

type signingCtx struct {
	Request            *http.Request
	Body               io.ReadSeeker
	SignedHeaderVals   http.Header
	UnsignedHeaderVals http.Header
	Query              url.Values

	credValues    credentials.AuthValue
	formattedTime string
	nonce         string

	contentHash       string
	canonicalHeaders  string
	pathQuery         string
	authHeaders       string
	signedAuthHeaders string
	signingData       string
	signingKey        string
}

func (ctx *signingCtx) build() error {
	ctx.buildTime()      // no deps
	ctx.buildNonce()     // no deps
	ctx.buildPathQuery() // no deps

	//	ctx.buildCanonicalHeaders() // depends on UnsignedHeaderVals
	ctx.buildSigningKey()  // depends on credValues and formattedTime
	ctx.buildContentHash() // no deps

	ctx.buildAuthHeaders() // depends on formattedTime and nonce
	ctx.buildSigningData() // depends on many things

	ctx.buildSignedAuthHeaders() // depends on like everything

	ctx.Request.Header.Set("Authorization", ctx.signedAuthHeaders)
	return nil
}

func (ctx *signingCtx) buildTime() {
	// format the timestamp for the akamai edgegrid api request
	local := time.FixedZone("GMT", 0)
	t := time.Now().In(local)
	ctx.formattedTime = fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d+0000",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func (ctx *signingCtx) buildNonce() {
	// generate a nonce for the akamai edgegrid api request
	uuid, _ := uuid.NewRandom()
	ctx.nonce = uuid.String()
}

// createSignature is the base64-encoding of the SHA–256 HMAC of the data to sign with the signing key.
func createSignature(data string, key string) string {
	k := []byte(key)
	h := hmac.New(sha256.New, k)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// signingKey is derived from the client secret.
// The signing key is computed as the base64 encoding of the SHA–256 HMAC of the timestamp string
// (the field value included in the HTTP authorization header described above) with the client secret as the key.
func (ctx *signingCtx) buildSigningKey() {
	key := createSignature(ctx.formattedTime, ctx.credValues.ClientSecret)
	ctx.signingKey = key
}

// buildSigningData formats the HTTP request to ensure its acceptance by Akamai.
//
// The data to sign includes the information from the HTTP request that is relevant to ensuring that the request is authentic.
// This data set comprised of the request data combined with the authorization header value (excluding the signature field,
// but including the ; right before the signature field).
//
// Documentation: https://developer.akamai.com/legacy/introduction/Client_Auth.html
func (ctx *signingCtx) buildSigningData() {
	dataSign := []string{
		ctx.Request.Method,
		ctx.Request.URL.Scheme,
		ctx.Request.URL.Host,
		ctx.pathQuery,
		ctx.canonicalHeaders,
		ctx.contentHash,
		ctx.authHeaders,
	}
	ctx.signingData = strings.Join(dataSign, "\t")
}

func (ctx *signingCtx) buildPathQuery() {
	if ctx.Request.URL.RawQuery == "" {
		ctx.pathQuery = ctx.Request.URL.Path
	}
	ctx.pathQuery = fmt.Sprintf("%s?%s", ctx.Request.URL.Path, ctx.Request.URL.RawQuery)
}

func (ctx *signingCtx) buildCanonicalHeaders() {
	var unsortedHeader []string
	var sortedHeader []string
	for k := range ctx.Request.Header {
		unsortedHeader = append(unsortedHeader, k)
	}
	sort.Strings(unsortedHeader)
	for _, k := range unsortedHeader {
		v := strings.TrimSpace(ctx.Request.Header.Get(k))
		sortedHeader = append(sortedHeader, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(stringMinifier(v))))
	}
	ctx.canonicalHeaders = strings.Join(sortedHeader, "\t")
}

// buildContentHash is the base64-encoded SHA–256 hash of the POST body.
// For any other request methods, this field is empty. But the tac separator (\t) must be included.
// The size of the POST body must be less than or equal to the value specified by the service.
// Any request that does not meet this criteria SHOULD be rejected during the signing process,
// as the request will be rejected by EdgeGrid.
func (ctx *signingCtx) buildContentHash() {
	var (
		preparedBody string
		bodyBytes    []byte
	)

	if ctx.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		preparedBody = string(bodyBytes)
	}

	if ctx.Request.Method == "POST" && len(preparedBody) > 0 {
		// MaxBody is set in edgegrid Go library, but wasn't found in docs. Set to 131072.
		if len(preparedBody) > 131072 {
			preparedBody = preparedBody[0:131072]
		}
		h := sha256.Sum256([]byte(preparedBody))
		ctx.contentHash = base64.StdEncoding.EncodeToString(h[:])
	}

}

func (ctx *signingCtx) buildAuthHeaders() {
	ctx.authHeaders = fmt.Sprintf("EG1-HMAC-SHA256 client_token=%s;access_token=%s;timestamp=%s;nonce=%s;",
		ctx.credValues.ClientToken,
		ctx.credValues.AccessToken,
		ctx.formattedTime,
		ctx.nonce,
	)
}

// buildSignedAuthHeaders puts it all together
func (ctx *signingCtx) buildSignedAuthHeaders() {
	signature := createSignature(ctx.signingData, ctx.signingKey)

	ctx.signedAuthHeaders = fmt.Sprintf("%ssignature=%s", ctx.authHeaders, signature)
}

func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}
