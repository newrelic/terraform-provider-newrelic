// +build unit

package newrelic

import (
	"strings"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nr "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
)

var (
	testThresholdLow  = 1.0
	testThresholdHigh = 10.9
)

func TestExpandNrqlAlertConditionInput(t *testing.T) {
	nrql := map[string]interface{}{
		"query":             "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'Dummy App'",
		"evaluation_offset": 3,
	}

	var criticalTerms []map[string]interface{}
	criticalTerms = append(criticalTerms, map[string]interface{}{
		"threshold":             1,
		"threshold_occurrences": alerts.ThresholdOccurrences.AtLeastOnce,
		"threshold_duration":    600,
		"operator":              alerts.AlertsNRQLConditionTermsOperatorTypes.ABOVE,
	})

	var warningTerms []map[string]interface{}
	warningTerms = append(warningTerms, map[string]interface{}{
		"threshold":             10.9,
		"threshold_occurrences": alerts.ThresholdOccurrences.AtLeastOnce,
		"threshold_duration":    660,
		"operator":              alerts.AlertsNRQLConditionTermsOperatorTypes.BELOW,
	})

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
		"baseline condition, requires baseline_direction attr": {
			Data: map[string]interface{}{
				"type": "baseline",
			},
			ExpectErr:    true,
			ExpectReason: "attribute `baseline_direction` is required for nrql alert conditions of type `baseline`",
		},
		"baseline condition, has baseline_direction attr": {
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
		"critical term": {
			Data: map[string]interface{}{
				"nrql":           []interface{}{nrql},
				"type":           "static",
				"value_function": "single_value",
				"critical":       criticalTerms,
			},
			ExpectErr:    false,
			ExpectReason: "",
			Expanded: func() *alerts.NrqlConditionInput {
				x := alerts.NrqlConditionInput{
					ValueFunction: &alerts.NrqlConditionValueFunctions.SingleValue,
				}
				x.Terms = []alerts.NrqlConditionTerm{
					{
						Threshold:            &testThresholdLow,
						ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             alerts.AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             alerts.NrqlConditionPriorities.Critical,
					},
				}

				return &x
			}(),
		},
		"critical and warning terms": {
			Data: map[string]interface{}{
				"nrql":           []interface{}{nrql},
				"type":           "static",
				"value_function": "single_value",
				"critical":       criticalTerms,
				"warning":        warningTerms,
			},
			ExpectErr:    false,
			ExpectReason: "",
			Expanded: func() *alerts.NrqlConditionInput {
				x := alerts.NrqlConditionInput{
					ValueFunction: &alerts.NrqlConditionValueFunctions.SingleValue,
				}
				x.Terms = []alerts.NrqlConditionTerm{
					{
						Threshold:            &testThresholdLow,
						ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    600,
						Operator:             alerts.AlertsNRQLConditionTermsOperatorTypes.ABOVE,
						Priority:             alerts.NrqlConditionPriorities.Critical,
					},
					{
						Threshold:            &testThresholdHigh,
						ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
						ThresholdDuration:    660,
						Operator:             alerts.AlertsNRQLConditionTermsOperatorTypes.BELOW,
						Priority:             alerts.NrqlConditionPriorities.Warning,
					},
				}

				return &x
			}(),
		},
		"aggregation window non-zero": {
			Data: map[string]interface{}{
				"nrql":               []interface{}{nrql},
				"aggregation_window": 60,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						AggregationWindow: &[]int{60}[0],
					},
				},
			},
		},
		"aggregation window nil": {
			Data: map[string]interface{}{
				"nrql":               []interface{}{nrql},
				"aggregation_window": nil,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{},
				},
			},
		},
		"gap filling static non-zero": {
			Data: map[string]interface{}{
				"nrql":        []interface{}{nrql},
				"fill_option": "static",
				"fill_value":  1.0,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						FillOption: &alerts.AlertsFillOptionTypes.STATIC,
						FillValue:  &[]float64{1.0}[0],
					},
				},
			},
		},
		"gap filling static zero": {
			Data: map[string]interface{}{
				"nrql":        []interface{}{nrql},
				"fill_option": "static",
				"fill_value":  0.0,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						FillOption: &alerts.AlertsFillOptionTypes.STATIC,
						FillValue:  &[]float64{0.0}[0],
					},
				},
			},
		},
		"gap filling static zero mixed case": {
			Data: map[string]interface{}{
				"nrql":        []interface{}{nrql},
				"fill_option": "sTaTiC",
				"fill_value":  0.0,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						FillOption: &alerts.AlertsFillOptionTypes.STATIC,
						FillValue:  &[]float64{0.0}[0],
					},
				},
			},
		},
		"gap filling last_value": {
			Data: map[string]interface{}{
				"nrql":        []interface{}{nrql},
				"fill_option": "last_value",
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						FillOption: &alerts.AlertsFillOptionTypes.LAST_VALUE,
					},
				},
			},
		},
		"gap filling none": {
			Data: map[string]interface{}{
				"nrql":        []interface{}{nrql},
				"fill_option": "none",
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Signal: &alerts.AlertsNrqlConditionSignal{
						FillOption: &alerts.AlertsFillOptionTypes.NONE,
					},
				},
			},
		},
		"expiration on": {
			Data: map[string]interface{}{
				"nrql":                           []interface{}{nrql},
				"expiration_duration":            120,
				"open_violation_on_expiration":   true,
				"close_violations_on_expiration": true,
			},
			Expanded: &alerts.NrqlConditionInput{
				NrqlConditionBase: alerts.NrqlConditionBase{
					Expiration: &alerts.AlertsNrqlConditionExpiration{
						ExpirationDuration:          &[]int{120}[0],
						CloseViolationsOnExpiration: true,
						OpenViolationOnExpiration:   true,
					},
				},
			},
		},
	}

	r := resourceNewRelicNrqlAlertCondition()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d := r.TestResourceData()

			for k, v := range tc.Data {
				if k == "critical" || k == "warning" {
					var terms []map[string]interface{}

					terms = append(terms, v.([]map[string]interface{})...)

					if err := d.Set(k, terms); err != nil {
						t.Fatalf("err: %s", err)
					}

				} else {
					if err := d.Set(k, v); err != nil {
						t.Fatalf("err: %s", err)
					}
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

				if len(tc.Expanded.Terms) > 0 {
					assert.Equal(t, tc.Expanded.Terms, expanded.Terms)
				}

				if tc.Expanded.Signal != nil {
					require.Equal(t, tc.Expanded.Signal, expanded.Signal)
				}

				if tc.Expanded.Expiration != nil {
					require.Equal(t, tc.Expanded.Expiration, expanded.Expiration)
				}
			}
		})
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
					Threshold:            &testThresholdLow,
					ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    600,
					Operator:             alerts.AlertsNRQLConditionTermsOperatorTypes.ABOVE,
					Priority:             alerts.NrqlConditionPriorities.Critical,
				},
				{
					Threshold:            &testThresholdHigh,
					ThresholdOccurrences: alerts.ThresholdOccurrences.AtLeastOnce,
					ThresholdDuration:    660,
					Operator:             alerts.AlertsNRQLConditionTermsOperatorTypes.BELOW,
					Priority:             alerts.NrqlConditionPriorities.Warning,
				},
			},
			ViolationTimeLimit: alerts.NrqlConditionViolationTimeLimits.OneHour,
			Signal: &alerts.AlertsNrqlConditionSignal{
				FillOption: &alerts.AlertsFillOptionTypes.LAST_VALUE,
			},
			Expiration: &alerts.AlertsNrqlConditionExpiration{
				ExpirationDuration:          &[]int{120}[0],
				CloseViolationsOnExpiration: true,
				OpenViolationOnExpiration:   true,
			},
		},
	}

	// Baseline
	nrqlConditionBaseline := nrqlCondition
	zero := 0.0
	nrqlConditionBaseline.Signal = &alerts.AlertsNrqlConditionSignal{
		FillOption: &alerts.AlertsFillOptionTypes.STATIC,
		FillValue:  &zero,
	}
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
	nrqlConditionOutlier.Signal = &alerts.AlertsNrqlConditionSignal{
		FillOption: &alerts.AlertsFillOptionTypes.NONE,
	}
	nrqlConditionOutlier.Expiration = nil

	conditions := []*alerts.NrqlAlertCondition{
		&nrqlConditionBaseline,
		&nrqlConditionStatic,
		&nrqlConditionOutlier,
	}

	// Use the API object above to construct a "user-configured" critical term.
	// This ensures that in cases where the user has one configured, but there
	// are two present in the API, we avoid crashing.
	// https://github.com/newrelic/terraform-provider-newrelic/issues/882
	term := nrqlCondition.Terms[0]
	testingTerms := make([]map[string]interface{}, 0)
	crit := map[string]interface{}{
		"operator": strings.ToLower(string(term.Operator)),
		// This is a critical term, so the priority is inferred.
		// "priority":              strings.ToLower(string(term.Priority)),
		"threshold":             term.Threshold,
		"threshold_duration":    term.ThresholdDuration,
		"threshold_occurrences": strings.ToLower(string(term.ThresholdOccurrences)),
	}
	testingTerms = append(testingTerms, crit)

	for i, condition := range conditions {
		d := r.TestResourceData()

		// Configure the critical testingTerms for one of the conditions.
		if i == 0 {
			err := d.Set("critical", testingTerms)
			require.NoError(t, err)
		}

		testAccountID, err := nr.GetTestAccountID()
		require.NoError(t, err)

		err = flattenNrqlAlertCondition(testAccountID, condition, d)
		require.NoError(t, err)

		require.Equal(t, 7654321, d.Get("policy_id").(int))
		require.Equal(t, testAccountID, d.Get("account_id").(int))

		criticalTerms := d.Get("critical").([]interface{})
		assert.Equal(t, 1, len(criticalTerms))
		assert.Equal(t, float64(1), criticalTerms[0].(map[string]interface{})["threshold"])
		assert.Equal(t, "at_least_once", criticalTerms[0].(map[string]interface{})["threshold_occurrences"])
		assert.Equal(t, 600, criticalTerms[0].(map[string]interface{})["threshold_duration"])
		assert.Equal(t, "above", criticalTerms[0].(map[string]interface{})["operator"])

		warningTerms := d.Get("warning").([]interface{})
		assert.Equal(t, 1, len(warningTerms))
		assert.Equal(t, float64(10.9), warningTerms[0].(map[string]interface{})["threshold"])
		assert.Equal(t, "at_least_once", warningTerms[0].(map[string]interface{})["threshold_occurrences"])
		assert.Equal(t, 660, warningTerms[0].(map[string]interface{})["threshold_duration"])
		assert.Equal(t, "below", warningTerms[0].(map[string]interface{})["operator"])

		// require.Equal(t, 1, d.Get("critical.threshold").(map[string]interface{}))

		switch condition.Type {
		case alerts.NrqlConditionTypes.Baseline:
			require.Equal(t, string(alerts.NrqlBaselineDirections.LowerOnly), d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("value_function").(string))
			require.Zero(t, d.Get("expected_groups").(int))
			require.Zero(t, d.Get("open_violation_on_group_overlap").(bool))
			require.Equal(t, 120, d.Get("expiration_duration").(int))
			require.True(t, d.Get("open_violation_on_expiration").(bool))
			require.True(t, d.Get("close_violations_on_expiration").(bool))
			require.Equal(t, "static", d.Get("fill_option").(string))
			require.Equal(t, 0.0, d.Get("fill_value").(float64))

		case alerts.NrqlConditionTypes.Static:
			require.Equal(t, string(alerts.NrqlConditionValueFunctions.Sum), d.Get("value_function").(string))
			require.Zero(t, d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("expected_groups").(int))
			require.Zero(t, d.Get("open_violation_on_group_overlap").(bool))
			require.Equal(t, 120, d.Get("expiration_duration").(int))
			require.True(t, d.Get("open_violation_on_expiration").(bool))
			require.True(t, d.Get("close_violations_on_expiration").(bool))
			require.Equal(t, "last_value", d.Get("fill_option").(string))
			require.Zero(t, d.Get("fill_value").(float64))

		case alerts.NrqlConditionTypes.Outlier:
			require.Equal(t, 2, d.Get("expected_groups").(int))
			require.Zero(t, d.Get("ignore_overlap").(bool))
			require.True(t, d.Get("open_violation_on_group_overlap").(bool))
			require.Zero(t, d.Get("baseline_direction").(string))
			require.Zero(t, d.Get("value_function").(string))
			require.Zero(t, d.Get("expiration_duration").(int))
			require.Zero(t, d.Get("open_violation_on_expiration").(bool))
			require.Zero(t, d.Get("close_violations_on_expiration").(bool))
			require.Equal(t, "none", d.Get("fill_option").(string))
			require.Zero(t, d.Get("fill_value").(float64))
		}
	}
}

