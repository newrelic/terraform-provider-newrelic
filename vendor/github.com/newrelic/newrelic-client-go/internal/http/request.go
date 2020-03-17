package http

import (
	"fmt"
	neturl "net/url"

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
	request      *retryablehttp.Request
}

// NewRequest creates a new Request struct.
func NewRequest(c Client, method string, url string, params interface{}, reqBody interface{}, value interface{}) (*Request, error) {
	// Make a copy of the client's config
	cfg := c.config

	req := &Request{
		method:       method,
		url:          url,
		params:       params,
		reqBody:      reqBody,
		value:        value,
		authStrategy: c.authStrategy,
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

// SetServiceName sets the service name for the request.
func (r *Request) SetServiceName(serviceName string) {
	serviceName = fmt.Sprintf("%s|%s", serviceName, defaultServiceName)
	r.SetHeader(defaultNewRelicRequestingServiceHeader, serviceName)
}

func (r *Request) makeURL() (*neturl.URL, error) {
	u, err := neturl.Parse(r.url)

	if err != nil {
		return nil, err
	}

	if u.Host != "" {
		return u, nil
	}

	u, err = neturl.Parse(r.config.BaseURL + u.Path)

	if err != nil {
		return nil, err
	}

	return u, err
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
