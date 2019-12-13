package internal

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type NewRelicClient struct {
	Client resty.Client
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
	Environment   Environment
}

// Environment specifies the New Relic environment to target.
type Environment int

const (
	// Production represents New Relic's US-based production deployment.
	Production = iota

	// EU represents New Relic's EU-based production deployment.
	EU

	// Staging represents New Relic's US-based staging deployment.  This is for internal use only.
	Staging
)

func (e Environment) String() string {
	return [...]string{"production", "eu", "staging"}[e]
}

var defaultBaseURLs = map[Environment]string{
	Production: "https://api.newrelic.com/v2",
	EU:         "https://api.eu.newrelic.com/v2",
	Staging:    "https://staging-api.newrelic.com/v2",
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
		Client: *client,
		pager:  config.Pager,
	}
}

func setHostURL(config Config, client *resty.Client) {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURLs[config.Environment]
	}

	client.SetHostURL(config.BaseURL)
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

// Get executes an HTTP GET request.
func (nr *NewRelicClient) Get(path string, params *map[string]string, result interface{}) error {
	req := nr.Client.R()

	if result != nil {
		req.SetResult(result)
	}

	if params != nil {
		req.SetQueryParams(*params)
	}

	nextPath := path

	for nextPath != "" {
		paging, err := nr.do(http.MethodGet, nextPath, req)

		if err != nil {
			return err
		}

		nextPath = paging.Next
	}

	return nil
}

// nolint
func (nr *NewRelicClient) Put(path string, body interface{}, result interface{}) error {
	req := nr.Client.R().
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
	req := nr.Client.R().
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

func (nr *NewRelicClient) do(method string, path string, req *resty.Request) (*Paging, error) {
	if req == nil {
		req = nr.Client.R()
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
