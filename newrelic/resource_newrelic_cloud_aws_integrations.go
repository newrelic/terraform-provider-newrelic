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
				Elem:        cloudAwsIntegrationAutoscalingSchemaElem(),
				MaxItems:    1,
			},
			"aws_app_sync": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Appsync integration",
				Elem:        cloudAwsIntegrationAppsyncSchemaElem(),
				MaxItems:    1,
			},
			"aws_athena": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Athena integration",
				Elem:        cloudAwsIntegrationAthenaSchemaElem(),
				MaxItems:    1,
			},
			"aws_cognito": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Cognito integration",
				Elem:        cloudAwsIntegrationCognitoSchemaElem(),
				MaxItems:    1,
			},
			"aws_connect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Connect integration",
				Elem:        cloudAwsIntegrationConnectSchemaElem(),
				MaxItems:    1,
			},
			"aws_direct_connect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Direct Connect integration",
				Elem:        cloudAwsIntegrationDirectConnectSchemaElem(),
				MaxItems:    1,
			},
			"aws_fsx": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Fsx integration",
				Elem:        cloudAwsIntegrationFsxSchemaElem(),
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

// function to add schema for autoscaling

func cloudAwsIntegrationAutoscalingSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for appsync

func cloudAwsIntegrationAppsyncSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for athena

func cloudAwsIntegrationAthenaSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for cognito

func cloudAwsIntegrationCognitoSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for connect

func cloudAwsIntegrationConnectSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for direct connect

func cloudAwsIntegrationDirectConnectSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for aws fsx

func cloudAwsIntegrationFsxSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

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

type enableDisableAwsIntegration struct {
	enableFunc  func([]interface{}, int)
	disableFunc func(int)
}

