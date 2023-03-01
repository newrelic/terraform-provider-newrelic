package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNewRelicCloudAwsIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAwsIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAwsIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAwsIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAwsIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the account in New Relic.",
			},
			"linked_account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the linked AWS account in New Relic",
			},
			"billing": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Billing integration",
				Elem:        cloudAwsIntegrationBillingSchemaElem(),
				MaxItems:    1,
			},
			"doc_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Billing integration",
				Elem:        cloudAwsIntegrationDocDBSchemaElem(),
				MaxItems:    1,
			},
			"cloudtrail": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "CloudTrail integration",
				Elem:        cloudAwsIntegrationCloudTrailSchemaElem(),
				MaxItems:    1,
			},
			"health": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Health integration",
				Elem:        cloudAwsIntegrationHealthSchemaElem(),
				MaxItems:    1,
			},
			"trusted_advisor": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Trusted Advisor integration",
				Elem:        cloudAwsIntegrationTrustedAdvisorSchemaElem(),
				MaxItems:    1,
			},
			"s3": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "S3 integration",
				Elem:        cloudAwsIntegrationS3SchemaElem(),
				MaxItems:    1,
			},
			"vpc": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "VPC integration",
				Elem:        cloudAwsIntegrationVpcSchemaElem(),
				MaxItems:    1,
			},
			"x_ray": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "X-Ray integration",
				Elem:        cloudAwsIntegrationXRaySchemaElem(),
				MaxItems:    1,
			},
		},
	}
}

func cloudAwsIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationDocDBSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

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

func cloudAwsIntegrationS3SchemaElem() *schema.Resource {
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
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
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
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	return &schema.Resource{
		Schema: s,
	}
}

func resourceNewRelicCloudAwsIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	cloudAwsIntegrationsInput, _ := expandCloudAwsIntegrationsInput(d)

	cloudAwsIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudAwsIntegrationsInput)
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

	if len(cloudAwsIntegrationsPayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}

	return nil
}

