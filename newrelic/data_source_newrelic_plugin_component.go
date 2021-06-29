package newrelic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicPluginComponent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicPluginComponentRead,

		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the plugin instance this component belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the plugin component.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the plugin component.",
			},
			"health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health status of the plugin component.",
			},
		},
	}
}

func dataSourceNewRelicPluginComponentRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use infrastructure integrations instead")
}
