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
	// Always return a slice (even if empty) for proper drift detection
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
	// Always return a slice (even if empty) for proper drift detection
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

	// Only set operating_system for HOST fleets (not for KUBERNETESCLUSTER)
	// Always set it for HOST fleets to detect drift, even if API returns empty value
	// Use string comparison to be defensive against type issues
	if strings.ToUpper(string(fleet.ManagedEntityType)) == "HOST" {
		osType := string(fleet.OperatingSystem.Type)
		// For HOST fleets, always set operating_system (even if empty)
		// This ensures drift detection works correctly
		if err := d.Set("operating_system", osType); err != nil {
			return err
		}
	}

	if err := d.Set("description", fleet.Description); err != nil {
		return err
	}

	// Always set tags (even if empty) to detect drift when tags are removed externally
	if err := d.Set("tags", flattenEntityManagementTags(fleet.Tags)); err != nil {
		return err
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

	// Only set operating_system for HOST fleets (not for KUBERNETESCLUSTER)
	// Always set it for HOST fleets to detect drift, even if API returns empty value
	// Use string comparison to be defensive against type issues
	if strings.ToUpper(string(fleet.ManagedEntityType)) == "HOST" {
		osType := string(fleet.OperatingSystem.Type)
		// For HOST fleets, always set operating_system (even if empty)
		// This ensures drift detection works correctly
		if err := d.Set("operating_system", osType); err != nil {
			return err
		}
	}

	if err := d.Set("description", fleet.Description); err != nil {
		return err
	}

	// Always set tags (even if empty) to detect drift when tags are removed externally
	if err := d.Set("tags", flattenFleetControlTags(fleet.Tags)); err != nil {
		return err
	}

	if err := d.Set("organization_id", organizationID); err != nil {
		return err
	}

	return nil
}
