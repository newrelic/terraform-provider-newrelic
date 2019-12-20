package alerts

import (
	"net/http"
	"net/http/httptest"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestAlerts(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}
