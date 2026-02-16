package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicCloudAwsEuSovereignIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsEuSovereignIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAwsEuSovereignIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAwsEuSovereignIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAwsEuSovereignIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the account in New Relic.",
			},
			"linked_account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the linked AWS EU Sovereign account in New Relic.",
			},
			// EU Sovereign only supports 5 integrations: billing, cloudtrail, xray
			"billing": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Billing integration",
				Elem:        cloudAwsEuSovereignIntegrationsBillingElem(),
			},
			"cloudtrail": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "CloudTrail integration",
				Elem:        cloudAwsEuSovereignIntegrationsCloudtrailElem(),
			},
			"x_ray": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "X-Ray integration",
				Elem:        cloudAwsEuSovereignIntegrationsXRayElem(),
			},
		},
	}
}

func resourceNewRelicCloudAwsEuSovereignIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	linkedAccountID := d.Get("linked_account_id").(int)

	configureInput := expandCloudAwsEuSovereignIntegrationsInput(d, linkedAccountID)

	payload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(payload.Errors) > 0 {
		for _, err := range payload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	if len(payload.Integrations) > 0 {
		d.SetId(strconv.Itoa(linkedAccountID))
	}

	return nil
}

func resourceNewRelicCloudAwsEuSovereignIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	flattenCloudAwsEuSovereignIntegrations(linkedAccount, accountID, d)

	return nil
}

func resourceNewRelicCloudAwsEuSovereignIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	linkedAccountID := d.Get("linked_account_id").(int)

	configureInput := expandCloudAwsEuSovereignIntegrationsInput(d, linkedAccountID)

	payload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(payload.Errors) > 0 {
		for _, err := range payload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	return resourceNewRelicCloudAwsEuSovereignIntegrationsRead(ctx, d, meta)
}

func resourceNewRelicCloudAwsEuSovereignIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	linkedAccountID := d.Get("linked_account_id").(int)

	disableInput := expandCloudAwsEuSovereignDisableIntegrationsInput(d, linkedAccountID)

	payload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(payload.Errors) > 0 {
		for _, err := range payload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	d.SetId("")

	return nil
}

// Billing integration schema
func cloudAwsEuSovereignIntegrationsBillingElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metrics_polling_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The data polling interval in seconds",
			},
		},
	}
}

// CloudTrail integration schema
func cloudAwsEuSovereignIntegrationsCloudtrailElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metrics_polling_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The data polling interval in seconds",
			},
			"aws_regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specify each AWS region that includes the resources that you want to monitor",
			},
		},
	}
}

// X-Ray integration schema
func cloudAwsEuSovereignIntegrationsXRayElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metrics_polling_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The data polling interval in seconds",
			},
			"aws_regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specify each AWS region that includes the resources that you want to monitor",
			},
		},
	}
}

// expandCloudAwsEuSovereignIntegrationsInput expands the schema data for configuring integrations
// EU Sovereign only supports: billing, cloudtrail, xray
func expandCloudAwsEuSovereignIntegrationsInput(d *schema.ResourceData, linkedAccountID int) cloud.CloudIntegrationsInput {
	awsEuSovereignInput := cloud.CloudAwsEuSovereignIntegrationsInput{}

	// Billing Integration
	if billingRaw := d.Get("billing"); billingRaw != nil {
		billingList := billingRaw.([]interface{})
		if len(billingList) > 0 {
			awsEuSovereignInput.Billing = expandCloudAwsEuSovereignIntegrationBilling(billingList, linkedAccountID)
		}
	}

	// CloudTrail Integration
	// Use Get instead of GetOk because GetOk returns false for empty blocks like `cloudtrail {}`
	if cloudtrailRaw := d.Get("cloudtrail"); cloudtrailRaw != nil {
		cloudtrailList := cloudtrailRaw.([]interface{})
		if len(cloudtrailList) > 0 {
			awsEuSovereignInput.Cloudtrail = expandCloudAwsEuSovereignIntegrationCloudtrail(cloudtrailList, linkedAccountID)
		}
	}

	// X-Ray Integration
	if xrayRaw := d.Get("x_ray"); xrayRaw != nil {
		xrayList := xrayRaw.([]interface{})
		if len(xrayList) > 0 {
			awsEuSovereignInput.AwsXray = expandCloudAwsEuSovereignIntegrationXRay(xrayList, linkedAccountID)
		}
	}

	input := cloud.CloudIntegrationsInput{
		AwsEuSovereign: awsEuSovereignInput,
	}

	return input
}