func TestExpandNrqlConditionTerm(t *testing.T) {

	cases := map[string]struct {
		ExpectErr     bool
		ExpectReason  string
		ConditionType string
		Priority      string
		Term          map[string]interface{}
		Expected      *alerts.NrqlConditionTerm
	}{
		"critical default priority": {
			Priority:      "critical",
			ConditionType: "static",
			Term: map[string]interface{}{
				"threshold":             10.9,
				"threshold_duration":    5,
				"threshold_occurrences": "ALL",
				"operator":              "equals",
			},
			Expected: &alerts.NrqlConditionTerm{
				Operator:             alerts.AlertsNRQLConditionTermsOperator("EQUALS"),
				Priority:             alerts.NrqlConditionPriority("CRITICAL"),
				Threshold:            &testThresholdHigh,
				ThresholdDuration:    5,
				ThresholdOccurrences: "ALL",
			},
		},
		"critical explicit priority": {
			Priority:      "critical",
			ConditionType: "static",
			Term: map[string]interface{}{
				"threshold":             10.9,
				"threshold_duration":    5,
				"threshold_occurrences": "ALL",
				"operator":              "equals",
				"priority":              "critical",
			},
			Expected: &alerts.NrqlConditionTerm{
				Operator:             alerts.AlertsNRQLConditionTermsOperator("EQUALS"),
				Priority:             alerts.NrqlConditionPriority("CRITICAL"),
				Threshold:            &testThresholdHigh,
				ThresholdDuration:    5,
				ThresholdOccurrences: "ALL",
			},
		},
		"warning priority passed at call": {
			Priority:      "warning",
			ConditionType: "static",
			Term: map[string]interface{}{
				"threshold":             10.9,
				"threshold_duration":    9,
				"threshold_occurrences": "ALL",
				"operator":              "equals",
			},
			Expected: &alerts.NrqlConditionTerm{
				Operator:             alerts.AlertsNRQLConditionTermsOperator("EQUALS"),
				Priority:             alerts.NrqlConditionPriority("WARNING"),
				Threshold:            &testThresholdHigh,
				ThresholdDuration:    9,
				ThresholdOccurrences: "ALL",
			},
		},
		"warning priority passed on the term": {
			Priority:      "",
			ConditionType: "static",
			Term: map[string]interface{}{
				"threshold":             10.9,
				"threshold_duration":    9,
				"threshold_occurrences": "ALL",
				"operator":              "equals",
				"priority":              "warning",
			},
			Expected: &alerts.NrqlConditionTerm{
				Operator:             alerts.AlertsNRQLConditionTermsOperator("EQUALS"),
				Priority:             alerts.NrqlConditionPriority("WARNING"),
				Threshold:            &testThresholdHigh,
				ThresholdDuration:    9,
				ThresholdOccurrences: "ALL",
			},
		},
		"critical priority passed on the term, and warning priority passed to the method": {
			Priority:      "warning",
			ConditionType: "static",
			Term: map[string]interface{}{
				"threshold":             10.9,
				"threshold_duration":    9,
				"threshold_occurrences": "ALL",
				"operator":              "equals",
				"priority":              "critical",
			},
			Expected: &alerts.NrqlConditionTerm{
				Operator:             alerts.AlertsNRQLConditionTermsOperator("EQUALS"),
				Priority:             alerts.NrqlConditionPriority("WARNING"),
				Threshold:            &testThresholdHigh,
				ThresholdDuration:    9,
				ThresholdOccurrences: "ALL",
			},
		},
	}

	for _, tc := range cases {
		expandedTerm, err := expandNrqlConditionTerm(tc.Term, tc.ConditionType, tc.Priority)
		if tc.ExpectErr {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			require.Nil(t, err)
		}

		if tc.Expected != nil {
			require.Equal(t, tc.Expected.Operator, expandedTerm.Operator)
			require.Equal(t, tc.Expected.Priority, expandedTerm.Priority)
			require.Equal(t, *tc.Expected.Threshold, *expandedTerm.Threshold)
			require.Equal(t, tc.Expected.ThresholdDuration, expandedTerm.ThresholdDuration)
			require.Equal(t, tc.Expected.ThresholdOccurrences, expandedTerm.ThresholdOccurrences)
		}
	}

}
