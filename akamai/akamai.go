package akamai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"

	"github.com/trussworks/akamai-sdk-go/akamai/credentials"
)

const (
	userAgent = "go-akamai"
)

// Client creates an Akamai client to make requests against the Akamai API.
type Client struct {
	// HTTP client used to make API calls.
	client *http.Client

	// User agent used when communicating with the API.
	UserAgent string

	// BaseURL contains the API URL.
	BaseURL *url.URL

	// Credentials object to use when signing requests.
	Credentials *credentials.Credentials

	// reuse a single struct rather than allocating one for each service on the heap
	common service

	// Services of the Akamai API.
	FastDNSv2 *FastDNSv2Service
}

type service struct {
	client *Client
}

// NewClient returns an Akamai API client.
// If no httpClient is provided, http.DefaultClient is used.
// The Akamai API uses a unique base URL that is generated for every API client.
// If this isn't set then there is no default URL we can fall back to and we
// have to return an error.
func NewClient(httpClient *http.Client, cc *credentials.Credentials) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// If no credentials are set, fall back to .edgerc file, as Akamai docs
	// all lean on the config file. Environment variables are available.
	if cc == nil {
		cc = credentials.NewSharedCredentials("", "default")
	}

	creds, err := cc.Get()
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Akamai authentication credentials: %v", err)
	}

	baseURL, err := url.Parse("https://" + creds.Host)
	if err != nil {
		return nil, err

	}

	// BaseURL needs a trailing slash for requests to be made
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path = baseURL.Path + "/"
	}

	c := &Client{
		client:      httpClient,
		BaseURL:     baseURL,
		Credentials: cc,
		UserAgent:   userAgent,
	}

	c.common.client = c
	c.FastDNSv2 = (*FastDNSv2Service)(&c.common)

	return c, nil
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// We need to sign the request. https://developer.akamai.com/legacy/introduction/Client_Auth.html
	signer := NewSigner(c.Credentials)
	if body != nil {
		b, err := ioutil.ReadAll(buf)
		if err != nil {
			return nil, err
		}
		signer.Sign(req, bytes.NewReader(b))

	} else {
		signer.Sign(req, nil)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response is an Akamai API response. It wraps http.Response and allows for us to add additional
// properties in the future.
type Response struct {
	*http.Response
}

// Do sends the API request and returns the API response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		// A nil ctx will cause a panic. Just use a background context.
		ctx = context.Background()
	}
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			return nil, e

		}

		return nil, err
	}

	defer resp.Body.Close()

	response := &Response{Response: resp}

	err = CheckResponse(resp)
	if err != nil {
		// AcceptedErrors are a special case. We return the response's payload.
		aerr, ok := err.(*AcceptedError)
		if ok {
			b, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return response, readErr
			}

			aerr.Raw = b
			return response, aerr
		}

		return response, err
	}

	// Do the actual copying into the interface
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// CheckResponse checks an API resonse for errors. If an error is found, it is returned.
// Errors are considered as anything outside of the 200 range of HTTP responses, with the exception
// being a 202 Accepted response.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return &AcceptedError{}
	}

	// if errors are within 200 range (but not 202, as above), then no error has occurred.
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// AcceptedError occurs when Akamai returns a 202 Accepted response. This means an asynchronous process
// has begun and is scheduled on the Akamai side.
// HTTP 202 is not an error, it's just used to indicate that the results are not ready yet, to check back soon.
type AcceptedError struct {
	// Response body's raw contents
	Raw []byte
}

func (*AcceptedError) Error() string {
	return "job scheduled with Akamai. check back later."
}

// Error wraps Akamai error responses.
// API error responses are outlined in the Akamai APIs, but aren't consistent and
// differ depending on service :( In order to deal with this we create an error type
// per API.
type Error struct {
	FastDNSv2 *FastDNSv2Error
}

// FastDNSv2Error is the error type of FastDNS v2 API.
type FastDNSv2Error struct {
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

func (e *FastDNSv2Error) Error() string {
	return fmt.Sprintf("%v.%v", e.Title, e.Detail)
}

// ErrorResponse holds the errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"`
	Errors   []Error        `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Errors)
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
