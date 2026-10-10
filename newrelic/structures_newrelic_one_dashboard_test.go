//go:build unit

package newrelic

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
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

func TestValidateDashboardVariableOptions(t *testing.T) {
	cases := map[string]struct {
		config      map[string]interface{}
		meta        interface{}
		expectedErr string
	}{
		"enum variable allows show_apply_action": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "enum_variable",
				"replacement_strategy": "default",
				"title":                "Enum Variable",
				"type":                 "enum",
				"options": []interface{}{
					map[string]interface{}{
						"show_apply_action": true,
					},
				},
			}),
		},
		"enum variable rejects ignore_time_range when explicitly set": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "enum_variable",
				"replacement_strategy": "default",
				"title":                "Enum Variable",
				"type":                 "enum",
				"options": []interface{}{
					map[string]interface{}{
						"ignore_time_range": true,
						"show_apply_action": true,
					},
				},
			}),
			expectedErr: "`ignore_time_range` in `options` can only be used with the variable type `nrql`",
		},
		"enum variable allows ignore_time_range false": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "enum_variable",
				"replacement_strategy": "default",
				"title":                "Enum Variable",
				"type":                 "enum",
				"options": []interface{}{
					map[string]interface{}{
						"ignore_time_range": false,
						"show_apply_action": true,
					},
				},
			}),
		},
		"enum variable rejects excluded when true": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "enum_variable",
				"replacement_strategy": "default",
				"title":                "Enum Variable",
				"type":                 "enum",
				"options": []interface{}{
					map[string]interface{}{
						"excluded":          true,
						"show_apply_action": true,
					},
				},
			}),
			expectedErr: "`excluded` in `options` can only be used with the variable type `nrql`",
		},
		"enum variable allows excluded false": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "enum_variable",
				"replacement_strategy": "default",
				"title":                "Enum Variable",
				"type":                 "enum",
				"options": []interface{}{
					map[string]interface{}{
						"excluded":          false,
						"show_apply_action": true,
					},
				},
			}),
		},
		"nrql variable allows nrql-only options": {
			config: testDashboardVariableValidationConfig(map[string]interface{}{
				"is_multi_selection":   true,
				"name":                 "nrql_variable",
				"replacement_strategy": "default",
				"title":                "NRQL Variable",
				"type":                 "nrql",
				"nrql_query": []interface{}{
					map[string]interface{}{
						"account_ids": []interface{}{12345},
						"query":       "FROM Transaction SELECT uniques(appName)",
					},
				},
				"options": []interface{}{
					map[string]interface{}{
						"excluded":          true,
						"ignore_time_range": true,
						"show_apply_action": true,
					},
				},
			}),
			meta: &ProviderConfig{AccountID: 12345},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := testDashboardDiff(tc.config, tc.meta)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
				return
			}

			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestExpandVariableOptions(t *testing.T) {
	cases := map[string]struct {
		options                  []interface{}
		variableType             dashboards.DashboardVariableType
		expectIgnoreTimeRangeNil bool
		expectIgnoreTimeRange    bool
		expectExcludedNil        bool
		expectExcluded           bool
		expectShowApplyActionNil bool
		expectShowApplyAction    bool
	}{
		"enum keeps only show_apply_action": {
			options: []interface{}{
				map[string]interface{}{
					"excluded":          true,
					"ignore_time_range": true,
					"show_apply_action": true,
				},
			},
			variableType:             dashboards.DashboardVariableTypeTypes.ENUM,
			expectIgnoreTimeRangeNil: true,
			expectExcludedNil:        true,
			expectShowApplyAction:    true,
		},
		"nrql keeps all supported options": {
			options: []interface{}{
				map[string]interface{}{
					"excluded":          true,
					"ignore_time_range": true,
					"show_apply_action": true,
				},
			},
			variableType:          dashboards.DashboardVariableTypeTypes.NRQL,
			expectIgnoreTimeRange: true,
			expectExcluded:        true,
			expectShowApplyAction: true,
		},
		"enum preserves explicit false apply button": {
			options: []interface{}{
				map[string]interface{}{
					"show_apply_action": false,
				},
			},
			variableType:             dashboards.DashboardVariableTypeTypes.ENUM,
			expectIgnoreTimeRangeNil: true,
			expectExcludedNil:        true,
			expectShowApplyAction:    false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out := expandVariableOptions(tc.options, tc.variableType)
			if assert.NotNil(t, out) {
				if tc.expectIgnoreTimeRangeNil {
					assert.Nil(t, out.IgnoreTimeRange)
				} else if assert.NotNil(t, out.IgnoreTimeRange) {
					assert.Equal(t, tc.expectIgnoreTimeRange, *out.IgnoreTimeRange)
				}

				if tc.expectExcludedNil {
					assert.Nil(t, out.Excluded)
				} else if assert.NotNil(t, out.Excluded) {
					assert.Equal(t, tc.expectExcluded, *out.Excluded)
				}

				if tc.expectShowApplyActionNil {
					assert.Nil(t, out.ShowApplyAction)
				} else if assert.NotNil(t, out.ShowApplyAction) {
					assert.Equal(t, tc.expectShowApplyAction, *out.ShowApplyAction)
				}
			}
		})
	}
}

func testDashboardDiff(config map[string]interface{}, meta interface{}) error {
	resourceConfig := terraform.NewResourceConfigRaw(config)
	_, err := resourceNewRelicOneDashboard().Diff(context.Background(), nil, resourceConfig, meta)
	return err
}

func testDashboardVariableValidationConfig(variable map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name": "test-dashboard",
		"page": []interface{}{
			map[string]interface{}{
				"name": "test-page",
			},
		},
		"variable": []interface{}{variable},
	}
}
