package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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

	isSyntheticMonitorCreated := d.Id() != ""

	_, runtimeTypeInConfig := d.GetChange(SyntheticsRuntimeTypeAttrLabel)
	isRuntimeTypeNil := runtimeTypeInConfig == ""

	_, runtimeTypeVersionInConfig := d.GetChange(SyntheticsRuntimeTypeVersionAttrLabel)
	isRuntimeTypeVersionNil := runtimeTypeVersionInConfig == ""

	_, useLegacyRuntimeInConfig := d.GetChange(SyntheticsUseLegacyRuntimeAttrLabel)
	useLegacyRuntime := useLegacyRuntimeInConfig == true

	_, monitorType := d.GetChange("type")
	isSimpleMonitor := monitorType == "SIMPLE"

	// in this first condition, we're trying to make sure 'use_unsupported_legacy_runtime' is only being used with the legacy runtime
	// and not with any sort of runtime values which signify the new runtime (since the intent of using this attribute
	// is to skip Terraform validation to create monitors in the legacy runtime after the August 26 EOL if exempt by the API)
	if useLegacyRuntime {
		if !syntheticMonitorConfigHasObsoleteRuntime(runtimeTypeInConfig, runtimeTypeVersionInConfig) &&
			(!isRuntimeTypeNil || !isRuntimeTypeVersionNil) {
			return []error{
				fmt.Errorf(
					`'%s' is intended to be used with legacy runtime values of '%s' and '%s' to skip restrictions on using
the legacy runtime, after the EOL of the Synthetics Legacy Runtime, if you have been granted an exemption to use the legacy runtime.
However, the configuration of the current resource seems to comprise values of these runtime attributes corresponding to the new runtime.
Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.
							`,
					SyntheticsUseLegacyRuntimeAttrLabel,
					SyntheticsRuntimeTypeAttrLabel,
					SyntheticsRuntimeTypeVersionAttrLabel,
					SyntheticsUseLegacyRuntimeAttrLabel,
				),
			}
		}
	}

	// apply further validation to block usage of the legacy runtime with monitors owing to the EOL, ONLY if 'use_unsupported_legacy_runtime' is false
	// (which it is, by default, unless made true by the customer) AND ONLY for a new monitor and not an existing one, as the August 26
	// Legacy Runtime EOL (the first phase) only applies to new monitors. Validation would include checking if runtime attribute are nil or if they
	// are not nil and comprise values corresponding to the legacy runtime; as they would lead to creating monitors in the legacy runtime either way,
	// and both of these cases in create requests of monitors would be blocked by the API via an error; which we're trying to reflect in Terraform.
	// also, this error scenario should not apply to SIMPLE Synthetic Monitors, as they do not support using runtime attributes.
	if !isSyntheticMonitorCreated && !useLegacyRuntime && !isSimpleMonitor {

		// if 'use_unsupported_legacy_runtime' is false (which it is, by default), check if 'runtime_type' and 'runtime_type_version' are nil
		// if either of these two runtime attributes are nil, throw a relevant error to explain that this is no longer allowed
		// owing to the Legacy Runtime EOL.
		if isRuntimeTypeNil {
			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel),
			)
		}

		if isRuntimeTypeVersionNil {
			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeVersionAttrLabel),
			)
		}

		// if 'use_unsupported_legacy_runtime' is false (which it is, by default) and neither 'runtime_type' nor 'runtime_type_version' is nil,
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
	return `with new monitors starting August 26, 2024;
creating new monitors with the legacy runtime/without the new runtime is no longer supported.
This is in relation with the upcoming Synthetics Legacy Runtime EOL on October 22, 2024; see this for more details: 
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
		return fmt.Errorf(`Cannot use 'devices', 'device_type', and 'device_orientation' simultaneously. 
	Use either 'devices' alone or both 'device_type' and 'device_orientation' fields together. 
	We recommend using the 'devices' field, as it allows you to select multiple combinations of device types and orientations.`)
	}

	if deviceTypeIsNil != deviceOrientationIsNil {
		return fmt.Errorf("you need to specify both 'device_type' and 'device_orientation' fields; you can't use just one of them")
	}
	return nil
}
