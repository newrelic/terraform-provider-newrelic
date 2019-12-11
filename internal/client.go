package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
	"github.com/tomnomnom/linkheader"
)

type NewRelicClient struct {
	Client *resty.Client
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
}

func NewClient(config Config) NewRelicClient {
	client := resty.New()

	setHostURL(config, client)
	setProxyURL(config, client)
	setHeaders(config, client)
	setTLSConfig(config, client)
	setDebug(config, client)
	setHTTPTransport(config, client)

	return NewRelicClient{
		Client: client,
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
func (nr *NewRelicClient) Get(path string, response interface{}) (string, error) {
	return nr.do(http.MethodGet, path, nil, response)
}

// nolint
func (nr *NewRelicClient) Put(path string, body interface{}, response interface{}) (string, error) {
	return nr.do(http.MethodPut, path, body, response)
}

// nolint
func (nr *NewRelicClient) Post(path string, body interface{}, response interface{}) (string, error) {
	return nr.do(http.MethodPost, path, body, response)
}

// nolint
func (nr *NewRelicClient) Delete(path string) (string, error) {
	return nr.do(http.MethodDelete, path, nil, nil)
}

// Do exectes an API request with the specified parameters.
func (nr *NewRelicClient) do(method string, path string, body interface{}, response interface{}) (string, error) {
	client := nr.Client.R().
		SetError(&ErrorResponse{}).
		SetHeader("Content-Type", "application/json")

	if body != nil {
		client = client.SetBody(body)
	}

	if response != nil {
		client = client.SetResult(response)
	}

	apiResponse, err := client.Execute(method, path)

	if err != nil {
		return "", err
	}

	nextPath := ""
	header := apiResponse.Header().Get("Link")
	if header != "" {
		links := linkheader.Parse(header)

		for _, link := range links.FilterByRel("next") {
			nextPath = link.URL
			break
		}
	}

	apiResponseBody := apiResponse.Body()
	if nextPath == "" && apiResponseBody != nil && len(apiResponseBody) > 0 {

		linksResponse := struct {
			Links struct {
				Next string `json:"next"`
			} `json:"links"`
		}{}

		err = json.Unmarshal(apiResponseBody, &linksResponse)
		if err != nil {
			return "", err
		}

		if linksResponse.Links.Next != "" {
			nextPath = linksResponse.Links.Next
		}
	}

	statusCode := apiResponse.StatusCode()
	statusClass := statusCode / 100 % 10

	if statusClass == 2 {
		return nextPath, nil
	}

	if statusCode == 404 {
		return "", ErrNotFound
	}

	rawError := apiResponse.Error()

	if rawError != nil {
		apiError := rawError.(*ErrorResponse)

		if apiError.Detail != nil {
			return "", apiError
		}
	}

	return "", fmt.Errorf("Unexpected status %v returned from API", apiResponse.StatusCode())
}