func expandCloudAwsIntegrationsInput(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	cloudAwsIntegration := cloud.CloudAwsIntegrationsInput{}
	cloudDisableAwsIntegration := cloud.CloudAwsDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}

	awsIntegrationMap := map[string]enableDisableAwsIntegration{
		"billing": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Billing = expandCloudAwsIntegrationBillingInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Billing = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"cloudtrail": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Cloudtrail = expandCloudAwsIntegrationCloudtrailInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"doc_db": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsDocdb = expandCloudAwsIntegrationDocDBInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsDocdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"health": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Health = expandCloudAwsIntegrationHealthInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Health = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"trusted_advisor": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Trustedadvisor = expandCloudAwsIntegrationTrustedAdvisorInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Trustedadvisor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"s3": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.S3 = expandCloudAwsIntegrationS3Input(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.S3 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"vpc": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Vpc = expandCloudAwsIntegrationVpcInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Vpc = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"x_ray": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsXray = expandCloudAwsIntegrationXRayInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsXray = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"sqs": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Sqs = expandCloudAwsIntegrationSqsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Sqs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"ebs": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Ebs = expandCloudAwsIntegrationEbsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Ebs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"alb": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Alb = expandCloudAwsIntegrationAlbInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Alb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"elasticache": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Elasticache = expandCloudAwsIntegrationElasticacheInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Elasticache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"api_gateway": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.APIgateway = expandCloudAwsIntegrationsAPIGatewayInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.APIgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"auto_scaling": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Autoscaling = expandCloudAwsIntegrationAutoscalingInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Autoscaling = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_app_sync": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsAppsync = expandCloudAwsIntegrationAppsyncInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsAppsync = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_athena": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsAthena = expandCloudAwsIntegrationAthenaInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsAthena = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_cognito": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsCognito = expandCloudAwsIntegrationCognitoInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsCognito = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_connect": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsConnect = expandCloudAwsIntegrationConnectInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsConnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_direct_connect": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsDirectconnect = expandCloudAwsIntegrationDirectconnectInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsDirectconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_fsx": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsFsx = expandCloudAwsIntegrationFsxInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsFsx = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
	}

	for key, fun := range awsIntegrationMap {
		if v, ok := d.GetOk(key); ok {
			fun.enableFunc(v.([]interface{}), linkedAccountID)
		} else if o, n := d.GetChange(key); len(n.([]interface{})) > len(o.([]interface{})) {
			fun.disableFunc(linkedAccountID)
		}
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

		if f, ok := in["fetch_extended_inventory"]; ok {
			s3Input.FetchExtendedInventory = f.(bool)
		}

		if f, ok := in["fetch_tags"]; ok {
			s3Input.FetchTags = f.(bool)
		}

		if tk, ok := in["tag_key"]; ok {
			s3Input.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			s3Input.TagValue = tv.(string)
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

func expandCloudAwsIntegrationSqsInput(h []interface{}, linkedAccountID int) []cloud.CloudSqsIntegrationInput {
	expanded := make([]cloud.CloudSqsIntegrationInput, len(h))

	for i, health := range h {
		var sqsInput cloud.CloudSqsIntegrationInput

		if health == nil {
			sqsInput.LinkedAccountId = linkedAccountID
			expanded[i] = sqsInput
			return expanded
		}

		in := health.(map[string]interface{})

		sqsInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			sqsInput.AwsRegions = regions
		}

		if f, ok := in["fetch_extended_inventory"]; ok {
			sqsInput.FetchExtendedInventory = f.(bool)
		}

		if f, ok := in["fetch_tags"]; ok {
			sqsInput.FetchTags = f.(bool)
		}

		if f, ok := in["queue_prefixes"]; ok {
			queuePrefixes := f.([]interface{})
			var prefixes []string

			for _, prefix := range queuePrefixes {
				prefixes = append(prefixes, prefix.(string))
			}
			sqsInput.QueuePrefixes = prefixes
		}

		if tk, ok := in["tag_key"]; ok {
			sqsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			sqsInput.TagValue = tv.(string)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			sqsInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = sqsInput
	}

	return expanded
}

func expandCloudAwsIntegrationEbsInput(h []interface{}, linkedAccountID int) []cloud.CloudEbsIntegrationInput {
	expanded := make([]cloud.CloudEbsIntegrationInput, len(h))

	for i, health := range h {
		var ebsInput cloud.CloudEbsIntegrationInput

		if health == nil {
			ebsInput.LinkedAccountId = linkedAccountID
			expanded[i] = ebsInput
			return expanded
		}

		in := health.(map[string]interface{})

		ebsInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			ebsInput.AwsRegions = regions
		}

		if f, ok := in["fetch_extended_inventory"]; ok {
			ebsInput.FetchExtendedInventory = f.(bool)
		}

		if tk, ok := in["tag_key"]; ok {
			ebsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			ebsInput.TagValue = tv.(string)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			ebsInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = ebsInput
	}

	return expanded
}

func expandCloudAwsIntegrationAlbInput(h []interface{}, linkedAccountID int) []cloud.CloudAlbIntegrationInput {
	expanded := make([]cloud.CloudAlbIntegrationInput, len(h))

	for i, health := range h {
		var albInput cloud.CloudAlbIntegrationInput

		if health == nil {
			albInput.LinkedAccountId = linkedAccountID
			expanded[i] = albInput
			return expanded
		}

		in := health.(map[string]interface{})

		albInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			albInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			albInput.MetricsPollingInterval = m.(int)
		}

		if f, ok := in["fetch_extended_inventory"]; ok {
			albInput.FetchExtendedInventory = f.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			albInput.FetchTags = ft.(bool)
		}

		if lb, ok := in["load_balancer_prefixes"]; ok {
			loadBalancerPrefixes := lb.([]interface{})
			var prefixes []string

			for _, prefix := range loadBalancerPrefixes {
				prefixes = append(prefixes, prefix.(string))
			}
			albInput.LoadBalancerPrefixes = prefixes
		}

		if tk, ok := in["tag_key"]; ok {
			albInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			albInput.TagValue = tv.(string)
		}

		expanded[i] = albInput
	}

	return expanded
}

func expandCloudAwsIntegrationElasticacheInput(h []interface{}, linkedAccountID int) []cloud.CloudElasticacheIntegrationInput {
	expanded := make([]cloud.CloudElasticacheIntegrationInput, len(h))

	for i, health := range h {
		var elasticacheInput cloud.CloudElasticacheIntegrationInput

		if health == nil {
			elasticacheInput.LinkedAccountId = linkedAccountID
			expanded[i] = elasticacheInput
			return expanded
		}

		in := health.(map[string]interface{})

		elasticacheInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elasticacheInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			elasticacheInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			elasticacheInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			elasticacheInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			elasticacheInput.TagValue = tv.(string)
		}
		expanded[i] = elasticacheInput
	}

	return expanded
}

// Expanding the api gateway
func expandCloudAwsIntegrationsAPIGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAPIgatewayIntegrationInput {
	expanded := make([]cloud.CloudAPIgatewayIntegrationInput, len(b))

	for i, apiGateway := range b {
		var apiGatewayInput cloud.CloudAPIgatewayIntegrationInput

		if apiGateway == nil {
			apiGatewayInput.LinkedAccountId = linkedAccountID
			expanded[i] = apiGatewayInput
			return expanded
		}

		in := apiGateway.(map[string]interface{})

		apiGatewayInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			apiGatewayInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			apiGatewayInput.MetricsPollingInterval = m.(int)
		}

		if sp, ok := in["stage_prefixes"]; ok {
			stagePrefixes := sp.([]interface{})
			var prefixes []string

			for _, prefix := range stagePrefixes {
				prefixes = append(prefixes, prefix.(string))
			}
			apiGatewayInput.StagePrefixes = prefixes

		}

		if tk, ok := in["tag_key"]; ok {
			apiGatewayInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			apiGatewayInput.TagValue = tv.(string)
		}

		expanded[i] = apiGatewayInput
	}

	return expanded
}

