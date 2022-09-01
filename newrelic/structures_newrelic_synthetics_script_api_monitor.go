package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
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
