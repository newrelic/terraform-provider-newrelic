package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var (
	testAPIKey    = "apiKey"
	testUserAgent = "userAgent"
)

// NewTestAPIClient returns a test NewRelicClient instance that is configured to communicate with a mock server.
func NewTestAPIClient(handler http.Handler) NewRelicClient {
	ts := httptest.NewServer(handler)

	c := NewClient(&config.Config{
		APIKey:    testAPIKey,
		BaseURL:   ts.URL,
		UserAgent: testUserAgent,
	})

	return c
}
