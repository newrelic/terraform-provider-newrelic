// +build unit

package synthetics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestDefaultEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{})

	actual := a.client.Config.BaseURL
	expected := "https://synthetics.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestEUEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.EU,
	})

	actual := a.client.Config.BaseURL
	expected := "https://synthetics.eu.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestStagingEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.Staging,
	})

	actual := a.client.Config.BaseURL
	expected := "https://staging-synthetics.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

// nolint
func newTestAlerts(handler http.Handler) Synthetics {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
		LogLevel:  "debug",
	})

	return c
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Synthetics {
	return newTestAlerts(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))
}
