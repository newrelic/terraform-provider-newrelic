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
			Operator:     "above",
			Priority:     "critical",
			Threshold:    1.5,
			TimeFunction: "all",
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
			TimeFunction: "all",
		},
	}

	expected := []map[string]interface{}{
		{
			"duration":      5,
			"operator":      "above",
			"priority":      "critical",
			"threshold":     1.5,
			"time_function": "all",
		},
	}

	flattened := flattenNrqlConditionTerms(expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}
