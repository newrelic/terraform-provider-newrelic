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

// TestFlattenVariableNRQLQuery_AccountIDsNormalization covers the read-side
// normalization that prevents drift for a defaulted account_ids: when the user
// leaves account_ids empty, the upstream value equal to the provider account ID is
// written back as an empty list (matching the empty config), while any other upstream
// value — or any explicitly configured value — is written through unchanged so real
// changes still surface as drift.
func TestFlattenVariableNRQLQuery_AccountIDsNormalization(t *testing.T) {
	const providerAccountID = 111111
	const subAccountID = 222222

	tests := []struct {
		name             string
		apiAccountIDs    []int
		configAccountIDs []interface{}
		wantAccountIDs   interface{}
	}{
		{
			name:             "empty config + upstream == provider default -> normalized to empty",
			apiAccountIDs:    []int{providerAccountID},
			configAccountIDs: nil,
			wantAccountIDs:   []int{},
		},
		{
			name:             "empty config + upstream differs (UI change) -> written through (drift)",
			apiAccountIDs:    []int{subAccountID},
			configAccountIDs: []interface{}{},
			wantAccountIDs:   []int{subAccountID},
		},
		{
			name:             "empty config + upstream has multiple accounts -> written through",
			apiAccountIDs:    []int{providerAccountID, subAccountID},
			configAccountIDs: nil,
			wantAccountIDs:   []int{providerAccountID, subAccountID},
		},
		{
			name:             "explicit config == provider default -> written through as-is",
			apiAccountIDs:    []int{providerAccountID},
			configAccountIDs: []interface{}{providerAccountID},
			wantAccountIDs:   []int{providerAccountID},
		},
		{
			name:             "explicit config sub-account -> written through",
			apiAccountIDs:    []int{subAccountID},
			configAccountIDs: []interface{}{subAccountID},
			wantAccountIDs:   []int{subAccountID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &entities.DashboardVariableNRQLQuery{
				AccountIDs: tt.apiAccountIDs,
				Query:      "FROM Transaction SELECT uniques(appName)",
			}
			out := flattenVariableNRQLQuery(in, tt.configAccountIDs, providerAccountID)
			assert.Len(t, out, 1)
			block := out[0].(map[string]interface{})
			assert.Equal(t, tt.wantAccountIDs, block["account_ids"])
		})
	}
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
