package internal

import (
	"net/http"
	"net/http/httptest"
)

func NewTestAPIClient(handler http.Handler) NewRelicClient {
	ts := httptest.NewServer(handler)

	c := NewClient(Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}
