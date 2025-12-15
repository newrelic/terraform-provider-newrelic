package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

func dataSourceNewRelicAlertCompoundCondition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAlertCompoundConditionRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID to operate on.",
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the compound alert condition.",
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the policy associated with this compound alert condition.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the compound alert condition.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the condition is enabled.",
			},
			"trigger_expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The logical expression defining when the alert compound condition triggers.",
			},
			"component_conditions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The component conditions that make up this alert compound condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the component NRQL condition.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias used in the trigger expression.",
						},
					},
				},
			},
			"facet_matching_behavior": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "How facets are handled when evaluating component conditions.",
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Runbook URL for the condition.",
			},
			"threshold_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration the trigger expression must be true before opening an incident.",
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the alert compound condition.",
			},
		},
	}
}

func dataSourceNewRelicAlertCompoundConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*ProviderConfig)
	client := cfg.NewClient
	accountID := selectAccountID(cfg, d)

	log.Printf("[INFO] Reading New Relic compound alert condition")

	conditionID := d.Get("id").(string)

	// Use efficient API filtering by ID - expects just the condition ID
	filter := &alerts.AlertsCompoundConditionFilterInput{
		Id: &alerts.AlertsCompoundConditionIDFilter{
			Eq: &conditionID,
		},
	}

	conditions, err := client.Alerts.SearchCompoundConditionsWithContext(ctx, accountID, filter, nil, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(conditions) == 0 {
		return diag.FromErr(fmt.Errorf("compound alert condition with ID '%s' not found", conditionID))
	}

	matchedCondition := conditions[0]
	d.SetId(conditionID)

	return diag.FromErr(flattenAlertCompoundCondition(accountID, matchedCondition, d))
}
