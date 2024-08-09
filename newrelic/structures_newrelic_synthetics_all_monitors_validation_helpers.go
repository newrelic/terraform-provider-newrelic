package newrelic

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateSyntheticMonitorRuntimeAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []error

	err := validateSyntheticMonitorLegacyRuntimeAttributesUponCreate(d)
	if err != nil {
		errorsList = append(errorsList, err...)
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

const SyntheticsRuntimeTypeAttributeLabel string = "runtime_type"
const SyntheticsRuntimeTypeVersionAttributeLabel string = "runtime_type_version"
const SyntheticsNodeLegacyRuntimeType string = "NODE_API"
const SyntheticsNodeLegacyRuntimeTypeVersion string = "10"
const SyntheticsChromeBrowserLegacyRuntimeType string = "CHROME_BROWSER"
const SyntheticsChromeBrowserLegacyRuntimeTypeVersion string = "72"

func validateSyntheticMonitorLegacyRuntimeAttributesUponCreate(d *schema.ResourceDiff) []error {
	var runtimeAttributesValidationErrors []error

	isSyntheticMonitorAlreadyCreated := d.Id() != ""
	rawConfiguration := d.GetRawConfig()

	isRuntimeTypeNotSpecifiedInConfiguration := rawConfiguration.GetAttr(SyntheticsRuntimeTypeAttributeLabel).IsNull()
	_, runtimeTypeInConfig := d.GetChange(SyntheticsRuntimeTypeAttributeLabel)

	isRuntimeTypeNullValue := runtimeTypeInConfig == ""

	// this would return true only if runtime_type_version is not specified in the configuration at all
	// and false, if runtime_type_version is specified either as an empty string "", or as any other non nil value (non-empty string)
	isRuntimeTypeVersionNotSpecifiedInConfiguration := rawConfiguration.GetAttr(SyntheticsRuntimeTypeVersionAttributeLabel).IsNull()
	_, runtimeTypeVersionInConfig := d.GetChange(SyntheticsRuntimeTypeVersionAttributeLabel)

	// this would return true both when `runtime_type_version` is not specified in the config and when `runtime_type_version` is specified as "" in the config
	// and false, if `runtime_type_version` has a non nil value (a non-empty string) as its value
	isRuntimeTypeVersionNullValue := runtimeTypeVersionInConfig == ""

	if !isSyntheticMonitorAlreadyCreated &&
		!isRuntimeTypeNotSpecifiedInConfiguration &&
		isRuntimeTypeNullValue {

		runtimeAttributesValidationErrors = append(
			runtimeAttributesValidationErrors,
			constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(SyntheticsRuntimeTypeAttributeLabel),
		)
	}

	if !isSyntheticMonitorAlreadyCreated &&
		!isRuntimeTypeVersionNotSpecifiedInConfiguration &&
		isRuntimeTypeVersionNullValue {

		runtimeAttributesValidationErrors = append(
			runtimeAttributesValidationErrors,
			constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(SyntheticsRuntimeTypeVersionAttributeLabel),
		)
	}

	if !isSyntheticMonitorAlreadyCreated &&
		!isRuntimeTypeNotSpecifiedInConfiguration &&
		!isRuntimeTypeVersionNotSpecifiedInConfiguration &&
		!isRuntimeTypeNullValue &&
		!isRuntimeTypeVersionNullValue {

		if runtimeTypeInConfig == SyntheticsNodeLegacyRuntimeType &&
			runtimeTypeVersionInConfig == SyntheticsNodeLegacyRuntimeTypeVersion {

			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
					SyntheticsRuntimeTypeAttributeLabel,
					SyntheticsRuntimeTypeVersionAttributeLabel,
					runtimeTypeInConfig.(string),
					runtimeTypeVersionInConfig.(string),
				),
			)
		}

		if runtimeTypeInConfig == SyntheticsChromeBrowserLegacyRuntimeType &&
			runtimeTypeVersionInConfig == SyntheticsChromeBrowserLegacyRuntimeTypeVersion {

			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
					SyntheticsRuntimeTypeAttributeLabel,
					SyntheticsRuntimeTypeVersionAttributeLabel,
					runtimeTypeInConfig.(string),
					runtimeTypeVersionInConfig.(string),
				),
			)
		}
	}

	// add any other validation errors as needed above this line

	if len(runtimeAttributesValidationErrors) > 0 {
		return runtimeAttributesValidationErrors
	}

	return nil
}

func constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(attributeName string) error {
	return fmt.Errorf(
		"`%s` can no longer be specified as an empty string \"\" %s",
		attributeName,
		constructSyntheticMonitorLegacyRuntimeAttributesValidationErrorUponCreate(),
	)
}

func constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
	runtimeTypeAttributeLabel string,
	runtimeTypeVersionAttributeLabel string,
	runtimeTypeInConfig string,
	runtimeTypeVersionInConfig string,
) error {
	return fmt.Errorf(
		"legacy runtime version `%s` can no longer be specified as the `%s` corresponding to the `%s` `%s` %s",
		runtimeTypeVersionInConfig,
		runtimeTypeVersionAttributeLabel,
		runtimeTypeAttributeLabel,
		runtimeTypeInConfig,
		constructSyntheticMonitorLegacyRuntimeAttributesValidationErrorUponCreate(),
	)
}

func constructSyntheticMonitorLegacyRuntimeAttributesValidationErrorUponCreate() string {
	return `with new monitors starting August 26, 2024;
creating new monitors with the legacy runtime is no longer supported.
This is in relation with the upcoming Synthetics Legacy Runtime EOL on October 22, 2024; see this for more details: 
https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm
`
}
