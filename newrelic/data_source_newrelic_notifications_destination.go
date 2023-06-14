package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicNotificationDestinationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The ID of the destination.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the destination.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account ID under which to put the destination.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("The type of the destination. One of: (%s).", strings.Join(listValidNotificationsDestinationTypes(), ", ")),
			},
			"property": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Notification destination property type.",
				Elem:        notificationsPropertySchema(),
			},
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the destination is active.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the destination.",
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
	var filters ai.AiNotificationsDestinationFilter
	sorter := notifications.AiNotificationsDestinationSorter{}

	nameValue, nameOk := d.Get("name").(string)
	if nameOk && nameValue != "" {
		filters = ai.AiNotificationsDestinationFilter{Name: nameValue}
	}
	idValue, idOk := d.Get("id").(string)
	if idOk && idValue != "" {
		filters = ai.AiNotificationsDestinationFilter{ID: idValue}
	}

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
		return diag.FromErr(fmt.Errorf("the id or name you provided does not match any New Relic notification destination"))
	}

	errors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return diag.FromErr(flattenNotificationDestinationDataSource(&destinationResponse.Entities[0], d))
}
