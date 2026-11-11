package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/maintenancewindow"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
	"github.com/stretchr/testify/assert"
)

func testMaintenanceWindowResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, resourceNewRelicMaintenanceWindow().Schema, raw)
}

func TestExpandMaintenanceWindowCreateInput(t *testing.T) {
	t.Parallel()

	d := testMaintenanceWindowResourceData(t, map[string]interface{}{
		"name":                 "tf-test-mw",
		"description":          "a test maintenance window",
		"scope_id":             "1067061",
		"scope_type":           "ACCOUNT",
		"start_time":           "2026-12-01T22:00:00",
		"duration":             "PT1H",
		"rrule":                "FREQ=MONTHLY;BYDAY=1SU",
		"timezone":             "America/Los_Angeles",
		"affected_entity_type": "SERVICE_LEVEL",
		"affected_entities":    []interface{}{"guid-1", "guid-2"},
	})

	got := expandMaintenanceWindowCreateInput(d)

	assert.Equal(t, "tf-test-mw", got.Name)
	assert.Equal(t, "a test maintenance window", got.Description)
	assert.Equal(t, maintenancewindow.MaintenanceWindowScopedReferenceInput{
		ID:   "1067061",
		Type: maintenancewindow.MaintenanceWindowEntityScopeTypes.ACCOUNT,
	}, got.Scope)
	assert.Equal(t, maintenancewindow.NaiveDateTime("2026-12-01T22:00:00"), got.StartTime)
	assert.Equal(t, workloads.Duration("PT1H"), got.Duration)
	assert.Equal(t, "FREQ=MONTHLY;BYDAY=1SU", got.Rrule)
	assert.Equal(t, "America/Los_Angeles", got.Timezone)
	assert.Equal(t, maintenancewindow.MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL, got.AffectedEntityType)
	assert.Equal(t, []string{"guid-1", "guid-2"}, got.AffectedEntities)
}

func TestExpandMaintenanceWindowUpdateInput(t *testing.T) {
	t.Parallel()

	d := testMaintenanceWindowResourceData(t, map[string]interface{}{
		"name":              "tf-test-mw-updated",
		"description":       "updated description",
		"start_time":        "2026-12-02T22:00:00",
		"duration":          "PT2H",
		"rrule":             "FREQ=WEEKLY",
		"timezone":          "UTC",
		"affected_entities": []interface{}{"guid-3"},
	})

	got := expandMaintenanceWindowUpdateInput(d)

	assert.Equal(t, "tf-test-mw-updated", got.Name)
	assert.Equal(t, "updated description", got.Description)
	assert.Equal(t, maintenancewindow.NaiveDateTime("2026-12-02T22:00:00"), got.StartTime)
	assert.Equal(t, workloads.Duration("PT2H"), got.Duration)
	assert.Equal(t, "FREQ=WEEKLY", got.Rrule)
	assert.Equal(t, "UTC", got.Timezone)
	assert.Equal(t, []string{"guid-3"}, got.AffectedEntities)
}

func TestExpandMaintenanceWindowAffectedEntities(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"a", "b"}, expandMaintenanceWindowAffectedEntities([]interface{}{"a", "b"}))
	assert.Equal(t, []string{}, expandMaintenanceWindowAffectedEntities([]interface{}{}))
}

func TestFlattenMaintenanceWindow(t *testing.T) {
	t.Parallel()

	d := testMaintenanceWindowResourceData(t, map[string]interface{}{})

	mw := &maintenancewindow.MaintenanceWindowConfigDefinition{
		ID:                 "12345",
		Name:               "tf-test-mw",
		Description:        "a test maintenance window",
		Scope:              maintenancewindow.MaintenanceWindowScopedReference{ID: "1067061", Type: maintenancewindow.MaintenanceWindowEntityScopeTypes.ACCOUNT},
		StartTime:          maintenancewindow.NaiveDateTime("2026-12-01T22:00:00"),
		Duration:           workloads.Duration("PT1H"),
		Rrule:              "FREQ=MONTHLY;BYDAY=1SU",
		Timezone:           "America/Los_Angeles",
		AffectedEntityType: maintenancewindow.MaintenanceWindowAffectedEntityTypeTypes.SERVICE_LEVEL,
		AffectedEntities:   []string{"guid-1", "guid-2"},
	}

	err := flattenMaintenanceWindow(mw, d)
	assert.NoError(t, err)

	assert.Equal(t, "tf-test-mw", d.Get("name"))
	assert.Equal(t, "a test maintenance window", d.Get("description"))
	assert.Equal(t, "1067061", d.Get("scope_id"))
	assert.Equal(t, "ACCOUNT", d.Get("scope_type"))
	assert.Equal(t, "2026-12-01T22:00:00", d.Get("start_time"))
	assert.Equal(t, "PT1H", d.Get("duration"))
	assert.Equal(t, "FREQ=MONTHLY;BYDAY=1SU", d.Get("rrule"))
	assert.Equal(t, "America/Los_Angeles", d.Get("timezone"))
	assert.Equal(t, "SERVICE_LEVEL", d.Get("affected_entity_type"))
	assert.Equal(t, []interface{}{"guid-1", "guid-2"}, d.Get("affected_entities"))
}
