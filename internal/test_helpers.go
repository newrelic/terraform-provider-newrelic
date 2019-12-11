package client

import (
	"net/http"
	"net/http/httptest"

	client "github.com/newrelic/newrelic-client-go/internal"
)

func NewTestAPIClient(handler http.Handler) *client.NewRelicClient {
	ts := httptest.NewServer(handler)

	c := client.NewClient(client.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return &c
}
