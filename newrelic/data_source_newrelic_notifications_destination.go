package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicNotificationDestinationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the destination.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account ID under which to put the destination.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the destination.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: fmt.Sprintf("The type of the destination. One of: (%s).", strings.Join(listValidNotificationsDestinationTypes(), ", ")),
			},
			"property": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Notification destination property type.",
				Elem:        notificationsPropertySchema(),
			},
			"auth_basic": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"auth_token"},
				Description:   "Basic username and password authentication credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"auth_token": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"auth_basic"},
				Description:   "Token authentication credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the destination is active.",
				Default:     true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the destination.",
			},
			"last_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time a notification was sent.",
			},
		},
	}
}

func dataSourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Notification Destination")

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)
	id := d.Id()
	filters := ai.AiNotificationsDestinationFilter{ID: id}
	sorter := notifications.AiNotificationsDestinationSorter{}

	destinationResponse, err := client.Notifications.GetDestinationsWithContext(updatedContext, accountID, "", filters, sorter)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if len(destinationResponse.Entities) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the id '%s' does not match any New Relic notification destination", id))
	}

	errors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return diag.FromErr(flattenNotificationDestination(&destinationResponse.Entities[0], d))
}
