package trenergy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL  = "https://core.tr.energy"
	testNileBaseURL = "https://nile-core.tr.energy"
	testNileApiKey  = "145|bt5CAoruwRJ7M7rUQMFuTV80HKyeqYyG1TiiPMIIb4c67526"
)

// Client is the main entry point for the SDK.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
	headers    http.Header
}

// Option serves as a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL allows overriding the default base URL.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		if u, err := url.Parse(rawURL); err == nil {
			c.baseURL = u
		}
	}
}

// WithHTTPClient allows providing a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithHeader allows adding a default header to all requests.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers.Set(key, value)
	}
}

// WithTestNet enables the testnet environment.
func WithTestNet() Option {
	return func(c *Client) {
		u, _ := url.Parse(testNileBaseURL)
		c.baseURL = u
		// Use test key if no key provided
		if c.apiKey == "" {
			c.apiKey = testNileApiKey
		}
	}
}

// NewClient creates a new Client instance.
func NewClient(apiKey string, opts ...Option) *Client {
	u, _ := url.Parse(defaultBaseURL)
	c := &Client{
		baseURL:    u,
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
		headers:    make(http.Header),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.headers.Set("Accept", "application/json")
	if c.apiKey != "" {
		c.headers.Set("Authorization", "Bearer "+c.apiKey)
	}

	return c
}

// NewRequest creates an HTTP request.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	var contentType string

	if body != nil {
		// handle multipart/form-data if needed, for now assuming JSON or form-urlencoded based on samples
		// actually samples show both JSON and FormData.
		// Let's check typical usage. Most creation endpoints seem to use FormData.
		// For simplicity, let's implement checking if body implements a specific interface or just default to JSON
		// But wait, samples used `form-data` for creating consumers.
		// We can handle this by checking if the body is `url.Values` for form-urlencoded or a map/struct for JSON.
		// Or helper methods can handle serialization.
		// Let's default to JSON here and allow helper methods to pass already encoded reader if needed?
		// Actually, standard is to handle JSON serialization here.
		// If we encounter form-data endpoints, we might need a separate way to signal that.

		// For now, let's stick to JSON default, and if we need form-data we can handle it in specific methods or support it here.
		// Given the `samples.json` shows strict key-value pairs in form-data, let's look at `consumer.go` later.
		// Actually, `samples.json` shows requests with `Content-Type: application/json` for some GETs? No, GETs don't have body.
		// POSTs in samples:
		// consumer create: form-data
		// consumer activate: no body?
		// consumer mass trx: form-data
		// It seems MANY write operations use form-data.

		// Let's implement a "Form" helper or just treat `body` as JSON unless it is of a specific type.

		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range c.headers {
		req.Header[k] = v
	}

	return req, nil
}

// Do executes the request and decodes the response into v.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Try to decode error response
		// For now just return fmt.Errorf
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("API error: %s (status: %d) body: %s", resp.Status, resp.StatusCode, string(bodyBytes))
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				return resp, err
			}
		}
	}

	return resp, nil
}

// Helper for sending generic requests
func (c *Client) sendRequest(ctx context.Context, method, path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	_, err = c.Do(req, v)
	return err
}

// Separate helper for form-data if needed.
func (c *Client) postForm(ctx context.Context, path string, data url.Values, v interface{}) error {
	rel, err := url.Parse(path)
	if err != nil {
		return err
	}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range c.headers {
		req.Header[k] = v
	}

	_, err = c.Do(req, v)
	return err
}

// Helper for multipart form data
func (c *Client) postMultipart(ctx context.Context, path string, fields map[string]string, v interface{}) error {
	rel, err := url.Parse(path)
	if err != nil {
		return err
	}
	u := c.baseURL.ResolveReference(rel)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for k, v := range c.headers {
		req.Header[k] = v
	}

	_, err = c.Do(req, v)
	return err
}
