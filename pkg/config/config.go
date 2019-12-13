package config

import (
	"crypto/tls"
	"net/http"
)

// Region specifies the New Relic environment to target.
type Region int

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
	Region        Region
}
