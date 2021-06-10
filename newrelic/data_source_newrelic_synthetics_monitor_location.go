package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func dataSourceNewRelicSyntheticsMonitorLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicSyntheticsMonitorLocationRead,
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The label of the Synthetics monitor location.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Synthetics monitor location.",
			},
			"high_security_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Represents if high security mode is enabled for the location. A value of true means that high security mode is enabled, and a value of false means it is disabled.",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Represents if this location is a private location. A value of true means that the location is private, and a value of false means it is public.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the Synthetics monitor location.",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading Synthetics monitor locations")

	label := d.Get("label").(string)
	locations, err := client.Synthetics.GetMonitorLocations()
	if err != nil {
		return diag.FromErr(err)
	}

	var location *synthetics.MonitorLocation
	for _, l := range locations {
		if l.Label == label {
			location = l
			break
		}
	}

	if location == nil {
		return diag.FromErr(fmt.Errorf("the label '%s' does not match any Synthetics monitor locations", label))
	}

	d.SetId(location.Name)
	_ = d.Set("name", location.Name)
	_ = d.Set("label", location.Label)
	_ = d.Set("high_security_mode", location.HighSecurityMode)
	_ = d.Set("private", location.Private)
	_ = d.Set("description", location.Description)

	return nil
}
