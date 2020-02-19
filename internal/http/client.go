package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/google/go-querystring/query"
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
	Client     *retryablehttp.Client
	Config     config.Config
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
		Client:     r,
		Config:     cfg,
		errorValue: &DefaultErrorResponse{},
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
	return c.do(http.MethodGet, url, queryParams, nil, respBody)
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

	reqBody, err := makeRequestBody(reqBody)
	if err != nil {
		return nil, err
	}

	return c.do(http.MethodPost, url, queryParams, reqBody, respBody)
}

// RawPost behaves the same as Post, but without marshaling the body into JSON before making the request.  This is required at least in the case of Syntheics Labels, since the POST doesn't handle JSON.
func (c *NewRelicClient) RawPost(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {

	switch val := reqBody.(type) {
	case []byte:
		return c.do(http.MethodPost, url, queryParams, reqBody, respBody)

	case string:
		requestBody := []byte(val)
		return c.do(http.MethodPost, url, queryParams, requestBody, respBody)

	default:
		return nil, errors.New("invalid request body")
	}

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

	reqBody, err := makeRequestBody(reqBody)
	if err != nil {
		return nil, err
	}

	return c.do(http.MethodPut, url, queryParams, reqBody, respBody)
}

// Delete represents an HTTP DELETE request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *NewRelicClient) Delete(url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	return c.do(http.MethodDelete, url, queryParams, nil, respBody)
}

func makeRequestBody(reqBody interface{}) (*bytes.Buffer, error) {
	b := bytes.NewBuffer([]byte{})
	if reqBody != nil {
		j, err := json.Marshal(reqBody)

		if err != nil {
			return nil, err
		}

		b = bytes.NewBuffer(j)
	}

	return b, nil
}

func (c *NewRelicClient) setHeaders(req *retryablehttp.Request) {
	req.Header.Set(defaultNewRelicRequestingServiceHeader, defaultServiceName)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Config.UserAgent)

	if c.Config.APIKey != "" {
		req.Header.Set("X-Api-Key", c.Config.APIKey)
	}

	if c.Config.PersonalAPIKey != "" {
		req.Header.Set("Api-Key", c.Config.PersonalAPIKey)
	}
}

func setQueryParams(req *retryablehttp.Request, params interface{}) error {
	if params == nil || len(req.URL.Query()) > 0 {
		return nil
	}

	q, err := query.Values(params)

	if err != nil {
		return err
	}

	req.URL.RawQuery = q.Encode()

	return nil
}

func (c *NewRelicClient) makeURL(url string) (*neturl.URL, error) {
	u, err := neturl.Parse(url)

	if err != nil {
		return nil, err
	}

	if u.Host != "" {
		return u, nil
	}

	u, err = neturl.Parse(c.Config.BaseURL + u.Path)

	if err != nil {
		return nil, err
	}

	return u, err
}

func (c *NewRelicClient) do(
	method string,
	url string,
	params interface{},
	reqBody interface{},
	value interface{},
) (*http.Response, error) {

	u, err := c.makeURL(url)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	err = setQueryParams(req, params)
	if err != nil {
		return nil, err
	}

	c.Config.GetLogger().Debug("performing request", "method", method, "url", req.URL)

	logHeaders, err := json.Marshal(req.Header)
	if err != nil {
		return nil, err
	}

	c.Config.GetLogger().Trace("request details", "headers", string(logHeaders), "body", reqBody)

	resp, retryErr := c.Client.Do(req)
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

	c.Config.GetLogger().Trace("request completed", "method", method, "url", req.URL, "status_code", resp.StatusCode, "headers", string(logHeaders), "body", string(body))

	errorValue := c.errorValue
	_ = json.Unmarshal(body, &errorValue)

	if !isResponseSuccess(resp) {
		return nil, nrErrors.NewUnexpectedStatusCode(resp.StatusCode, c.errorValue.Error())
	}

	if errorValue.Error() != "" {
		return nil, errors.New(errorValue.Error())
	}

	if value == nil {
		return resp, nil
	}

	jsonErr := json.Unmarshal(body, value)
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
