package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
)

func dataSourceNewRelicPlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicPluginRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Plugins")

	guid := d.Get("guid").(string)

	params := plugins.ListPluginsParams{
		GUID: guid,
	}

	ps, err := client.Plugins.ListPlugins(&params)
	if err != nil {
		return err
	}

	var plugin *plugins.Plugin

	for _, p := range ps {
		if p.GUID == guid {
			plugin = p
			break
		}
	}

	if plugin == nil {
		return fmt.Errorf("the GUID '%s' does not match any New Relic plugins", guid)
	}

	d.SetId(strconv.Itoa(plugin.ID))
	d.Set("id", plugin.ID)

	return nil
}
