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

func resourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationDestinationCreate,
		ReadContext:   resourceNewRelicNotificationDestinationRead,
		UpdateContext: resourceNewRelicNotificationDestinationUpdate,
		DeleteContext: resourceNewRelicNotificationDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The account ID under which to put the destination.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "(Required) The name of the destination.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidNotificationsDestinationTypes(), false),
				Description:  fmt.Sprintf("(Required) The type of the destination. One of: (%s).", strings.Join(listValidNotificationsDestinationTypes(), ", ")),
			},
			"property": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Notification destination property type.",
				Elem:        notificationsPropertySchema(),
			},
			"auth_basic": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				ExactlyOneOf: []string{"auth_basic", "auth_token"},
				Description:  "Basic username and password authentication credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"auth_token": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				ExactlyOneOf: []string{"auth_basic", "auth_token"},
				Description:  "Token authentication credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"token": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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

			// Computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the destination.",
			},
			"is_user_authenticated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the user is authenticated with the destination.",
			},
			"last_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time a notification was sent.",
			},
		},
	}
}

func resourceNewRelicNotificationDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestination(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic notification destinationResponse %s", destinationInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.AiNotificationsCreateDestinationWithContext(updatedContext, accountID, *destinationInput)
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildAiNotificationsErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	d.SetId(destinationResponse.Destination.ID)

	return resourceNewRelicNotificationDestinationRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic notification destinationResponse %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	filters := ai.AiNotificationsDestinationFilter{ID: d.Id()}
	sorter := notifications.AiNotificationsDestinationSorter{}
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.GetDestinationsWithContext(updatedContext, accountID, "", filters, sorter)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if len(destinationResponse.Entities) == 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] notification destinationResponse.Entities response is empty"))
	}

	errors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return diag.FromErr(flattenNotificationDestination(&destinationResponse.Entities[0], d))
}

func resourceNewRelicNotificationDestinationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestinationUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.AiNotificationsUpdateDestinationWithContext(updatedContext, accountID, *destinationInput, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildAiNotificationsErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return resourceNewRelicNotificationDestinationRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationDestinationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic notification destination %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	destinationResponse, err := client.Notifications.AiNotificationsDeleteDestinationWithContext(updatedContext, accountID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Validation function to validate allowed destination types
func listValidNotificationsDestinationTypes() []string {
	return []string{
		string(notifications.AiNotificationsDestinationTypeTypes.WEBHOOK),
		string(notifications.AiNotificationsDestinationTypeTypes.EMAIL),
		string(notifications.AiNotificationsDestinationTypeTypes.SERVICE_NOW),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION),
		string(notifications.AiNotificationsDestinationTypeTypes.JIRA),
	}
}
