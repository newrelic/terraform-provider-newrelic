//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
	"github.com/stretchr/testify/require"
)

func TestIsDimensionalMetricsProviderValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		provider string
		isDM     bool
		expected bool
	}{
		{name: "not DM, any provider is valid", provider: "aws", isDM: false, expected: true},
		{name: "DM with gcp provider is valid", provider: "gcp", isDM: true, expected: true},
		{name: "DM with gcp provider is valid case-insensitive", provider: "GCP", isDM: true, expected: true},
		{name: "DM with non-gcp provider is invalid", provider: "aws", isDM: true, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, isDimensionalMetricsProviderValid(tc.provider, tc.isDM))
		})
	}
}

func TestFindCloudLinkedAccount(t *testing.T) {
	t.Parallel()

	accounts := []cloud.CloudLinkedAccount{
		{ID: 1, NrAccountId: 100, Name: "legacy-account", HasDimensionalMetrics: false},
		{ID: 2, NrAccountId: 100, Name: "dm-account", HasDimensionalMetrics: true},
		{ID: 3, NrAccountId: 200, Name: "dm-account", HasDimensionalMetrics: true},
	}

	t.Run("matches legacy account when isDM is false", func(t *testing.T) {
		result := findCloudLinkedAccount(accounts, 100, "legacy-account", false)
		require.NotNil(t, result)
		require.Equal(t, 1, result.ID)
	})

	t.Run("matches dimensional metrics account when isDM is true", func(t *testing.T) {
		result := findCloudLinkedAccount(accounts, 100, "dm-account", true)
		require.NotNil(t, result)
		require.Equal(t, 2, result.ID)
	})

	t.Run("does not match a same-named account of the wrong kind", func(t *testing.T) {
		result := findCloudLinkedAccount(accounts, 100, "dm-account", false)
		require.Nil(t, result)
	})

	t.Run("matches name case-insensitively", func(t *testing.T) {
		result := findCloudLinkedAccount(accounts, 200, "DM-Account", true)
		require.NotNil(t, result)
		require.Equal(t, 3, result.ID)
	})

	t.Run("returns nil when account id does not match", func(t *testing.T) {
		result := findCloudLinkedAccount(accounts, 999, "dm-account", true)
		require.Nil(t, result)
	})

	t.Run("returns nil for empty account list", func(t *testing.T) {
		result := findCloudLinkedAccount(nil, 100, "dm-account", true)
		require.Nil(t, result)
	})
}
