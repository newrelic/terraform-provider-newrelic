package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type SyntheticsMonitorBase struct {
	//Locations     synthetics.SyntheticsLocationsInput
	Name          string
	Period        synthetics.SyntheticsMonitorPeriod
	Status        synthetics.SyntheticsMonitorStatus
	Tags          []synthetics.SyntheticsTag
	URI           string                                   // Move URI outside of base (does not apply to all monitors)
	CustomHeaders []synthetics.SyntheticsCustomHeaderInput // Move CustomHeaders outside of base (does not apply to all monitors)
}

// nolint:revive
var SyntheticsMonitorTypes = struct {
	SIMPLE         SyntheticsMonitorType
	BROWSER        SyntheticsMonitorType
	SCRIPT_API     SyntheticsMonitorType
	SCRIPT_BROWSER SyntheticsMonitorType
}{
	SIMPLE:         "SIMPLE",
	BROWSER:        "BROWSER",
	SCRIPT_API:     "SCRIPT_API",
	SCRIPT_BROWSER: "SCRIPT_BROWSER",
}

//validation function to validate monitor period
func listValidSyntheticsMonitorPeriods() []string {
	return []string{
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY),
	}
}

//validate func to validate monitor status
func listValidSyntheticsMonitorStatuses() []string {
	return []string{
		string(synthetics.SyntheticsMonitorStatusTypes.DISABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.ENABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.MUTED),
	}
}

//validate func to validate monitor type.
func listValidSyntheticsScriptMonitorTypes() []string {
	return []string{
		string(SyntheticsMonitorTypes.SCRIPT_API),
		string(SyntheticsMonitorTypes.SCRIPT_BROWSER),
	}
}

//func to build input for script API monitor creation.
func buildSyntheticsScriptAPIMonitorInput(d *schema.ResourceData) synthetics.SyntheticsCreateScriptAPIMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateScriptAPIMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}

	if attr, ok := d.GetOk("locations_private"); ok {
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

//function to build input for script Browser monitor creation.
func buildSyntheticsScriptBrowserMonitorInput(d *schema.ResourceData) synthetics.SyntheticsCreateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)
	input := synthetics.SyntheticsCreateScriptBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}
	if v, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = v.(bool)
	}
	if attr, ok := d.GetOk("locations_private"); ok {
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

//function to build input for script API monitor update function
func buildSyntheticsScriptAPIMonitorUpdateInput(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptAPIMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)
	input := synthetics.SyntheticsUpdateScriptAPIMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}
	if v, ok := d.GetOk("locations_private"); ok {
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

//function to build input for script browser monitor update function
func buildSyntheticsScriptBrowserUpdateInput(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)
	input := synthetics.SyntheticsUpdateScriptBrowserMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Script: d.Get("script").(string),
	}
	if v, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = v.(bool)
	}
	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandSyntheticsPublicLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("locations_private"); ok {
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

//func to expand synthetics private locations.
func expandSyntheticsPrivateLocations(locations []interface{}) []synthetics.SyntheticsPrivateLocationInput {
	locationsOut := make([]synthetics.SyntheticsPrivateLocationInput, len(locations))

	for i, v := range locations {
		pl := v.(map[string]interface{})
		locationsOut[i].GUID = pl["guid"].(int) // This should be a string I think since it's a scalar `ID`
		if v, ok := pl["vse_password"]; ok {
			locationsOut[i].VsePassword = synthetics.SecureValue(v.(string))
		}
	}
	return locationsOut
}

//function to expand synthetics public locations.
func expandSyntheticsPublicLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))

	for i, v := range locations {
		locationsOut[i] = v.(string)
	}

	return locationsOut
}

//function to expand monitor base for monitors input.
func expandSyntheticsMonitorBase(d *schema.ResourceData) SyntheticsMonitorBase {
	inputBase := SyntheticsMonitorBase{}

	name := d.Get("name")
	inputBase.Name = name.(string)

	status := d.Get("status")
	inputBase.Status = synthetics.SyntheticsMonitorStatus(status.(string))

	period := d.Get("period")
	inputBase.Period = synthetics.SyntheticsMonitorPeriod(period.(string))

	if uri, ok := d.GetOk("uri"); ok {
		inputBase.URI = uri.(string)
	}

	if tags, ok := d.GetOk("tag"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	if headers, ok := d.GetOk("custom_headers"); ok {
		inputBase.CustomHeaders = expandSyntheticsCustomHeaders(headers.(*schema.Set).List())
	}

	return inputBase
}

//function to expand custom headers
func expandSyntheticsCustomHeaders(headers []interface{}) []synthetics.SyntheticsCustomHeaderInput {
	output := make([]synthetics.SyntheticsCustomHeaderInput, len(headers))
	for i, v := range headers {
		header := v.(map[string]interface{})
		expanded := synthetics.SyntheticsCustomHeaderInput{
			Name:  header["name"].(string),
			Value: header["value"].(string),
		}
		output[i] = expanded
	}
	return output
}

//function to expand simple public locations
func expandSyntheticsSimplePublicLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))
	for i, v := range locations {
		locationsOut[i] = v.(string)
	}
	return locationsOut
}

//function to expand simple private locations
func expandSyntheticsSimplePrivateLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))
	for i, v := range locations {
		locationsOut[i] = v.(string)
	}
	return locationsOut
}

//function to expand synthetics tags
func expandSyntheticsTags(tags []interface{}) []synthetics.SyntheticsTag {
	out := make([]synthetics.SyntheticsTag, len(tags))
	for i, v := range tags {
		tag := v.(map[string]interface{})
		expanded := synthetics.SyntheticsTag{
			Key:    tag["key"].(string),
			Values: expandSyntheticsTagValues(tag["values"].([]interface{})),
		}
		out[i] = expanded
	}
	return out
}

//function to expand synthetics tag values
func expandSyntheticsTagValues(v []interface{}) []string {
	values := make([]string, len(v))
	for i, value := range v {
		values[i] = value.(string)
	}
	return values
}
