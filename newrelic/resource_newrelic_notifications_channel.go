package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicNotificationChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationChannelCreate,
		ReadContext:   resourceNewRelicNotificationChannelRead,
		UpdateContext: resourceNewRelicNotificationChannelUpdate,
		DeleteContext: resourceNewRelicNotificationChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "(Required) The name of the channel.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The id of the destination.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidNotificationsChannelTypes(), false),
				Description:  fmt.Sprintf("(Required) The type of the channel. One of: (%s).", strings.Join(listValidNotificationsChannelTypes(), ", ")),
			},
			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidNotificationsProductTypes(), false),
				Description:  fmt.Sprintf("(Required) The type of the channel product. One of: (%s).", strings.Join(listValidNotificationsProductTypes(), ", ")),
			},

			// Optional
			"property": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Notification channel property type.",
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
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the channel is active.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The account id of the channel.",
			},

			// Computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the channel.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the channel.",
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

	d.SetId(channelResponse.Channel.ID)

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

func resourceNewRelicNotificationChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput, err := expandNotificationChannelUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	channelID := d.Get("id").(string)
	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	_, err = client.Notifications.AiNotificationsUpdateChannelWithContext(updatedContext, accountID, *updateInput, channelID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicNotificationChannelRead(updatedContext, d, meta)
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

// Validation function to validate allowed channel types
func listValidNotificationsChannelTypes() []string {
	return []string{
		string(notifications.AiNotificationsChannelTypeTypes.WEBHOOK),
		string(notifications.AiNotificationsChannelTypeTypes.EMAIL),
		string(notifications.AiNotificationsChannelTypeTypes.SERVICENOW_INCIDENTS),
		string(notifications.AiNotificationsChannelTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION),
		string(notifications.AiNotificationsChannelTypeTypes.PAGERDUTY_SERVICE_INTEGRATION),
	}
}

// Validation function to validate allowed product types
func listValidNotificationsProductTypes() []string {
	return []string{
		string(notifications.AiNotificationsProductTypes.ALERTS),
		string(notifications.AiNotificationsProductTypes.DISCUSSIONS),
		string(notifications.AiNotificationsProductTypes.ERROR_TRACKING),
		string(notifications.AiNotificationsProductTypes.NTFC),
		string(notifications.AiNotificationsProductTypes.SHARING),
		string(notifications.AiNotificationsProductTypes.PD),
		string(notifications.AiNotificationsProductTypes.IINT),
	}
}
