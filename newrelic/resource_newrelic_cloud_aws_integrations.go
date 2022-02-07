package newrelic

import (
	"context"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNewRelicCloudAwsIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAwsIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAwsIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAwsIntegrationsDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
			"linked_account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the linked AWS account in New Relic",
			},
			"billing": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Billing integration",
				Elem:        cloudAwsIntegrationBillingSchemaElem(),
			},
			"cloudtrail": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "CloudTrail integration",
				Elem:        cloudAwsIntegrationCloudTrailSchemaElem(),
			},
			"health": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Health integration",
				Elem:        cloudAwsIntegrationHealthSchemaElem(),
			},
			"trusted_advisor": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Trusted Advisor integration",
				Elem:        cloudAwsIntegrationTrustedAdvisorSchemaElem(),
			},
			"vpc": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "VPC integration",
				Elem:        cloudAwsIntegrationVpcSchemaElem(),
			},
			"x_ray": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "X-Ray integration",
				Elem:        cloudAwsIntegrationXRaySchemaElem(),
			},
		},
	}
}

func cloudAwsIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"integration_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The ID of the AWS integration",
		},
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds.",
		},
	}
}

func cloudAwsIntegrationBillingSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationCloudTrailSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationHealthSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationTrustedAdvisorSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationVpcSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
	}

	s["fetch_nat_gateway"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specify if NAT gateway should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.",
	}

	s["fetch_vpn"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specify if VPN should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.",
	}

	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
	}

	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationXRaySchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func resourceNewRelicCloudAwsIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	cloudAwsIntegrationsInput := expandCloudAwsIntegrationsInput(d)

	cloudAwsIntegrationsPayload, err := client.Cloud.CloudConfigureIntegration(accountID, cloudAwsIntegrationsInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudAwsIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudAwsIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return resourceNewRelicCloudAwsIntegrationsRead(ctx, d, meta)
}

func expandCloudAwsIntegrationsInput(d *schema.ResourceData) cloud.CloudIntegrationsInput {
	cloudAwsIntegration := cloud.CloudAwsIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}

	if b, ok := d.GetOk("billing"); ok {
		cloudAwsIntegration.Billing = expandCloudAwsIntegrationBillingInput(b.(map[string]interface{}), linkedAccountID)
	}

	if c, ok := d.GetOk("cloudtrail"); ok {
		cloudAwsIntegration.Cloudtrail = expandCloudAwsIntegrationCloudtrailInput(c.(map[string]interface{}), linkedAccountID)
	}

	if h, ok := d.GetOk("health"); ok {
		cloudAwsIntegration.Health = expandCloudAwsIntegrationHealthInput(h.(map[string]interface{}), linkedAccountID)
	}

	if t, ok := d.GetOk("trusted_advisor"); ok {
		cloudAwsIntegration.Trustedadvisor = expandCloudAwsIntegrationTrustedAdvisorInput(t.(map[string]interface{}), linkedAccountID)
	}

	if v, ok := d.GetOk("vpc"); ok {
		cloudAwsIntegration.Vpc = expandCloudAwsIntegrationVpcInput(v.(map[string]interface{}), linkedAccountID)
	}

	if x, ok := d.GetOk("x_ray"); ok {
		cloudAwsIntegration.AwsXray = expandCloudAwsIntegrationXRayInput(x.(map[string]interface{}), linkedAccountID)
	}

	input := cloud.CloudIntegrationsInput{
		Aws: cloudAwsIntegration,
	}

	return input
}

func expandCloudAwsIntegrationBillingInput(b map[string]interface{}, linkedAccountID int) []cloud.CloudBillingIntegrationInput {
	var billingInput cloud.CloudBillingIntegrationInput

	billingInput.LinkedAccountId = linkedAccountID

	if m, ok := b["metrics_polling_interval"]; ok {
		billingInput.MetricsPollingInterval = m.(int)
	}

	return []cloud.CloudBillingIntegrationInput{billingInput}
}

func expandCloudAwsIntegrationCloudtrailInput(c map[string]interface{}, linkedAccountID int) []cloud.CloudCloudtrailIntegrationInput {
	var cloudtrailInput cloud.CloudCloudtrailIntegrationInput

	cloudtrailInput.LinkedAccountId = linkedAccountID

	if a, ok := c["aws_regions"]; ok {
		cloudtrailInput.AwsRegions = a.([]string)
	}

	if m, ok := c["metrics_polling_interval"]; ok {
		cloudtrailInput.MetricsPollingInterval = m.(int)
	}

	return []cloud.CloudCloudtrailIntegrationInput{cloudtrailInput}
}

func expandCloudAwsIntegrationHealthInput(h map[string]interface{}, linkedAccountID int) []cloud.CloudHealthIntegrationInput {
	var healthInput cloud.CloudHealthIntegrationInput

	healthInput.LinkedAccountId = linkedAccountID

	if m, ok := h["metrics_polling_interval"]; ok {
		healthInput.MetricsPollingInterval = m.(int)
	}

	return []cloud.CloudHealthIntegrationInput{healthInput}
}

