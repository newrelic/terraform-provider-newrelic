// +build integration

package alerts

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAlertConditions(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: Read
	readResult, err := testReadAlertConditions(t, client)

	require.NoError(t, err)

	log.Println(readResult)
}

func testReadAlertConditions(t *testing.T, client Alerts) ([]*AlertCondition, error) {
	result, err := client.ListAlertConditions(593566)

	return result, err
}
