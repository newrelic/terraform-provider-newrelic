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

// NewRelicClient represents a client for communicating with the New Relic APIs.
type NewRelicClient struct {
	// Client represents the underlying HTTP client.
	Client *retryablehttp.Client

	// Config is the HTTP client configuration.
	Config config.Config

	// AuthStrategy allows us to use multiple authentication methods for API calls
	AuthStrategy RequestAuthorizer

	errorValue ErrorResponse
}

// NewClient is used to create a new instance of NewRelicClient.
func NewClient(cfg config.Config) NewRelicClient {
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
		cfg.BaseURL = region.DefaultBaseURLs[region.Parse(cfg.Region)]
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

	return NewRelicClient{
		Client:       r,
		Config:       cfg,
		errorValue:   &DefaultErrorResponse{},
		AuthStrategy: &ClassicV2Authorizer{},
	}
}

// SetErrorValue is used to unmarshal error body responses in JSON format.
func (c *NewRelicClient) SetErrorValue(v ErrorResponse) *NewRelicClient {
	c.errorValue = v
	return c
}

// Get represents an HTTP GET request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *NewRelicClient) Get(
	url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, url, queryParams, nil, respBody)
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
func (c *NewRelicClient) Post(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPost, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// RawPost behaves the same as Post, but without marshaling the body into JSON before making the request.
// This is required at least in the case of Synthetics Labels, since the POST doesn't handle JSON.
func (c *NewRelicClient) RawPost(
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

	req, err := c.NewRequest(http.MethodPost, url, queryParams, requestBody, respBody)
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
func (c *NewRelicClient) Put(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPut, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Delete represents an HTTP DELETE request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *NewRelicClient) Delete(url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodDelete, url, queryParams, nil, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// NewRequest creates a new Request struct.
func (c *NewRelicClient) NewRequest(method string, url string, params interface{}, reqBody interface{}, value interface{}) (*Request, error) {
	// Make a copy of the client's config
	cfg := c.Config

	req := &Request{
		method:       method,
		url:          url,
		params:       params,
		reqBody:      reqBody,
		value:        value,
		authStrategy: c.AuthStrategy,
	}

	req.config = cfg

	u, err := req.makeURL()
	if err != nil {
		return nil, err
	}

	var r *retryablehttp.Request
	if reqBody != nil {
		if _, ok := reqBody.([]byte); !ok {
			reqBody, err = makeRequestBodyReader(reqBody)
			if err != nil {
				return nil, err
			}
		}

		r, err = retryablehttp.NewRequest(req.method, u.String(), reqBody)
		if err != nil {
			return nil, err
		}
	} else {
		r, err = retryablehttp.NewRequest(req.method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	req.request = r

	req.SetHeader(defaultNewRelicRequestingServiceHeader, defaultServiceName)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("User-Agent", cfg.UserAgent)

	return req, nil
}

// Do initiates an HTTP request as configured by the passed Request struct.
func (c *NewRelicClient) Do(req *Request) (*http.Response, error) {
	r, err := req.makeRequest()
	if err != nil {
		return nil, err
	}

	c.Config.GetLogger().Debug("performing request", "method", req.method, "url", r.URL)

	logHeaders, err := json.Marshal(r.Header)
	if err != nil {
		return nil, err
	}

	c.Config.GetLogger().Trace("request details", "headers", string(logHeaders), "body", req.reqBody)

	resp, retryErr := c.Client.Do(r)
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

	c.Config.GetLogger().Trace("request completed", "method", req.method, "url", r.URL, "status_code", resp.StatusCode, "headers", string(logHeaders), "body", string(body))

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
