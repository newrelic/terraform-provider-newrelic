package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{}

	simpleBrowserMonitorInput.Name = inputBase.Name
	simpleBrowserMonitorInput.Period = inputBase.Period
	simpleBrowserMonitorInput.Status = inputBase.Status
	simpleBrowserMonitorInput.Tags = inputBase.Tags

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

	if v, ok := d.GetOk("script_language"); ok {
		simpleBrowserMonitorInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	return simpleBrowserMonitorInput
}

func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{}

	simpleBrowserMonitorUpdateInput.Name = inputBase.Name
	simpleBrowserMonitorUpdateInput.Period = inputBase.Period
	simpleBrowserMonitorUpdateInput.Status = inputBase.Status
	simpleBrowserMonitorUpdateInput.Tags = inputBase.Tags

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

	if v, ok := d.GetOk("script_language"); ok {
		simpleBrowserMonitorUpdateInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		simpleBrowserMonitorUpdateInput.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		simpleBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	return simpleBrowserMonitorUpdateInput
}
