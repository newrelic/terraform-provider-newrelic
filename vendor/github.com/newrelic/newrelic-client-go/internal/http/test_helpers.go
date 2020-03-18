package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var (
	testPersonalAPIKey = "apiKey"
	testAdminAPIKey    = "adminAPIKey"
	testUserAgent      = "userAgent"
)

// NewTestAPIClient returns a test Client instance that is configured to communicate with a mock server.
func NewTestAPIClient(handler http.Handler) Client {
	ts := httptest.NewServer(handler)

	c := NewClient(config.Config{
		PersonalAPIKey: testPersonalAPIKey,
		AdminAPIKey:    testAdminAPIKey,
		BaseURL:        ts.URL,
		UserAgent:      testUserAgent,
	})

	return c
}
