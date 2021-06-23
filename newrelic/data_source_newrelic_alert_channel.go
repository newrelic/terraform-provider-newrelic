package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func dataSourceNewRelicAlertChannel() *schema.Resource {
	validAlertChannelTypes := make([]string, 0, len(alertChannelTypes))
	for k := range alertChannelTypes {
		validAlertChannelTypes = append(validAlertChannelTypes, k)
	}

	return &schema.Resource{
		ReadContext: dataSourceNewRelicAlertChannelRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert channel in New Relic.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("The type of channel. One of: (%s).", strings.Join(validAlertChannelTypes, ", ")),
			},
			"policy_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "A list of policy IDs associated with the alert channel.",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alert channel configuration.",
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

func dataSourceNewRelicAlertChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Alert Channels")

	channels, err := client.Alerts.ListChannelsWithContext(ctx)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic alert channel", name))
	}

	return diag.FromErr(flattenAlertChannelDataSource(channel, d))
}
