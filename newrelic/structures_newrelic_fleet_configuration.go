package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// Fleet configuration resource expand/flatten functions

// flattenFleetConfiguration flattens a CreateConfigurationResponse into Terraform state
func flattenFleetConfiguration(config *fleetcontrol.CreateConfigurationResponse, d *schema.ResourceData, organizationID string) error {
	// Set configuration_id
	if err := d.Set("configuration_id", config.ConfigurationEntityGUID); err != nil {
		return err
	}

	// Set version from the configuration version entity
	if err := d.Set("version", config.ConfigurationVersion.ConfigurationVersionNumber); err != nil {
		return err
	}

	// Set organization_id
	if err := d.Set("organization_id", organizationID); err != nil {
		return err
	}

	// Note: name is already set from the input, no need to set it again here

	return nil
}
