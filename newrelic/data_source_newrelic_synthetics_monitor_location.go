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

		Schema: map[string]*schema.Schema{ //required, optional, computed
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The label of the synthetics monitor.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the synthetics monitor in New Relic.",
			},
			"highSecurityMode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorLocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic synthetics monitors location")

	label := d.Get("label").(string)                          //label
	locations, err := client.Synthetics.GetMonitorLocations() //Get all locations

	if err != nil {
		return err
	}

	var location *synthetics.MonitorLocation //Filtering - to find matching label
	for _, l := range locations {
		if l.Name == label {
			location = l
			break
		}
	}

	if location == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic monitors", label)
	}

	d.SetId(location.Name)
	d.Set("name", location.Name)
	d.Set("label", location.Label)
	d.Set("highSecurityMode", location.HighSecurityMode)
	d.Set("private", location.Private)
	d.Set("description", location.Description)

	return nil
}
