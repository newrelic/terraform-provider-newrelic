package newrelic

import (
	"fmt"
	"log"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicSyntheticsMonitorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_KEY", nil),
				Sensitive:   true,
			},
			"max_check": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {

	apiKey := d.Get("api_key").(string)
	maxCheck := d.Get("max_check").(int)

	conf := func(s *synthetics.Client) {
		s.APIKey = apiKey
	}

	syntheticsClient, _ := synthetics.NewClient(conf)

	log.Printf("[INFO] Reading New Relic synthetics monitors")

	offset := 0
	max := 100
	monitors, err := syntheticsClient.GetAllMonitors(uint(offset), uint(max))
	var monitor *synthetics.ExtendedMonitor
	name := d.Get("name").(string)
	for offset <= maxCheck {
		if len(monitors.Monitors) > 0 && err == nil {
			mon := *monitors

			for i := 0; i < len(monitors.Monitors); i++ {
				if mon.Monitors[i].Name == name {
					monitor = mon.Monitors[i]
					break
				}
			}
		} else {
			break
		}

		offset = offset + 100
		monitors, err = syntheticsClient.GetAllMonitors(uint(offset), uint(max))
	}
	if err != nil {
		return err
	}

	if monitor == nil {
		return fmt.Errorf("The name '%s' does not match any New Relic monitors.", name)
	}

	d.SetId(monitor.ID)
	d.Set("name", monitor.Name)
	d.Set("monitor_id", monitor.ID)

	return nil
}
