//go:build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Every interval the schema advertises must actually pass validation, so the allowed set
// and the ValidateFunc cannot drift apart.
func TestCloudAwsInfra360PollingIntervalsAreAllAccepted(t *testing.T) {
	validateFunc := cloudAwsInfra360Elem().Schema["metrics_polling_interval"].ValidateFunc

	require.ElementsMatch(t, []int{28800, 43200, 86400}, cloudAwsInfra360PollingIntervals)

	for _, interval := range cloudAwsInfra360PollingIntervals {
		warns, errs := validateFunc(interval, "metrics_polling_interval")

		require.Empty(t, warns)
		require.Empty(t, errs, "advertised interval %d must be accepted", interval)
	}
}

func TestCloudAwsInfra360PollingIntervalRejectsOthers(t *testing.T) {
	validateFunc := cloudAwsInfra360Elem().Schema["metrics_polling_interval"].ValidateFunc

	for _, interval := range []int{300, 6000, 57600, 100000} {
		_, errs := validateFunc(interval, "metrics_polling_interval")

		require.NotEmpty(t, errs, "interval %d must be rejected", interval)
	}
}

// The shared base is used by many integration blocks that allow different intervals, so it
// must stay unvalidated.
func TestCloudAwsIntegrationSchemaBaseIsNotValidated(t *testing.T) {
	require.Nil(t, cloudAwsIntegrationSchemaBaseExtended()["metrics_polling_interval"].ValidateFunc)
}
