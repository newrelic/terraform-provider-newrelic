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
		config       map[string]interface{}
		meta         interface{}
		expectedErr  string
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
