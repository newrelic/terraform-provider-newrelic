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

func TestBaseURLs(t *testing.T) {
	t.Parallel()

	pairs := map[region.Region]string{
		region.US:      "https://synthetics.newrelic.com/synthetics/api/v3",
		region.EU:      "https://synthetics.eu.newrelic.com/synthetics/api/v3",
		region.Staging: "https://staging-synthetics.newrelic.com/synthetics/api/v3",
	}

	// Default should be region.US
	a := New(config.Config{})
	assert.Equal(t, pairs[region.US], a.client.Config.BaseURL)

	for k, v := range pairs {
		b := New(config.Config{
			Region: k.String(),
		})

		assert.Equal(t, v, b.client.Config.BaseURL)
	}
}

// TestError checks that messages are reported in the correct
// order by going through priorities backwards
func TestError(t *testing.T) {
	t.Parallel()

	// Default empty
	e := ErrorResponse{}
	assert.Equal(t, "", e.Error())

	// 3rd Messages concat
	e.Messages = []ErrorDetail{
		{Message: "detail message"},
		{Message: "another detail"},
	}
	assert.Equal(t, "detail message, another detail", e.Error())

	// 2nd Message
	e.Message = "message"
	assert.Equal(t, "message", e.Error())

	// 1st Server Error Message
	e.ServerErrorMessage = "server message"
	assert.Equal(t, "server message", e.Error())

}

// nolint
func newTestClient(handler http.Handler) Synthetics {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:            "abc123",
		SyntheticsBaseURL: ts.URL,
		UserAgent:         "newrelic/newrelic-client-go",
		LogLevel:          "debug",
	})

	return c
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Synthetics {
	return newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))
}

func TestSetSyntheticsBaseURL(t *testing.T) {
	a := New(config.Config{
		SyntheticsBaseURL: "http://localhost",
	})

	assert.Equal(t, "http://localhost", a.client.Config.BaseURL)
}
