package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
	"log"
	"time"
)

func resourceNewRelicDataPartition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicDataPartitionCreate,
		ReadContext:   resourceNewRelicDataPartitionRead,
		UpdateContext: resourceNewRelicDataPartitionUpdate,
		DeleteContext: resourceNewRelicDataPartitionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "The account id associated with the data partition rule.",
				Computed:    true,
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the data partition rule.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether or not this data partition rule is enabled.",
				Required:    true,
			},
			"nrql": {
				Type:        schema.TypeString,
				Description: "The NRQL to match events for this data partition rule. Logs matching this criteria will be routed to the specified data partition.",
				Optional:    true,
			},
			"retention_policy": {
				Type:         schema.TypeString,
				Description:  "The retention policy of the data partition data.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(listValidDataPartitionRuleRetentionPolicyType(), false),
			},
			"target_data_partition": {
				Type:        schema.TypeString,
				Description: "The name of the data partition where logs will be allocated once the rule is enabled.",
				Required:    true,
				ForceNew:    true,
			},
			"matching_criteria": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The matching criteria of the data partition rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The attribute name against which this matching condition will be evaluated.",
						},
						"matching_expression": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The matching expression of the data partition rule definition.",
						},
						"matching_method": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The matching method of the data partition rule definition.",
							ValidateFunc: validation.StringInSlice(listValidDataPartitionRuleMatchingOperator(), false),
						},
					},
				},
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this data partition rule is deleted. Deleting a data partition rule does not delete the already persisted data. This data will be retained for a given period of time specified in the retention policy field.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func listValidDataPartitionRuleRetentionPolicyType() []string {
	return []string{
		string(logconfigurations.LogConfigurationsDataPartitionRuleRetentionPolicyTypeTypes.SECONDARY),
		string(logconfigurations.LogConfigurationsDataPartitionRuleRetentionPolicyTypeTypes.STANDARD),
	}
}

func listValidDataPartitionRuleMatchingOperator() []string {
	return []string{
		string(logconfigurations.LogConfigurationsDataPartitionRuleMatchingOperatorTypes.EQUALS),
		string(logconfigurations.LogConfigurationsDataPartitionRuleMatchingOperatorTypes.LIKE),
	}
}

func expandDataPartitionUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsUpdateDataPartitionRuleInput {
	updateInp := logconfigurations.LogConfigurationsUpdateDataPartitionRuleInput{
		ID: d.Id(),
	}
	if e, ok := d.GetOk("enabled"); ok {
		updateInp.Enabled = e.(bool)
	}

	if e, ok := d.GetOk("description"); ok {
		updateInp.Description = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		updateInp.NRQL = logconfigurations.NRQL(e.(string))
	}

	if v, ok := d.GetOk("matching_criteria"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			updateInp.MatchingCriteria = &logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      cfg["attribute_name"].(string),
				MatchingExpression: cfg["matching_expression"].(string),
				MatchingMethod:     logconfigurations.LogConfigurationsDataPartitionRuleMatchingOperator(cfg["matching_method"].(string)),
			}
		}
	}

	return updateInp
}

// Create the data partition rule

// Read the data partition rule

// Update the data partition rule

// Delete the data partition rule

func resourceNewRelicDataPartitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := logconfigurations.LogConfigurationsCreateDataPartitionRuleInput{
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
	}

	if e, ok := d.GetOk("nrql"); ok {
		createInput.NRQL = logconfigurations.NRQL(e.(string))
	}

	if e, ok := d.GetOk("target_data_partition"); ok {
		createInput.TargetDataPartition = logconfigurations.LogConfigurationsLogDataPartitionName(e.(string))
	}

	if e, ok := d.GetOk("retention_policy"); ok {
		createInput.RetentionPolicy = logconfigurations.LogConfigurationsDataPartitionRuleRetentionPolicyType(e.(string))
	}

	if v, ok := d.GetOk("matching_criteria"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			createInput.MatchingCriteria = &logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteriaInput{
				AttributeName:      cfg["attribute_name"].(string),
				MatchingExpression: cfg["matching_expression"].(string),
				MatchingMethod:     logconfigurations.LogConfigurationsDataPartitionRuleMatchingOperator(cfg["matching_method"].(string)),
			}
		}
	}

	log.Printf("[INFO] Creating New Relic Data Partition Rule  %s", createInput.TargetDataPartition)

	created, err := client.Logconfigurations.LogConfigurationsCreateDataPartitionRuleWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var apiDiags diag.Diagnostics

	if created.Errors != nil {
		for _, err := range created.Errors {
			apiDiags = append(apiDiags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Message,
				Detail:   string(err.Type),
			})
		}
		return apiDiags
	}

	if created == nil {
		return diag.Errorf("err: data partition rule create result wasn't returned or rule was not created.")
	}

	ruleID := created.Rule.ID

	d.SetId(ruleID)

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		rules, err := client.Logconfigurations.GetDataPartitionRulesWithContext(ctx, accountID)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		for _, v := range *rules {
			if v.ID == ruleID && !v.Deleted {
				return nil
			}
		}
		return resource.RetryableError(fmt.Errorf("data partition rule was not created"))
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourceNewRelicDataPartitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ruleID := d.Id()
	rule, err := getDataPartitionByID(ctx, client, accountID, ruleID)

	if err != nil || rule == nil || rule.Deleted {
		d.SetId("")
		return nil
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("target_data_partition", rule.TargetDataPartition)
	_ = d.Set("nrql", rule.NRQL)
	_ = d.Set("retention_policy", rule.RetentionPolicy)
	_ = d.Set("deleted", rule.Deleted)

	if rule.MatchingCriteria != (logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteria{}) {
		mc := map[string]interface{}{
			"attribute_name":      rule.MatchingCriteria.AttributeName,
			"matching_expression": rule.MatchingCriteria.MatchingExpression,
			"matching_method":     string(rule.MatchingCriteria.MatchingOperator),
		}
		if err := d.Set("matching_criteria", []interface{}{mc}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNewRelicDataPartitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandDataPartitionUpdateInput(d)

	log.Printf("[INFO] Updating New Relic Data Partition Rule %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)

	updated, err := client.Logconfigurations.LogConfigurationsUpdateDataPartitionRuleWithContext(ctx, accountID, updateInput)

	if err != nil {
		return diag.FromErr(err)
	}

	var apiDiags diag.Diagnostics

	if updated.Errors != nil {
		for _, err := range updated.Errors {
			apiDiags = append(apiDiags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Message,
				Detail:   string(err.Type),
			})
		}
		return apiDiags
	}

	return nil
}

func resourceNewRelicDataPartitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic Data Partition Rule id %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)
	expressionID := d.Id()

	_, err := client.Logconfigurations.LogConfigurationsDeleteDataPartitionRuleWithContext(ctx, accountID, expressionID)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getDataPartitionByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*logconfigurations.LogConfigurationsDataPartitionRule, error) {
	rules, err := client.Logconfigurations.GetDataPartitionRulesWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range *rules {
		if v.ID == ruleID {
			return &v, nil
		}
	}
	return nil, errors.New("data partition rule not found")
}

func expandDataPartitionRuleCreateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsCreateDataPartitionRuleInput {
	input := logconfigurations.LogConfigurationsCreateDataPartitionRuleInput{
		Enabled: d.Get("enabled").(bool),
	}

	if e, ok := d.GetOk("description"); ok {
		input.Description = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		input.NRQL = logconfigurations.NRQL(e.(string))
	}

	if e, ok := d.GetOk("target_data_partition"); ok {
		input.TargetDataPartition = logconfigurations.LogConfigurationsLogDataPartitionName(e.(string))
	}

	if e, ok := d.GetOk("retention_policy"); ok {
		input.RetentionPolicy = logconfigurations.LogConfigurationsDataPartitionRuleRetentionPolicyType(e.(string))
	}

	if v, ok := d.GetOk("matching_criteria"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			input.MatchingCriteria = expandDataPartitionRuleMatchingCriteriaInput(items[0].(map[string]interface{}))
		}
	}

	return input
}

func expandDataPartitionRuleUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsUpdateDataPartitionRuleInput {
	input := logconfigurations.LogConfigurationsUpdateDataPartitionRuleInput{
		ID:      d.Id(),
		Enabled: d.Get("enabled").(bool),
	}

	if e, ok := d.GetOk("description"); ok {
		input.Description = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		input.NRQL = logconfigurations.NRQL(e.(string))
	}

	if v, ok := d.GetOk("matching_criteria"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			input.MatchingCriteria = expandDataPartitionRuleMatchingCriteriaInput(items[0].(map[string]interface{}))
		}
	}

	return input
}

func expandDataPartitionRuleMatchingCriteriaInput(cfg map[string]interface{}) *logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteriaInput {
	input := &logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteriaInput{}

	if v, ok := cfg["attribute_name"]; ok {
		input.AttributeName = v.(string)
	}

	if v, ok := cfg["matching_expression"]; ok {
		input.MatchingExpression = v.(string)
	}

	if v, ok := cfg["matching_method"]; ok {
		input.MatchingMethod = logconfigurations.LogConfigurationsDataPartitionRuleMatchingOperator(v.(string))
	}

	return input
}

func flattenDataPartitionRule(rule *logconfigurations.LogConfigurationsDataPartitionRule, d *schema.ResourceData) error {
	if rule == nil {
		return nil
	}

	_ = d.Set("description", rule.Description)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("target_data_partition", string(rule.TargetDataPartition))
	_ = d.Set("nrql", string(rule.NRQL))
	_ = d.Set("retention_policy", string(rule.RetentionPolicy))
	_ = d.Set("deleted", rule.Deleted)

	if rule.MatchingCriteria != (logconfigurations.LogConfigurationsDataPartitionRuleMatchingCriteria{}) {
		mc := map[string]interface{}{
			"attribute_name":      rule.MatchingCriteria.AttributeName,
			"matching_expression": rule.MatchingCriteria.MatchingExpression,
			"matching_method":     string(rule.MatchingCriteria.MatchingOperator),
		}
		if err := d.Set("matching_criteria", []interface{}{mc}); err != nil {
			return fmt.Errorf("[DEBUG] Error setting `matching_criteria`: %v", err)
		}
	}

	return nil
}
