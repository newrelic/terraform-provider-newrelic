package newrelic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateSyntheticMonitorAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []error

	err := validateSyntheticMonitorLegacyRuntimeAttributesOnCreate(d)
	if err != nil {
		errorsList = append(errorsList, err...)
	}

	_, monitorType := d.GetChange("type")
	if monitorType != nil {
		isBrowserMonitor := strings.Contains(monitorType.(string), "BROWSER")
		if isBrowserMonitor {
			err := validateDevicesFields(d)
			if err != nil {
				errorsList = append(errorsList, err)
			}
		}
	}

	if len(errorsList) == 0 {
		return nil
	}

	errorsString := "the following validation errors have been identified with the configuration of the synthetic monitor: \n"

	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}

	return errors.New(errorsString)
}

func validateSyntheticMonitorLegacyRuntimeAttributesOnCreate(d *schema.ResourceDiff) []error {
	var runtimeAttributesValidationErrors []error

	_, runtimeTypeInConfig := d.GetChange(SyntheticsRuntimeTypeAttrLabel)
	isRuntimeTypeNil := runtimeTypeInConfig == ""

	_, runtimeTypeVersionInConfig := d.GetChange(SyntheticsRuntimeTypeVersionAttrLabel)
	isRuntimeTypeVersionNil := runtimeTypeVersionInConfig == ""

	_, monitorType := d.GetChange("type")
	isSimpleMonitor := monitorType == "SIMPLE"

	// Validation would include checking if runtime attribute values comprise values corresponding to the legacy runtime;
	// create/update requests of monitors would be blocked by the API via an error; which we're trying to reflect in Terraform.
	// also, this error scenario should not apply to SIMPLE Synthetic Monitors, as they do not support using runtime attributes.
	if !isSimpleMonitor {

		// if neither 'runtime_type' nor 'runtime_type_version' is nil,
		// check if either of these attributes are not nil and actually comprise values which signify the legacy runtime instead,
		// NODE_API 10 or CHROME_BROWSER 72; in which case, a similar error explaining the restriction would be thrown.
		if !isRuntimeTypeNil && !isRuntimeTypeVersionNil {
			if syntheticMonitorConfigHasObsoleteRuntime(runtimeTypeInConfig, runtimeTypeVersionInConfig) {
				runtimeAttributesValidationErrors = append(
					runtimeAttributesValidationErrors,
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						runtimeTypeInConfig.(string),
						runtimeTypeVersionInConfig.(string),
					),
				)
			}
		}
	}

	// add any other validation errors as needed above this line

	if len(runtimeAttributesValidationErrors) > 0 {
		return runtimeAttributesValidationErrors
	}

	return nil
}

func buildSyntheticsLegacyEmptyRuntimeError(attributeName string) error {
	return fmt.Errorf(
		"attribute `%s` is required to be specified with new runtime values, %s",
		attributeName,
		buildSyntheticsLegacyRuntimeValidationError(),
	)
}

func buildSyntheticsLegacyObsoleteRuntimeError(
	runtimeTypeAttrLabel string,
	runtimeTypeVersionAttrLabel string,
	runtimeTypeInConfig string,
	runtimeTypeVersionInConfig string,
) error {
	return fmt.Errorf(
		"legacy runtime version `%s` can no longer be specified as the `%s` corresponding to the `%s` `%s` %s",
		runtimeTypeVersionInConfig,
		runtimeTypeVersionAttrLabel,
		runtimeTypeAttrLabel,
		runtimeTypeInConfig,
		buildSyntheticsLegacyRuntimeValidationError(),
	)
}

func buildSyntheticsLegacyRuntimeValidationError() string {
	return `with new and existing monitors starting October 22, 2024;
creating and updating monitors comprising legacy runtime values/without the new runtime is *no longer supported*.
This is in relation with the Synthetics Legacy Runtime EOL which has taken effect on October 22, 2024; see the following for more details: 
https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm
https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide
`
}

// The following function will validate the device fields at the Terraform plan, ensuring that the user specifies
// either the devices field alone or both the device_type and device_orientation fields
func validateDevicesFields(d *schema.ResourceDiff) error {
	rawConfiguration := d.GetRawConfig()

	// GetAttr func will get below fields corresponding values from raw configuration that is from terraform configuration file
	devicesIsNil := rawConfiguration.GetAttr("devices").IsNull()
	deviceTypeIsNil := rawConfiguration.GetAttr("device_type").IsNull()
	deviceOrientationIsNil := rawConfiguration.GetAttr("device_orientation").IsNull()

	if !devicesIsNil && !(deviceTypeIsNil && deviceOrientationIsNil) {
		return fmt.Errorf(`cannot use 'devices', 'device_type', and 'device_orientation' simultaneously 
	use either 'devices' alone or both 'device_type' and 'device_orientation' fields together 
	we recommend using the 'devices' field, as it allows you to select multiple combinations of device types and orientations`)
	}
	return nil
}
