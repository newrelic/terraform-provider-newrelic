package newrelic

import (
	"crypto/tls"
	"net/http"

	"github.com/newrelic/newrelic-client-go/internal"
)

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

// Config contains all the configuration data for the API Client.
type Config struct {
	APIKey        string
	BaseURL       string
	ProxyURL      string
	Debug         bool
	TLSConfig     *tls.Config
	UserAgent     string
	HTTPTransport http.RoundTripper
	Environment   Environment
}

func (c *Config) ToInternal() internal.Config {
	return internal.Config{
		APIKey:      c.APIKey,
		BaseURL:     c.BaseURL,
		ProxyURL:    c.ProxyURL,
		Debug:       c.Debug,
		TLSConfig:   c.TLSConfig,
		UserAgent:   c.UserAgent,
		Environment: internal.Environment(c.Environment),
	}
}
