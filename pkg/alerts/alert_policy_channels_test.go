package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func newMockResponse(
	t *testing.T,
	mockJsonResponse string,
	statusCode int,
) Alerts {
	return newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJsonResponse))

		require.NoError(t, err)
	}))
}
