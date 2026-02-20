package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

func dataSourceNewRelicAlertPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAlertPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert policy in New Relic.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID to operate on.",
			},
			"incident_preference": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rollup strategy for the policy, which can be `PER_POLICY`, `PER_CONDITION`, or `PER_CONDITION_AND_TARGET`.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the policy was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the policy was last updated.",
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the alert policy in New Relic.",
			},
		},
	}
}

func dataSourceNewRelicAlertPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*ProviderConfig)
	client := cfg.NewClient
	log.Printf("[INFO] Reading New Relic Alert Policies")

	name := d.Get("name").(string)
	accountID := selectAccountID(cfg, d)

	params := alerts.AlertsPoliciesSearchCriteriaInput{}

	policies, err := client.Alerts.QueryPolicySearchWithContext(ctx, accountID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy *alerts.AlertsPolicy

	for _, c := range policies {
		if strings.EqualFold(c.Name, name) {
			policy = c

			break
		}
	}

	if policy == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic alert policy", name))
	}

	d.SetId(policy.ID)

	return diag.FromErr(flattenAlertPolicyWithEntityGUID(ctx, client, policy, d, accountID))
}
