package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	synthetics "github.com/newrelic/newrelic-client-go/pkg/synthetics"
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

	name := d.Get("name").(string)
	monitors, err := client.ListMonitors()

	if err != nil {
		return err
	}

	var monitor *synthetics.Monitor
	for _, m := range monitors {
		if m.Name == name {
			monitor = &m
			break
		}
	}

	if monitor == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic monitors", name)
	}

	d.SetId(monitor.ID)
	d.Set("name", monitor.Name)
	d.Set("monitor_id", monitor.ID)

	return nil
}
