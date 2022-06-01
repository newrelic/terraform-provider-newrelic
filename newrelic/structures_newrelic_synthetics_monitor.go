package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type SyntheticsMonitorBase struct {
	Locations synthetics.SyntheticsLocationsInput
	Name      string
	Period    synthetics.SyntheticsMonitorPeriod
	Status    synthetics.SyntheticsMonitorStatus
	Tags      []synthetics.SyntheticsTag
	Uri       string
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
		inputBase.Uri = uri.(string)
	}

	if tags, ok := d.GetOk("tags"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	locations := d.Get("locations").(*schema.Set).List()
	inputBase.Locations = expandSyntheticsLocations(locations)

	return inputBase
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

func getSyntheticsMonitorPeriodTypesAsStrings() []string {
	return []string{
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE),
	}
}
