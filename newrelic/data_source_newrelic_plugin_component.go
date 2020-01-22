package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
)

func dataSourceNewRelicPluginComponent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicPluginComponentRead,

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

func dataSourceNewRelicPluginComponentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Plugin Components")

	pluginID := d.Get("plugin_id").(int)
	name := d.Get("name").(string)

	params := plugins.ListComponentsParams{
		PluginID: pluginID,
		Name:     name,
	}

	components, err := client.Plugins.ListComponents(&params)
	if err != nil {
		return err
	}

	var component *plugins.Component

	for _, c := range components {
		if c.Name == name {
			component = c
			break
		}
	}

	if component == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic plugin components", name)
	}

	flattenPluginsComponent(component, d)

	return nil
}

func flattenPluginsComponent(component *plugins.Component, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(component.ID))
	d.Set("id", component.ID)
	d.Set("name", component.Name)
	d.Set("health_status", component.HealthStatus)
}
