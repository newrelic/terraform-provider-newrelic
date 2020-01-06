package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint
func newTestPluginsClient(handler http.Handler) Plugins {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

// nolint
func newMockResponse(
	t *testing.T,
	mockJSONResponse string,
	statusCode int,
) Plugins {
	return newTestPluginsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJSONResponse))

		require.NoError(t, err)
	}))
}

// nolint
func newIntegrationTestClient(t *testing.T) Plugins {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey: apiKey,
	})
}

var (
	testPluginJSON = `{
		"id": 999,
		"name": "Redis",
		"guid": "net.jondoe.newrelic_redis_plugin",
		"publisher": "Jon Doe",
		"summary_metrics": [
			{
				"id": 123,
				"name": "Connected Clients",
				"metric": "Component/Connection/Clients[connections]",
				"value_function": "average_value",
				"thresholds": {
					"caution": null,
					"critical": null
				}
			},
			{
				"id": 124,
				"name": "Rejected Connections",
				"metric": "Component/ConnectionRate/Rejected[connections]",
				"value_function": "average_value",
				"thresholds": {
					"caution": null,
					"critical": null
				}
			}
		]
	}`

	testPlugin = Plugin{
		ID:        999,
		Name:      "Redis",
		GUID:      "net.jondoe.newrelic_redis_plugin",
		Publisher: "Jon Doe",
		SummaryMetrics: []SummaryMetric{
			{
				ID:            123,
				Name:          "Connected Clients",
				Metric:        "Component/Connection/Clients[connections]",
				ValueFunction: "average_value",
				Thresholds:    MetricThreshold{},
			},
			{
				ID:            124,
				Name:          "Rejected Connections",
				Metric:        "Component/ConnectionRate/Rejected[connections]",
				ValueFunction: "average_value",
				Thresholds:    MetricThreshold{},
			},
		},
	}
)

func TestListPlugins(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"plugins": [%s]}`, testPluginJSON)
	client := newMockResponse(t, responseJSON, http.StatusOK)

	expected := []*Plugin{
		&testPlugin,
	}

	actual, err := client.ListPlugins(nil)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestListPluginsWithParams(t *testing.T) {
	t.Parallel()

	guidFilter := "net.jondoe.newrelic_redis_plugin"
	idsFilter := "999"

	client := newTestPluginsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[guid]")
		require.Equal(t, guidFilter, name)

		ids := values.Get("filter[ids]")
		require.Equal(t, idsFilter, ids)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(fmt.Sprintf(`{"plugins": [%s]}`, testPluginJSON)))

		require.NoError(t, err)
	}))

	params := ListPluginsParams{
		GUID: guidFilter,
		IDs:  []int{999},
	}

	expected := []*Plugin{
		&testPlugin,
	}

	actual, err := client.ListPlugins(&params)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
