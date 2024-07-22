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

func validateSyntheticMonitorLegacyRuntimeAttributesUponCreate(d *schema.ResourceDiff) []error {
	var runtimeAttributesValidationErrors []error
	runtimeTypeAttributeLabel := "runtime_type"
	runtimeTypeVersionAttributeLabel := "runtime_type_version"

	isSyntheticMonitorAlreadyCreated := d.Id() != ""
	rawConfiguration := d.GetRawConfig()

	isRuntimeTypeNotSpecifiedInConfiguration := rawConfiguration.GetAttr(runtimeTypeAttributeLabel).IsNull()
	_, runtimeTypeInConfig := d.GetChange(runtimeTypeAttributeLabel)
	isRuntimeTypeNullValue := runtimeTypeInConfig == ""

	// this would return true only if runtime_type_version is not specified in the configuration at all
	// and false, if runtime_type_version is specified either as an empty string "", or as any other non nil value (non-empty string)
	isRuntimeTypeVersionNotSpecifiedInConfiguration := rawConfiguration.GetAttr(runtimeTypeVersionAttributeLabel).IsNull()
	_, runtimeTypeVersionInConfig := d.GetChange(runtimeTypeAttributeLabel)

	// this would return true both when `runtime_type_version` is not specified in the config and when `runtime_type_version` is specified as "" in the config
	// and false, if `runtime_type_version` has a non nil value (a non-empty string) as its value
	isRuntimeTypeVersionNullValue := runtimeTypeVersionInConfig == ""

	if !isSyntheticMonitorAlreadyCreated && !isRuntimeTypeNotSpecifiedInConfiguration && isRuntimeTypeNullValue {
		runtimeAttributesValidationErrors = append(
			runtimeAttributesValidationErrors,
			constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(runtimeTypeAttributeLabel),
		)
	}

	if !isSyntheticMonitorAlreadyCreated && !isRuntimeTypeVersionNotSpecifiedInConfiguration && isRuntimeTypeVersionNullValue {
		runtimeAttributesValidationErrors = append(
			runtimeAttributesValidationErrors,
			constructSyntheticMonitorLegacyRuntimeAttributesEmptyValidationErrorUponCreate(runtimeTypeVersionAttributeLabel),
		)
	}

	/////////////////////////////
	// THIS IS WIP, not working
	/////////////////////////////

	if !isSyntheticMonitorAlreadyCreated && !isRuntimeTypeNotSpecifiedInConfiguration && !isRuntimeTypeVersionNotSpecifiedInConfiguration && !isRuntimeTypeNullValue && !isRuntimeTypeVersionNullValue {
		if runtimeTypeInConfig == "NODE_API" && runtimeTypeVersionInConfig == "10" {
			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
					runtimeTypeAttributeLabel,
					runtimeTypeVersionAttributeLabel,
					runtimeTypeInConfig.(string),
					runtimeTypeVersionInConfig.(string),
				),
			)
		}

		if runtimeTypeInConfig == "CHROME_BROWSER" && runtimeTypeVersionInConfig == "72" {
			runtimeAttributesValidationErrors = append(
				runtimeAttributesValidationErrors,
				constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
					runtimeTypeAttributeLabel,
					runtimeTypeVersionAttributeLabel,
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
	errorString := `can no longer be specified as an empty string with new monitors starting August 26, 2024; creating new monitors with the legacy runtime is no longer supported.
This is in relation with the upcoming Synthetics Legacy Runtime EOL on October 22, 2024; see this for more details: 
https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm`
	return fmt.Errorf("`%s` %s", attributeName, errorString)
}

func constructSyntheticMonitorLegacyRuntimeAttributesObsoleteValidationErrorUponCreate(
	runtimeTypeAttributeLabel string,
	runtimeTypeVersionAttributeLabel string,
	runtimeTypeInConfig string,
	runtimeTypeVersionInConfig string,
) error {
	errorString := ` with new monitors starting August 26, 2024; creating new monitors with the legacy runtime is no longer supported.
This is in relation with the upcoming Synthetics Legacy Runtime EOL on October 22, 2024; see this for more details: 
https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm`
	return fmt.Errorf(
		"legacy runtime version `%s` can no longer be specified as the `%s` corresponding to the `%s` `%s` %s",
		runtimeTypeVersionInConfig,
		runtimeTypeVersionAttributeLabel,
		runtimeTypeAttributeLabel,
		runtimeTypeInConfig,
		errorString,
	)
}
