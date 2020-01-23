package alerts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TODO: This is used by incidents_test.go still, need to refactor
// nolint
func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:                "abc123",
		BaseURL:               ts.URL,
		InfrastructureBaseURL: ts.URL,
		UserAgent:             "newrelic/newrelic-client-go",
	})

	return c
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Alerts {
	ts := mock.NewMockServer(t, mockJSONResponse, statusCode)

	return New(config.Config{
		APIKey:                "abc123",
		BaseURL:               ts.URL,
		InfrastructureBaseURL: ts.URL,
		UserAgent:             "newrelic/newrelic-client-go",
	})
}

// nolint
func newIntegrationTestClient(t *testing.T) Alerts {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})
}

func TestSetBaseURL(t *testing.T) {
	a := New(config.Config{
		BaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.client.Config.BaseURL)
}

func TestSetInfrastructureBaseURL(t *testing.T) {
	a := New(config.Config{
		InfrastructureBaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.infraClient.Config.BaseURL)
}
