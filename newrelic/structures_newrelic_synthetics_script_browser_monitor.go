package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func buildSyntheticsScriptBrowserMonitorInput(d *schema.ResourceData) synthetics.SyntheticsCreateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateScriptBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v.(bool) {
		e := v.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(attr.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("script_language"); ok {
		input.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		input.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		input.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	return input
}

func buildSyntheticsScriptBrowserUpdateInput(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateScriptBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v.(bool) {
		e := v.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("script_language"); ok {
		input.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		input.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		input.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	return input
}
