//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAiNotificationsResponseErrors_Empty(t *testing.T) {
	result := buildAiNotificationsResponseErrors([]ai.AiNotificationsResponseError{})
	assert.Empty(t, result)
}

func TestBuildAiNotificationsResponseErrors_SingleError(t *testing.T) {
	errors := []ai.AiNotificationsResponseError{
		{
			Description: "test error description",
			Type:        ai.AiNotificationsErrorType("INVALID_PARAMETER"),
		},
	}

	result := buildAiNotificationsResponseErrors(errors)
	require.Len(t, result, 1)
	assert.Equal(t, diag.Error, result[0].Severity)
	assert.Equal(t, "INVALID_PARAMETER: test error description", result[0].Summary)
}

func TestBuildAiNotificationsResponseErrors_MultipleErrors(t *testing.T) {
	errors := []ai.AiNotificationsResponseError{
		{
			Description: "first error",
			Type:        ai.AiNotificationsErrorType("INVALID_PARAMETER"),
		},
		{
			Description: "second error",
			Type:        ai.AiNotificationsErrorType("NOT_FOUND"),
		},
	}

	result := buildAiNotificationsResponseErrors(errors)
	require.Len(t, result, 2)
	assert.Equal(t, "INVALID_PARAMETER: first error", result[0].Summary)
	assert.Equal(t, "NOT_FOUND: second error", result[1].Summary)
}

func TestBuildAiNotificationsResponseErrors_Nil(t *testing.T) {
	result := buildAiNotificationsResponseErrors(nil)
	assert.Empty(t, result)
}

func TestBuildAiNotificationsErrors_WithFieldErrors(t *testing.T) {
	errors := []ai.AiNotificationsError{
		{
			Details: "validation failed",
			Fields: []ai.AiNotificationsFieldError{
				{Field: "name", Message: "is required"},
			},
		},
	}

	result := buildAiNotificationsErrors(errors)
	require.Len(t, result, 1)
	assert.Equal(t, diag.Error, result[0].Severity)
	assert.Equal(t, "validation failed", result[0].Detail)
}

func TestBuildAiNotificationsErrors_WithDependencyErrors(t *testing.T) {
	errors := []ai.AiNotificationsError{
		{
			Name:         "constraint_name",
			Dependencies: []string{"dep1", "dep2"},
		},
	}

	result := buildAiNotificationsErrors(errors)
	require.Len(t, result, 1)
	assert.Equal(t, diag.Error, result[0].Severity)
	assert.Contains(t, result[0].Summary, "constraint_name")
	assert.Contains(t, result[0].Summary, "dep1")
}

func TestBuildAiNotificationsErrors_WithResponseError(t *testing.T) {
	errors := []ai.AiNotificationsError{
		{
			Description: "something went wrong",
			Details:     "details here",
			Type:        ai.AiNotificationsErrorType("SERVER_ERROR"),
		},
	}

	result := buildAiNotificationsErrors(errors)
	require.Len(t, result, 1)
	assert.Equal(t, diag.Error, result[0].Severity)
	assert.Equal(t, "SERVER_ERROR: something went wrong", result[0].Summary)
	assert.Equal(t, "details here", result[0].Detail)
}

func TestBuildAiNotificationsErrors_Empty(t *testing.T) {
	result := buildAiNotificationsErrors([]ai.AiNotificationsError{})
	assert.Empty(t, result)
}

func TestBuildEntityScopeInput_WithOrganizationScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	d := r.TestResourceData()

	err := d.Set("scope", []interface{}{
		map[string]interface{}{
			"type": "ORGANIZATION",
			"id":   "org-uuid-123",
		},
	})
	require.NoError(t, err)

	scope := buildEntityScopeInput(d, 12345)

	assert.Equal(t, notifications.AiNotificationsEntityScopeTypeInput("ORGANIZATION"), scope.Type)
	assert.Equal(t, "org-uuid-123", scope.ID)
}

func TestBuildEntityScopeInput_WithAccountScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	d := r.TestResourceData()

	err := d.Set("scope", []interface{}{
		map[string]interface{}{
			"type": "ACCOUNT",
			"id":   "98765",
		},
	})
	require.NoError(t, err)

	scope := buildEntityScopeInput(d, 12345)

	assert.Equal(t, notifications.AiNotificationsEntityScopeTypeInput("ACCOUNT"), scope.Type)
	assert.Equal(t, "98765", scope.ID)
}

func TestBuildEntityScopeInput_DefaultsToAccountScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	d := r.TestResourceData()

	scope := buildEntityScopeInput(d, 12345)

	assert.Equal(t, notifications.AiNotificationsEntityScopeTypeInputTypes.ACCOUNT, scope.Type)
	assert.Equal(t, "12345", scope.ID)
}

func TestListValidNotificationsScopeTypes(t *testing.T) {
	scopeTypes := listValidNotificationsScopeTypes()

	require.Len(t, scopeTypes, 2)
	assert.Contains(t, scopeTypes, "ACCOUNT")
	assert.Contains(t, scopeTypes, "ORGANIZATION")
}

func TestCreateMonitoringProperty(t *testing.T) {
	prop := createMonitoringProperty()

	assert.Equal(t, "source", prop.Key)
	assert.Equal(t, "terraform", prop.Value)
	assert.Equal(t, "terraform-source-internal", prop.Label)
}
