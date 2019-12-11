package newrelic

import (
	"crypto/tls"
	"net/http"

	"github.com/newrelic/newrelic-client-go/internal"
)

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

func (c *Config) ToInternal() internal.Config {
	return internal.Config{
		APIKey:    c.APIKey,
		BaseURL:   c.BaseURL,
		ProxyURL:  c.ProxyURL,
		Debug:     c.Debug,
		TLSConfig: c.TLSConfig,
		UserAgent: c.UserAgent,
	}
}