func expandCloudAwsIntegrationsInput(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	cloudAwsIntegration := cloud.CloudAwsIntegrationsInput{}
	cloudDisableAwsIntegration := cloud.CloudAwsDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}

	if b, ok := d.GetOk("billing"); ok {
		cloudAwsIntegration.Billing = expandCloudAwsIntegrationBillingInput(b.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("billing"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.Billing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if c, ok := d.GetOk("cloudtrail"); ok {
		cloudAwsIntegration.Cloudtrail = expandCloudAwsIntegrationCloudtrailInput(c.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("cloudtrail"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if t, ok := d.GetOk("doc_db"); ok {
		cloudAwsIntegration.AwsDocdb = expandCloudAwsIntegrationDocDBInput(t.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("doc_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.AwsDocdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if h, ok := d.GetOk("health"); ok {
		cloudAwsIntegration.Health = expandCloudAwsIntegrationHealthInput(h.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("health"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.Health = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if t, ok := d.GetOk("trusted_advisor"); ok {
		cloudAwsIntegration.Trustedadvisor = expandCloudAwsIntegrationTrustedAdvisorInput(t.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("trusted_advisor"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.Trustedadvisor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if t, ok := d.GetOk("s3"); ok {
		cloudAwsIntegration.S3 = expandCloudAwsIntegrationS3Input(t.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("s3"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.S3 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("vpc"); ok {
		cloudAwsIntegration.Vpc = expandCloudAwsIntegrationVpcInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("vpc"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.Vpc = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if x, ok := d.GetOk("x_ray"); ok {
		cloudAwsIntegration.AwsXray = expandCloudAwsIntegrationXRayInput(x.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("x_ray"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsIntegration.AwsXray = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	configureInput := cloud.CloudIntegrationsInput{
		Aws: cloudAwsIntegration,
	}

	disableInput := cloud.CloudDisableIntegrationsInput{
		Aws: cloudDisableAwsIntegration,
	}

	return configureInput, disableInput
}

func expandCloudAwsIntegrationBillingInput(b []interface{}, linkedAccountID int) []cloud.CloudBillingIntegrationInput {
	expanded := make([]cloud.CloudBillingIntegrationInput, len(b))

	for i, billing := range b {
		var billingInput cloud.CloudBillingIntegrationInput

		if billing == nil {
			billingInput.LinkedAccountId = linkedAccountID
			expanded[i] = billingInput
			return expanded
		}

		in := billing.(map[string]interface{})

		billingInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			billingInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = billingInput
	}

	return expanded
}

func expandCloudAwsIntegrationCloudtrailInput(c []interface{}, linkedAccountID int) []cloud.CloudCloudtrailIntegrationInput {
	expanded := make([]cloud.CloudCloudtrailIntegrationInput, len(c))

	for i, cloudtrail := range c {
		var cloudtrailInput cloud.CloudCloudtrailIntegrationInput

		if cloudtrail == nil {
			cloudtrailInput.LinkedAccountId = linkedAccountID
			expanded[i] = cloudtrailInput
			return expanded
		}

		in := cloudtrail.(map[string]interface{})

		cloudtrailInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			cloudtrailInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			cloudtrailInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = cloudtrailInput
	}

	return expanded
}

func expandCloudAwsIntegrationHealthInput(h []interface{}, linkedAccountID int) []cloud.CloudHealthIntegrationInput {
	expanded := make([]cloud.CloudHealthIntegrationInput, len(h))

	for i, health := range h {
		var healthInput cloud.CloudHealthIntegrationInput

		if health == nil {
			healthInput.LinkedAccountId = linkedAccountID
			expanded[i] = healthInput
			return expanded
		}

		in := health.(map[string]interface{})

		healthInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			healthInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = healthInput
	}

	return expanded
}

func expandCloudAwsIntegrationDocDBInput(h []interface{}, linkedAccountID int) []cloud.CloudAwsDocdbIntegrationInput {
	expanded := make([]cloud.CloudAwsDocdbIntegrationInput, len(h))

	for i, docDb := range h {
		var docDbInput cloud.CloudAwsDocdbIntegrationInput

		if docDb == nil {
			docDbInput.LinkedAccountId = linkedAccountID
			expanded[i] = docDbInput
			return expanded
		}

		in := docDb.(map[string]interface{})

		docDbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			docDbInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = docDbInput
	}

	return expanded
}

func expandCloudAwsIntegrationTrustedAdvisorInput(t []interface{}, linkedAccountID int) []cloud.CloudTrustedadvisorIntegrationInput {
	expanded := make([]cloud.CloudTrustedadvisorIntegrationInput, len(t))

	for i, trustedAdvisor := range t {
		var trustedAdvisorInput cloud.CloudTrustedadvisorIntegrationInput

		if trustedAdvisor == nil {
			trustedAdvisorInput.LinkedAccountId = linkedAccountID
			expanded[i] = trustedAdvisorInput
			return expanded
		}

		in := trustedAdvisor.(map[string]interface{})

		trustedAdvisorInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			trustedAdvisorInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = trustedAdvisorInput
	}

	return expanded
}

func expandCloudAwsIntegrationS3Input(t []interface{}, linkedAccountID int) []cloud.CloudS3IntegrationInput {
	expanded := make([]cloud.CloudS3IntegrationInput, len(t))

	for i, s3 := range t {
		var s3Input cloud.CloudS3IntegrationInput

		if s3 == nil {
			s3Input.LinkedAccountId = linkedAccountID
			expanded[i] = s3Input
			return expanded
		}

		in := s3.(map[string]interface{})

		s3Input.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			s3Input.MetricsPollingInterval = m.(int)
		}

		expanded[i] = s3Input
	}

	return expanded
}

func expandCloudAwsIntegrationVpcInput(v []interface{}, linkedAccountID int) []cloud.CloudVpcIntegrationInput {
	expanded := make([]cloud.CloudVpcIntegrationInput, len(v))

	for i, vpc := range v {
		var vpcInput cloud.CloudVpcIntegrationInput

		if vpc == nil {
			vpcInput.LinkedAccountId = linkedAccountID
			expanded[i] = vpcInput
			return expanded
		}

		in := vpc.(map[string]interface{})

		vpcInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			vpcInput.AwsRegions = regions
		}

		if nat, ok := in["fetch_nat_gateway"]; ok {
			vpcInput.FetchNatGateway = nat.(bool)
		}

		if vpn, ok := in["fetch_vpn"]; ok {
			vpcInput.FetchVpn = vpn.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			vpcInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			vpcInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			vpcInput.TagValue = tv.(string)
		}
		expanded[i] = vpcInput
	}

	return expanded
}

func expandCloudAwsIntegrationXRayInput(x []interface{}, linkedAccountID int) []cloud.CloudAwsXrayIntegrationInput {
	expanded := make([]cloud.CloudAwsXrayIntegrationInput, len(x))

	for i, xray := range x {
		var xrayInput cloud.CloudAwsXrayIntegrationInput

		if xray == nil {
			xrayInput.LinkedAccountId = linkedAccountID
			expanded[i] = xrayInput
			return expanded
		}

		in := xray.(map[string]interface{})

		xrayInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			xrayInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			xrayInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = xrayInput
	}

	return expanded
}

func resourceNewRelicCloudAwsIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func flattenCloudAwsBillingIntegration(in *cloud.CloudBillingIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsCloudTrailIntegration(in *cloud.CloudCloudtrailIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsHealthIntegration(in *cloud.CloudHealthIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsTrustedAdvisorIntegration(in *cloud.CloudTrustedadvisorIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsVpcIntegration(in *cloud.CloudVpcIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["fetch_nat_gateway"] = in.FetchNatGateway
	out["fetch_vpn"] = in.FetchVpn
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

func flattenCloudAwsXRayIntegration(in *cloud.CloudAwsXrayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func resourceNewRelicCloudAwsIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	configureInput, disableInput := expandCloudAwsIntegrationsInput(d)

	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)

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

	cloudAwsIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, configureInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(cloudAwsIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudAwsIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

func resourceNewRelicCloudAwsIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	deleteInput := buildDeleteInput(d)

	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)

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

func buildDeleteInput(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudDisableAwsIntegration := cloud.CloudAwsDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}

	if _, ok := d.GetOk("billing"); ok {
		cloudDisableAwsIntegration.Billing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("cloudtrail"); ok {
		cloudDisableAwsIntegration.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("doc_db"); ok {
		cloudDisableAwsIntegration.AwsDocdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("health"); ok {
		cloudDisableAwsIntegration.Health = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("trusted_advisor"); ok {
		cloudDisableAwsIntegration.Trustedadvisor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("vpc"); ok {
		cloudDisableAwsIntegration.Vpc = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("s3"); ok {
		cloudDisableAwsIntegration.S3 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("x_ray"); ok {
		cloudDisableAwsIntegration.AwsXray = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	deleteInput := cloud.CloudDisableIntegrationsInput{
		Aws: cloudDisableAwsIntegration,
	}

	return deleteInput
}
