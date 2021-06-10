package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
)

func dataSourceNewRelicPlugin() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "`newrelic_plugin` has been deprecated and will not be supported as of June 16, 2021",
		ReadContext:        dataSourceNewRelicPluginRead,
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

func dataSourceNewRelicPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Plugins")

	guid := d.Get("guid").(string)

	params := plugins.ListPluginsParams{
		GUID: guid,
	}

	ps, err := client.Plugins.ListPlugins(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	var plugin *plugins.Plugin

	for _, p := range ps {
		if p.GUID == guid {
			plugin = p
			break
		}
	}

	if plugin == nil {
		return diag.FromErr(fmt.Errorf("the GUID '%s' does not match any New Relic plugins", guid))
	}

	id := strconv.Itoa(plugin.ID)
	d.SetId(id)
	_ = d.Set("id", id)

	return nil
}
