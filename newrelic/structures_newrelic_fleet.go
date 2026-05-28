package newrelic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// flattenEntityManagementTags converts EntityManagementTag slice to "key:v1,v2" strings.
func flattenEntityManagementTags(tags []fleetcontrol.EntityManagementTag) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag.Values) > 0 {
			result = append(result, fmt.Sprintf("%s:%s", tag.Key, strings.Join(tag.Values, ",")))
		}
	}
	return result
}

// flattenFleetControlTags converts FleetControlTag slice to "key:v1,v2" strings.
func flattenFleetControlTags(tags []fleetcontrol.FleetControlTag) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag.Values) > 0 {
			result = append(result, fmt.Sprintf("%s:%s", tag.Key, strings.Join(tag.Values, ",")))
		}
	}
	return result
}

// applyFleetEntityToState writes the fleet entity fields shared by all fleet
// resource CRUD responses into Terraform state. operatingSystem is only set
// when non-empty and the managed entity type is HOST.
func applyFleetEntityToState(d *schema.ResourceData, name, managedEntityType, operatingSystem, description string, tags []string, organizationID string) error {
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("managed_entity_type", managedEntityType); err != nil {
		return err
	}
	if strings.ToUpper(managedEntityType) == "HOST" && operatingSystem != "" {
		if err := d.Set("operating_system", operatingSystem); err != nil {
			return err
		}
	}
	if err := d.Set("description", description); err != nil {
		return err
	}
	if err := d.Set("tags", tags); err != nil {
		return err
	}
	if err := d.Set("organization_id", organizationID); err != nil {
		return err
	}
	return nil
}

// flattenFleetEntity writes an EntityManagementFleetEntity (Read response) into state.
func flattenFleetEntity(fleet *fleetcontrol.EntityManagementFleetEntity, d *schema.ResourceData, organizationID string) error {
	return applyFleetEntityToState(d,
		fleet.Name,
		string(fleet.ManagedEntityType),
		string(fleet.OperatingSystem.Type),
		fleet.Description,
		flattenEntityManagementTags(fleet.Tags),
		organizationID,
	)
}

// flattenFleetControlEntity writes a FleetControlFleetEntityResult (Create/Update response) into state.
func flattenFleetControlEntity(fleet *fleetcontrol.FleetControlFleetEntityResult, d *schema.ResourceData, organizationID string) error {
	return applyFleetEntityToState(d,
		fleet.Name,
		string(fleet.ManagedEntityType),
		string(fleet.OperatingSystem.Type),
		fleet.Description,
		flattenFleetControlTags(fleet.Tags),
		organizationID,
	)
}
