package newrelic

import (
	"context"
	"strconv"

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
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
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
			"sqs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SQS integration",
				Elem:        cloudAwsIntegrationSqsSchemaElem(),
				MaxItems:    1,
			},
			"ebs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "EBS integration",
				Elem:        cloudAwsIntegrationEbsSchemaElem(),
				MaxItems:    1,
			},
			"alb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "ALB integration",
				Elem:        cloudAwsIntegrationAlbSchemaElem(),
				MaxItems:    1,
			},
			"elasticache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Elasticache integration",
				Elem:        cloudAwsIntegrationElasticacheSchemaElem(),
				MaxItems:    1,
			},
			"api_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "API Gateway integration",
				Elem:        cloudAwsIntegrationAPIGatewaySchemaElem(),
				MaxItems:    1,
			},
			"auto_scaling": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AutoScaling integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_app_sync": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Appsync integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_athena": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Athena integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_cognito": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Cognito integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_connect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Connect integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_direct_connect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Direct Connect integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_fsx": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Fsx integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_glue": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Glue integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_kinesis_analytics": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Kinesis Analytics integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_media_convert": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Media Convert integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_media_package_vod": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Media PackageVod integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_mq": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Mq integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_msk": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Msk integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_neptune": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Neptune integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_qldb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Qldb integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_route53resolver": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Route53resolver integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_states": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws states integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_transit_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Transit Gateway integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_waf": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Waf integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"aws_wafv2": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Wafv2 integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"cloudfront": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cloudfront integration",
				Elem:        cloudAwsIntegrationCloudfrontSchemaElem(),
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

func cloudAwsIntegrationSchemaBaseExtended() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"aws_regions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify each AWS region that includes the resources that you want to monitor.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
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

func cloudAwsIntegrationS3SchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
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

func cloudAwsIntegrationSqsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["queue_prefixes"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specify each name or prefix for the Queues that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
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

func cloudAwsIntegrationEbsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
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

func cloudAwsIntegrationAlbSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["load_balancer_prefixes"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specify each name or prefix for the LBs that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAwsIntegrationElasticacheSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each AWS region that includes the resources that you want to monitor.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for api gateway

func cloudAwsIntegrationAPIGatewaySchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["stage_prefixes"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
	}
	return &schema.Resource{
		Schema: s,
	}
}

// function to add common schema for various services.

func cloudAwsIntegrationCommonSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for cloudfront

func cloudAwsIntegrationCloudfrontSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["fetch_lambdas_at_edge"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if Lambdas@Edge should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
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