// Expanding the auto scaling
func expandCloudAwsIntegrationAutoscalingInput(b []interface{}, linkedAccountID int) []cloud.CloudAutoscalingIntegrationInput {
	expanded := make([]cloud.CloudAutoscalingIntegrationInput, len(b))

	for i, autoScaling := range b {
		var autoScalingInput cloud.CloudAutoscalingIntegrationInput

		if autoScaling == nil {
			autoScalingInput.LinkedAccountId = linkedAccountID
			expanded[i] = autoScalingInput
			return expanded
		}

		in := autoScaling.(map[string]interface{})

		autoScalingInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			autoScalingInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			autoScalingInput.AwsRegions = regions
		}

		expanded[i] = autoScalingInput
	}

	return expanded
}

// Expanding the aws app sync
func expandCloudAwsIntegrationAppsyncInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsAppsyncIntegrationInput {
	expanded := make([]cloud.CloudAwsAppsyncIntegrationInput, len(b))

	for i, appSync := range b {
		var appSyncInput cloud.CloudAwsAppsyncIntegrationInput

		if appSync == nil {
			appSyncInput.LinkedAccountId = linkedAccountID
			expanded[i] = appSyncInput
			return expanded
		}

		in := appSync.(map[string]interface{})

		appSyncInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			appSyncInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			appSyncInput.AwsRegions = regions
		}

		expanded[i] = appSyncInput
	}

	return expanded
}

// Expanding the aws athena
func expandCloudAwsIntegrationAthenaInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsAthenaIntegrationInput {
	expanded := make([]cloud.CloudAwsAthenaIntegrationInput, len(b))

	for i, athena := range b {
		var athenaInput cloud.CloudAwsAthenaIntegrationInput

		if athena == nil {
			athenaInput.LinkedAccountId = linkedAccountID
			expanded[i] = athenaInput
			return expanded
		}

		in := athena.(map[string]interface{})

		athenaInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			athenaInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			athenaInput.AwsRegions = regions
		}

		expanded[i] = athenaInput
	}

	return expanded
}

