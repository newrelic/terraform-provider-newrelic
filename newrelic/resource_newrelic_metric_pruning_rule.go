package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pruningrules"
)

func resourceNewRelicMetricPruningRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicMetricPruningRuleCreate,
		ReadContext:   resourceNewRelicMetricPruningRuleRead,
		DeleteContext: resourceNewRelicMetricPruningRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account ID in which the pruning rule is created.",
			},
			"nrql": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The NRQL query that identifies the metric attributes to prune. Must select specific attributes from Metric (e.g. `SELECT collector.name FROM Metric WHERE metricName = 'my.metric'`).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A human-readable description of the pruning rule.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the pruning rule assigned by New Relic.",
			},
		},
	}
}

func resourceNewRelicMetricPruningRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	input := []pruningrules.NRQLDropRulesCreateDropRuleInput{
		{
			Action:      pruningrules.NRQLDropRulesActionTypes.DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES,
			NRQL:        d.Get("nrql").(string),
			Description: d.Get("description").(string),
		},
	}

	log.Printf("[INFO] Creating New Relic metric pruning rule for account %d", accountID)

	created, err := client.Pruningrules.NRQLDropRulesCreateWithContext(ctx, accountID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil || len(created.Failures) > 0 {
		var diags diag.Diagnostics
		for _, f := range created.Failures {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  string(f.Error.Reason),
				Detail:   f.Error.Description,
			})
		}
		return diags
	}

	rule := created.Successes[0]
	d.SetId(fmt.Sprintf("%d:%s", rule.AccountID, rule.ID))

	return resourceNewRelicMetricPruningRuleRead(ctx, d, meta)
}

func resourceNewRelicMetricPruningRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic metric pruning rule %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := client.Pruningrules.GetListWithContext(ctx, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	if rules.Error.Reason == pruningrules.NRQLDropRulesErrorReasonTypes.RULE_NOT_FOUND {
		d.SetId("")
		return nil
	}

	for _, r := range rules.Rules {
		if r.ID == ruleID {
			if err := d.Set("account_id", accountID); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("rule_id", r.ID); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("nrql", r.NRQL); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("description", r.Description); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	// Rule not found — removed outside of Terraform.
	d.SetId("")
	return nil
}

func resourceNewRelicMetricPruningRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic metric pruning rule %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleted, err := client.Pruningrules.NRQLDropRulesDeleteWithContext(ctx, accountID, []string{ruleID})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(deleted.Failures) > 0 {
		var diags diag.Diagnostics
		for _, f := range deleted.Failures {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  string(f.Error.Reason),
				Detail:   f.Error.Description,
			})
		}
		return diags
	}

	return nil
}
