package http

import (
	"fmt"

	"github.com/google/go-querystring/query"
	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Request represents a configurable HTTP request.
type Request struct {
	method       string
	url          string
	params       interface{}
	reqBody      interface{}
	value        interface{}
	config       config.Config
	authStrategy RequestAuthorizer
	errorValue   ErrorResponse
	request      *retryablehttp.Request
}

// NewRequest creates a new Request struct.
func (c *Client) NewRequest(method string, url string, params interface{}, reqBody interface{}, value interface{}) (*Request, error) {
	var err error

	req := &Request{
		method:       method,
		url:          url,
		params:       params,
		reqBody:      reqBody,
		value:        value,
		authStrategy: c.authStrategy,
		errorValue:   c.errorValue,
	}

	// FIXME: We should remove this requirement on the request
	// Make a copy of the client's config
	cfg := c.config
	req.config = cfg

	if reqBody != nil {
		if _, ok := reqBody.([]byte); !ok {
			reqBody, err = makeRequestBodyReader(reqBody)
			if err != nil {
				return nil, err
			}
		}
	}

	req.request, err = retryablehttp.NewRequest(req.method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.SetHeader(defaultNewRelicRequestingServiceHeader, cfg.ServiceName)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("User-Agent", cfg.UserAgent)

	return req, nil
}

// SetHeader sets a header on the underlying request.
func (r *Request) SetHeader(key string, value string) {
	r.request.Header.Set(key, value)
}

// SetAuthStrategy sets the authentication strategy for the request.
func (r *Request) SetAuthStrategy(ra RequestAuthorizer) {
	r.authStrategy = ra
}

// SetErrorValue sets the error object for the request.
func (r *Request) SetErrorValue(e ErrorResponse) {
	r.errorValue = e
}

// SetServiceName sets the service name for the request.
func (r *Request) SetServiceName(serviceName string) {
	serviceName = fmt.Sprintf("%s|%s", serviceName, defaultServiceName)
	r.SetHeader(defaultNewRelicRequestingServiceHeader, serviceName)
}

func (r *Request) makeRequest() (*retryablehttp.Request, error) {
	r.authStrategy.AuthorizeRequest(r, &r.config)

	err := r.setQueryParams()
	if err != nil {
		return nil, err
	}

	return r.request, nil
}

func (r *Request) setQueryParams() error {
	if r.params == nil || len(r.request.URL.Query()) > 0 {
		return nil
	}

	q, err := query.Values(r.params)

	if err != nil {
		return err
	}

	r.request.URL.RawQuery = q.Encode()

	return nil
}
