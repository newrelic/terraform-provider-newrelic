package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/internal/version"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

const (
	defaultNewRelicRequestingServiceHeader = "NewRelic-Requesting-Services"
	defaultServiceName                     = "newrelic-client-go"
	defaultTimeout                         = time.Second * 30
	defaultRetryMax                        = 3
)

var (
	defaultUserAgent = fmt.Sprintf("newrelic/%s/%s (https://github.com/newrelic/%s)", defaultServiceName, version.Version, defaultServiceName)
)

// Client represents a client for communicating with the New Relic APIs.
type Client struct {
	// client represents the underlying HTTP client.
	client *retryablehttp.Client

	// config is the HTTP client configuration.
	config config.Config

	// authStrategy allows us to use multiple authentication methods for API calls
	authStrategy RequestAuthorizer

	errorValue ErrorResponse
}

// NewClient is used to create a new instance of Client.
func NewClient(cfg config.Config) Client {
	c := http.Client{
		Timeout: defaultTimeout,
	}

	if cfg.Timeout != nil {
		c.Timeout = *cfg.Timeout
	}

	if cfg.HTTPTransport != nil {
		if transport, ok := (cfg.HTTPTransport).(*http.Transport); ok {
			c.Transport = transport
		}
	} else {
		c.Transport = http.DefaultTransport
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = region.DefaultBaseURLs[region.Parse(string(cfg.Region))]
	}

	if cfg.NerdGraphBaseURL == "" {
		cfg.NerdGraphBaseURL = region.NerdGraphBaseURLs[region.Parse(string(cfg.Region))]
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = defaultUserAgent
	}

	// Either set or append the library name
	if cfg.ServiceName == "" {
		cfg.ServiceName = defaultServiceName
	} else {
		cfg.ServiceName = fmt.Sprintf("%s|%s", cfg.ServiceName, defaultServiceName)
	}

	r := retryablehttp.NewClient()
	r.HTTPClient = &c
	r.RetryMax = defaultRetryMax
	r.CheckRetry = RetryPolicy

	// Disable logging in go-retryablehttp since we are logging requests directly here
	r.Logger = nil

	return Client{
		authStrategy: &ClassicV2Authorizer{},
		client:       r,
		config:       cfg,
		errorValue:   &DefaultErrorResponse{},
	}
}

// SetAuthStrategy is used to set the default auth strategy for this client
// which can be overridden per request
func (c *Client) SetAuthStrategy(da RequestAuthorizer) {
	c.authStrategy = da
}

// SetErrorValue is used to unmarshal error body responses in JSON format.
func (c *Client) SetErrorValue(v ErrorResponse) *Client {
	c.errorValue = v
	return c
}

// Get represents an HTTP GET request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Get(
	url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := NewRequest(*c, http.MethodGet, url, queryParams, nil, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Post represents an HTTP POST request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The reqBody argument will be marshaled to JSON from the type provided and included in the request body.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Post(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := NewRequest(*c, http.MethodPost, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// RawPost behaves the same as Post, but without marshaling the body into JSON before making the request.
// This is required at least in the case of Synthetics Labels, since the POST doesn't handle JSON.
func (c *Client) RawPost(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {

	var requestBody []byte

	switch val := reqBody.(type) {
	case []byte:
		requestBody = val
	case string:
		requestBody = []byte(val)
	default:
		return nil, errors.New("invalid request body")
	}

	req, err := NewRequest(*c, http.MethodPost, url, queryParams, requestBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Put represents an HTTP PUT request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The reqBody argument will be marshaled to JSON from the type provided and included in the request body.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Put(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := NewRequest(*c, http.MethodPut, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Delete represents an HTTP DELETE request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Delete(url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := NewRequest(*c, http.MethodDelete, url, queryParams, nil, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Do initiates an HTTP request as configured by the passed Request struct.
func (c *Client) Do(req *Request) (*http.Response, error) {
	r, err := req.makeRequest()
	if err != nil {
		return nil, err
	}

	c.config.GetLogger().Debug("performing request", "method", req.method, "url", r.URL)

	logHeaders, err := json.Marshal(r.Header)
	if err != nil {
		return nil, err
	}

	c.config.GetLogger().Trace("request details", "headers", string(logHeaders), "body", req.reqBody)

	resp, retryErr := c.client.Do(r)
	if retryErr != nil {
		return nil, retryErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &nrErrors.NotFound{}
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	logHeaders, err = json.Marshal(resp.Header)
	if err != nil {
		return nil, err
	}

	c.config.GetLogger().Trace("request completed", "method", req.method, "url", r.URL, "status_code", resp.StatusCode, "headers", string(logHeaders), "body", string(body))

	errorValue := c.errorValue.New()
	_ = json.Unmarshal(body, &errorValue)

	if !isResponseSuccess(resp) {
		return nil, nrErrors.NewUnexpectedStatusCode(resp.StatusCode, errorValue.Error())
	}

	if errorValue.Error() != "" {
		return nil, errors.New(errorValue.Error())
	}

	if req.value == nil {
		return resp, nil
	}

	jsonErr := json.Unmarshal(body, req.value)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return resp, nil
}

// Ensures the response status code falls within the
// status codes that are commonly considered successful.
func isResponseSuccess(resp *http.Response) bool {
	statusCode := resp.StatusCode

	return statusCode >= http.StatusOK && statusCode <= 299
}

func makeRequestBodyReader(reqBody interface{}) (*bytes.Buffer, error) {
	if reqBody == nil {
		return nil, nil
	}

	j, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(j)

	return b, nil
}

// GetBaseURL returns the BaseURL used by the client
func (c *Client) GetBaseURL() string {
	return c.config.BaseURL
}

// Query runs a graphQL query.
func (c *Client) Query(query string, vars map[string]interface{}, respBody interface{}) error {
	graphqlReqBody := &graphQLRequest{
		Query:     query,
		Variables: vars,
	}

	graphqlRespBody := &graphQLResponse{
		Data: respBody,
	}

	req, err := NewRequest(*c, http.MethodPost, c.config.NerdGraphBaseURL, nil, graphqlReqBody, graphqlRespBody)
	if err != nil {
		return err
	}

	req.SetAuthStrategy(&NerdGraphAuthorizer{})
	c.SetErrorValue(&graphQLErrorResponse{})

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	return nil
}
