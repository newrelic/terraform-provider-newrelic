package newrelic

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
	"log"
)

//Defining the schema of the New Relic Obfuscation rule
func resourceNewRelicObfuscationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicObfuscationRuleCreate,
		ReadContext:   resourceNewRelicObfuscationRuleRead,
		UpdateContext: resourceNewRelicObfuscationRuleUpdate,
		DeleteContext: resourceNewRelicObfuscationRuleDelete,
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
			"filter": {
				Type:        schema.TypeString,
				Description: "NRQL for determining whether a given log record should have obfuscation actions applied.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of expression.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the rule should be applied or not to incoming data.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of rule.",
				Optional:    true,
			},
			"action": {
				Type:        schema.TypeSet,
				Description: "Actions for the rule. The actions will be applied in the order specified by this list.",
				Required:    true,
				Elem:        ObfuscationRuleActionInputSchemaElem(),
			},
			"action_id": {
				Type:        schema.TypeString,
				Description: "The id of the obfuscation action.",
				Computed:    true,
			},
		},
	}
}

//ObfuscationRuleActionInputSchemaElem returns the schema of the actions of the New Relic obfuscation rule
func ObfuscationRuleActionInputSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Attribute names for action. An empty list applies the action to all the attributes.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"expression_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expression Id for action.",
			},
			"method": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Obfuscation method to use.",
				ValidateFunc: validation.StringInSlice(listValidLogConfigurationsObfuscationMethod(), false),
			},
		},
	}

}

// A function to validate the methods for the actions of New Relic obfuscation rule
func listValidLogConfigurationsObfuscationMethod() []string {
	return []string{
		string(logconfigurations.LogConfigurationsObfuscationMethodTypes.HASH_SHA256),
		string(logconfigurations.LogConfigurationsObfuscationMethodTypes.MASK),
	}
}

//Create the obfuscation rule
func resourceNewRelicObfuscationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := logconfigurations.LogConfigurationsCreateObfuscationRuleInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Filter:      logconfigurations.NRQL(d.Get("filter").(string)),
	}

	if e, ok := d.GetOk("action"); ok {
		createInput.Actions = expandObfuscationRuleActionInput(e.(*schema.Set).List())
	}

	created, err := client.Logconfigurations.LogConfigurationsCreateObfuscationRuleWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: obfuscation rule create result wasn't returned or rule was not created.")
	}

	ruleID := created.ID
	d.SetId(ruleID)

	return resourceNewRelicObfuscationRuleRead(ctx, d, meta)
}

//Read the obfuscation Rule
func resourceNewRelicObfuscationRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ruleID := d.Id()
	rule, err := getObfuscationRuleByID(ctx, client, accountID, ruleID)

	if err != nil && rule == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", rule.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("filter", rule.Filter); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", rule.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("action", flattenActions(&rule.Actions)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenActions(actions *[]logconfigurations.LogConfigurationsObfuscationAction) []map[string]interface{} {
	var flatActions []map[string]interface{}
	for _, v := range *actions {
		m := map[string]interface{}{
			"expression_id": v.Expression.ID,
			"attribute":     v.Attributes,
			"method":        v.Method,
			//"action_id":     v.ID,
		}
		flatActions = append(flatActions, m)
	}
	return flatActions
}

//Update the obfuscation Rule
func resourceNewRelicObfuscationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandObfuscationRuleUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One obfuscation Rule %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)

	_, err := client.Logconfigurations.LogConfigurationsUpdateObfuscationRuleWithContext(ctx, accountID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicObfuscationRuleRead(ctx, d, meta)
}

func expandObfuscationRuleUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsUpdateObfuscationRuleInput {
	updateInp := logconfigurations.LogConfigurationsUpdateObfuscationRuleInput{
		ID: d.Id(),
		//Actions are required field
	}

	if e, ok := d.GetOk("name"); ok {
		updateInp.Name = e.(string)
	}

	if e, ok := d.GetOk("enabled"); ok {
		updateInp.Enabled = e.(bool)
	}

	if e, ok := d.GetOk("description"); ok {
		updateInp.Description = e.(string)
	}

	if e, ok := d.GetOk("filter"); ok {
		updateInp.Filter = logconfigurations.NRQL(e.(string))
	}

	if e, ok := d.GetOk("action"); ok {
		updateInp.Actions = expandObfuscationRuleActionUpdateInput(e.(*schema.Set).List())

	}

	return updateInp
}

//Delete the obfuscation Rule
func resourceNewRelicObfuscationRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic obfuscation Rule id %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)
	ruleID := d.Id()

	_, err := client.Logconfigurations.LogConfigurationsDeleteObfuscationRuleWithContext(ctx, accountID, ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandObfuscationRuleActionInput(list []interface{}) []logconfigurations.LogConfigurationsCreateObfuscationActionInput {

	actionsList := make([]logconfigurations.LogConfigurationsCreateObfuscationActionInput, len(list))

	for i, v := range list {
		setActions := v.(map[string]interface{})
		action := logconfigurations.LogConfigurationsCreateObfuscationActionInput{
			Method:       logconfigurations.LogConfigurationsObfuscationMethod(setActions["method"].(string)),
			Attributes:   expandAttributes(setActions["attribute"].(*schema.Set).List()),
			ExpressionId: setActions["expression_id"].(string),
		}
		actionsList[i] = action
	}
	return actionsList
}

func expandObfuscationRuleActionUpdateInput(list []interface{}) []logconfigurations.LogConfigurationsUpdateObfuscationActionInput {

	actionsList := make([]logconfigurations.LogConfigurationsUpdateObfuscationActionInput, len(list))

	for i, v := range list {
		setActions := v.(map[string]interface{})
		action := logconfigurations.LogConfigurationsUpdateObfuscationActionInput{
			Method:       logconfigurations.LogConfigurationsObfuscationMethod(setActions["method"].(string)),
			Attributes:   expandAttributes(setActions["attribute"].(*schema.Set).List()),
			ExpressionId: setActions["expression_id"].(string),
		}
		actionsList[i] = action
	}
	return actionsList
}

func expandAttributes(attributeList []interface{}) []string {
	actionsOut := make([]string, len(attributeList))

	for i, v := range attributeList {
		actionsOut[i] = v.(string)
	}
	return actionsOut
}

func getObfuscationRuleByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*logconfigurations.LogConfigurationsObfuscationRule, error) {
	rules, err := client.Logconfigurations.GetObfuscationRulesWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range *rules {
		if v.ID == ruleID {
			return &v, nil
		}
	}
	return nil, errors.New("obfuscation rule not found")

}
