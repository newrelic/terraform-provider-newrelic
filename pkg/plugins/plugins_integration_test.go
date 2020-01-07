// +build integration

package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationPlugins(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: List
	listResult, err := client.ListPlugins(nil)

	require.NoError(t, err)

	if len(listResult) == 0 {
		t.Skip("Skipping `GetPlugin()` integration test due to zero plugins found")
		return
	}

	// Test: Get
	qp := GetPluginParams{
		Detailed: true,
	}
	getResult, err := client.GetPlugin(listResult[0].ID, &qp)

	require.NoError(t, err)
	require.NotNil(t, getResult)
	require.NotNil(t, getResult.Details)
}
