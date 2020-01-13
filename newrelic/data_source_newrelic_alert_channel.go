package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func dataSourceNewRelicAlertChannel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicAlertChannelRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicAlertChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Alert Channels")

	channels, err := client.Alerts.ListAlertChannels()
	if err != nil {
		return err
	}

	var channel *alerts.AlertChannel
	name := d.Get("name").(string)

	for _, c := range channels {
		if strings.EqualFold(c.Name, name) {
			channel = &c
			break
		}
	}

	if channel == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert channel", name)
	}

	d.SetId(strconv.Itoa(channel.ID))
	d.Set("policy_ids", channel.Links.PolicyIDs)

	flattenAlertChannel(channel, d)

	return nil
}
