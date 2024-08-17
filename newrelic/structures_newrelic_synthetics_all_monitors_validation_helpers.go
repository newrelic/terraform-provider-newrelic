package newrelic

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateSyntheticMonitorRuntimeAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []error

	err := validateSyntheticMonitorLegacyRuntimeAttributesOnCreate(d)
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

const SyntheticsRuntimeTypeAttrLabel string = "runtime_type"
const SyntheticsRuntimeTypeVersionAttrLabel string = "runtime_type_version"
const SyntheticsUseLegacyRuntimeAttrLabel string = "use_legacy_runtime_unsupported"
const SyntheticsNodeLegacyRuntimeType string = "NODE_API"
const SyntheticsNodeLegacyRuntimeTypeVersion string = "10"
const SyntheticsChromeBrowserLegacyRuntimeType string = "CHROME_BROWSER"
const SyntheticsChromeBrowserLegacyRuntimeTypeVersion string = "72"

func validateSyntheticMonitorLegacyRuntimeAttributesOnCreate(d *schema.ResourceDiff) []error {
	var runtimeAttributesValidationErrors []error

	isSyntheticMonitorCreated := d.Id() != ""

	_, runtimeTypeInConfig := d.GetChange(SyntheticsRuntimeTypeAttrLabel)
	isRuntimeTypeNil := runtimeTypeInConfig == ""

	_, runtimeTypeVersionInConfig := d.GetChange(SyntheticsRuntimeTypeVersionAttrLabel)
	isRuntimeTypeVersionNil := runtimeTypeVersionInConfig == ""

	_, useLegacyRuntimeInConfig := d.GetChange(SyntheticsUseLegacyRuntimeAttrLabel)
	useLegacyRuntime := useLegacyRuntimeInConfig == true

	// if !isSyntheticMonitorCreated is a condition that needs to exist only until October 22, 2024 as the first phase of Legacy Runtime EOL changes
	// aim at restricting new monitors from using the legacy runtime. For the release on October 22, 2024; this condition may be discarded.
	if !isSyntheticMonitorCreated && !useLegacyRuntime {
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
		"attribute `%s` is required to be specified with new runtime values, \"\" %s",
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

func syntheticMonitorConfigHasObsoleteRuntime(
	runtimeTypeInConfig interface{},
	runtimeTypeVersionInConfig interface{},
) bool {
	return (runtimeTypeInConfig == SyntheticsNodeLegacyRuntimeType && runtimeTypeVersionInConfig == SyntheticsNodeLegacyRuntimeTypeVersion) || (runtimeTypeInConfig == SyntheticsChromeBrowserLegacyRuntimeType && runtimeTypeVersionInConfig == SyntheticsChromeBrowserLegacyRuntimeTypeVersion)
}

func buildSyntheticsLegacyRuntimeValidationError() string {
	return `with new monitors starting August 26, 2024;
creating new monitors with the legacy runtime/without the new runtime is no longer supported.
This is in relation with the upcoming Synthetics Legacy Runtime EOL on October 22, 2024; see this for more details: 
https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm
https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide
`
}
