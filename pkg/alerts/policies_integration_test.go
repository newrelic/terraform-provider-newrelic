// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPolicy(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	testIntegrationPolicyNameRandStr := nr.RandSeq(5)
	policy := Policy{
		IncidentPreference: "PER_POLICY",
		Name:               fmt.Sprintf("test-alert-policy-%s", testIntegrationPolicyNameRandStr),
	}

	// Test: Create
	createResult, err := client.CreatePolicy(policy)

	require.NoError(t, err)
	require.NotNil(t, createResult)

	// Test: Read
	readResult, err := client.GetPolicy(createResult.ID)

	require.NoError(t, err)
	require.NotNil(t, readResult)

	// Test: Update
	createResult.Name = fmt.Sprintf("test-alert-policy-updated-%s", testIntegrationPolicyNameRandStr)
	updateResult, err := client.UpdatePolicy(*createResult)

	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := client.DeletePolicy(updateResult.ID)

	require.NoError(t, err)
	require.NotNil(t, *deleteResult)
}
