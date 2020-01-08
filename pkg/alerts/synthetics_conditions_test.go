// +build unit

package alerts

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testPolicyID            = 12345
	testSyntheticsCondition = SyntheticsCondition{
		Name:       "Synthetics Condition",
		RunbookURL: "https://example.com/runbook.md",
		MonitorID:  "12345678-1234-1234-1234-1234567890ab",
		Enabled:    true,
	}
	testSyntheticsConditionJson = `
	{
		"name": "Synthetics Condition",
		"runbook_url": "https://example.com/runbook.md",
		"monitor_id": "12345678-1234-1234-1234-1234567890ab",
		"enabled": true
	}`
)

func TestListSyntheticsConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "synthetics_conditions": [%s] }`, testSyntheticsConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := []*SyntheticsCondition{&testSyntheticsCondition}
	actual, err := alerts.ListSyntheticsConditions(testPolicyID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestCreateSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "synthetics_condition": %s }`, testSyntheticsConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := alerts.CreateSyntheticsCondition(testPolicyID, testSyntheticsCondition)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, &testSyntheticsCondition, actual)
}

func TestUpdateSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "synthetics_condition": %s }`, testSyntheticsConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := alerts.UpdateSyntheticsCondition(testSyntheticsCondition)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, &testSyntheticsCondition, actual)
}

func TestDeleteSyntheticsCondition(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "synthetics_condition": %s }`, testSyntheticsConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	actual, err := alerts.DeleteSyntheticsCondition(testSyntheticsCondition.ID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, &testSyntheticsCondition, actual)
}
