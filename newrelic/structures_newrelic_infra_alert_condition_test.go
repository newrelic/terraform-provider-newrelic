// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

func TestExpandInfraAlertThreshold(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"duration":      5,
			"value":         1.5,
			"time_function": "all",
		},
	}

	expected := alerts.InfrastructureConditionThreshold{
		Duration: 5,
		Function: "all",
		Value:    1.5,
	}

	expanded := expandInfraAlertThreshold(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, &expected, expanded)
}

func TestFlattenInfraAlertThreshold(t *testing.T) {
	expanded := alerts.InfrastructureConditionThreshold{
		Duration: 5,
		Function: "all",
		Value:    1.5,
	}

	expected := []interface{}{
		map[string]interface{}{
			"duration":      5,
			"value":         1.5,
			"time_function": "all",
		},
	}

	flattened := flattenAlertThreshold(&expanded)

	require.NotNil(t, flattened)
	require.Equal(t, expected, flattened)
}
