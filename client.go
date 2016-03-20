package tmdb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)


// Client may be used to make requests to the TMDB WebService APIs
type Client struct {
httpClient        *http.Client
apiKey            string
baseURL           string
requestsPerSecond int
rateLimiter       chan int
}

// ClientOption is the type of constructor options for NewClient(...).
type ClientOption func(*Client) error

var defaultRequestsPerSecond = 4

// NewClient constructs a new Client which can make requests to the TMDB WebService APIs.
func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{requestsPerSecond: defaultRequestsPerSecond}
	WithHTTPClient(&http.Client{})(c)
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.apiKey == ""  {
		return nil, errors.New("tmdb: API Key missing")
	}

	// Implement a bursty rate limiter.
	// Allow up to 1 second worth of requests to be made at once.
	c.rateLimiter = make(chan int, c.requestsPerSecond)
	// Prefill rateLimiter with 1 seconds worth of requests.
	for i := 0; i < c.requestsPerSecond; i++ {
		c.rateLimiter <- 1
	}
	go func() {
		// Refill rateLimiter continuously
		for _ = range time.Tick(time.Second / time.Duration(c.requestsPerSecond)) {
			c.rateLimiter <- 1
		}
	}()

	return c, nil
}

// WithHTTPClient configures a TMDB API client with a http.Client to make requests over.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = c
		return nil
	}
}

// WithAPIKey configures a TMDB API client with an API Key
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

type apiConfig struct {
	host            string
	path            string
}

type apiRequest interface {
	params() url.Values
}

func (c *Client) get(config *apiConfig, apiReq apiRequest) (*http.Response, error) {
	select {
	case <-c.rateLimiter:
	// Execute request.
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	q, err := c.generateQuery(config.path, apiReq.params())
	if err != nil {
		return nil, err
	}
	url := host+config.path+"?"+q
	return http.Get(url)
}

func (c *Client) getJSON(config *apiConfig, apiReq apiRequest, resp interface{}) error {
	httpResp, err := c.get(config, apiReq)
	if err != nil {
		return err
	}

	if (httpResp.StatusCode/100 != 2) {
		return errors.New(httpResp.Status)
	}

	defer httpResp.Body.Close()

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *Client) generateQuery(path string, q url.Values) (string, error) {
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
		return q.Encode(), nil
	}
	return "", errors.New("tmdb: API Key missing")
}
