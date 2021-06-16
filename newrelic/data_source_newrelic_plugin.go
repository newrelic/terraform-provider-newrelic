package newrelic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicPlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicPluginRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the plugin in New Relic.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the installed plugin instance.",
			},
		},
	}
}

func dataSourceNewRelicPluginRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("plugins have reached end of life, use infrastructure integrations instead")
}
