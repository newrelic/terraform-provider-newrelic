package newrelic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// Fleet resource expand/flatten functions

// flattenEntityManagementTags converts EntityManagementTag array to Terraform-friendly format
func flattenEntityManagementTags(tags []fleetcontrol.EntityManagementTag) []string {
	if len(tags) == 0 {
		return nil
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag.Values) > 0 {
			tagStr := fmt.Sprintf("%s:%s", tag.Key, strings.Join(tag.Values, ","))
			result = append(result, tagStr)
		}
	}
	return result
}

// flattenFleetControlTags converts FleetControlTag array to Terraform-friendly format
func flattenFleetControlTags(tags []fleetcontrol.FleetControlTag) []string {
	if len(tags) == 0 {
		return nil
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag.Values) > 0 {
			tagStr := fmt.Sprintf("%s:%s", tag.Key, strings.Join(tag.Values, ","))
			result = append(result, tagStr)
		}
	}
	return result
}

// flattenFleetEntity flattens a fleet entity from EntityManagement into Terraform state
func flattenFleetEntity(fleet *fleetcontrol.EntityManagementFleetEntity, d *schema.ResourceData, organizationID string) error {
	if err := d.Set("name", fleet.Name); err != nil {
		return err
	}

	if err := d.Set("managed_entity_type", string(fleet.ManagedEntityType)); err != nil {
		return err
	}

	if fleet.OperatingSystem.Type != "" {
		if err := d.Set("operating_system", string(fleet.OperatingSystem.Type)); err != nil {
			return err
		}
	}

	if err := d.Set("description", fleet.Description); err != nil {
		return err
	}


	if len(fleet.Tags) > 0 {
		if err := d.Set("tags", flattenEntityManagementTags(fleet.Tags)); err != nil {
			return err
		}
	}

	if err := d.Set("organization_id", organizationID); err != nil {
		return err
	}

	return nil
}

// flattenFleetControlEntity flattens a FleetControlFleetEntityResult into Terraform state
func flattenFleetControlEntity(fleet *fleetcontrol.FleetControlFleetEntityResult, d *schema.ResourceData, organizationID string) error {
	if err := d.Set("name", fleet.Name); err != nil {
		return err
	}

	if err := d.Set("managed_entity_type", string(fleet.ManagedEntityType)); err != nil {
		return err
	}

	if fleet.OperatingSystem.Type != "" {
		if err := d.Set("operating_system", string(fleet.OperatingSystem.Type)); err != nil {
			return err
		}
	}

	if err := d.Set("description", fleet.Description); err != nil {
		return err
	}


	if len(fleet.Tags) > 0 {
		if err := d.Set("tags", flattenFleetControlTags(fleet.Tags)); err != nil {
			return err
		}
	}

	if err := d.Set("organization_id", organizationID); err != nil {
		return err
	}

	return nil
}
