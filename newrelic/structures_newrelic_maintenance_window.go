package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/maintenancewindow"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func expandMaintenanceWindowCreateInput(d *schema.ResourceData) maintenancewindow.MaintenanceWindowConfigCreateInput {
	return maintenancewindow.MaintenanceWindowConfigCreateInput{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		StartTime:          maintenancewindow.NaiveDateTime(d.Get("start_time").(string)),
		Duration:           workloads.Duration(d.Get("duration").(string)),
		Rrule:              d.Get("rrule").(string),
		Timezone:           d.Get("timezone").(string),
		AffectedEntityType: maintenancewindow.MaintenanceWindowAffectedEntityType(d.Get("affected_entity_type").(string)),
		AffectedEntities:   expandMaintenanceWindowAffectedEntities(d.Get("affected_entities").([]interface{})),
		Scope: maintenancewindow.MaintenanceWindowScopedReferenceInput{
			ID:   d.Get("scope_id").(string),
			Type: maintenancewindow.MaintenanceWindowEntityScope(d.Get("scope_type").(string)),
		},
	}
}

func expandMaintenanceWindowUpdateInput(d *schema.ResourceData) maintenancewindow.MaintenanceWindowConfigUpdateInput {
	return maintenancewindow.MaintenanceWindowConfigUpdateInput{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		StartTime:        maintenancewindow.NaiveDateTime(d.Get("start_time").(string)),
		Duration:         workloads.Duration(d.Get("duration").(string)),
		Rrule:            d.Get("rrule").(string),
		Timezone:         d.Get("timezone").(string),
		AffectedEntities: expandMaintenanceWindowAffectedEntities(d.Get("affected_entities").([]interface{})),
	}
}

func expandMaintenanceWindowAffectedEntities(cfg []interface{}) []string {
	entities := make([]string, len(cfg))
	for i, e := range cfg {
		entities[i] = e.(string)
	}
	return entities
}

func flattenMaintenanceWindow(mw *maintenancewindow.MaintenanceWindowConfigDefinition, d *schema.ResourceData) error {
	if err := d.Set("name", mw.Name); err != nil {
		return err
	}
	if err := d.Set("description", mw.Description); err != nil {
		return err
	}
	if err := d.Set("scope_id", mw.Scope.ID); err != nil {
		return err
	}
	if err := d.Set("scope_type", string(mw.Scope.Type)); err != nil {
		return err
	}
	if err := d.Set("start_time", string(mw.StartTime)); err != nil {
		return err
	}
	if err := d.Set("duration", string(mw.Duration)); err != nil {
		return err
	}
	if err := d.Set("rrule", mw.Rrule); err != nil {
		return err
	}
	if err := d.Set("timezone", mw.Timezone); err != nil {
		return err
	}
	if err := d.Set("affected_entity_type", string(mw.AffectedEntityType)); err != nil {
		return err
	}
	if err := d.Set("affected_entities", mw.AffectedEntities); err != nil {
		return err
	}

	return nil
}
