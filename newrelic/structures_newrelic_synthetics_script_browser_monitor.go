package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsScriptBrowserMonitorInput(d *schema.ResourceData) synthetics.SyntheticsCreateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateScriptBrowserMonitorInput{
		Name:            inputBase.Name,
		Period:          inputBase.Period,
		Status:          inputBase.Status,
		Tags:            inputBase.Tags,
		Script:          d.Get("script").(string),
		AdvancedOptions: synthetics.SyntheticsScriptBrowserMonitorAdvancedOptionsInput{},
	}

	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(attr.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	input.Browsers = &typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	input.Devices = &typedDevices

	sciptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		input.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			input.Runtime.ScriptLanguage = sciptLang.(string)
		}

		if runtimeTypeOk {
			input.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			input.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	enableScreenshot := d.Get("enable_screenshot_on_failure_and_script").(bool)
	if enableScreenshot {
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &enableScreenshot
	}

	deviceType, deviceTypeOk := d.GetOk("device_type")
	deviceOrientation, deviceOrientationOk := d.GetOk("device_orientation")
	if deviceTypeOk || deviceOrientationOk {
		input.AdvancedOptions.DeviceEmulation = &synthetics.SyntheticsDeviceEmulationInput{}

		if deviceTypeOk {
			input.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(deviceType.(string))
		}

		if deviceOrientationOk {
			input.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(deviceOrientation.(string))
		}
	}

	return input
}

func buildSyntheticsScriptBrowserUpdateInput(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateScriptBrowserMonitorInput{
		Name:            inputBase.Name,
		Period:          inputBase.Period,
		Status:          inputBase.Status,
		Tags:            inputBase.Tags,
		Script:          d.Get("script").(string),
		AdvancedOptions: synthetics.SyntheticsScriptBrowserMonitorAdvancedOptionsInput{},
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(v.(*schema.Set).List())
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	input.Browsers = &typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	input.Devices = &typedDevices

	sciptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")
	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		input.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			input.Runtime.ScriptLanguage = sciptLang.(string)
		}

		if runtimeTypeOk {
			input.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			input.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	enableScreenshot := d.Get("enable_screenshot_on_failure_and_script").(bool)
	if enableScreenshot {
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &enableScreenshot
	}

	deviceType, deviceTypeOk := d.GetOk("device_type")
	deviceOrientation, deviceOrientationOk := d.GetOk("device_orientation")
	if deviceTypeOk || deviceOrientationOk {
		input.AdvancedOptions.DeviceEmulation = &synthetics.SyntheticsDeviceEmulationInput{}

		if deviceTypeOk {
			input.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(deviceType.(string))
		}

		if deviceOrientationOk {
			input.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(deviceOrientation.(string))
		}
	}

	return input
}
