package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"auth_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringInSlice([]string{"BASIC"}, false),
						},
						"auth_username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"base_url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"headers": {
							Type:      schema.TypeMap,
							Elem:      &schema.Schema{Type: schema.TypeString},
							Optional:  true,
							Sensitive: true,
						},
						"key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"include_json_attachment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"payload": {
							Type:      schema.TypeMap,
							Elem:      &schema.Schema{Type: schema.TypeString},
							Sensitive: true,
							Optional:  true,
						},
						"payload_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"recipients": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
						},
						"route_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"service_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"teams": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNewRelicAlertChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Alert Channels")

	channels, err := client.Alerts.ListChannels()
	if err != nil {
		return err
	}

	var channel *alerts.Channel
	name := d.Get("name").(string)

	for _, c := range channels {
		if strings.EqualFold(c.Name, name) {
			channel = c
			break
		}
	}

	if channel == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert channel", name)
	}

	d.SetId(strconv.Itoa(channel.ID))
	d.Set("policy_ids", channel.Links.PolicyIDs)

	err = flattenAlertChannel(channel, d)

	if err != nil {
		return err
	}

	return nil
}
