package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestExpandPluginsAlertConditionEntities(t *testing.T) {
	flattened := []interface{}{123, 456}
	expected := []string{"123", "456"}

	expanded := expandPluginsConditionEntities(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenPluginsAlertConditionEntities(t *testing.T) {
	expanded := []string{"123", "456"}
	expected := []int{123, 456}

	flattened, err := flattenPluginsConditionEntities(expanded)

	require.NoError(t, err)
	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}

func TestExpandPluginsConditionTerms(t *testing.T) {
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

	expanded := expandPluginsConditionTerms(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}

func TestFlattenPluginsConditionTerms(t *testing.T) {
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

	flattened := flattenPluginsConditionTerms(expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}
