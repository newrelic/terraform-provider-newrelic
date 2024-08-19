//go:build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Step Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Step Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Step Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "100"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Step Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Step Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Step Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Step Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Simple Browser Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Simple Browser Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Simple Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "100"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Simple Browser Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Simple Browser Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Simple Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Simple Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedBrowserMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Scripted Browser Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
					"SCRIPT_BROWSER",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Scripted Browser Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
					"SCRIPT_BROWSER",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Scripted Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
					"SCRIPT_BROWSER",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedBrowserMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "100"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Scripted Browser Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
					"SCRIPT_BROWSER",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted Browser Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
					"SCRIPT_BROWSER",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
					"SCRIPT_BROWSER",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted Browser Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
					"SCRIPT_BROWSER",
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedAPIMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Scripted API Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
					"SCRIPT_API",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Scripted API Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
					"SCRIPT_API",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Scripted API Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
					"SCRIPT_API",
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsScriptedAPIMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Scripted API Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
					"SCRIPT_API",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted API Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
					"SCRIPT_API",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted API Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
					"SCRIPT_API",
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Scripted API Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
					"SCRIPT_API",
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Broken Links Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Broken Links Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Broken Links Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Broken Links Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Broken Links Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Broken Links Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Broken Links Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributes_Errors(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Cert Check Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Cert Check Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyEmptyRuntimeError(SyntheticsRuntimeTypeAttrLabel).Error(),
				),
			},

			// create a Cert Check Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see an error as there exists no use_legacy_runtime_unsupported in the configuration
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					false,
				),
				ExpectError: regexp.MustCompile(
					buildSyntheticsLegacyObsoleteRuntimeError(
						SyntheticsRuntimeTypeAttrLabel,
						SyntheticsRuntimeTypeVersionAttrLabel,
						SyntheticsChromeBrowserRuntimeType,
						SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					).Error(),
				),
			},
		},
	})
}

func TestAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributes_ErrorsSkippedByUseLegacyRuntime(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	newRuntimeVersion := "10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{

			// create a Cert Check Monitor with runtime attributes in the configuration as empty strings (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Cert Check Monitor with no runtime attributes in the configuration at all (i.e. Legacy Runtime)
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					"",
					"",
					false,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Cert Check Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime
			// the expected outcome is to see NO error as use_legacy_runtime_unsupported is now added to the config with the value 'true'
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					SyntheticsChromeBrowserLegacyRuntimeTypeVersion,
					true,
					true,
				),
				ExpectError:        nil,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},

			// create a Cert Check Monitor with runtime attributes comprising non nil values, but corresponding to the legacy runtime (the version though, is not of the legacy runtime)
			// the expected outcome is to see an error as use_legacy_runtime_unsupported is now added to the config with the value 'true' and we're trying to use runtime_type_version with the new runtime
			{
				Config: testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
					rName,
					SyntheticsChromeBrowserRuntimeType,
					newRuntimeVersion,
					true,
					true,
				),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`Please use '%s' only with runtime attributes with values corresponding to the legacy runtime.`,
						SyntheticsUseLegacyRuntimeAttrLabel,
					),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsStepMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_step_monitor" "legacy_synthetics_step_monitor" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			steps {
				ordinal = 0
				type    = "NAVIGATE"
				values  = ["https://google.com"]
			}			
			`+testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
		runtimeType,
		runtimeTypeVersion,
		runtimeAttributesExistInConfig,
		useLegacyRuntimeUnsupportedInConfig,
	)+
		`
		}
`,
		name,
	)
}

func testAccNewRelicSyntheticsSimpleBrowserMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_monitor" "legacy_synthetics_monitor" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			type									= "BROWSER"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			`+testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
		runtimeType,
		runtimeTypeVersion,
		runtimeAttributesExistInConfig,
		useLegacyRuntimeUnsupportedInConfig,
	)+
		`
		}
`,
		name,
	)
}

func testAccNewRelicSyntheticsScriptedMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
	scriptType string,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_script_monitor" "legacy_synthetics_script_monitor" {
			name                                    = "%[1]s"
			period                                  = "EVERY_DAY"
			type									= "%[2]s"
			status                                  = "ENABLED"
			locations_public                        = ["US_WEST_2"]
			enable_screenshot_on_failure_and_script = true
			script									= "console.log('');"
			`+testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
		runtimeType,
		runtimeTypeVersion,
		runtimeAttributesExistInConfig,
		useLegacyRuntimeUnsupportedInConfig,
	)+
		`
		}
`,
		name,
		scriptType,
	)
}

func testAccNewRelicSyntheticsBrokenLinksMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_broken_links_monitor" "legacy_synthetics_broken_links_monitor" {
			  name                 = "%[1]s"
			  uri                  = "https://www.one.example.com"
			  locations_public     = ["AP_SOUTH_1"]
			  period               = "EVERY_6_HOURS"
			  status               = "ENABLED"
			  tag {
				key    = "some_key"
				values = ["some_value"]
			  }
			`+testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
		runtimeType,
		runtimeTypeVersion,
		runtimeAttributesExistInConfig,
		useLegacyRuntimeUnsupportedInConfig,
	)+
		`
			}
`,
		name,
	)
}

func testAccNewRelicSyntheticsCertCheckMonitor_CreateWithLegacyRuntimeAttributesConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_cert_check_monitor" "legacy_synthetics_cert_check_monitor" {
			  name                   = "%[1]s"
			  domain                 = "www.example.com"
			  locations_public       = ["AP_SOUTH_1"]
			  certificate_expiration = "10"
			  period                 = "EVERY_6_HOURS"
			  status                 = "ENABLED"
			  tag {
				key    = "some_key"
				values = ["some_value"]
			  }
			`+testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
		runtimeType,
		runtimeTypeVersion,
		runtimeAttributesExistInConfig,
		useLegacyRuntimeUnsupportedInConfig,
	)+
		`
			}
`,
		name,
	)
}

func testAccNewRelicSyntheticsAllMonitors_ConfigureRuntimeAttributesInConfig(
	runtimeType string,
	runtimeTypeVersion string,
	runtimeAttributesExistInConfig bool,
	useLegacyRuntimeUnsupportedInConfig bool,
) string {
	if runtimeAttributesExistInConfig && !useLegacyRuntimeUnsupportedInConfig {
		return fmt.Sprintf(`
		runtime_type = "%s"
		runtime_type_version = "%s"

`,
			runtimeType,
			runtimeTypeVersion)
	} else if !runtimeAttributesExistInConfig && useLegacyRuntimeUnsupportedInConfig {
		return fmt.Sprintf(`
	%s = true
`,
			SyntheticsUseLegacyRuntimeAttrLabel)
	} else if runtimeAttributesExistInConfig && useLegacyRuntimeUnsupportedInConfig {
		runtimeAttributesString := fmt.Sprintf(`
		runtime_type = "%s"
		runtime_type_version = "%s"

`,
			runtimeType,
			runtimeTypeVersion)

		useLegacyRuntimeUnsupportedString := fmt.Sprintf(`
	%s = true
`,
			SyntheticsUseLegacyRuntimeAttrLabel)
		return runtimeAttributesString + useLegacyRuntimeUnsupportedString
	} else {
		return ""
	}

}
