// +build unit

package synthetics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestDefaultEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{})

	assert.Equal(t, BaseURLs[region.US], a.client.Config.BaseURL)
}

func TestUSEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: "US",
	})

	assert.Equal(t, BaseURLs[region.US], a.client.Config.BaseURL)
}

func TestEUEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: "EU",
	})

	assert.Equal(t, BaseURLs[region.EU], a.client.Config.BaseURL)
}

func TestStagingEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: "Staging",
	})

	assert.Equal(t, BaseURLs[region.Staging], a.client.Config.BaseURL)
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