// Expanding the aws cognito
func expandCloudAwsIntegrationCognitoInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsCognitoIntegrationInput {
	expanded := make([]cloud.CloudAwsCognitoIntegrationInput, len(b))

	for i, cognito := range b {
		var cognitoInput cloud.CloudAwsCognitoIntegrationInput

		if cognito == nil {
			cognitoInput.LinkedAccountId = linkedAccountID
			expanded[i] = cognitoInput
			return expanded
		}

		in := cognito.(map[string]interface{})

		cognitoInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			cognitoInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			cognitoInput.AwsRegions = regions
		}

		expanded[i] = cognitoInput
	}

	return expanded
}

// Expanding the aws connect
func expandCloudAwsIntegrationConnectInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsConnectIntegrationInput {
	expanded := make([]cloud.CloudAwsConnectIntegrationInput, len(b))

	for i, connect := range b {
		var connectInput cloud.CloudAwsConnectIntegrationInput

		if connect == nil {
			connectInput.LinkedAccountId = linkedAccountID
			expanded[i] = connectInput
			return expanded
		}

		in := connect.(map[string]interface{})

		connectInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			connectInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			connectInput.AwsRegions = regions
		}

		expanded[i] = connectInput
	}

	return expanded
}

// Expanding the aws direct connect
func expandCloudAwsIntegrationDirectconnectInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsDirectconnectIntegrationInput {
	expanded := make([]cloud.CloudAwsDirectconnectIntegrationInput, len(b))

	for i, directConnect := range b {
		var directConnectInput cloud.CloudAwsDirectconnectIntegrationInput

		if directConnect == nil {
			directConnectInput.LinkedAccountId = linkedAccountID
			expanded[i] = directConnectInput
			return expanded
		}

		in := directConnect.(map[string]interface{})

		directConnectInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			directConnectInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			directConnectInput.AwsRegions = regions
		}

		expanded[i] = directConnectInput
	}

	return expanded
}

