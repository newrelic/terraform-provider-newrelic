package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func dataSourceNewRelicComponent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicComponentRead,

		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicComponentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Components")

	pluginID := d.Get("plugin_id").(int)
	components, err := client.ListComponents(pluginID)
	if err != nil {
		return err
	}

	var component *newrelic.Component
	name := d.Get("name").(string)

	for _, c := range components {
		if c.Name == name {
			component = &c
			break
		}
	}

	if component == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic components", name)
	}

	d.SetId(strconv.Itoa(component.ID))
	d.Set("id", component.ID)
	d.Set("name", component.Name)
	d.Set("health_status", component.HealthStatus)

	return nil
}
