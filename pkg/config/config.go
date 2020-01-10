package config

import (
	"net/http"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
)

// DefaultBaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var DefaultBaseURLs = map[region.Region]string{
	region.US:      "https://api.newrelic.com/v2",
	region.EU:      "https://api.eu.newrelic.com/v2",
	region.Staging: "https://staging-api.newrelic.com/v2",
}

// Config contains all the configuration data for the API Client.
type Config struct {
	BaseURL       string
	APIKey        string
	Timeout       *time.Duration
	HTTPTransport *http.RoundTripper
	UserAgent     string
	Region        region.Region

	// LogLevel can be one of the following values:
	// "panic", "fatal", "error", "warn", "info", "debug", "trace"
	LogLevel string
	LogJSON  bool
	Logger   logging.Logger
}

// GetLogger returns a logger instance based on the config values.
func (c *Config) GetLogger() logging.Logger {
	if c.Logger != nil {
		return c.Logger
	}

	l := logging.NewStructuredLogger().
		SetDefaultFields(map[string]string{"newrelic-client-go": version.Version}).
		LogJSON(c.LogJSON).
		SetLogLevel(c.LogLevel)

	c.Logger = l
	return l
}
