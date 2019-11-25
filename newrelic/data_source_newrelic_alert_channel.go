package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
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
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Alert Channels")

	channels, err := client.ListAlertChannels()
	if err != nil {
		return err
	}

	var channel *newrelic.AlertChannel
	name := d.Get("name").(string)

	for _, c := range channels {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(name)) {
			channel = &c
			break
		}
	}

	if channel == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert channel", name)
	}

	d.SetId(strconv.Itoa(channel.ID))
	d.Set("name", channel.Name)
	d.Set("type", channel.Type)
	d.Set("policy_ids", channel.Links.PolicyIDs)
	d.Set("configuration", channel.Configuration)

	return nil
}
