package newrelic

import (
	"errors"
	"net/http"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

// NewRelic is a collection of New Relic APIs.
type NewRelic struct {
	Alerts     alerts.Alerts
	APM        apm.APM
	Dashboards dashboards.Dashboards
	Plugins    plugins.Plugins
	Synthetics synthetics.Synthetics
}

// New returns a collection of New Relic APIs.
func New(apiKey string, opts ...ConfigOption) (*NewRelic, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey required")
	}

	config := config.Config{
		APIKey: apiKey,
	}

	// Loop through config options
	for _, fn := range opts {
		if nil != fn {
			if err := fn(&config); err != nil {
				return nil, err
			}
		}
	}

	nr := &NewRelic{
		Alerts:     alerts.New(config),
		APM:        apm.New(config),
		Dashboards: dashboards.New(config),
		Plugins:    plugins.New(config),
		Synthetics: synthetics.New(config),
	}

	return nr, nil
}

// ConfigOption configures the Config when provided to NewApplication.
type ConfigOption func(*config.Config) error

// ConfigRegion sets the New Relic Region this client will use
func ConfigRegion(region string) ConfigOption {
	return func(cfg *config.Config) error {
		cfg.Region = region
		return nil
	}
}

// ConfigHTTPTimeout sets the timeout for HTTP requests
func ConfigHTTPTimeout(t time.Duration) ConfigOption {
	return func(cfg *config.Config) error {
		var timeout = &t
		cfg.Timeout = timeout
		return nil
	}
}

// ConfigHTTPTransport sets the HTTP Transporter
func ConfigHTTPTransport(transport *http.RoundTripper) ConfigOption {
	return func(cfg *config.Config) error {
		if transport != nil {
			cfg.HTTPTransport = transport
			return nil
		}

		return errors.New("HTTP Transport can not be nil")
	}
}

// ConfigUserAgent sets the HTTP UserAgent for API requests
func ConfigUserAgent(ua string) ConfigOption {
	return func(cfg *config.Config) error {
		if ua != "" {
			cfg.UserAgent = ua
			return nil
		}

		return errors.New("user-agent can not be empty")
	}
}

// ConfigBaseURL sets the Base URL used to make requests
func ConfigBaseURL(url string) ConfigOption {
	return func(cfg *config.Config) error {
		if url != "" {
			cfg.BaseURL = url
			return nil
		}

		return errors.New("base URL can not be empty")
	}
}
