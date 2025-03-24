package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type ClientConfig func(*Client) error
type apiRequest interface {
	Params() url.Values
}
type apiHeader interface {
	Headers() map[string]string
}
type ApiConfig struct {
	Host  string
	Path  string
	FPath string
}

var defaultRequestsPerSecond = 10

// Client may be used to make requests to the Google Maps WebService APIs
type Client struct {
	httpClient        *http.Client
	apiKey            string
	baseURL           string
	requestsPerSecond int
	rateLimiter       *rate.Limiter
}
type HttpConfig struct {
	Timeout   int32
	Transport http.Transport
}

func NewClient(configs ...ClientConfig) (*Client, error) {
	gc := &Client{
		requestsPerSecond: defaultRequestsPerSecond,
	}
	WithHTTPClient(&http.Client{})(gc)
	for _, config := range configs {
		err := config(gc)
		if err != nil {
			return nil, err
		}
	}
	if gc.apiKey == "" {
		return nil, errors.New("maps: API Key missing")
	}

	if gc.requestsPerSecond > 0 {
		// configure go token bucket rate limiter module
		gc.rateLimiter = rate.NewLimiter(rate.Limit(gc.requestsPerSecond), gc.requestsPerSecond)
	}
	return gc, nil
}

// WithHTTPClient configures a Maps API client with a http.Client to make requests
//
//	with transport layer configured to retry
func WithHTTPClient(c *http.Client) ClientConfig {
	return func(client *Client) error {
		c.Transport = &RetryRoundTripper{
			Transport:   http.DefaultTransport,
			MaxRetries:  3,                  // Max 3 retries
			RetryDelay:  2 * time.Second,    // 2 seconds delay between retries
			ShouldRetry: defaultShouldRetry, // Use default retry policy
		}
		client.httpClient = c
		return nil
	}
}
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := c.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req.WithContext(ctx))
}
func (c *Client) get(ctx context.Context, config *ApiConfig, apiReq apiRequest, apiHeader apiHeader) (*http.Response, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	if c.baseURL == "" {
		c.baseURL = config.Host + config.Path // get api base url concatenating host and path
	}
	req, err := http.NewRequest("GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	if apiHeader != nil { // if you want to set header in GET request, pass it from the caller
		headers := apiHeader.Headers()
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if apiReq != nil {
		q, err := c.setAuthQueryParam(apiReq.Params())
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = q
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}
	return resp, nil
}
func (c *Client) JsonGet(ctx context.Context, config *ApiConfig, apiReq apiRequest, apiHeader apiHeader, resp interface{}) error {
	httpResp, err := c.get(ctx, config, apiReq, apiHeader)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != http.StatusOK {
		return HttpError{Status: httpResp.StatusCode}
	}
	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(resp)

	return err
}
func (c *Client) JsonPost(ctx context.Context, config *ApiConfig, apiReq interface{}, apiHeader apiHeader, resp interface{}) error {
	httpResp, err := c.post(ctx, config, apiReq, apiHeader)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return HttpError{Status: httpResp.StatusCode}
	}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return err
}
func (c *Client) post(ctx context.Context, config *ApiConfig, apiReq interface{}, apiHeader apiHeader) (*http.Response, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	if c.baseURL == "" {
		c.baseURL = config.Host + config.Path // get api base url concatenating host and path
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	headers := apiHeader.Headers()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.do(ctx, req)
}

func (c *Client) setAuthQueryParam(q url.Values) (string, error) {
	if c.apiKey != "" {
		q.Set("key", c.apiKey)
		return q.Encode(), nil
	}
	return "", errors.New("maps: API Key missing")
}

// Add API Key configures a Maps API client with an API Key
func AddAPIKey(apiKey string) ClientConfig {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}
func (c *Client) waitRateLimit(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	return c.rateLimiter.Wait(ctx)
}

// WithBaseURL configures a Maps API client with a custom base url
func WithBaseURL(baseURL string) ClientConfig {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}
func WithRateLimit(requestsPerSecond int) ClientConfig {
	return func(c *Client) error {
		c.requestsPerSecond = requestsPerSecond
		return nil
	}
}

// RetryRoundTripper is a custom RoundTripper that retries failed requests.
type RetryRoundTripper struct {
	Transport   http.RoundTripper
	MaxRetries  int
	RetryDelay  time.Duration
	ShouldRetry func(*http.Response, error) bool
}

// RoundTrip executes an HTTP request with retry logic.
func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response

	// Retry loop
	for i := 0; i < r.MaxRetries; i++ {
		resp, err = r.Transport.RoundTrip(req)
		if err == nil && r.ShouldRetry(resp, err) {
			// If the response should be retried, wait and retry
			time.Sleep(r.RetryDelay)
			fmt.Printf("Retry attempt: %d\n", i)
			continue
		}
		break
	}
	return resp, err
}

// Default retry policy: Retry on server errors (5xx) or network errors
func defaultShouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		// Network or transport-level error, retry
		return true
	}
	if resp != nil && resp.StatusCode >= 500 {
		// Retry on 5xx server errors
		return true
	}
	return false
}
