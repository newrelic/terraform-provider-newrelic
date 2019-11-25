package newrelic

import (
	"fmt"
	"log"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicSyntheticsMonitorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Reading New Relic synthetics monitors")

	offset := 0
	max := 100
	monitors, err := client.GetAllMonitors(uint(offset), uint(max))
	var monitor *synthetics.ExtendedMonitor
	name := d.Get("name").(string)
	for monitors != nil {
		if len(monitors.Monitors) > 0 && err == nil {
			mon := *monitors

			for i := 0; i < len(monitors.Monitors); i++ {
				if mon.Monitors[i].Name == name {
					monitor = mon.Monitors[i]
					break
				}
			}
		}

		offset = offset + 100
		monitors, err = client.GetAllMonitors(uint(offset), uint(max))

		if len(monitors.Monitors) == 0 {
			monitors = nil
		}
	}
	if err != nil {
		return err
	}

	if monitor == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic monitors", name)
	}

	d.SetId(monitor.ID)
	d.Set("name", monitor.Name)
	d.Set("monitor_id", monitor.ID)

	return nil
}
