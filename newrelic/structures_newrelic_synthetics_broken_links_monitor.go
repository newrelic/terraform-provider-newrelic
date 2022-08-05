package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func buildSyntheticsBrokenLinksMonitorCreateInput(d *schema.ResourceData) *synthetics.SyntheticsCreateBrokenLinksMonitorInput {
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
	return &input
}

func buildSyntheticsBrokenLinksMonitorUpdateInput(d *schema.ResourceData) *synthetics.SyntheticsUpdateBrokenLinksMonitorInput {
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

	return &input
}
