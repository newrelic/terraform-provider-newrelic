package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsScriptAPIMonitorInput(d *schema.ResourceData) synthetics.SyntheticsCreateScriptAPIMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateScriptAPIMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}

	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(attr.(*schema.Set).List())
	}
	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(attr.(*schema.Set).List())
	}

	_, scriptLangOk := d.GetOk("script_language")
	_, runtimeTypeOk := d.GetOk("runtime_type")
	_, runtimeVersionOk := d.GetOk("runtime_type_version")

	input.Runtime = &synthetics.SyntheticsRuntimeInput{}
	if scriptLangOk || runtimeTypeOk || runtimeVersionOk {
		if v, ok := d.GetOk("script_language"); ok {
			input.Runtime.ScriptLanguage = v.(string)
		}
		if v, ok := d.GetOk("runtime_type"); ok {
			input.Runtime.RuntimeType = v.(string)
		}
		if v, ok := d.GetOk("runtime_type_version"); ok {
			input.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
		}
	} else {
		input.Runtime.RuntimeType = ""
		input.Runtime.RuntimeTypeVersion = ""
	}

	return input
}

func buildSyntheticsScriptAPIMonitorUpdateInput(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptAPIMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateScriptAPIMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}

	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}

	_, scriptLangOk := d.GetOk("script_language")
	_, runtimeTypeOk := d.GetOk("runtime_type")
	_, runtimeVersionOk := d.GetOk("runtime_type_version")

	input.Runtime = &synthetics.SyntheticsRuntimeInput{}
	if scriptLangOk || runtimeTypeOk || runtimeVersionOk {
		if v, ok := d.GetOk("script_language"); ok {
			input.Runtime.ScriptLanguage = v.(string)
		}
		if v, ok := d.GetOk("runtime_type"); ok {
			input.Runtime.RuntimeType = v.(string)
		}
		if v, ok := d.GetOk("runtime_type_version"); ok {
			input.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
		}
	} else {
		input.Runtime.RuntimeType = ""
		input.Runtime.RuntimeTypeVersion = ""
	}

	return input
}