// expandCloudAwsEuSovereignDisableIntegrationsInput expands the schema data for disabling integrations
func expandCloudAwsEuSovereignDisableIntegrationsInput(d *schema.ResourceData, linkedAccountID int) cloud.CloudDisableIntegrationsInput {
	awsEuSovereignInput := cloud.CloudAwsEuSovereignDisableIntegrationsInput{}

	// Use Get instead of GetOk because GetOk returns false for empty blocks
	if v := d.Get("billing"); v != nil && len(v.([]interface{})) > 0 {
		awsEuSovereignInput.Billing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v := d.Get("cloudtrail"); v != nil && len(v.([]interface{})) > 0 {
		awsEuSovereignInput.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v := d.Get("x_ray"); v != nil && len(v.([]interface{})) > 0 {
		awsEuSovereignInput.AwsXray = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	input := cloud.CloudDisableIntegrationsInput{
		AwsEuSovereign: awsEuSovereignInput,
	}

	return input
}

// flattenCloudAwsEuSovereignIntegrations flattens the integrations data from the API into the schema
func flattenCloudAwsEuSovereignIntegrations(linkedAccount *cloud.CloudLinkedAccount, accountID int, d *schema.ResourceData) {
	_ = d.Set("account_id", accountID)
	_ = d.Set("linked_account_id", linkedAccount.ID)

	for _, i := range linkedAccount.Integrations {
		switch t := i.(type) {
		case *cloud.CloudBillingIntegration:
			_ = d.Set("billing", flattenCloudAwsEuSovereignIntegrationBilling(t))
		case *cloud.CloudCloudtrailIntegration:
			_ = d.Set("cloudtrail", flattenCloudAwsEuSovereignIntegrationCloudtrail(t))
		case *cloud.CloudAwsXrayIntegration:
			_ = d.Set("x_ray", flattenCloudAwsEuSovereignIntegrationXRay(t))
		}
	}
}

// Billing expand/flatten
func expandCloudAwsEuSovereignIntegrationBilling(b []interface{}, linkedAccountID int) []cloud.CloudBillingIntegrationInput {
	expanded := make([]cloud.CloudBillingIntegrationInput, len(b))

	for i, billing := range b {
		var billingInput cloud.CloudBillingIntegrationInput

		if billing == nil {
			billingInput.LinkedAccountId = linkedAccountID
			expanded[i] = billingInput
			return expanded
		}

		cfg := billing.(map[string]interface{})
		billingInput.LinkedAccountId = linkedAccountID

		if v, ok := cfg["metrics_polling_interval"]; ok {
			billingInput.MetricsPollingInterval = v.(int)
		}

		expanded[i] = billingInput
	}

	return expanded
}

func flattenCloudAwsEuSovereignIntegrationBilling(t *cloud.CloudBillingIntegration) []interface{} {
	result := make(map[string]interface{})
	result["metrics_polling_interval"] = t.MetricsPollingInterval
	return []interface{}{result}
}

// CloudTrail expand/flatten
func expandCloudAwsEuSovereignIntegrationCloudtrail(b []interface{}, linkedAccountID int) []cloud.CloudCloudtrailIntegrationInput {
	expanded := make([]cloud.CloudCloudtrailIntegrationInput, len(b))

	for i, cloudtrail := range b {
		var cloudtrailInput cloud.CloudCloudtrailIntegrationInput

		if cloudtrail == nil {
			cloudtrailInput.LinkedAccountId = linkedAccountID
			expanded[i] = cloudtrailInput
			return expanded
		}

		cfg := cloudtrail.(map[string]interface{})
		cloudtrailInput.LinkedAccountId = linkedAccountID

		if v, ok := cfg["metrics_polling_interval"]; ok {
			cloudtrailInput.MetricsPollingInterval = v.(int)
		}

		if v, ok := cfg["aws_regions"]; ok {
			awsRegions := v.([]interface{})
			var regions []string
			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			cloudtrailInput.AwsRegions = regions
		}

		expanded[i] = cloudtrailInput
	}

	return expanded
}

func flattenCloudAwsEuSovereignIntegrationCloudtrail(t *cloud.CloudCloudtrailIntegration) []interface{} {
	result := make(map[string]interface{})
	result["metrics_polling_interval"] = t.MetricsPollingInterval
	result["aws_regions"] = t.AwsRegions
	return []interface{}{result}
}

// X-Ray expand/flatten
func expandCloudAwsEuSovereignIntegrationXRay(b []interface{}, linkedAccountID int) []cloud.CloudAwsXrayIntegrationInput {
	expanded := make([]cloud.CloudAwsXrayIntegrationInput, len(b))

	for i, xray := range b {
		var xrayInput cloud.CloudAwsXrayIntegrationInput

		if xray == nil {
			xrayInput.LinkedAccountId = linkedAccountID
			expanded[i] = xrayInput
			return expanded
		}

		cfg := xray.(map[string]interface{})
		xrayInput.LinkedAccountId = linkedAccountID

		if v, ok := cfg["metrics_polling_interval"]; ok {
			xrayInput.MetricsPollingInterval = v.(int)
		}

		if v, ok := cfg["aws_regions"]; ok {
			awsRegions := v.([]interface{})
			var regions []string
			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			xrayInput.AwsRegions = regions
		}

		expanded[i] = xrayInput
	}

	return expanded
}

func flattenCloudAwsEuSovereignIntegrationXRay(t *cloud.CloudAwsXrayIntegration) []interface{} {
	result := make(map[string]interface{})
	result["metrics_polling_interval"] = t.MetricsPollingInterval
	result["aws_regions"] = t.AwsRegions
	return []interface{}{result}
}
