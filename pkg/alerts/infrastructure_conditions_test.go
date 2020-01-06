// +build unit

package alerts

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
	"github.com/stretchr/testify/require"
)

var (
	testInfrastructureConditionPolicyId  = 111111
	testInfrastructureConditionTimestamp = serialization.EpochTime(time.Unix(1490996713872, 0))
	testInfrastructureCriticalThreshold  = InfrastructureConditionThreshold{
		Duration: 6,
		Function: "all",
		Value:    12.3,
	}
	testInfrastructureWarningThreshold = InfrastructureConditionThreshold{
		Duration: 6,
		Function: "all",
		Value:    10,
	}

	testInfrastructureCondition = InfrastructureCondition{
		Comparison:   "equal",
		CreatedAt:    &testInfrastructureConditionTimestamp,
		Critical:     &testInfrastructureCriticalThreshold,
		Enabled:      true,
		ID:           13890,
		Name:         "Java is running",
		PolicyID:     testInfrastructureConditionPolicyId,
		ProcessWhere: "(commandName = 'java')",
		Type:         "infra_process_running",
		UpdatedAt:    &testInfrastructureConditionTimestamp,
		Warning:      &testInfrastructureWarningThreshold,
		Where:        "(hostname LIKE '%cassandra%')",
	}
	testInfrastructureConditionJson = `
		{
			"type":"infra_process_running",
			"name":"Java is running",
			"enabled":true,
			"where_clause":"(hostname LIKE '%cassandra%')",
			"id":13890,
			"created_at_epoch_millis":1490996713872,
			"updated_at_epoch_millis":1490996713872,
			"policy_id":111111,
			"comparison":"equal",
			"critical_threshold":{
				"value":12.3,
				"duration_minutes":6,
				"time_function": "all"
			},
			"warning_threshold": {
				"value": 10,
				"duration_minutes": 6,
				"time_function": "all"
			},
			"process_where_clause":"(commandName = 'java')"
		}`
)

func TestListInfrastructureConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":[%s] }`, testInfrastructureConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := []InfrastructureCondition{testInfrastructureCondition}

	actual, err := alerts.ListInfrastructureConditions(testInfrastructureConditionPolicyId)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestGetInfrastructureConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testInfrastructureConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testInfrastructureCondition

	actual, err := alerts.GetInfrastructureCondition(testInfrastructureCondition.ID)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestCreateInfrastructureConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testInfrastructureConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testInfrastructureCondition

	actual, err := alerts.CreateInfrastructureCondition(testInfrastructureCondition)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestUpdateInfrastructureConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testInfrastructureConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	expected := &testInfrastructureCondition

	actual, err := alerts.UpdateInfrastructureCondition(testInfrastructureCondition)

	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func TestDeleteInfrastructureConditions(t *testing.T) {
	t.Parallel()
	respJSON := fmt.Sprintf(`{ "data":%s }`, testInfrastructureConditionJson)
	alerts := newMockResponse(t, respJSON, http.StatusOK)

	err := alerts.DeleteInfrastructureCondition(testInfrastructureCondition.ID)

	require.NoError(t, err)
}
