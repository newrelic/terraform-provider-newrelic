package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/eventstometrics"
)

func resourceNewRelicEventsToMetricsRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicEventsToMetricsRuleCreate,
		ReadContext:   resourceNewRelicEventsToMetricsRuleRead,
		UpdateContext: resourceNewRelicEventsToMetricsRuleUpdate,
		DeleteContext: resourceNewRelicEventsToMetricsRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Account with the event and where the metrics will be put.",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the rule. This must be unique within an account.",
			},
			"nrql": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Explains how to create metrics from events.",
			},
			"description": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Provides additional information about the rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "True means this rule is enabled. False means the rule is currently not creating metrics.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id, uniquely identifying the rule.",
			},
		},
	}
}

func resourceNewRelicEventsToMetricsRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient

	createInput := []eventstometrics.EventsToMetricsCreateRuleInput{
		{
			AccountID:   d.Get("account_id").(int),
			Description: d.Get("description").(string),
			Name:        d.Get("name").(string),
			NRQL:        d.Get("nrql").(string),
		},
	}

	rules, err := client.EventsToMetrics.CreateRulesWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := rules[0]

	id := fmt.Sprintf("%d:%s", rule.AccountID, rule.ID)

	d.SetId(id)

	if enabled, ok := d.GetOk("enabled"); ok {
		updateInput := []eventstometrics.EventsToMetricsUpdateRuleInput{
			{
				AccountID: rule.AccountID,
				RuleId:    rule.ID,
				Enabled:   enabled.(bool),
			},
		}

		_, err := client.EventsToMetrics.UpdateRulesWithContext(ctx, updateInput)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNewRelicEventsToMetricsRuleRead(ctx, d, meta)
}

func resourceNewRelicEventsToMetricsRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic entity tags for entity guid %s", d.Id())

	accountID, ruleID, err := getEventsToMetricsRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := client.EventsToMetrics.GetRuleWithContext(ctx, accountID, ruleID)

	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule_id", ruleID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", rule.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nrql", rule.NRQL); err != nil {
		return diag.FromErr(err)
	}

	_, ok := d.GetOk("enabled")
	if ok {
		if err := d.Set("enabled", rule.Enabled); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNewRelicEventsToMetricsRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Update")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Updating New Relic events to metric rules")

	accountID, ruleID, err := getEventsToMetricsRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateInput := []eventstometrics.EventsToMetricsUpdateRuleInput{
		{
			AccountID: accountID,
			RuleId:    ruleID,
			Enabled:   d.Get("enabled").(bool),
		},
	}

	_, err = client.EventsToMetrics.UpdateRulesWithContext(ctx, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicEventsToMetricsRuleRead(ctx, d, meta)
}

func resourceNewRelicEventsToMetricsRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Delete")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity tags from entity guid %s", d.Id())

	accountID, ruleID, err := getEventsToMetricsRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInput := []eventstometrics.EventsToMetricsDeleteRuleInput{
		{
			AccountID: accountID,
			RuleId:    ruleID,
		},
	}

	_, err = client.EventsToMetrics.DeleteRulesWithContext(ctx, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getEventsToMetricsRuleIDs(id string) (int, string, error) {
	strIDs := strings.Split(id, ":")

	if len(strIDs) != 2 {
		return 0, "", errors.New("could not parse events to metrics rule IDs")
	}

	accountID, err := strconv.Atoi(strIDs[0])
	if err != nil {
		return 0, "", err
	}

	return accountID, strIDs[1], nil
}
