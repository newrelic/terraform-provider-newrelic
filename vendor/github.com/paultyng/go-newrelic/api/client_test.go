package api

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
)

func newTestAPIClient(handler http.Handler) *Client {
	ts := httptest.NewServer(handler)

	c := New(Config{
		APIKey:  "123456",
		BaseURL: ts.URL,
		Debug:   false,
	})

	return &c
}

func newTestAPIClientTLSConfig(handler http.Handler) *Client {
	ts := httptest.NewServer(handler)

	tlsCfg := &tls.Config{}
	tlsCfg.InsecureSkipVerify = true

	c := New(Config{
		APIKey:    "123456",
		BaseURL:   ts.URL,
		Debug:     false,
		TLSConfig: tlsCfg,
	})

	return &c
}
