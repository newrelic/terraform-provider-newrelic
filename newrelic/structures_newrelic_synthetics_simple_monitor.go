package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func buildSyntheticsSimpleMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleMonitorInput := synthetics.SyntheticsCreateSimpleMonitorInput{}

	simpleMonitorInput.Name = inputBase.Name
	simpleMonitorInput.Period = inputBase.Period
	simpleMonitorInput.Status = inputBase.Status
	simpleMonitorInput.Tags = inputBase.Tags

	if v, ok := d.GetOk("custom_header"); ok {
		simpleMonitorInput.AdvancedOptions.CustomHeaders = expandSyntheticsCustomHeaders(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("uri"); ok {
		simpleMonitorInput.Uri = v.(string)
	}

	if v, ok := d.GetOk("locations_public"); ok {
		simpleMonitorInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleMonitorInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
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

func buildSyntheticsSimpleMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleMonitorInput{}

	simpleMonitorUpdateInput.Name = inputBase.Name
	simpleMonitorUpdateInput.Period = inputBase.Period
	simpleMonitorUpdateInput.Status = inputBase.Status
	simpleMonitorUpdateInput.Tags = inputBase.Tags

	if v, ok := d.GetOk("locations_public"); ok {
		simpleMonitorUpdateInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleMonitorUpdateInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
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
