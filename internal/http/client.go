package http

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/newrelic/newrelic-client-go/internal/version"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// NewRelicClient is the internal client for communicating with the New Relic APIs.
type NewRelicClient struct {
	Client resty.Client
	pager  Pager
}

// NewClient is used to create a new instance of the NewRelicClient type.
func NewClient(config config.Config) NewRelicClient {
	client := resty.New()

	setHostURL(config, client)
	setProxyURL(config, client)
	setHeaders(config, client)
	setTLSConfig(config, client)
	setDebug(config, client)
	setHTTPTransport(config, client)

	client.SetError(&RestyErrorResponse{})

	c := NewRelicClient{
		Client: *client,
	}

	c.pager = &LinkHeaderPager{}

	return c
}

// SetPager allows for use of different pagination implementations.
func (n NewRelicClient) SetPager(pager Pager) NewRelicClient {
	n.pager = pager
	return n
}

// SetError allows for registering different well-known error response structures.
func (n NewRelicClient) SetError(err interface{}) NewRelicClient {
	n.Client.SetError(err)
	return n
}

func setHostURL(config config.Config, client *resty.Client) {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURLs[config.Region]
	}

	client.SetHostURL(config.BaseURL)
}

func setProxyURL(config config.Config, client *resty.Client) {
	proxyURL := config.ProxyURL
	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}
}

func setHeaders(config config.Config, client *resty.Client) {
	userAgent := config.UserAgent
	if userAgent == "" {
		userAgent = fmt.Sprintf("newrelic/newrelic-client-go/%s (https://github.com/newrelic/newrelic-client-go)", version.Version)
	}

	client.SetHeaders(map[string]string{
		"X-Api-Key":  config.APIKey,
		"User-Agent": userAgent,
	})
}

func setTLSConfig(config config.Config, client *resty.Client) {
	if config.TLSConfig != nil {
		client.SetTLSClientConfig(config.TLSConfig)
	}
}

func setDebug(config config.Config, client *resty.Client) {
	if config.Debug {
		client.SetDebug(true)
	}
}

func setHTTPTransport(config config.Config, client *resty.Client) {
	if config.HTTPTransport != nil {
		client.SetTransport(config.HTTPTransport)
	}
}

// Get executes an HTTP GET request.
func (n *NewRelicClient) Get(path string, params *map[string]string, result interface{}) error {
	req := n.Client.R()

	if result != nil {
		req.SetResult(result)
	}

	if params != nil {
		req.SetQueryParams(*params)
	}

	nextPath := path

	for nextPath != "" {
		paging, err := n.do(http.MethodGet, nextPath, req)

		if err != nil {
			return err
		}

		nextPath = paging.Next
	}

	return nil
}

// Put executes an HTTP PUT request.
// nolint
func (n *NewRelicClient) Put(path string, body interface{}, result interface{}) error {
	req := n.Client.R().
		SetBody(body).
		SetResult(result)

	_, err := n.do(http.MethodPut, path, req)

	if err != nil {
		return err
	}

	return nil
}

// Post executes an HTTP POST request.
// nolint
func (n *NewRelicClient) Post(path string, body interface{}, result interface{}) error {
	req := n.Client.R().
		SetBody(body).
		SetResult(result)

	_, err := n.do(http.MethodPost, path, req)

	if err != nil {
		return err
	}

	return nil
}

// Delete executes an HTTP DELETE request.
// nolint
func (n *NewRelicClient) Delete(path string) error {
	_, err := n.do(http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return nil
}

func (n *NewRelicClient) do(method string, path string, req *resty.Request) (*Paging, error) {
	if req == nil {
		req = n.Client.R()
	}

	req.SetHeader("Content-Type", "application/json")

	apiResponse, err := req.Execute(method, path)

	if err != nil {
		return nil, err
	}

	paging := n.pager.Parse(apiResponse.RawResponse)

	if err != nil {
		return nil, err
	}

	statusCode := apiResponse.StatusCode()
	statusClass := statusCode / 100 % 10

	if statusClass == 2 {
		return &paging, nil
	}

	if statusCode == 404 {
		return nil, ErrNotFound
	}

	rawError := apiResponse.Error()

	if rawError != nil {
		apiError := rawError.(*RestyErrorResponse)

		if apiError.Detail != nil {
			return nil, apiError
		}
	}

	return nil, fmt.Errorf("unexpected status %v returned from API", apiResponse.StatusCode())
}
