//go:build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The Azure Auto Discovery integration performs an inventory sweep rather than a metric
// poll, and accepts only three intervals: 8, 12 and 16 hours. These tests pin that set at
// the schema level so an unsupported value fails during plan instead of apply.
func TestCloudAzureAutoDiscoveryPollingIntervalValidation(t *testing.T) {
	validateFunc := cloudAzureAutoDiscoveryElem().Schema["metrics_polling_interval"].ValidateFunc
	require.NotNil(t, validateFunc, "auto_discovery must validate metrics_polling_interval")

	testCases := []struct {
		name     string
		interval int
		valid    bool
	}{
		{name: "8 hours, the default", interval: 28800, valid: true},
		{name: "12 hours", interval: 43200, valid: true},
		{name: "16 hours", interval: 57600, valid: true},
		{name: "24 hours, allowed for AWS but not Azure", interval: 86400, valid: false},
		{name: "between two allowed values", interval: 36000, valid: false},
		{name: "one second below the smallest allowed value", interval: 28799, valid: false},
		{name: "1 hour, the metric integration default", interval: 3600, valid: false},
		{name: "5 minutes", interval: 300, valid: false},
		{name: "zero", interval: 0, valid: false},
		{name: "negative", interval: -500, valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			warns, errs := validateFunc(testCase.interval, "metrics_polling_interval")

			require.Empty(t, warns)
			if testCase.valid {
				require.Empty(t, errs)
			} else {
				require.NotEmpty(t, errs)
			}
		})
	}
}

// Every interval the schema advertises must actually pass validation, so the allowed set
// and the ValidateFunc cannot drift apart.
func TestCloudAzureAutoDiscoveryPollingIntervalsAreAllAccepted(t *testing.T) {
	validateFunc := cloudAzureAutoDiscoveryElem().Schema["metrics_polling_interval"].ValidateFunc

	require.ElementsMatch(t, []int{28800, 43200, 57600}, cloudAzureAutoDiscoveryPollingIntervals)

	for _, interval := range cloudAzureAutoDiscoveryPollingIntervals {
		warns, errs := validateFunc(interval, "metrics_polling_interval")

		require.Empty(t, warns)
		require.Empty(t, errs, "advertised interval %d must be accepted", interval)
	}
}

// The shared base is used by all of the metric integration blocks, each of which allows a
// different set of intervals, so it must stay unvalidated. This guards against the Auto
// Discovery restriction being moved onto it, which would break blocks such as monitor.
func TestCloudAzureIntegrationSchemaBaseIsNotValidated(t *testing.T) {
	require.Nil(
		t,
		cloudAzureIntegrationSchemaBase()["metrics_polling_interval"].ValidateFunc,
		"the shared base is used by integrations with different allowed intervals and must not enforce the Auto Discovery set",
	)
}
