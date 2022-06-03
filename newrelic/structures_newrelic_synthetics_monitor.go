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
	Uri           string
	CustomHeaders []synthetics.SyntheticsCustomHeaderInput
}

//1, 5, 10, 15, 30, 60, 360, 720, 1440
func periodConvIntToString(v interface{}) synthetics.SyntheticsMonitorPeriod {
	var output synthetics.SyntheticsMonitorPeriod
	switch v.(int) {
	case 1:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE
	case 5:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES
	case 10:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES
	case 15:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES
	case 30:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES
	case 60:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR
	case 360:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS
	case 720:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS
	case 1440:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY
	}
	return output
}

func expandSyntheticsMonitorBase(d *schema.ResourceData) SyntheticsMonitorBase {
	inputBase := SyntheticsMonitorBase{}

	name := d.Get("name")
	inputBase.Name = name.(string)

	status := d.Get("status")
	inputBase.Status = synthetics.SyntheticsMonitorStatus(status.(string))

	period := d.Get("frequency")
	inputBase.Period = periodConvIntToString(period)

	if uri, ok := d.GetOk("uri"); ok {
		inputBase.Uri = uri.(string)
	}

	if tags, ok := d.GetOk("tags"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	headers := d.Get("custom_headers")
	inputBase.CustomHeaders = expandSyntheticsCustomHeaders(headers.(*schema.Set).List())

	locations := d.Get("locations").(*schema.Set).List()
	inputBase.Locations = expandSyntheticsLocations(locations)

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
