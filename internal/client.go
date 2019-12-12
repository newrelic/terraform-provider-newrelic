package internal

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type NewRelicClient struct {
	client resty.Client
	pager  Pager
}

// Config contains all the configuration data for the API Client.
type Config struct {
	APIKey        string
	BaseURL       string
	ProxyURL      string
	Debug         bool
	TLSConfig     *tls.Config
	UserAgent     string
	HTTPTransport http.RoundTripper
	Pager         Pager
}

func NewClient(config Config) NewRelicClient {
	client := resty.New()

	setHostURL(config, client)
	setProxyURL(config, client)
	setHeaders(config, client)
	setTLSConfig(config, client)
	setDebug(config, client)
	setHTTPTransport(config, client)

	if config.Pager == nil {
		config.Pager = &LinkHeaderPager{}
	}

	return NewRelicClient{
		client: *client,
		pager:  config.Pager,
	}
}

func setHostURL(config Config, client *resty.Client) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.newrelic.com/v2"
	}

	client.SetHostURL(baseURL)
}

func setProxyURL(config Config, client *resty.Client) {
	proxyURL := config.ProxyURL
	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}
}

func setHeaders(config Config, client *resty.Client) {
	userAgent := config.UserAgent
	if userAgent == "" {
		userAgent = fmt.Sprintf("newrelic/newrelic-client-go/%s (https://github.com/newrelic/newrelic-client-go)", "VERSION")
	}

	client.SetHeaders(map[string]string{
		"X-Api-Key":  config.APIKey,
		"User-Agent": userAgent,
	})
}

func setTLSConfig(config Config, client *resty.Client) {
	if config.TLSConfig != nil {
		client.SetTLSClientConfig(config.TLSConfig)
	}
}

func setDebug(config Config, client *resty.Client) {
	if config.Debug {
		client.SetDebug(true)
	}
}

func setHTTPTransport(config Config, client *resty.Client) {
	if config.HTTPTransport != nil {
		client.SetTransport(config.HTTPTransport)
	}
}

// nolint
func (nr *NewRelicClient) Get(path string, params *map[string]string, result interface{}) error {
	req := nr.client.R()

	if result != nil {
		req.SetResult(result)
	}

	if params != nil {
		req.SetQueryParams(*params)
	}

	nextPath := path

	for nextPath != "" {
		paging, err := nr.do(http.MethodGet, path, req)

		if err != nil {
			return err
		}

		nextPath = paging.Next
	}

	return nil
}

// nolint
func (nr *NewRelicClient) Put(path string, body interface{}, result interface{}) error {
	req := nr.client.R().
		SetBody(body).
		SetResult(result)

	_, err := nr.do(http.MethodPut, path, req)

	if err != nil {
		return err
	}

	return nil
}

// nolint
func (nr *NewRelicClient) Post(path string, body interface{}, result interface{}) error {
	req := nr.client.R().
		SetBody(body).
		SetResult(result)

	_, err := nr.do(http.MethodPost, path, req)

	if err != nil {
		return err
	}

	return nil
}

// nolint
func (nr *NewRelicClient) Delete(path string) error {
	_, err := nr.do(http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	return nil
}

// Do exectes an API request with the specified parameters.
func (nr *NewRelicClient) do(method string, path string, req *resty.Request) (*Paging, error) {
	if req == nil {
		req = nr.client.R()
	}

	req.SetError(&ErrorResponse{}).
		SetHeader("Content-Type", "application/json")

	apiResponse, err := req.Execute(method, path)

	if err != nil {
		return nil, err
	}

	paging := nr.pager.Parse(apiResponse)

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
		apiError := rawError.(*ErrorResponse)

		if apiError.Detail != nil {
			return nil, apiError
		}
	}

	return nil, fmt.Errorf("Unexpected status %v returned from API", apiResponse.StatusCode())
}
