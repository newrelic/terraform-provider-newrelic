package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestExpandAlertConditionEntities(t *testing.T) {
	flattened := []interface{}{123, 456}
	expected := []string{"123", "456"}

	expanded := expandAlertConditionEntities(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenAlertConditionEntities(t *testing.T) {
	expanded := []string{"123", "456"}
	expected := []int{123, 456}

	flattened, err := flattenAlertConditionEntities(&expanded)

	require.NoError(t, err)
	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}

func TestExpandAlertConditionTerms(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"duration":      123,
			"operator":      "operator",
			"priority":      "priority",
			"threshold":     123.456,
			"time_function": alerts.TimeFunctionTypes.All,
		},
	}

	expected := []alerts.ConditionTerm{
		{
			Duration:     123,
			Operator:     "operator",
			Priority:     "priority",
			Threshold:    123.456,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
	}

	expanded := expandAlertConditionTerms(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenAlertConditionTerms(t *testing.T) {
	expanded := []alerts.ConditionTerm{
		{
			Duration:     123,
			Operator:     "operator",
			Priority:     "priority",
			Threshold:    123.456,
			TimeFunction: alerts.TimeFunctionTypes.All,
		},
	}

	expected := []map[string]interface{}{
		{
			"duration":      123,
			"operator":      "operator",
			"priority":      "priority",
			"threshold":     123.456,
			"time_function": alerts.TimeFunctionTypes.All,
		},
	}

	flattened := flattenAlertConditionTerms(&expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}
