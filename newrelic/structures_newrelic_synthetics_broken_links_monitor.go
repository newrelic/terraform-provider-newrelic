package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsBrokenLinksMonitorCreateInput(d *schema.ResourceData) (*synthetics.SyntheticsCreateBrokenLinksMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateBrokenLinksMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(attr.(*schema.Set).List())
	}
	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}
	if v, ok := d.GetOk("uri"); ok {
		input.Uri = v.(string)
	}

	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if runtimeTypeOk || runtimeTypeVersionOk {
		if !(runtimeTypeOk && runtimeTypeVersionOk) {
			return &input, fmt.Errorf("both `runtime_type` and `runtime_type_version` are to be specified")
		} else {
			r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
				RuntimeType:        runtimeType.(string),
				RuntimeTypeVersion: synthetics.SemVer(runtimeTypeVersion.(string)),
			}
			input.Runtime = &r
		}

	} else {
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{}
		input.Runtime = &r
	}

	return &input, nil
}

func buildSyntheticsBrokenLinksMonitorUpdateInput(d *schema.ResourceData) (*synthetics.SyntheticsUpdateBrokenLinksMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateBrokenLinksMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(attr.(*schema.Set).List())
	}
	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.(*schema.Set).List())
	}
	if v, ok := d.GetOk("uri"); ok {
		input.Uri = v.(string)
	}
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if runtimeTypeOk || runtimeTypeVersionOk {
		if !(runtimeTypeOk && runtimeTypeVersionOk) {
			return &input, fmt.Errorf("both `runtime_type` and `runtime_type_version` are to be specified")
		} else {
			r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
				RuntimeType:        runtimeType.(string),
				RuntimeTypeVersion: synthetics.SemVer(runtimeTypeVersion.(string)),
			}
			input.Runtime = &r
		}

	} else {
		r := synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{}
		input.Runtime = &r
	}

	return &input, nil
}
