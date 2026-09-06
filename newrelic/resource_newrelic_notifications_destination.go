package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicNotificationDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotificationDestinationCreate,
		ReadContext:   resourceNewRelicNotificationDestinationRead,
		UpdateContext: resourceNewRelicNotificationDestinationUpdate,
		DeleteContext: resourceNewRelicNotificationDestinationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNewRelicNotificationDestinationImport,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"scope"},
				Description:   "The account ID under which to put the destination.",
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
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"auth_token", "auth_custom_header"},
				Description:   "Basic username and password authentication credentials.",
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
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"auth_basic", "auth_custom_header"},
				Description:   "Token authentication credentials.",
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
			"auth_custom_header": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"auth_basic", "auth_token"},
				Description:   "Custom header based authentication",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
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
			"last_sent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time a notification was sent.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination entity GUID",
			},
			"secure_url": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "URL in secure format",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secure_suffix": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"scope": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
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
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNewRelicNotificationDestinationV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateStateNewRelicNotificationDestinationV0toV1,
				Version: 0,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(16 * time.Second),
			Update: schema.DefaultTimeout(16 * time.Second),
		},
	}
}

func resourceNewRelicNotificationDestinationV0() *schema.Resource {
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
				Optional:    true,
				Computed:    true,
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
				Type:          schema.TypeList,
				Optional:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{"auth_token"},
				Description:   "Basic username and password authentication credentials.",
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
				Type:          schema.TypeList,
				Optional:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{"auth_basic"},
				Description:   "Token authentication credentials.",
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
		SchemaVersion: 0,
	}
}

func resourceNewRelicNotificationDestinationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestination(d)
	destinationInput.Properties = append(destinationInput.Properties, createMonitoringProperty())
	if err != nil {
		return diag.FromErr(err)
	}

	if isOAuth2SlackType(destinationInput.Type) {
		return diag.FromErr(fmt.Errorf("a destination with '%s' type cannot be created via terraform", destinationInput.Type))
	}

	log.Printf("[INFO] Creating New Relic notification destinationResponse %s", destinationInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	scope := buildEntityScopeInput(d, accountID)
	destinationResponse, err := client.Notifications.AiNotificationsCreateDestinationWithContext(updatedContext, &accountID, *destinationInput, scope)
	if err != nil {
		diagErr := diag.FromErr(err)
		newDiagErr := diag.Diagnostics{
			diag.Diagnostic{
				Severity: diagErr[0].Severity,
				Summary:  diagErr[0].Summary,
				Detail:   "NOTICE: fields are statically typed. Make sure all fields are of the correct type",
			},
		}
		return newDiagErr
	}

	// Always persist the destination ID to state if the backend created it, even if the
	// response also contains errors. This prevents orphaned resources when the API
	// creates the destination successfully but also returns errors in the response.
	if destinationResponse.Destination.ID != "" {
		d.SetId(destinationResponse.Destination.ID)
	}

	errors := buildAiNotificationsErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return resourceNewRelicNotificationDestinationRead(updatedContext, d, meta)
}

func resourceNewRelicNotificationDestinationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic notification destinationResponse %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	filters := ai.AiNotificationsDestinationFilter{ID: d.Id()}
	updatedContext := updateContextWithAccountID(ctx, accountID)

	scope := buildEntityScopeInput(d, accountID)

	var destinationResponse *notifications.AiNotificationsDestinationsResponse
	var err error
	if scope.Type == notifications.AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION {
		destinationResponse, err = client.Notifications.GetDestinationsWithContextOrganization(updatedContext, nil, &filters, nil)
	} else {
		scopeID, atoiErr := strconv.Atoi(scope.ID)
		if atoiErr != nil {
			return diag.FromErr(atoiErr)
		}
		destinationResponse, err = client.Notifications.GetDestinationsWithContextAccount(updatedContext, scopeID, nil, &filters, nil)
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
		return nil
	}

	errors := buildAiNotificationsResponseErrors(destinationResponse.Errors)
	if len(errors) > 0 {
		return errors
	}

	return diag.FromErr(flattenNotificationDestination(&destinationResponse.Entities[0], d))
}

// resourceNewRelicNotificationDestinationImport supports importing with a composite ID format:
//
//	<destination_id>:<scope_type>:<scope_id> (e.g., "abc-123:ORGANIZATION:org-uuid")
//	<destination_id> (defaults to ACCOUNT scope with provider's account ID)
func resourceNewRelicNotificationDestinationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 3)

	switch len(parts) {
	case 1:
		// Plain ID — backward compatible, defaults to account scope
	case 3:
		scopeType := parts[1]
		scopeID := parts[2]

		if scopeType != "ORGANIZATION" && scopeType != "ACCOUNT" {
			return nil, fmt.Errorf("invalid scope type %q, must be ORGANIZATION or ACCOUNT", scopeType)
		}

		d.SetId(parts[0])
		if err := d.Set("scope", []interface{}{
			map[string]interface{}{
				"type": scopeType,
				"id":   scopeID,
			},
		}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid import ID format: expected <id> or <id>:<scope_type>:<scope_id>, got %q", d.Id())
	}

	return []*schema.ResourceData{d}, nil
}

func resourceNewRelicNotificationDestinationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	destinationInput, err := expandNotificationDestinationUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	destinationType := notifications.AiNotificationsDestinationType(d.Get("type").(string))

	if isOAuth2SlackType(destinationType) {
		return diag.FromErr(fmt.Errorf("a destination with '%s' type cannot be updated via terraform", destinationType))
	}

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	scope := buildEntityScopeInput(d, accountID)

	// Use the regular update method (scope is passed in mutation for org-scoped)
	destinationResponse, err := client.Notifications.AiNotificationsUpdateDestinationWithContext(updatedContext, &accountID, *destinationInput, d.Id(), scope)
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

	scope := buildEntityScopeInput(d, accountID)

	destinationResponse, err := client.Notifications.AiNotificationsDeleteDestinationWithContext(updatedContext, &accountID, d.Id(), scope)
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
		string(notifications.AiNotificationsDestinationTypeTypes.SERVICE_NOW_APP),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION),
		string(notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION),
		string(notifications.AiNotificationsDestinationTypeTypes.JIRA),
		string(notifications.AiNotificationsDestinationTypeTypes.SLACK),
		string(notifications.AiNotificationsDestinationTypeTypes.SLACK_COLLABORATION),
		string(notifications.AiNotificationsDestinationTypeTypes.SLACK_LEGACY),
		string(notifications.AiNotificationsDestinationTypeTypes.MOBILE_PUSH),
		string(notifications.AiNotificationsDestinationTypeTypes.EVENT_BRIDGE),
		string(notifications.AiNotificationsDestinationTypeTypes.MICROSOFT_TEAMS),
		string(notifications.AiNotificationsDestinationTypeTypes.WORKFLOW_AUTOMATION),
	}
}

// Validation function to OAuth2 slack types
func isOAuth2SlackType(destinationType notifications.AiNotificationsDestinationType) bool {
	return destinationType == notifications.AiNotificationsDestinationTypeTypes.SLACK || destinationType == notifications.AiNotificationsDestinationTypeTypes.SLACK_COLLABORATION
}
