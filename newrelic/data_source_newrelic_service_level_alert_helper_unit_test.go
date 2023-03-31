//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateAlertThreshold(t *testing.T) {
	var sloPeriod = 28
	var sloTarget = 99.9
	var toleratedBudgetConsumption = 2.0
	var evaluationPeriod = 60

	threshold := calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

	require.NotNil(t, threshold)
	require.Equal(t, 1.3439999999999237, threshold)
}
