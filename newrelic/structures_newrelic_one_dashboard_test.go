//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestExpandDashboardBillboardThreshold(t *testing.T) {
	dashboard := entities.DashboardWidget{
		ID: "abcde",
		Visualization: entities.DashboardWidgetVisualization{
			ID: "viz.billboard",
		},
		RawConfiguration: []byte(`
		{
			"facet": {
				"showOtherSeries": false
			},
		 	"nrqlQueries": [
				{
					"accountId": 1606862,
					"query": "FROM Transaction SELECT average(duration) WHERE appName = 'WebPortal' "
				}
			],
			"platformOptions": {
				"ignoreTimeRange": false
			},
		  	"thresholds": [
				{
					"alertSeverity": "WARNING",
					"value": 1
				},
				{
					"alertSeverity": "CRITICAL",
					"value": 2
				}
		  	]
		}
		`),
	}
	widgetType, out := flattenDashboardWidget(&dashboard, "abcde")
	assert.Equal(t, "widget_billboard", widgetType)
	assert.Contains(t, out, "nrql_query")
	assert.Contains(t, out, "critical")
	assert.Contains(t, out, "warning")
	assert.Equal(t, out["critical"], "2")
	assert.Equal(t, out["warning"], "1")
}

func TestExpandDashboardBillboardThresholdNullValue(t *testing.T) {
	dashboard := entities.DashboardWidget{
		ID: "abcde",
		Visualization: entities.DashboardWidgetVisualization{
			ID: "viz.billboard",
		},
		RawConfiguration: []byte(`
		{
			"facet": {
				"showOtherSeries": false
			},
		 	"nrqlQueries": [
				{
					"accountId": 1606862,
					"query": "FROM Transaction SELECT average(duration) WHERE appName = 'WebPortal' "
				}
			],
			"platformOptions": {
				"ignoreTimeRange": false
			},
		  	"thresholds": [
				{
					"alertSeverity": "WARNING",
					"value": null
				},
				{
					"alertSeverity": "CRITICAL",
					"value": 2
				}
		  	]
		}
		`),
	}
	widgetType, out := flattenDashboardWidget(&dashboard, "abcde")
	assert.Equal(t, "widget_billboard", widgetType)
	assert.Contains(t, out, "nrql_query")
	assert.Contains(t, out, "critical")
	assert.NotContains(t, out, "warning")
	assert.Equal(t, out["critical"], "2")
}
