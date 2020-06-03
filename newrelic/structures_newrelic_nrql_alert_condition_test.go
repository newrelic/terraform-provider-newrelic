package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestExpandNrqlConditionTerms(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"duration":      5,
			"operator":      "above",
			"priority":      "critical",
			"threshold":     1.5,
			"time_function": "all",
		},
	}

	expected := []alerts.ConditionTerm{
		{
			Duration:     5,
			Operator:     alerts.OperatorTypes.Above,
			Priority:     alerts.PriorityTypes.Critical,
			Threshold:    1.5,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
	}

	expanded := expandNrqlConditionTerms(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenNrql(t *testing.T) {
	expanded := alerts.NrqlQuery{
		Query:      "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip",
		SinceValue: "3",
	}

	expected := []interface{}{map[string]interface{}{
		"query":       "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip",
		"since_value": "3",
	}}

	flattened := flattenNrqlQuery(expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}

func TestFlattenNrqlConditionTerms(t *testing.T) {
	expanded := []alerts.ConditionTerm{
		{
			Duration:     5,
			Operator:     "above",
			Priority:     "critical",
			Threshold:    1.5,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
	}

	expected := []map[string]interface{}{
		{
			"duration":      5,
			"operator":      alerts.OperatorTypes.Above,
			"priority":      alerts.PriorityTypes.Critical,
			"threshold":     1.5,
			"time_function": alerts.TimeFunctionTypes.All,
		},
	}

	flattened := flattenNrqlConditionTerms(expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}

func TestExpandNrqlAlertConditionInput(t *testing.T) {

	nrql := map[string]interface{}{
		"query":             "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'Dummy App'",
		"evaluation_offset": 3,
	}

	expectedNrql := &alerts.NrqlConditionInput{}
	expectedNrql.Nrql.Query = nrql["query"].(string)
	expectedNrql.Nrql.EvaluationOffset = nrql["evaluation_offset"].(int)

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *alerts.NrqlConditionInput
	}{
		"invalid nrql": {
			Data:         map[string]interface{}{},
			ExpectErr:    true,
			ExpectReason: "one of `since_value` or `evaluation_offset` must be configured for block `nrql`",
		},
		"valid nrql": {
			Data: map[string]interface{}{
				"nrql": []interface{}{nrql},
			},
			ExpectErr:    false,
			ExpectReason: "",
			Expanded:     expectedNrql,
		},
		"basline condition, requires baseline_direction attr": {
			Data: map[string]interface{}{
				"type": "baseline",
			},
			ExpectErr:    true,
			ExpectReason: "attribute `baseline_direction` is required for nrql alert conditions of type `baseline`",
		},
		"basline condition, has baseline_direction attr": {
			Data: map[string]interface{}{
				"nrql":               []interface{}{nrql},
				"type":               "baseline",
				"baseline_direction": "lower_only",
			},
			ExpectErr:    false,
			ExpectReason: "",
			Expanded: &alerts.NrqlConditionInput{
				BaselineDirection: &alerts.NrqlBaselineDirections.LowerOnly,
			},
		},
		"static condition, requires value_function attr": {
			Data: map[string]interface{}{
				"nrql": []interface{}{nrql},
				"type": "static",
			},
			ExpectErr:    true,
			ExpectReason: "attribute `value_function` is required for nrql alert conditions of type `static`",
		},
		"static condition, has value_function attr": {
			Data: map[string]interface{}{
				"nrql":           []interface{}{nrql},
				"type":           "static",
				"value_function": "single_value",
			},
			ExpectErr:    false,
			ExpectReason: "",
			Expanded: &alerts.NrqlConditionInput{
				ValueFunction: &alerts.NrqlConditionValueFunctions.SingleValue,
			},
		},
	}

	r := resourceNewRelicNrqlAlertCondition()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded, err := expandNrqlAlertConditionInput(d)

		if tc.ExpectErr {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			require.Nil(t, err)
		}

		if tc.Expanded != nil {

			// Static conditions specific
			if tc.Expanded.ValueFunction != nil {
				require.Equal(t, *tc.Expanded.ValueFunction, *expanded.ValueFunction)
			}

			// Baseline conditions specific
			if tc.Expanded.BaselineDirection != nil {
				require.Equal(t, *tc.Expanded.BaselineDirection, *expanded.BaselineDirection)
			}

			if tc.Expanded.Nrql.Query != "" {
				require.Equal(t, tc.Expanded.Nrql.Query, expanded.Nrql.Query)
			}

			if tc.Expanded.Nrql.EvaluationOffset > 0 {
				require.Equal(t, tc.Expanded.Nrql.EvaluationOffset, expanded.Nrql.EvaluationOffset)
			}
		}
	}

}
