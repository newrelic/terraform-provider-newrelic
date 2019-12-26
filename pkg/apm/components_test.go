// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testComponent     = Component{}
	testComponentJSON = `{}`
)

func TestListComponents(t *testing.T) {
	apm := newMockResponse(t, testComponentJSON, http.StatusOK)
	c, err := apm.ListComponents(nil)

	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestListComponentsWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "componentName"
	expectedIDs := "123,456"
	expectedPluginID := "1234"
	expectedHealthStatus := "true"

	apm := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", recieved: "%s"`, expectedName, name)
		}

		ids := values.Get("filter[ids]")
		if ids != expectedIDs {
			t.Errorf(`expected ID filter "%s", recieved: "%s"`, expectedIDs, ids)
		}

		pluginID := values.Get("filter[plugin_id]")
		if pluginID != expectedPluginID {
			t.Errorf(`expected plugin ID filter "%s", recieved: "%s"`, expectedPluginID, pluginID)
		}

		healthStatus := values.Get("health_status")
		if healthStatus != expectedHealthStatus {
			t.Errorf(`expected health status filter "%s", recieved: "%s"`, expectedHealthStatus, healthStatus)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"applications":[]}`))

		require.NoError(t, err)
	}))

	params := ListComponentsParams{
		IDs:          []int{123, 456},
		PluginID:     1234,
		Name:         expectedName,
		HealthStatus: true,
	}

	c, err := apm.ListComponents(&params)

	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestShowComponents(t *testing.T) {
	apm := newMockResponse(t, testComponentJSON, http.StatusOK)
	c, err := apm.GetComponent(testComponent.ID)

	require.NoError(t, err)
	require.NotNil(t, c)
}
