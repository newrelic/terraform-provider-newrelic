// +build integration

package apm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationKeyTransactions(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: List
	listResult, err := client.ListKeyTransactions(nil)

	require.NoError(t, err)

	if len(listResult) == 0 {
		t.Skip("Skipping `GetKeyTransaction` integration test due to zero key transactions found")
		return
	}

	// Test: Get
	getResult, err := client.GetKeyTransaction(listResult[0].ID)

	require.NoError(t, err)
	require.NotNil(t, getResult)
}
