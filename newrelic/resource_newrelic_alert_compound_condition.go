package newrelic

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicAlertCompoundCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAlertCompoundConditionCreate,
		ReadContext:   resourceNewRelicAlertCompoundConditionRead,
		UpdateContext: resourceNewRelicAlertCompoundConditionUpdate,
		DeleteContext: resourceNewRelicAlertCompoundConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID for managing your compound alert conditions.",
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy where this condition should be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the compound alert condition.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not to enable the alert condition.",
			},
			"trigger_expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expression that defines how component condition evaluations are combined. Valid operators are 'AND', 'OR', 'NOT'. For more complex expressions, use parentheses. Simple example: 'A AND B'. Complex example: 'A AND (B OR C) AND NOT D'.",
			},
			"component_conditions": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    2,
				Description: "The list of NRQL conditions to be combined. Each component condition must be enabled. Must have at least 2 components.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the existing NRQL alert condition to use as a component.",
						},
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The identifier that will be used in the compound alert condition's trigger_expression (e.g., 'A', 'B', 'C').",
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`),
								"alias must start with a letter and contain only letters, numbers, and underscores",
							),
						},
					},
				},
				// Custom hash to ensure uniqueness by alias
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					// Use alias for uniqueness
					return schema.HashString(m["alias"].(string))
				},
			},
			"facet_matching_behavior": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "FACETS_IGNORED",
				Description:  "How the compound condition will take into account the component conditions' facets during evaluation. Valid values: 'FACETS_IGNORED' (default) - facets are not taken into consideration when determining when the compound alert condition activates; 'FACETS_MATCH' - the compound alert condition will activate only when shared facets have matching values.",
				ValidateFunc: validation.StringInSlice([]string{"FACETS_MATCH", "FACETS_IGNORED"}, false),
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
			},
			"threshold_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The duration, in seconds, that the trigger expression must be true before the compound alert condition will activate. Between 30-1440 seconds.",
				ValidateFunc: validation.All(
					validation.IntBetween(30, 1440),
				),
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the alert compound condition in New Relic.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceNewRelicAlertCompoundConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	policyID := strconv.Itoa(d.Get("policy_id").(int))

	// Build create input
	conditionInput, err := expandAlertCompoundConditionCreateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic compound alert condition %s via NerdGraph API", conditionInput.Name)

	// Call API
	condition, err := client.Alerts.CreateCompoundConditionWithContext(ctx, accountID, policyID, *conditionInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if condition == nil {
		return diag.Errorf("error creating compound alert condition: response was nil")
	}

	// Set ID to just the condition ID
	d.SetId(condition.ID)

	// Flatten response into state
	return diag.FromErr(flattenAlertCompoundCondition(accountID, condition, d))
}

func resourceNewRelicAlertCompoundConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic compound alert condition %s", d.Id())

	conditionID := d.Id()

	// Search for the specific alert compound condition
	filter := &alerts.AlertsCompoundConditionFilterInput{
		Id: &alerts.AlertsCompoundConditionIDFilter{
			Eq: &conditionID,
		},
	}

	conditions, err := client.Alerts.SearchCompoundConditionsWithContext(ctx, accountID, filter, nil, nil)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			log.Printf("[WARN] Compound alert condition %s not found, removing from state", conditionID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if len(conditions) == 0 {
		log.Printf("[WARN] Compound alert condition %s not found, removing from state", conditionID)
		d.SetId("")
		return nil
	}

	condition := conditions[0]

	return diag.FromErr(flattenAlertCompoundCondition(accountID, condition, d))
}

func resourceNewRelicAlertCompoundConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	conditionID := d.Id()

	// Build update input
	conditionInput, err := expandAlertCompoundConditionUpdateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating New Relic compound alert condition %s", conditionID)

	// Call API
	_, err = client.Alerts.UpdateCompoundConditionWithContext(ctx, accountID, conditionID, *conditionInput)
	if err != nil {
		return diag.FromErr(err)
	}

	// Re-read to get updated state
	return resourceNewRelicAlertCompoundConditionRead(ctx, d, meta)
}

func resourceNewRelicAlertCompoundConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	conditionID := d.Id()

	log.Printf("[INFO] Deleting New Relic compound alert condition %s", conditionID)

	_, err := client.Alerts.DeleteCompoundConditionWithContext(ctx, accountID, conditionID)
	if err != nil {
		// Check if condition was already deleted
		if _, ok := err.(*errors.NotFound); ok {
			log.Printf("[WARN] Compound alert condition %s already deleted", conditionID)
			return nil
		}
		// Check for GraphQL errors indicating the condition doesn't exist
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "does not exist") {
			log.Printf("[WARN] Compound alert condition %s not found during delete", conditionID)
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
