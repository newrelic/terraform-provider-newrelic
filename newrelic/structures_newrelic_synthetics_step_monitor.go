package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsStepMonitorCreateInput(d *schema.ResourceData) (*synthetics.SyntheticsCreateStepMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateStepMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Steps:  expandSyntheticsMonitorSteps(d.Get("steps").([]interface{})),
	}

	if attr, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandPrivateLocations(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		v := attr.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &v
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	input.Browsers = &typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	input.Devices = &typedDevices

	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if runtimeTypeOk || runtimeTypeVersionOk {
		if !(runtimeTypeOk && runtimeTypeVersionOk) {
			return &input, fmt.Errorf("both `runtime_type` and `runtime_type_version` are to be specified")
		}
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
			RuntimeType:        runtimeType.(string),
			RuntimeTypeVersion: synthetics.SemVer(runtimeTypeVersion.(string)),
		}
		input.Runtime = &r

	} else {
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{}
		input.Runtime = &r
	}

	return &input, nil
}

func buildSyntheticsStepMonitorUpdateInput(d *schema.ResourceData) (*synthetics.SyntheticsUpdateStepMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateStepMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
		Steps:  expandSyntheticsMonitorSteps(d.Get("steps").([]interface{})),
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandPrivateLocations(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		v := attr.(bool)
		input.AdvancedOptions.EnableScreenshotOnFailureAndScript = &v
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	input.Browsers = &typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	input.Devices = &typedDevices

	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if runtimeTypeOk || runtimeTypeVersionOk {
		if !(runtimeTypeOk && runtimeTypeVersionOk) {
			return &input, fmt.Errorf("both `runtime_type` and `runtime_type_version` are to be specified")
		}
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
			RuntimeType:        runtimeType.(string),
			RuntimeTypeVersion: synthetics.SemVer(runtimeTypeVersion.(string)),
		}
		input.Runtime = &r

	} else {
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{}
		input.Runtime = &r
	}

	return &input, nil
}

func expandSyntheticsMonitorSteps(steps []interface{}) []synthetics.SyntheticsStepInput {
	stepsOut := []synthetics.SyntheticsStepInput{}

	for _, s := range steps {
		st := s.(map[string]interface{})

		stepsOut = append(stepsOut, synthetics.SyntheticsStepInput{
			Ordinal: st["ordinal"].(int),
			Type:    synthetics.SyntheticsStepType(st["type"].(string)),
			Values:  expandStringSlice(st["values"].([]interface{})),
		})
	}

	return stepsOut
}

func expandPrivateLocations(locations []interface{}) []synthetics.SyntheticsPrivateLocationInput {
	pl := []synthetics.SyntheticsPrivateLocationInput{}

	for _, v := range locations {
		loc := v.(map[string]interface{})
		pl = append(pl, synthetics.SyntheticsPrivateLocationInput{
			GUID:        loc["guid"].(string),
			VsePassword: synthetics.SecureValue(loc["vse_password"].(string)),
		})
	}

	return pl
}

func flattenSyntheticsMonitorSteps(stepsIn []synthetics.SyntheticsStep) []map[string]interface{} {
	steps := []map[string]interface{}{}

	for _, s := range stepsIn {
		step := map[string]interface{}{
			"ordinal": s.Ordinal,
			"type":    string(s.Type),
			"values":  s.Values,
		}

		steps = append(steps, step)
	}

	return steps
}
