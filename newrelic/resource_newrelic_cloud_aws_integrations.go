package newrelic

import (
	"context"
	"strconv"
	"strings"

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
				ForceNew:    true,
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
				Description: "Doc DB integration",
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
			"aws_auto_discovery": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Aws Auto Discovery Integration",
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
			"dynamodb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dynamo DB integration",
				Elem:        cloudAwsIntegrationDynamoDBSchemaElem(),
				MaxItems:    1,
			},
			"ec2": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ec2 integration",
				Elem:        cloudAwsIntegrationEc2SchemaElem(),
				MaxItems:    1,
			},
			"ecs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ecs integration",
				Elem:        cloudAwsIntegrationEcsSchemaElem(),
				MaxItems:    1,
			},
			"efs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Efs integration",
				Elem:        cloudAwsIntegrationEfsSchemaElem(),
				MaxItems:    1,
			},
			"elasticbeanstalk": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Elastic Bean Stalk integration",
				Elem:        cloudAwsIntegrationElasticBeanStalkSchemaElem(),
				MaxItems:    1,
			},
			"elasticsearch": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Elastic Search integration",
				Elem:        cloudAwsIntegrationElasticSearchSchemaElem(),
				MaxItems:    1,
			},
			"elb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Elb integration",
				Elem:        cloudAwsIntegrationElbSchemaElem(),
				MaxItems:    1,
			},
			"emr": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Emr integration",
				Elem:        cloudAwsIntegrationEmrSchemaElem(),
				MaxItems:    1,
			},
			"iam": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Iam integration",
				Elem:        cloudAwsIntegrationIamSchemaElem(),
				MaxItems:    1,
			},
			"iot": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Iot integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"kinesis": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Kinesis integration",
				Elem:        cloudAwsIntegrationKinesisSchemaElem(),
				MaxItems:    1,
			},
			"kinesis_firehose": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Kinesis Firehose integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"lambda": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Lambda integration",
				Elem:        cloudAwsIntegrationLambdaSchemaElem(),
				MaxItems:    1,
			},
			"rds": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rds integration",
				Elem:        cloudAwsIntegrationRdsSchemaElem(),
				MaxItems:    1,
			},
			"redshift": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Redshift integration",
				Elem:        cloudAwsIntegrationRedshiftSchemaElem(),
				MaxItems:    1,
			},
			"route53": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Route53 integration",
				Elem:        cloudAwsIntegrationRoute53SchemaElem(),
				MaxItems:    1,
			},
			"ses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ses integration",
				Elem:        cloudAwsIntegrationCommonSchemaElem(),
				MaxItems:    1,
			},
			"sns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Sns integration",
				Elem:        cloudAwsIntegrationSnsSchemaElem(),
				MaxItems:    1,
			},
			"security_hub": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security Hub integration",
				Elem:        cloudAwsIntegrationSecurityHubSchemaElem(),
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

// function to add common schema for various services.

func cloudAwsIntegrationCommonSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	return &schema.Resource{
		Schema: s,
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

func cloudAwsIntegrationSchemaBaseExtendedOne() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"aws_regions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify each AWS region that includes the resources that you want to monitor.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"fetch_tags": {
			Type:        schema.TypeBool,
			Description: "Specify if tags and the extended inventory should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
			Optional:    true,
		},
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds.",
		},
		"tag_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		},
		"tag_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.",
		},
	}
}

// function to add schema for Billing integration

func cloudAwsIntegrationBillingSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for Cloud Trail integration

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

// function to add schema for Health integration

func cloudAwsIntegrationHealthSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for Trutsted Advisor integration

func cloudAwsIntegrationTrustedAdvisorSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for S3 integration

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

// function to add schema for VPC integration

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

// function to add schema for XRay integration

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

// function to add schema for SQS integration

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

// function to add schema for Ebs integration

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

// function to add schema for Alb integration

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

// function to add schema for Elasticache integration

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

// function to add schema for api gateway integration

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

// function to add schema for cloudfront integration

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

// function to add schema for DynamoDB integration

func cloudAwsIntegrationDynamoDBSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for Ec2 integration

func cloudAwsIntegrationEc2SchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	s["duplicate_ec2_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specify if the old legacy metadata and tag names have to be kept, it will consume more ingest data size",
	}
	s["fetch_ip_addresses"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specify if IP addresses of ec2 instance should be collected",
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

// function to add schema for ecs integration.

func cloudAwsIntegrationEcsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for efs integration.

func cloudAwsIntegrationEfsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for elastic bean stalk integration.

func cloudAwsIntegrationElasticBeanStalkSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for elastic search integration.

func cloudAwsIntegrationElasticSearchSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	s["fetch_nodes"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if metrics should be collected for nodes. Turning it on will increase the number of API calls made to CloudWatch.",
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

// function to add schema for Elb integration.

func cloudAwsIntegrationElbSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

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

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for emr integration.

func cloudAwsIntegrationEmrSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for iam integration.

func cloudAwsIntegrationIamSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

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

// function to add schema for kinesis integration.

func cloudAwsIntegrationKinesisSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	s["fetch_shards"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if Shards should be monitored. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for lambda integration.

func cloudAwsIntegrationLambdaSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for rds integration.

func cloudAwsIntegrationRdsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtendedOne()

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for redshift integration.

func cloudAwsIntegrationRedshiftSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

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

// function to add schema for route53 integration.

func cloudAwsIntegrationRoute53SchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBase()

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for sns integration.

func cloudAwsIntegrationSnsSchemaElem() *schema.Resource {
	s := cloudAwsIntegrationSchemaBaseExtended()

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
	}

	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for SecurityHub integration

func cloudAwsIntegrationSecurityHubSchemaElem() *schema.Resource {
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
