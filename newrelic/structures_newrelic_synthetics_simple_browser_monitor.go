package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		AdvancedOptions: synthetics.SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			DeviceEmulation: &synthetics.SyntheticsDeviceEmulationInput{},
		},
	}

	if v, ok := d.GetOk("custom_header"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.CustomHeaders = expandSyntheticsCustomHeaders(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("uri"); ok {
		simpleBrowserMonitorInput.Uri = v.(string)
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if v, ok := d.GetOk("locations_public"); ok {
		simpleBrowserMonitorInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleBrowserMonitorInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.ResponseValidationText = v.(string)
	}

	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleBrowserMonitorInput.AdvancedOptions.UseTlsValidation = &vs
	}

	sciptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		simpleBrowserMonitorInput.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			simpleBrowserMonitorInput.Runtime.ScriptLanguage = sciptLang.(string)
		}

		if runtimeTypeOk {
			simpleBrowserMonitorInput.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	if v, ok := d.GetOk("device_orientation"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(v.(string))
	}

	if v, ok := d.GetOk("device_type"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(v.(string))
	}

	return simpleBrowserMonitorInput
}

func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		AdvancedOptions: synthetics.SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{
			DeviceEmulation: &synthetics.SyntheticsDeviceEmulationInput{},
		},
	}

	if v, ok := d.GetOk("custom_header"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.CustomHeaders = expandSyntheticsCustomHeaders(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("uri"); ok {
		simpleBrowserMonitorUpdateInput.Uri = v.(string)
	}

	if v, ok := d.GetOk("locations_public"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if v, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.ResponseValidationText = v.(string)
	}

	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleBrowserMonitorUpdateInput.AdvancedOptions.UseTlsValidation = &vs
	}

	sciptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		simpleBrowserMonitorUpdateInput.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			simpleBrowserMonitorUpdateInput.Runtime.ScriptLanguage = sciptLang.(string)
		}

		if runtimeTypeOk {
			simpleBrowserMonitorUpdateInput.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			simpleBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	if v, ok := d.GetOk("device_orientation"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(v.(string))
	}

	if v, ok := d.GetOk("device_type"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(v.(string))
	}

	return simpleBrowserMonitorUpdateInput
}
