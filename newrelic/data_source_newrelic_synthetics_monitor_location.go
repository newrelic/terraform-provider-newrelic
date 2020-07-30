package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func dataSourceNewRelicSyntheticsMonitorLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicSyntheticsMonitorLocationRead,

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
				Description: "The high security mode for the Synthetics monitor location",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The private setting for the Synthetics monitor location",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description mode for the Synthetics monitor location",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorLocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading Synthetics monitors location")

	label := d.Get("label").(string)
	locations, err := client.Synthetics.GetMonitorLocations()

	if err != nil {
		return err
	}

	var location *synthetics.MonitorLocation
	for _, l := range locations {
		if l.Label == label {
			location = l
			break
		}
	}

	if location == nil {
		return fmt.Errorf("the label '%s' does not match any Synthetics monitor locations", label)
	}

	d.SetId(location.Name)
	d.Set("name", location.Name)
	d.Set("label", location.Label)
	d.Set("high_security_mode", location.HighSecurityMode)
	d.Set("private", location.Private)
	d.Set("description", location.Description)

	return nil
}
