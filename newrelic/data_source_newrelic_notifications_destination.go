package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
)

func dataSourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicNotificationDestinationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name", "exact_name"},
				Description:  "The ID of the destination.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination entity GUID",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name", "exact_name"},
				Description:  "The name of the destination. Uses a contains match.",
			},
			"exact_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name", "exact_name"},
				Description:  "The exact name of the destination. Uses an exact match.",
			},
			"account_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"scope"},
				Description:   "The account ID under which the particular destination belong to.",
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
			"secure_url": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Description: "URL in secure format",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"scope": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"account_id"},
				Description:   "Scope of the destination",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(listValidNotificationsScopeTypes(), false),
							Description:  fmt.Sprintf("(Required) The scope type of the destination. One of: (%s).", strings.Join(listValidNotificationsScopeTypes(), ", ")),
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the Scope (Organization UUID for ORGANIZATION scope, Account ID for ACCOUNT scope)",
						},
					},
				},
			},
		},
	}
}

func dataSourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Notification Destination")

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	scope := buildEntityScopeInput(d, accountID)
	updatedContext := updateContextWithAccountID(ctx, accountID)
	var filters ai.AiNotificationsDestinationFilter
	sorter := notifications.AiNotificationsDestinationSorter{}

	nameValue, nameOk := d.Get("name").(string)
	if nameOk && nameValue != "" {
		filters = ai.AiNotificationsDestinationFilter{Name: nameValue}
	}
	exactNameValue, exactNameOk := d.Get("exact_name").(string)
	if exactNameOk && exactNameValue != "" {
		filters = ai.AiNotificationsDestinationFilter{ExactName: exactNameValue}
	}
	idValue, idOk := d.Get("id").(string)
	if idOk && idValue != "" {
		filters = ai.AiNotificationsDestinationFilter{ID: idValue}
	}

	return getDestinationWithScope(updatedContext, client, &filters, &sorter, idValue, nameValue, exactNameValue, scope, d)
}

// getDestinationWithScope handles retrieval for all scope types (ACCOUNT and ORGANIZATION)
// Calls GetDestinationsAccount when scope is nil or ACCOUNT, GetDestinationsOrganization when ORGANIZATION
func getDestinationWithScope(
	updatedContext context.Context,
	client *nr.NewRelic,
	filters *ai.AiNotificationsDestinationFilter,
	sorter *notifications.AiNotificationsDestinationSorter,
	idValue, nameValue, exactNameValue string,
	scope *notifications.AiNotificationsEntityScopeInput,
	d *schema.ResourceData,
) diag.Diagnostics {
	var destinationResponse *notifications.AiNotificationsDestinationsResponse
	var err error

	if scope.Type == notifications.AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION {
		destinationResponse, err = client.Notifications.GetDestinationsWithContextOrganization(updatedContext, nil, filters, sorter)
	} else {
		scopeID, atoiErr := strconv.Atoi(scope.ID)
		if atoiErr != nil {
			return diag.FromErr(atoiErr)
		}
		destinationResponse, err = client.Notifications.GetDestinationsWithContextAccount(updatedContext, scopeID, nil, filters, sorter)
	}

	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if len(destinationResponse.Entities) == 0 {
		d.SetId("")
		if err := getNotificationDestinationNotFoundError(idValue, nameValue, exactNameValue); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	respErrors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(respErrors) > 0 {
		return respErrors
	}

	return diag.FromErr(flattenNotificationDestinationDataSource(&destinationResponse.Entities[0], *scope, d))
}

// getNotificationDestinationNotFoundError returns an appropriate error message based on which filter attribute was provided
func getNotificationDestinationNotFoundError(idValue, nameValue, exactNameValue string) error {
	filterAttributes := []struct {
		value string
		name  string
	}{
		{idValue, "id"},
		{nameValue, "name"},
		{exactNameValue, "exact_name"},
	}

	for _, attr := range filterAttributes {
		if attr.value != "" {
			return fmt.Errorf("the %s provided does not match any New Relic notification destination", attr.name)
		}
	}
	return nil
}
