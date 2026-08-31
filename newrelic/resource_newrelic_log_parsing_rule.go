package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
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
			"source": {
				Type:        schema.TypeString,
				Description: "Source of the parsing rule.",
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"NO_CODE",
					"NO_CODE_WRITE_YOUR_OWN",
					"WRITE_YOUR_OWN",
				}, false),
			},
		},
	}
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
	return nil, nrErrors.NewNotFound("parsing rule not found")
}

func getLogParsingRuleByName(ctx context.Context, client *newrelic.NewRelic, accountID int, name string) (*logconfigurations.LogConfigurationsParsingRule, error) {
	rules, err := client.Logconfigurations.GetParsingRulesWithContext(ctx, accountID)
	if rules == nil && err != nil && err.Error() == "resource not found" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for _, v := range *rules {
		if v.Description == name && !v.Deleted {
			return v, nil
		}
	}
	return nil, nil
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

	if e, ok := d.GetOk("source"); ok {
		updateInp.Source = logconfigurations.LogConfigurationsParsingRuleSource(e.(string))
	}

	return updateInp
}

var _ = log.Printf
var _ = diag.FromErr
func resourceNewRelicLogParsingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := expandLogParsingRuleUpdateInput(d)

	var diags diag.Diagnostics

	e := d.Get("name")
	rule, err := getLogParsingRuleByName(ctx, client, accountID, e.(string))
	if (rule != nil && err != nil) || (rule == nil && err != nil) {
		return diag.FromErr(err)
	}

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

	return append(diags, resourceNewRelicLogParsingRuleRead(ctx, d, meta)...)
}

func resourceNewRelicLogParsingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ruleID := d.Id()
	rule, err := getLogParsingRuleByID(ctx, client, accountID, ruleID)

	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", accountID)
	_ = d.Set("name", rule.Description)
	_ = d.Set("attribute", rule.Attribute)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("grok", rule.Grok)
	_ = d.Set("lucene", rule.Lucene)
	_ = d.Set("nrql", string(rule.NRQL))
	_ = d.Set("deleted", rule.Deleted)
	_ = d.Set("source", string(rule.Source))

	return nil
}

func resourceNewRelicLogParsingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	if e, ok := d.GetOk("name"); ok {
		if o, n := d.GetChange("name"); o != n {
			rule, err := getLogParsingRuleByName(ctx, client, accountID, e.(string))
			if (rule != nil && err != nil) || (rule == nil && err != nil) {
				return diag.FromErr(err)
			}
		}
	}

	updateInput := expandLogParsingRuleUpdateInput(d)

	log.Printf("[INFO] Updating New Relic log parsing rule %s", d.Id())

	ruleID := d.Id()

	var diags diag.Diagnostics
	if e, ok := d.GetOk("matched"); ok {
		if !e.(bool) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "The grok pattern is not tested against log lines from the New Relic",
			})
		}
	}

	_, err := client.Logconfigurations.LogConfigurationsUpdateParsingRuleWithContext(ctx, accountID, ruleID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return append(diags, resourceNewRelicLogParsingRuleRead(ctx, d, meta)...)
}

func resourceNewRelicLogParsingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic log parsing rule id %s", d.Id())

	accountID := selectAccountID(providerConfig, d)
	ruleID := d.Id()

	_, err := client.Logconfigurations.LogConfigurationsDeleteParsingRuleWithContext(ctx, accountID, ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}