package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
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
				Optional: true,
			},
		},
	}
}

func dataSourceNewRelicPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Plugins")

	plugins, err := client.ListPlugins()
	if err != nil {
		return err
	}

	var plugin *newrelic.Plugin
	guid := d.Get("guid").(string)

	for _, p := range plugins {
		if p.GUID == guid {
			plugin = &p
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
