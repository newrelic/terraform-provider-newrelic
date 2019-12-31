// +build unit

package apm

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testTimestamp, _ = time.Parse(time.RFC3339, time.RFC3339)

	testComponent = Component{
		ID:   212222915,
		Name: "Redis",
		SummaryMetrics: []SummaryMetric{
			{
				ID:            190174,
				Name:          "Connected Clients",
				Metric:        "Component/Connection/Clients[connections]",
				ValueFunction: "average_value",
				Thresholds: MetricThreshold{
					Caution:  3,
					Critical: 0,
				},
			},
		},
	}

	testComponentJSON = `{
		"id": 212222915,
		"name": "Redis",
		"summary_metrics": [
			{
				"id": 190174,
				"name": "Connected Clients",
				"metric": "Component/Connection/Clients[connections]",
				"value_function": "average_value",
				"thresholds": {
					"caution": 3,
					"critical": 0
				}
			}
		]
	}`

	testComponentMetric = ComponentMetric{
		Name: "Component/Memory/RSS[bytes]",
		Values: []string{
			"average_value",
			"total_value",
			"max_value",
			"min_value",
			"standard_deviation",
			"rate",
			"count",
		},
	}

	testComponentMetricJSON = `{
		"name": "Component/Memory/RSS[bytes]",
		"values": [
			"average_value",
			"total_value",
			"max_value",
			"min_value",
			"standard_deviation",
			"rate",
			"count"
		]
	}`

	testMetricData = Metric{
		Name: "Component/Memory/RSS[bytes]",
		Timeslices: []MetricTimeslice{
			{
				From: &testTimestamp,
				To:   &testTimestamp,
				Values: map[string]float64{
					"average_value":      10,
					"total_value":        10,
					"max_value":          10,
					"min_value":          10,
					"standard_deviation": 10,
					"rate":               10,
					"count":              10,
				},
			},
		},
	}

	testMetricDataJSON = fmt.Sprintf(`{
        "name": "Component/Memory/RSS[bytes]",
        "timeslices": [
			{
				"from": "%[1]s",
				"to": "%[1]s",
				"values": {
					"average_value": 10,
					"total_value": 10,
					"max_value": 10,
					"min_value": 10,
					"standard_deviation": 10,
					"rate": 10,
					"count": 10
				}
			}
		]
	}`, testTimestamp.Format(time.RFC3339))
)

func TestListComponents(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{"components": [%s]}`, testComponentJSON)
	apm := newMockResponse(t, responseJSON, http.StatusOK)
	c, err := apm.ListComponents(nil)

	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, testComponent, c[0])
}

func TestListComponentsWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "componentName"
	expectedIDs := "123,456"
	expectedPluginID := "1234"
	expectedHealthStatus := "true"

	apm := newTestAPMClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		if name != expectedName {
			t.Errorf(`expected name filter "%s", received: "%s"`, expectedName, name)
		}

		ids := values.Get("filter[ids]")
		if ids != expectedIDs {
			t.Errorf(`expected ID filter "%s", received: "%s"`, expectedIDs, ids)
		}

		pluginID := values.Get("filter[plugin_id]")
		if pluginID != expectedPluginID {
			t.Errorf(`expected plugin ID filter "%s", received: "%s"`, expectedPluginID, pluginID)
		}

		healthStatus := values.Get("health_status")
		if healthStatus != expectedHealthStatus {
			t.Errorf(`expected health status filter "%s", received: "%s"`, expectedHealthStatus, healthStatus)
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

func TestGetComponent(t *testing.T) {
	responseJSON := fmt.Sprintf(`{"component": %s}`, testComponentJSON)
	apm := newMockResponse(t, responseJSON, http.StatusOK)
	c, err := apm.GetComponent(testComponent.ID)

	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, testComponent, *c)
}

func TestListComponentMetrics(t *testing.T) {
	responseJSON := fmt.Sprintf(`{"metrics": [%s]}`, testComponentMetricJSON)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	m, err := apm.ListComponentMetrics(testComponent.ID, nil)

	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, testComponentMetric, m[0])
}

func TestGetComponentMetricData(t *testing.T) {
	responseJSON := fmt.Sprintf(`{
		"metric_data": {
			 "metrics": [%s]
		}
	}`, testMetricDataJSON)
	apm := newMockResponse(t, responseJSON, http.StatusOK)
	m, err := apm.GetComponentMetricData(testComponent.ID, nil)

	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, testMetricData, m[0])
}

func TestGetComponentMetricDataWithParams(t *testing.T) {
	expectedNames := "componentName"
	expectedValues := "123"
	expectedTo := testTimestamp.Format(time.RFC3339)
	expectedFrom := testTimestamp.Format(time.RFC3339)
	expectedPeriod := "30"
	expectedSummarize := "true"
	expectedRaw := "true"

	apm := newTestAPMClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		names := values.Get("names[]")
		if names != expectedNames {
			t.Errorf(`expected names filter "%s", received: "%s"`, expectedNames, names)
		}

		v := values.Get("values[]")
		if v != expectedValues {
			t.Errorf(`expected values filter "%s", received: "%s"`, expectedValues, v)
		}

		from := values.Get("from")
		if from != expectedFrom {
			t.Errorf(`expected from param "%s", received: "%s"`, expectedFrom, from)
		}

		to := values.Get("to")
		if to != expectedTo {
			t.Errorf(`expected to param "%s", received: "%s"`, expectedTo, to)
		}

		period := values.Get("period")
		if period != expectedPeriod {
			t.Errorf(`expected period param "%s", received: "%s"`, expectedPeriod, period)
		}

		raw := values.Get("raw")
		if raw != expectedRaw {
			t.Errorf(`expected raw param "%s", received: "%s"`, expectedRaw, raw)
		}

		summarize := values.Get("summarize")
		if summarize != expectedSummarize {
			t.Errorf(`expected summarize param "%s", received: "%s"`, expectedSummarize, summarize)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"applications":[]}`))

		require.NoError(t, err)
	}))

	params := GetComponentMetricDataParams{
		Names:     []string{"componentName"},
		Values:    []string{"123"},
		From:      &testTimestamp,
		To:        &testTimestamp,
		Period:    30,
		Summarize: true,
		Raw:       true,
	}
	_, err := apm.GetComponentMetricData(testComponent.ID, &params)

	require.NoError(t, err)
}
