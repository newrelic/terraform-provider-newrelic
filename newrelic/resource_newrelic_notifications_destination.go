package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

var notificationsDestinationTypes = map[notifications.DestinationType][]string{
	"EMAIL":      {},
	"PAGER_DUTY": {},
	"WEBHOOK":    {},
}

func resourceNewRelicNotificationDestination() *schema.Resource {
	validNotificationDestinationTypes := make([]string, 0, len(notificationsDestinationTypes))
	for k := range notificationsDestinationTypes {
		validNotificationDestinationTypes = append(validNotificationDestinationTypes, string(k))
	}

	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationDestinationCreate,
		ReadContext:   resourceNewRelicNotificationDestinationRead,
		DeleteContext: resourceNewRelicNotificationDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The name of the destination.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validNotificationDestinationTypes, false),
				Description:  fmt.Sprintf("(Required) The type of the destination. One of: (%s).", strings.Join(validNotificationDestinationTypes, ", ")),
			},
			"properties": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of notification destination property types.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification destination property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification destination property value.",
						},
						"label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification destination property label.",
						},
						"display_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification destination property display key.",
						},
					},
				},
			},
			"auth": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "A set of key-value pairs to represent a Notification destination auth.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Sensitive:   true,
			},
		},
	}
}

func resourceNewRelicNotificationDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestinationInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic notification destination %s", destinationInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destination, err := client.Notifications.CreateDestinationMutationWithContext(updatedContext, accountID, *destinationInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(destination.ID))

	return resourceNewRelicNotificationDestinationRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic notification destination %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destination, err := client.Notifications.GetDestinationWithContext(updatedContext, accountID, notifications.UUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenNotificationDestination(destination, d))
}

func resourceNewRelicNotificationDestinationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic notification destination %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	if _, err := client.Notifications.DeleteDestinationMutationWithContext(updatedContext, accountID, notifications.UUID(d.Id())); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