func expandCloudAwsIntegrationTrustedAdvisorInput(t map[string]interface{}, linkedAccountID int) []cloud.CloudTrustedadvisorIntegrationInput {
	var trustedAdvisorInput cloud.CloudTrustedadvisorIntegrationInput

	trustedAdvisorInput.LinkedAccountId = linkedAccountID

	if m, ok := t["metrics_polling_interval"]; ok {
		trustedAdvisorInput.MetricsPollingInterval = m.(int)
	}

	return []cloud.CloudTrustedadvisorIntegrationInput{trustedAdvisorInput}
}

func expandCloudAwsIntegrationVpcInput(v map[string]interface{}, linkedAccountID int) []cloud.CloudVpcIntegrationInput {
	var vpcInput cloud.CloudVpcIntegrationInput

	vpcInput.LinkedAccountId = linkedAccountID

	if a, ok := v["aws_regions"]; ok {
		vpcInput.AwsRegions = a.([]string)
	}

	if nat, ok := v["fetch_nat_gateway"]; ok {
		vpcInput.FetchNatGateway = nat.(bool)
	}

	if vpn, ok := v["fetch_vpn"]; ok {
		vpcInput.FetchVpn = vpn.(bool)
	}

	if m, ok := v["metrics_polling_interval"]; ok {
		vpcInput.MetricsPollingInterval = m.(int)
	}

	if tk, ok := v["tag_key"]; ok {
		vpcInput.TagKey = tk.(string)
	}

	if tv, ok := v["tag_value"]; ok {
		vpcInput.TagKey = tv.(string)
	}

	return []cloud.CloudVpcIntegrationInput{vpcInput}
}

func expandCloudAwsIntegrationXRayInput(x map[string]interface{}, linkedAccountID int) []cloud.CloudAwsXrayIntegrationInput {
	var xrayInput cloud.CloudAwsXrayIntegrationInput

	xrayInput.LinkedAccountId = linkedAccountID

	if a, ok := x["aws_regions"]; ok {
		xrayInput.AwsRegions = a.([]string)
	}

	if m, ok := x["metrics_polling_interval"]; ok {
		xrayInput.MetricsPollingInterval = m.(int)
	}

	return []cloud.CloudAwsXrayIntegrationInput{xrayInput}
}

func resourceNewRelicCloudAwsIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID := d.Get("linked_account_id").(int)

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenCloudAwsLinkedAccount(d, linkedAccount)

	return nil
}

func flattenCloudAwsLinkedAccount(d *schema.ResourceData, linkedAccount *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", linkedAccount.NrAccountId)
	_ = d.Set("linked_account_id", linkedAccount.ID)

	for _, i := range linkedAccount.Integrations {
		switch t := i.(type) {
		case *cloud.CloudBillingIntegration:
			_ = d.Set("billing", flattenCloudAwsBillingIntegration(t))
		case *cloud.CloudCloudtrailIntegration:
			_ = d.Set("cloudtrail", flattenCloudAwsCloudTrailIntegration(t))
		case *cloud.CloudHealthIntegration:
			_ = d.Set("health", flattenCloudAwsHealthIntegration(t))
		case *cloud.CloudTrustedadvisorIntegration:
			_ = d.Set("trusted_advisor", flattenCloudAwsTrustedAdvisorIntegration(t))
		case *cloud.CloudVpcIntegration:
			_ = d.Set("vpc", flattenCloudAwsVpcIntegration(t))
		case *cloud.CloudAwsXrayIntegration:
			_ = d.Set("x_ray", flattenCloudAwsXRayIntegration(t))
		}
	}
}

func flattenCloudAwsBillingIntegration(in *cloud.CloudBillingIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	return out
}

func flattenCloudAwsCloudTrailIntegration(in *cloud.CloudCloudtrailIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	return out
}

func flattenCloudAwsHealthIntegration(in *cloud.CloudHealthIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	return out
}

func flattenCloudAwsTrustedAdvisorIntegration(in *cloud.CloudTrustedadvisorIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	return out
}

func flattenCloudAwsVpcIntegration(in *cloud.CloudVpcIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["aws_regions"] = in.AwsRegions
	out["fetch_nat_gateway"] = in.FetchNatGateway
	out["fetch_vpn"] = in.FetchVpn
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	return out
}

func flattenCloudAwsXRayIntegration(in *cloud.CloudAwsXrayIntegration) map[string]interface{} {
	out := make(map[string]interface{})

	out["integration_id"] = in.ID
	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	return out
}

func resourceNewRelicCloudAwsIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	cloudAwsIntegrationsInput := expandCloudAwsIntegrationsInput(d)

	cloudAwsIntegrationsPayload, err := client.Cloud.CloudConfigureIntegration(accountID, cloudAwsIntegrationsInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudAwsIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudAwsIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return resourceNewRelicCloudAwsIntegrationsRead(ctx, d, meta)
}

func resourceNewRelicCloudAwsIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}

	disableIntegrationsInput := cloud.CloudDisableIntegrationsInput{
		Aws: cloud.CloudAwsDisableIntegrationsInput{
			Billing:        []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
			Cloudtrail:     []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
			Health:         []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
			Trustedadvisor: []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
			Vpc:            []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
			AwsXray:        []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}},
		},
	}

	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableIntegrationsInput)

	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudDisableIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudDisableIntegrationsPayload.Errors {
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