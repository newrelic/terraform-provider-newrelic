package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// NewTestAPIClient returns a test NewRelicClient instance that is configured to communicate with a mock server.
func NewTestAPIClient(handler http.Handler) NewRelicClient {
	ts := httptest.NewServer(handler)

	c := NewClient(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}
