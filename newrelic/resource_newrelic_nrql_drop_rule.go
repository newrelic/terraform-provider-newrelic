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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/nrqldroprules"
)

func resourceNewRelicNRQLDropRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNRQLDropRuleCreate,
		ReadContext:   resourceNewRelicNRQLDropRuleRead,
		DeleteContext: resourceNewRelicNRQLDropRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Account with the NRQL drop rule will be put.",
			},
			"action": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"drop_data", "drop_attributes", "drop_attributes_from_metric_aggregates"}, false),
				Description:  "The drop rule action (drop_data, drop_attributes, or drop_attributes_from_metric_aggregates).",
			},
			"nrql": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Explains which data to apply the drop rule to.",
			},
			"description": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Provides additional information about the rule.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id, uniquely identifying the rule.",
			},
		},
	}
}

func resourceNewRelicNRQLDropRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := []nrqldroprules.NRQLDropRulesCreateDropRuleInput{
		{
			Description: d.Get("description").(string),
			Action:      nrqldroprules.NRQLDropRulesAction(strings.ToUpper(d.Get("action").(string))),
			NRQL:        d.Get("nrql").(string),
		},
	}

	created, err := client.Nrqldroprules.NRQLDropRulesCreateWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil || len(created.Successes) == 0 {
		return diag.Errorf("err: drop rule create result wasn't returned. Validate the action value or NRQL query.")
	}

	rule := created.Successes[0]

	id := fmt.Sprintf("%d:%s", rule.AccountID, rule.ID)

	d.SetId(id)

	return resourceNewRelicNRQLDropRuleRead(ctx, d, meta)
}

func resourceNewRelicNRQLDropRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic NRQL Drop Rule for %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := getNRQLDropRuleByID(ctx, client, accountID, ruleID)

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

	if err := d.Set("action", strings.ToLower(string(rule.Action))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nrql", rule.NRQL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicNRQLDropRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Delete")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity tags from entity guid %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInput := []string{ruleID}

	_, err = client.Nrqldroprules.NRQLDropRulesDeleteWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func parseNRQLDropRuleIDs(id string) (int, string, error) {
	strIDs := strings.Split(id, ":")

	if len(strIDs) != 2 {
		return 0, "", errors.New("could not parse drop rule IDs")
	}

	accountID, err := strconv.Atoi(strIDs[0])
	if err != nil {
		return 0, "", err
	}

	return accountID, strIDs[1], nil
}

func getNRQLDropRuleByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*nrqldroprules.NRQLDropRulesDropRule, error) {
	rules, err := client.Nrqldroprules.GetListWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range rules.Rules {
		if v.ID == ruleID {
			return &v, nil
		}
	}
	return nil, errors.New("drop rule not found")
}
