package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

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

func TestFlattenNrqlAlertCondition(t *testing.T) {
	r := resourceNewRelicNrqlAlertCondition()

	nrqlCondition := alerts.NrqlAlertCondition{
		ID:       "1234567",
		PolicyID: "7654321",
		NrqlConditionBase: alerts.NrqlConditionBase{
			Description: "description test",
			Enabled:     true,
			Name:        "name-test",
			Nrql: alerts.NrqlConditionQuery{
				Query:            "SELECT average(duration) from Transaction where appName='Dummy App'",
				EvaluationOffset: 3,
			},
			RunbookURL: "test.com",
			Terms: []alerts.NrqlConditionTerm{
				{
					Threshold:            1,
					ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             alerts.NrqlConditionOperators.Above,
					Priority:             alerts.NrqlConditionPriorities.Critical,
				},
			},
			ViolationTimeLimit: alerts.NrqlConditionViolationTimeLimits.OneHour,
		},
	}

	// Baseline
	nrqlConditionBaseline := nrqlCondition
	nrqlConditionBaseline.Type = alerts.NrqlConditionTypes.Baseline
	nrqlConditionBaseline.BaselineDirection = &alerts.NrqlBaselineDirections.LowerOnly

	// Static
	nrqlConditionStatic := nrqlCondition
	nrqlConditionStatic.Type = alerts.NrqlConditionTypes.Static
	nrqlConditionStatic.ValueFunction = &alerts.NrqlConditionValueFunctions.Sum

	// Outlier
	expectedGroups := 2
	openViolationOnOverlap := true
	nrqlConditionOutlier := nrqlCondition
	nrqlConditionOutlier.Type = "OUTLIER"
	nrqlConditionOutlier.ExpectedGroups = &expectedGroups
	nrqlConditionOutlier.OpenViolationOnGroupOverlap = &openViolationOnOverlap

	conditions := []*alerts.NrqlAlertCondition{
		&nrqlConditionBaseline,
		&nrqlConditionStatic,
		&nrqlConditionOutlier,
	}

	for _, condition := range conditions {
		d := r.TestResourceData()
		err := flattenNrqlAlertCondition(nr.TestAccountID, condition, d)
		require.NoError(t, err)

		require.Equal(t, 7654321, d.Get("policy_id").(int))
		require.Equal(t, nr.TestAccountID, d.Get("account_id").(int))

		switch condition.Type {
		case alerts.NrqlConditionTypes.Baseline:
			require.Equal(t, string(alerts.NrqlBaselineDirections.LowerOnly), d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("value_function").(string))
			require.Zero(t, d.Get("expected_groups").(int))
			require.Zero(t, d.Get("open_violation_on_group_overlap").(bool))

		case alerts.NrqlConditionTypes.Static:
			require.Equal(t, string(alerts.NrqlConditionValueFunctions.Sum), d.Get("value_function").(string))
			require.Zero(t, d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("expected_groups").(int))
			require.Zero(t, d.Get("open_violation_on_group_overlap").(bool))

		case alerts.NrqlConditionTypes.Outlier:
			require.Equal(t, 2, d.Get("expected_groups").(int))
			require.True(t, d.Get("open_violation_on_group_overlap").(bool))
			require.Zero(t, d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("value_function").(string))
		}
	}
}
