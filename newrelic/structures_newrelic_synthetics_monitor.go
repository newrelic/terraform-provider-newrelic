package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type SyntheticsMonitorBase struct {
	Locations     synthetics.SyntheticsLocationsInput
	Name          string
	Period        synthetics.SyntheticsMonitorPeriod
	Status        synthetics.SyntheticsMonitorStatus
	Tags          []synthetics.SyntheticsTag
	URI           string                                   // Move URI outside of base (does not apply to all monitors)
	CustomHeaders []synthetics.SyntheticsCustomHeaderInput // Move CustomHeaders outside of base (does not apply to all monitors)
}

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

func listValidSyntheticsMonitorStatuses() []string {
	return []string{
		string(synthetics.SyntheticsMonitorStatusTypes.DISABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.ENABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.MUTED),
	}
}

func listValidSyntheticsScriptMonitorTypes() []string {
	return []string{
		string(SyntheticsMonitorTypes.SCRIPT_API),
		string(SyntheticsMonitorTypes.SCRIPT_BROWSER),
	}
}

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

func expandSyntheticsPrivateLocations(locations []interface{}) []synthetics.SyntheticsPrivateLocationInput {
	locationsOut := make([]synthetics.SyntheticsPrivateLocationInput, len(locations))

	for i, v := range locations {
		pl := v.(map[string]interface{})
		locationsOut[i].GUID = pl["guid"].(int) // This should be a string I think since it's a scalar `ID`

		// if v, ok := pl["vse_password"]; ok {
		// 	locationsOut[i].VsePassword = synthetics.SecureValue(v.(string))
		// }
	}

	return locationsOut
}

func expandSyntheticsPublicLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))

	for i, v := range locations {
		locationsOut[i] = v.(string)
	}

	return locationsOut
}

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

	if tags, ok := d.GetOk("tags"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	if headers, ok := d.GetOk("custom_headers"); ok {
		inputBase.CustomHeaders = expandSyntheticsCustomHeaders(headers.(*schema.Set).List())
	}

	// Locations can't be in the "base" due to different data structures in different monitor types
	publicLocations := d.Get("locations_public").(*schema.Set).List()
	inputBase.Locations = expandSyntheticsLocations(publicLocations)

	//privateLocations:=d.Get("locations_private").(*schema.Set).List()

	return inputBase
}

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

func expandSyntheticsLocations(locations []interface{}) synthetics.SyntheticsLocationsInput {
	locationsOut := make([]string, len(locations))

	for i, v := range locations {
		locationsOut[i] = v.(string)
	}

	return synthetics.SyntheticsLocationsInput{
		Public: locationsOut,
		// What about private?
	}
}

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

func expandSyntheticsTagValues(v []interface{}) []string {
	values := make([]string, len(v))
	for i, value := range v {
		values[i] = value.(string)
	}
	return values
}

// func flattenSyntheticsScriptMonitor(v *entities.EntityInterface, d *schema.ResourceData) diag.Diagnostics {
// 	switch e := (*v).(type) {
// 	case *entities.SyntheticMonitorEntityOutline:
// 		_ = d.Set("guid", string(e.GUID))
// 		_ = d.Set("name", e.Name)
// 		_ = d.Set("type", string(e.MonitorType))
// 		// _ = d.Set("period", e.Period) // entity.Period is NOT the same as synthetics.Period ugh
// 		// _ = d.Set("period", e.Period)
// 	}

// 	return nil
// }
