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
		string(SyntheticsMonitorTypes.BROWSER),
		string(SyntheticsMonitorTypes.SIMPLE),
		string(SyntheticsMonitorTypes.SCRIPT_API),
		string(SyntheticsMonitorTypes.SCRIPT_BROWSER),
	}
}

//func to build the input to create simple browser monitor
func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{}

	simpleBrowserMonitorInput.AdvancedOptions.CustomHeaders = inputBase.CustomHeaders
	simpleBrowserMonitorInput.Name = inputBase.Name
	simpleBrowserMonitorInput.Period = inputBase.Period
	simpleBrowserMonitorInput.Status = inputBase.Status
	simpleBrowserMonitorInput.Tags = inputBase.Tags
	simpleBrowserMonitorInput.Uri = inputBase.URI

	if v := d.Get("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}
	if v, ok := d.GetOk("location_public"); ok {
		simpleBrowserMonitorInput.Locations.Public = expandSyntheticsSimplePublicLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		simpleBrowserMonitorInput.Locations.Private = expandSyntheticsSimplePrivateLocations(v.(*schema.Set).List())
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

//func to build input to create simple monitor
func buildSyntheticsSimpleMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleMonitorInput {

	inputBase := expandSyntheticsMonitorBase(d)

	simpleMonitorInput := synthetics.SyntheticsCreateSimpleMonitorInput{}

	simpleMonitorInput.AdvancedOptions.CustomHeaders = inputBase.CustomHeaders
	simpleMonitorInput.Name = inputBase.Name
	simpleMonitorInput.Period = inputBase.Period
	simpleMonitorInput.Status = inputBase.Status
	simpleMonitorInput.Tags = inputBase.Tags
	simpleMonitorInput.Uri = inputBase.URI

	if v, ok := d.GetOk("location_public"); ok {
		simpleMonitorInput.Locations.Public = expandSyntheticsSimplePublicLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		simpleMonitorInput.Locations.Private = expandSyntheticsSimplePrivateLocations(v.(*schema.Set).List())
	}
	if v := d.Get("treat_redirect_as_failure"); v != nil {
		t := v.(bool)
		simpleMonitorInput.AdvancedOptions.RedirectIsFailure = &t
	}
	if v, ok := d.GetOk("validation_string"); ok {
		simpleMonitorInput.AdvancedOptions.ResponseValidationText = v.(string)
	}
	if v := d.Get("bypass_head_request"); v != nil {
		b := v.(bool)
		simpleMonitorInput.AdvancedOptions.ShouldBypassHeadRequest = &b
	}
	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleMonitorInput.AdvancedOptions.UseTlsValidation = &vs
	}
	return simpleMonitorInput
}

//func to build input to update simple browser monitor
func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleBrowserMonitorInput {
	simpleBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{}

	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput.AdvancedOptions.CustomHeaders = inputBase.CustomHeaders
	simpleBrowserMonitorUpdateInput.Name = inputBase.Name
	simpleBrowserMonitorUpdateInput.Period = inputBase.Period
	simpleBrowserMonitorUpdateInput.Status = inputBase.Status
	simpleBrowserMonitorUpdateInput.Tags = inputBase.Tags
	simpleBrowserMonitorUpdateInput.Uri = inputBase.URI

	if v, ok := d.GetOk("location_public"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Public = expandSyntheticsSimplePublicLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Private = expandSyntheticsSimplePrivateLocations(v.(*schema.Set).List())
	}
	if v, _ := d.GetOk("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}
	if v, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.ResponseValidationText = v.(string)
	}
	if v, _ := d.GetOk("verify_ssl"); v != nil {
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

//func to build input to update simple monitor
func buildSyntheticsSimpleMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleMonitorInput {

	simpleMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleMonitorInput{}

	inputBase := expandSyntheticsMonitorBase(d)

	simpleMonitorUpdateInput.AdvancedOptions.CustomHeaders = inputBase.CustomHeaders
	simpleMonitorUpdateInput.Name = inputBase.Name
	simpleMonitorUpdateInput.Period = inputBase.Period
	simpleMonitorUpdateInput.Status = inputBase.Status
	simpleMonitorUpdateInput.Tags = inputBase.Tags
	simpleMonitorUpdateInput.Uri = inputBase.URI

	if v, ok := d.GetOk("location_public"); ok {
		simpleMonitorUpdateInput.Locations.Public = expandSyntheticsSimplePublicLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		simpleMonitorUpdateInput.Locations.Private = expandSyntheticsSimplePrivateLocations(v.(*schema.Set).List())
	}
	if v := d.Get("treat_redirect_as_failure"); v != nil {
		i := v.(bool)
		simpleMonitorUpdateInput.AdvancedOptions.RedirectIsFailure = &i
	}
	if v, ok := d.GetOk("validation_string"); ok {
		simpleMonitorUpdateInput.AdvancedOptions.ResponseValidationText = v.(string)
	}
	if v := d.Get("bypass_head_request"); v != nil {
		b := v.(bool)
		simpleMonitorUpdateInput.AdvancedOptions.ShouldBypassHeadRequest = &b
	}
	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleMonitorUpdateInput.AdvancedOptions.UseTlsValidation = &vs
	}
	return simpleMonitorUpdateInput
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

	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(attr.(*schema.Set).List())
	}
	if attr, ok := d.GetOk("location_public"); ok {
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

	if v := d.Get("enable_screenshot_on_failure_and_script"); v.(bool) {
		e := v.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}
	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(attr.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_public"); ok {
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

	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsPrivateLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_public"); ok {
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
	if v := d.Get("enable_screenshot_on_failure_and_script"); v.(bool) {
		e := v.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}
	if v, ok := d.GetOk("location_public"); ok {
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
	if headers, ok := d.GetOk("custom_header"); ok {
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
