package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/newrelic/newrelic-client-go/internal/version"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

const (
	defaultTimeout  = time.Second * 30
	defaultRetryMax = 3
)

var (
	defaultUserAgent = fmt.Sprintf("newrelic/newrelic-client-go/%s (https://github.com/newrelic/newrelic-client-go)", version.Version)
	defaultBaseURLs  = map[config.RegionType]string{
		config.Region.US:      "https://api.newrelic.com/v2",
		config.Region.EU:      "https://api.eu.newrelic.com/v2",
		config.Region.Staging: "https://staging-api.newrelic.com/v2",
	}
)

type ReplacementClient struct {
	Client     *retryablehttp.Client
	Config     config.ReplacementConfig
	errorValue ErrorResponse
}

func NewReplacementClient(config config.ReplacementConfig) ReplacementClient {
	c := http.Client{
		Timeout: defaultTimeout,
	}

	if config.Timeout != nil {
		c.Timeout = *config.Timeout
	}

	if config.HTTPTransport != nil {
		if transport, ok := (*config.HTTPTransport).(*http.Transport); ok {
			c.Transport = transport
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURLs[config.Region]
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	r := retryablehttp.NewClient()
	r.HTTPClient = &c
	r.RetryMax = defaultRetryMax
	r.CheckRetry = RetryPolicy

	return ReplacementClient{
		Client:     r,
		Config:     config,
		errorValue: &DefaultErrorResponse{},
	}
}

func (c *ReplacementClient) SetErrorValue(v ErrorResponse) *ReplacementClient {
	c.errorValue = v
	return c
}

func (c *ReplacementClient) Get(url string, params map[string]string, reqBody interface{}, value interface{}) (*http.Response, error) {
	return c.do(http.MethodGet, url, params, reqBody, value)
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

func (c *ReplacementClient) setHeaders(req *retryablehttp.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.Config.APIKey)
	req.Header.Set("User-Agent", c.Config.UserAgent)
}

func setQueryParams(req *retryablehttp.Request, params map[string]string) {
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}
}

func (c *ReplacementClient) makeURL(url string) (*neturl.URL, error) {
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

func (c *ReplacementClient) do(method string, url string, params map[string]string, reqBody interface{}, value interface{}) (*http.Response, error) {
	reqBody, err := makeRequestBody(reqBody)

	if err != nil {
		return nil, err
	}

	u, err := c.makeURL(url)

	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(method, u.String(), reqBody)

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	setQueryParams(req, params)

	resp, retryErr := c.Client.Do(req)

	if retryErr != nil {
		return nil, retryErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &ErrorNotFound{}
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	if resp.StatusCode != http.StatusOK {
		errorValue := c.errorValue
		_ = json.Unmarshal(body, &errorValue)

		return nil, &ErrorUnexpectedStatusCode{
			statusCode: resp.StatusCode,
			err:        c.errorValue.Error()}
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
