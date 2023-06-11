package newrelic

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
)

func resourceNewRelicLogParsingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicLogParsingRuleCreate,
		ReadContext:   resourceNewRelicLogParsingRuleRead,
		UpdateContext: resourceNewRelicLogParsingRuleUpdate,
		DeleteContext: resourceNewRelicLogParsingRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "The account id associated with the obfuscation expression.",
				Computed:    true,
				Optional:    true,
			},
			"attribute": {
				Type:        schema.TypeString,
				Description: "The parsing rule will apply to value of this attribute. If field is not provided, value will default to message.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A description of what this parsing rule represents.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether or not this rule is enabled.",
				Required:    true,
			},
			"grok": {
				Type:        schema.TypeString,
				Description: "The Grok of what to parse.",
				Required:    true,
			},
			"lucene": {
				Type:        schema.TypeString,
				Description: "The Lucene to match events to the parsing rule.",
				Required:    true,
			},
			"nrql": {
				Type:        schema.TypeString,
				Description: "The NRQL to match events to the parsing rule.",
				Required:    true,
			},
			"deleted": {
				Type:        schema.TypeBool,
				Description: "Whether or not this rule is deleted.",
				Computed:    true,
			},
			"matched": {
				Type:        schema.TypeBool,
				Description: "Whether the Grok pattern matched.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Create the obfuscation expression
func resourceNewRelicLogParsingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := logconfigurations.LogConfigurationsParsingRuleConfiguration{
		Enabled: d.Get("enabled").(bool),
		Grok:    d.Get("grok").(string),
		Lucene:  d.Get("lucene").(string),
		NRQL:    logconfigurations.NRQL(d.Get("nrql").(string)),
	}
	var diags diag.Diagnostics
	if e, ok := d.GetOk("attribute"); ok {
		createInput.Attribute = e.(string)
	}

	e := d.Get("name")
	rule, err := getLogParsingRuleByName(ctx, client, accountID, e.(string))
	if (rule != nil && err != nil) || (rule == nil && err != nil) {
		return diag.FromErr(err)
	}
	createInput.Description = e.(string)
	if d.Get("matched") == false {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "The grok pattern is not tested against log lines from the New Relic",
		})
	}

	created, err := client.Logconfigurations.LogConfigurationsCreateParsingRuleWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: rule not created.")
	}

	parsingRuleID := created.Rule.ID

	d.SetId(parsingRuleID)
	return diags
}

// Read the obfuscation expression
func resourceNewRelicLogParsingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ruleID := d.Id()
	rule, err := getLogParsingRuleByID(ctx, client, accountID, ruleID)

	if err != nil && rule == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("attribute", rule.Attribute); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", rule.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("grok", rule.Grok); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lucene", rule.Lucene); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nrql", rule.NRQL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deleted", rule.Deleted); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Update the obfuscation expression
func resourceNewRelicLogParsingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	accountID := selectAccountID(meta.(*ProviderConfig), d)
	if e, ok := d.GetOk("name"); ok {
		rule, err := getLogParsingRuleByName(ctx, client, accountID, e.(string))
		if (rule != nil && err != nil) || (rule == nil && err != nil) {
			return diag.FromErr(err)
		}
	}

	updateInput := expandLogParsingRuleUpdateInput(d)

	log.Printf("[INFO] Updating New Relic logging parsing rule %s", d.Id())

	ruleID := d.Id()

	var diags diag.Diagnostics
	if e, ok := d.GetOk("matched"); ok {
		if !e.(bool) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "The grok pattern is not tested against log lines from the New Relic",
			})
		}
		return diags
	}
	_, err := client.Logconfigurations.LogConfigurationsUpdateParsingRuleWithContext(ctx, accountID, ruleID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceNewRelicLogParsingRuleRead(ctx, d, meta)
}

func expandLogParsingRuleUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsParsingRuleConfiguration {
	updateInp := logconfigurations.LogConfigurationsParsingRuleConfiguration{}

	if e, ok := d.GetOk("attribute"); ok {
		updateInp.Attribute = e.(string)
	}

	if e, ok := d.GetOk("enabled"); ok {
		updateInp.Enabled = e.(bool)
	}

	if e, ok := d.GetOk("name"); ok {
		updateInp.Description = e.(string)
	}

	if e, ok := d.GetOk("grok"); ok {
		updateInp.Grok = e.(string)
	}

	if e, ok := d.GetOk("lucene"); ok {
		updateInp.Lucene = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		updateInp.NRQL = logconfigurations.NRQL(e.(string))
	}

	return updateInp
}

// Delete the logging parsing rule
func resourceNewRelicLogParsingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic logging parsing rule id %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)
	ruleID := d.Id()

	_, err := client.Logconfigurations.LogConfigurationsDeleteParsingRuleWithContext(ctx, accountID, ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getLogParsingRuleByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*logconfigurations.LogConfigurationsParsingRule, error) {
	rules, err := client.Logconfigurations.GetParsingRulesWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range *rules {
		if v.ID == ruleID {
			return v, nil
		}
	}
	return nil, errors.New("parsing rule not found")

}

func getLogParsingRuleByName(ctx context.Context, client *newrelic.NewRelic, accountID int, name string) (*logconfigurations.LogConfigurationsParsingRule, error) {
	rules, err := client.Logconfigurations.GetParsingRulesWithContext(ctx, accountID)
	if rules == nil && err.Error() == "resource not found" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for _, v := range *rules {
		if v.Description == name && !v.Deleted {
			return v, errors.New("name is already in use by another rule")
		}
	}
	return nil, nil

}
