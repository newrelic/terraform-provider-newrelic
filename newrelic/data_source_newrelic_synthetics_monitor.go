package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func dataSourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicSyntheticsMonitorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the synthetics monitor in New Relic.",
			},
			"monitor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the synthetics monitor.",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic synthetics monitors")

	name := d.Get("name").(string)
	monitors, err := client.Synthetics.ListMonitors()
	if err != nil {
		return diag.FromErr(err)
	}

	var monitor *synthetics.Monitor
	for _, m := range monitors {
		if m.Name == name {
			monitor = m
			break
		}
	}

	if monitor == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic monitors", name))
	}

	d.SetId(monitor.ID)
	_ = d.Set("name", monitor.Name)
	_ = d.Set("monitor_id", monitor.ID)

	return nil
}
