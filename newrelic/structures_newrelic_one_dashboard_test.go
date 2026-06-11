//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestExpandVariableNRQLQuery_DefaultsToProviderAccountID(t *testing.T) {
	providerAccountID := 12345

	// account_ids omitted (empty slice) — should default to provider account
	in := []interface{}{
		map[string]interface{}{
			"account_ids": []interface{}{},
			"query":       "FROM Transaction SELECT uniques(appName)",
		},
	}
	out := expandVariableNRQLQuery(in, providerAccountID)
	assert.Equal(t, []int{providerAccountID}, out.AccountIDs)
}

func TestExpandVariableNRQLQuery_ExplicitAccountIDs(t *testing.T) {
	providerAccountID := 12345

	in := []interface{}{
		map[string]interface{}{
			"account_ids": []interface{}{111, 222},
			"query":       "FROM Transaction SELECT uniques(appName)",
		},
	}
	out := expandVariableNRQLQuery(in, providerAccountID)
	assert.Equal(t, []int{111, 222}, out.AccountIDs)
}

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
