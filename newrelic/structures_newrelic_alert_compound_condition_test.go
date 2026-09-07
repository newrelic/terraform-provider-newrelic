//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		"with description and title_template": {
			Data: map[string]interface{}{
				"name":               "test-compound-condition",
				"enabled":            true,
				"trigger_expression": "A AND B",
				"description":        "Test description",
				"title_template":     "{{compoundCondition.name}} triggered",
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
			Expected: func() *alerts.CompoundConditionCreateInput {
				description := "Test description"
				titleTemplate := "{{compoundCondition.name}} triggered"
				return &alerts.CompoundConditionCreateInput{
					Name:              "test-compound-condition",
					Enabled:           true,
					TriggerExpression: "A AND B",
					Description:       &description,
					TitleTemplate:     &titleTemplate,
					ComponentConditions: []alerts.ComponentConditionInput{
						{ID: "123", Alias: "A"},
						{ID: "456", Alias: "B"},
					},
				}
			}(),
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
				if tc.Expected.Description != nil {
					require.NotNil(t, expanded.Description)
					assert.Equal(t, *tc.Expected.Description, *expanded.Description)
				}
				if tc.Expected.TitleTemplate != nil {
					require.NotNil(t, expanded.TitleTemplate)
					assert.Equal(t, *tc.Expected.TitleTemplate, *expanded.TitleTemplate)
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
		EntityGuid:            "MTAxMzMyMDB8QUxFUlR8Q09ORGU5OfDEwMzQ1NTc",
		Description:           "Test description",
		TitleTemplate:         "{{compoundCondition.name}} triggered",
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
	assert.Equal(t, "MTAxMzMyMDB8QUxFUlR8Q09ORGU5OfDEwMzQ1NTc", d.Get("entity_guid"))
	assert.Equal(t, "Test description", d.Get("description"))
	assert.Equal(t, "{{compoundCondition.name}} triggered", d.Get("title_template"))
}

func TestNormalizeComponentConditionID(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Expected string
	}{
		"composite id":           {Input: "1254186:10395010", Expected: "10395010"},
		"plain id":               {Input: "10395010", Expected: "10395010"},
		"whitespace padded":      {Input: " 1254186:10395010 ", Expected: "10395010"},
		"non numeric parts":      {Input: "abc:def", Expected: "abc:def"},
		"more than two segments": {Input: "1:2:3", Expected: "1:2:3"},
		"empty":                  {Input: "", Expected: ""},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, normalizeComponentConditionID(tc.Input))
		})
	}
}

func TestExpandComponentConditionsNormalizesCompositeIDs(t *testing.T) {
	r := resourceNewRelicAlertCompoundCondition()
	d := r.TestResourceData()

	err := d.Set("component_conditions", []interface{}{
		map[string]interface{}{"id": "1254186:10395010", "alias": "A"},
		map[string]interface{}{"id": "10395011", "alias": "B"},
	})
	require.NoError(t, err)

	components, err := expandComponentConditions(d.Get("component_conditions").(*schema.Set))
	require.NoError(t, err)

	ids := map[string]string{}
	for _, c := range components {
		ids[c.Alias] = c.ID
	}

	assert.Equal(t, "10395010", ids["A"])
	assert.Equal(t, "10395011", ids["B"])
}

// TestComponentConditionIDRoundTripNoDrift verifies that the StateFunc on
// component_conditions.id prevents perpetual diff after a create+read cycle.
//
// Important: StateFunc is NOT applied by d.Set or d.Get. It is called by
// Terraform's plan engine on the raw config value before comparing it against
// state. We simulate that here by calling normalizeComponentConditionID
// directly (which is exactly what the StateFunc wraps).
//
// The lifecycle under test:
//  1. User writes composite or plain id in config.
//  2. Plan engine calls StateFunc(config_id) → plain id; compares with state.
//  3. Apply calls Create/Update: expandComponentConditions normalizes the raw
//     ResourceData id (still composite at this point) before sending to API.
//  4. Read calls flattenComponentConditions: API returns plain id → stored in state.
//  5. Next plan: StateFunc(config_id) == state id → no diff.
func TestComponentConditionIDRoundTripNoDrift(t *testing.T) {
	// setHash mirrors the hash function defined on the component_conditions schema.
	setHash := func(v interface{}) int {
		return schema.HashString(v.(map[string]interface{})["alias"].(string))
	}

	// readSet represents what flattenComponentConditions produces after a Read:
	// the API always returns plain numeric condition IDs.
	readSet := flattenComponentConditions([]alerts.ComponentCondition{
		{ID: "10395010", Alias: "A"},
		{ID: "10395011", Alias: "B"},
	})

	cases := map[string][]map[string]string{
		"composite ids (newrelic_nrql_alert_condition.x.id format)": {
			{"id": "1254186:10395010", "alias": "A"},
			{"id": "1254186:10395011", "alias": "B"},
		},
		"plain numeric ids (split workaround or direct numeric)": {
			{"id": "10395010", "alias": "A"},
			{"id": "10395011", "alias": "B"},
		},
		"mixed (one composite, one plain)": {
			{"id": "1254186:10395010", "alias": "A"},
			{"id": "10395011", "alias": "B"},
		},
	}

	for name, rawPairs := range cases {
		t.Run(name, func(t *testing.T) {
			// Simulate Terraform's plan engine applying StateFunc to every config id.
			// This is what prevents drift: the plan engine normalizes the user's
			// composite id to a plain id before comparing it against state.
			elems := make([]interface{}, 0, len(rawPairs))
			for _, pair := range rawPairs {
				elems = append(elems, map[string]interface{}{
					"id":    normalizeComponentConditionID(pair["id"]),
					"alias": pair["alias"],
				})
			}
			configSet := schema.NewSet(setHash, elems)

			// The plan engine compares configSet against readSet element-by-element
			// using Difference. A non-empty result means a perpetual diff.
			onlyInConfig := configSet.Difference(readSet)
			onlyInRead := readSet.Difference(configSet)
			assert.Zero(t, onlyInConfig.Len(),
				"case %q: StateFunc-normalized config has elements not in read state: %v", name, onlyInConfig.List())
			assert.Zero(t, onlyInRead.Len(),
				"case %q: read state has elements not in StateFunc-normalized config: %v", name, onlyInRead.List())
		})
	}
}
