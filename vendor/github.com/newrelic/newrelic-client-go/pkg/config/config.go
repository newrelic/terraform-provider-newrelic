package config

import (
	"net/http"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/version"
)

// Config contains all the configuration data for the API Client.
type Config struct {
	// APIKey to authenticate API requests
	APIKey string

	// PersonalAPIKey to authenticate API requests
	PersonalAPIKey string

	// Region of the New Relic platform to use
	// Valid values are: US, EU
	Region string

	// HTTP
	Timeout       *time.Duration
	HTTPTransport *http.RoundTripper
	UserAgent     string
	BaseURL       string

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
