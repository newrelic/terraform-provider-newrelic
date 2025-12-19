//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandAlertCompoundConditionCreateInput(t *testing.T) {
	r := resourceNewRelicAlertCompoundCondition()

	cases := map[string]struct {
		Data     map[string]interface{}
		Expected *alerts.CompoundConditionCreateInput
	}{
		"basic compound condition": {
			Data: map[string]interface{}{
				"name":               "test-compound-condition",
				"enabled":            true,
				"trigger_expression": "A AND B",
				"component_conditions": []interface{}{
					map[string]interface{}{
						"id":    "123",
						"alias": "A",
					},
					map[string]interface{}{
						"id":    "456",
						"alias": "B",
					},
				},
			},
			Expected: &alerts.CompoundConditionCreateInput{
				Name:              "test-compound-condition",
				Enabled:           true,
				TriggerExpression: "A AND B",
				ComponentConditions: []alerts.ComponentConditionInput{
					{
						ID:    "123",
						Alias: "A",
					},
					{
						ID:    "456",
						Alias: "B",
					},
				},
			},
		},
		"with optional fields": {
			Data: map[string]interface{}{
				"name":                    "test-compound-condition",
				"enabled":                 false,
				"trigger_expression":      "(A AND B) OR C",
				"runbook_url":             "https://example.com/runbook",
				"threshold_duration":      120,
				"facet_matching_behavior": "FACETS_MATCH",
				"component_conditions": []interface{}{
					map[string]interface{}{
						"id":    "123",
						"alias": "A",
					},
					map[string]interface{}{
						"id":    "456",
						"alias": "B",
					},
					map[string]interface{}{
						"id":    "789",
						"alias": "C",
					},
				},
			},
			Expected: func() *alerts.CompoundConditionCreateInput {
				runbookURL := "https://example.com/runbook"
				thresholdDuration := 120
				facetBehavior := "FACETS_MATCH"
				return &alerts.CompoundConditionCreateInput{
					Name:                  "test-compound-condition",
					Enabled:               false,
					TriggerExpression:     "(A AND B) OR C",
					RunbookURL:            &runbookURL,
					ThresholdDuration:     &thresholdDuration,
					FacetMatchingBehavior: &facetBehavior,
					ComponentConditions: []alerts.ComponentConditionInput{
						{
							ID:    "123",
							Alias: "A",
						},
						{
							ID:    "456",
							Alias: "B",
						},
						{
							ID:    "789",
							Alias: "C",
						},
					},
				}
			}(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d := r.TestResourceData()

			for k, v := range tc.Data {
				if err := d.Set(k, v); err != nil {
					t.Fatalf("err: %s", err)
				}
			}

			expanded, err := expandAlertCompoundConditionCreateInput(d)
			require.NoError(t, err)

			if tc.Expected != nil {
				assert.Equal(t, tc.Expected.Name, expanded.Name)
				assert.Equal(t, tc.Expected.Enabled, expanded.Enabled)
				assert.Equal(t, tc.Expected.TriggerExpression, expanded.TriggerExpression)
				assert.Equal(t, len(tc.Expected.ComponentConditions), len(expanded.ComponentConditions))

				if tc.Expected.RunbookURL != nil {
					require.NotNil(t, expanded.RunbookURL)
					assert.Equal(t, *tc.Expected.RunbookURL, *expanded.RunbookURL)
				}
				if tc.Expected.ThresholdDuration != nil {
					require.NotNil(t, expanded.ThresholdDuration)
					assert.Equal(t, *tc.Expected.ThresholdDuration, *expanded.ThresholdDuration)
				}
			}
		})
	}
}

func TestFlattenAlertCompoundCondition(t *testing.T) {
	r := resourceNewRelicAlertCompoundCondition()
	testAccountID := 123456

	condition := &alerts.CompoundCondition{
		ID:                    "test-id",
		Name:                  "test-compound-condition",
		Enabled:               true,
		PolicyID:              "987654",
		TriggerExpression:     "A AND B",
		RunbookURL:            "https://example.com/runbook",
		ThresholdDuration:     120,
		FacetMatchingBehavior: "FACETS_IGNORED",
		ComponentConditions: []alerts.ComponentCondition{
			{
				ID:    "123",
				Alias: "A",
			},
			{
				ID:    "456",
				Alias: "B",
			},
		},
	}

	d := r.TestResourceData()
	// The flatten function doesn't set ID - that's done by the resource Read function
	d.SetId(condition.ID)

	err := flattenAlertCompoundCondition(testAccountID, condition, d)
	require.NoError(t, err)

	assert.Equal(t, "test-id", d.Id())
	assert.Equal(t, "test-compound-condition", d.Get("name"))
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 987654, d.Get("policy_id"))
	assert.Equal(t, "A AND B", d.Get("trigger_expression"))
	assert.Equal(t, "https://example.com/runbook", d.Get("runbook_url"))
	assert.Equal(t, 120, d.Get("threshold_duration"))
	assert.Equal(t, "FACETS_IGNORED", d.Get("facet_matching_behavior"))
	assert.Equal(t, testAccountID, d.Get("account_id"))
}