// Expanding the aws fsx
func expandCloudAwsIntegrationFsxInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsFsxIntegrationInput {
	expanded := make([]cloud.CloudAwsFsxIntegrationInput, len(b))

	for i, fsx := range b {
		var fsxInput cloud.CloudAwsFsxIntegrationInput

		if fsx == nil {
			fsxInput.LinkedAccountId = linkedAccountID
			expanded[i] = fsxInput
			return expanded
		}

		in := fsx.(map[string]interface{})

		fsxInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			fsxInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			fsxInput.AwsRegions = regions
		}

		expanded[i] = fsxInput
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

// TODO: Reduce the cyclomatic complexity of this func
// nolint: gocyclo
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
		case *cloud.CloudSqsIntegration:
			_ = d.Set("sqs", flattenCloudAwsSqsIntegration(t))
		case *cloud.CloudEbsIntegration:
			_ = d.Set("ebs", flattenCloudAwsEbsIntegration(t))
		case *cloud.CloudAlbIntegration:
			_ = d.Set("alb", flattenCloudAwsAlbIntegration(t))
		case *cloud.CloudElasticacheIntegration:
			_ = d.Set("elasticache", flattenCloudAwsElasticacheIntegration(t))
		case *cloud.CloudS3Integration:
			_ = d.Set("s3", flattenCloudAwsS3Integration(t))
		case *cloud.CloudAwsDocdbIntegration:
			_ = d.Set("doc_db", flattenCloudAwsDocDBIntegration(t))
		case *cloud.CloudAPIgatewayIntegration:
			_ = d.Set("api_gateway", flattenCloudAwsAPIGatewayIntegration(t))
		case *cloud.CloudAutoscalingIntegration:
			_ = d.Set("auto_scaling", flattenCloudAwsAutoScalingIntegration(t))
		case *cloud.CloudAwsAppsyncIntegration:
			_ = d.Set("aws_app_sync", flattenCloudAwsAppsyncIntegration(t))
		case *cloud.CloudAwsAthenaIntegration:
			_ = d.Set("aws_athena", flattenCloudAwsAthenaIntegration(t))
		case *cloud.CloudAwsCognitoIntegration:
			_ = d.Set("aws_cognito", flattenCloudAwsCognitoIntegration(t))
		case *cloud.CloudAwsConnectIntegration:
			_ = d.Set("aws_connect", flattenCloudAwsConnectIntegration(t))
		case *cloud.CloudAwsDirectconnectIntegration:
			_ = d.Set("aws_direct_connect", flattenCloudAwsDirectconnectIntegration(t))
		case *cloud.CloudAwsFsxIntegration:
			_ = d.Set("aws_fsx", flattenCloudAwsFsxIntegration(t))
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

func flattenCloudAwsSqsIntegration(in *cloud.CloudSqsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["queue_prefixes"] = in.QueuePrefixes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsEbsIntegration(in *cloud.CloudEbsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

func flattenCloudAwsAlbIntegration(in *cloud.CloudAlbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["load_balancer_prefixes"] = in.LoadBalancerPrefixes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

func flattenCloudAwsElasticacheIntegration(in *cloud.CloudElasticacheIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["aws_regions"] = in.AwsRegions
	out["fetch_tags"] = in.FetchTags
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

func flattenCloudAwsS3Integration(in *cloud.CloudS3Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

func flattenCloudAwsDocDBIntegration(in *cloud.CloudAwsDocdbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for api gateway
func flattenCloudAwsAPIGatewayIntegration(in *cloud.CloudAPIgatewayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["stage_prefixes"] = in.StagePrefixes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for autoscaling
func flattenCloudAwsAutoScalingIntegration(in *cloud.CloudAutoscalingIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for appsync
func flattenCloudAwsAppsyncIntegration(in *cloud.CloudAwsAppsyncIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for athena
func flattenCloudAwsAthenaIntegration(in *cloud.CloudAwsAthenaIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for cognito
func flattenCloudAwsCognitoIntegration(in *cloud.CloudAwsCognitoIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for connect
func flattenCloudAwsConnectIntegration(in *cloud.CloudAwsConnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for direct connect
func flattenCloudAwsDirectconnectIntegration(in *cloud.CloudAwsDirectconnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for fsx
func flattenCloudAwsFsxIntegration(in *cloud.CloudAwsFsxIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

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

// TODO: Reduce the cyclomatic complexity of this func
// nolint: gocyclo
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

	if _, ok := d.GetOk("sqs"); ok {
		cloudDisableAwsIntegration.Sqs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("ebs"); ok {
		cloudDisableAwsIntegration.Ebs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("alb"); ok {
		cloudDisableAwsIntegration.Alb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("elasticache"); ok {
		cloudDisableAwsIntegration.Elasticache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("api_gateway"); ok {
		cloudDisableAwsIntegration.APIgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("auto_scaling"); ok {
		cloudDisableAwsIntegration.Autoscaling = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_app_sync"); ok {
		cloudDisableAwsIntegration.AwsAppsync = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_athena"); ok {
		cloudDisableAwsIntegration.AwsAthena = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_cognito"); ok {
		cloudDisableAwsIntegration.AwsCognito = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_connect"); ok {
		cloudDisableAwsIntegration.AwsConnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_direct_connect"); ok {
		cloudDisableAwsIntegration.AwsDirectconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_fsx"); ok {
		cloudDisableAwsIntegration.AwsFsx = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	deleteInput := cloud.CloudDisableIntegrationsInput{
		Aws: cloudDisableAwsIntegration,
	}

	return deleteInput
}
