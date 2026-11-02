//go:build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataPartitionRuleRetentionPolicyValidation(t *testing.T) {
	t.Parallel()

	validate := resourceNewRelicDataPartition().Schema["retention_policy"].ValidateFunc

	for _, value := range []string{"SECONDARY", "STANDARD", "STANDARD_LIVE_ARCHIVE"} {
		_, errs := validate(value, "retention_policy")
		require.Empty(t, errs, "expected %q to be a valid retention policy", value)
	}

	for _, value := range []string{"", "standard", "LIVE_ARCHIVE"} {
		_, errs := validate(value, "retention_policy")
		require.NotEmpty(t, errs, "expected %q to be rejected", value)
	}
}
