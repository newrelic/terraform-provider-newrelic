package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

var notificationsChannelTypes = map[notifications.AiNotificationsChannelType][]string{
	"EMAIL":                         {},
	"SERVICE_NOW":                   {},
	"PAGERDUTY_ACCOUNT_INTEGRATION": {},
	"PAGERDUTY_SERVICE_INTEGRATION": {},
	"WEBHOOK":                       {},
}

var notificationsChannelProductTypes = map[notifications.AiNotificationsProduct][]string{
	"ALERTS":         {},
	"DISCUSSIONS":    {},
	"ERROR_TRACKING": {},
	"IINT":           {},
	"NTFC":           {},
	"PD":             {},
	"SHARING":        {},
}

func resourceNewRelicNotificationChannel() *schema.Resource {
	validNotificationChannelTypes := make([]string, 0, len(notificationsChannelTypes))
	for k := range notificationsChannelTypes {
		validNotificationChannelTypes = append(validNotificationChannelTypes, string(k))
	}

	validNotificationChannelProductTypes := make([]string, 0, len(notificationsChannelTypes))
	for k := range notificationsChannelTypes {
		validNotificationChannelProductTypes = append(validNotificationChannelProductTypes, string(k))
	}

	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationChannelCreate,
		ReadContext:   resourceNewRelicNotificationChannelRead,
		DeleteContext: resourceNewRelicNotificationChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The name of the channel.",
			},
			"destinationId": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The id of the destination.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validNotificationChannelTypes, false),
				Description:  fmt.Sprintf("(Required) The type of the channel. One of: (%s).", strings.Join(validNotificationChannelTypes, ", ")),
			},
			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validNotificationChannelProductTypes, false),
				Description:  fmt.Sprintf("(Required) The type of the channel product. One of: (%s).", strings.Join(validNotificationChannelProductTypes, ", ")),
			},
			"properties": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of notification channel property types.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification channel property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification channel property value.",
						},
						"label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification channel property label.",
						},
						"display_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification channel property display key.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicNotificationChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	channelInput, err := expandNotificationChannel(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic notification channelResponse %s", channelInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	channelResponse, err := client.Notifications.AiNotificationsCreateChannelWithContext(updatedContext, accountID, *channelInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(channelResponse.Channel.ID))

	return resourceNewRelicNotificationChannelRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic notification channelResponse %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	filters := ai.AiNotificationsChannelFilter{ID: d.Id()}
	sorter := notifications.AiNotificationsChannelSorter{}
	updatedContext := updateContextWithAccountID(ctx, accountID)

	channelResponse, err := client.Notifications.GetChannelsWithContext(updatedContext, accountID, "", filters, sorter)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenNotificationChannel(&channelResponse.Entities[0], d))
}

func resourceNewRelicNotificationChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic notification channel %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	if _, err := client.Notifications.AiNotificationsDeleteChannelWithContext(updatedContext, accountID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
