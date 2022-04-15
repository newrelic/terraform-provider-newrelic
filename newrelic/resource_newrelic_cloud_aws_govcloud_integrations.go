package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNewRelicAwsGovCloudIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsGovCloudIntegrationsCreate,
		ReadContext:   resourceNewRelicAwsGovCloudIntegrationsRead,
		UpdateContext: resourceNewRelicAwsGovCloudIntegrationsUpdate,
		DeleteContext: resourceNewRelicAwsGovCloudIntegrationsDelete,
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
				Description: "The ID of the linked AwsGovCloud account in New Relic",
			},

			// list of resources in AwsGov cloud for integrations

			"alb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The alb integration",
				Elem:        AwsGovCloudIntegrationAlbElem(),
				MaxItems:    1,
			},
			"api_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The api gateway integration",
				Elem:        AwsGovCloudIntegrationApiGatewayElem(),
				MaxItems:    1,
			},
			"auto_scaling": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The auto scaling integration",
				Elem:        AwsGovCloudIntegrationAutoScalingElem(),
				MaxItems:    1,
			},
			"aws_direct_connect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The aws direct connect integration",
				Elem:        AwsGovCloudIntegrationAwsDirectConnectElem(),
				MaxItems:    1,
			},
			"aws_states": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The aws states integration",
				Elem:        AwsGovCloudIntegrationAwsStatesElem(),
				MaxItems:    1,
			},
			"cloudtrail": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The cloudtrail integration",
				Elem:        AwsGovCloudIntegrationCloudTrailElem(),
				MaxItems:    1,
			},
			"dynamo_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The dynamo DB integration",
				Elem:        AwsGovCloudIntegrationDynamodbElem(),
				MaxItems:    1,
			},
			"ebs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The ebs integration",
				Elem:        AwsGovCloudIntegrationEbsElem(),
				MaxItems:    1,
			},
			"ec2": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The ec2 integration",
				Elem:        AwsGovCloudIntegrationEc2Elem(),
				MaxItems:    1,
			},
			"elastic_search": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The elastic search integration",
				Elem:        AwsGovCloudIntegrationElasticSearchElem(),
				MaxItems:    1,
			},
			"elb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The elb integration",
				Elem:        AwsGovCloudIntegrationElbElem(),
				MaxItems:    1,
			},
			"emr": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The emr integration",
				Elem:        AwsGovCloudIntegrationEmrElem(),
				MaxItems:    1,
			},
			"iam": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The iam integration",
				Elem:        AwsGovCloudIntegrationIamElem(),
				MaxItems:    1,
			},
			"lambda": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The lambda integration",
				Elem:        AwsGovCloudIntegrationLambdaElem(),
				MaxItems:    1,
			},
			"rds": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The rds integration",
				Elem:        AwsGovCloudIntegrationRdsElem(),
				MaxItems:    1,
			},
			"red_shift": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The redshift integration",
				Elem:        AwsGovCloudIntegrationRedshiftElem(),
				MaxItems:    1,
			},
			"route53": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The route53 integration",
				Elem:        AwsGovCloudIntegrationRoute53Elem(),
				MaxItems:    1,
			},
			"s3": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The s3 integration",
				Elem:        AwsGovCloudIntegrationS3Elem(),
				MaxItems:    1,
			},
			"sns": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The sns integration",
				Elem:        AwsGovCloudIntegrationSnsElem(),
				MaxItems:    1,
			},
			"sqs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The sqs integration",
				Elem:        AwsGovCloudIntegrationSqsElem(),
				MaxItems:    1,
			},
		},
	}
}
func AwsGovCloudIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds",
		},
	}
}

func AwsGovCloudIntegrationAlbElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}

	s["load_balancer_prefixes"] = &schema.Schema{
		Type:        schema.TypeString,
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
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

func AwsGovCloudIntegrationApiGatewayElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["stage_prefixes"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

func AwsGovCloudIntegrationAutoScalingElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

func AwsGovCloudIntegrationAwsDirectConnectElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationAwsStatesElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationCloudTrailElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationDynamodbElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationEbsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}

	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationEc2Elem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_ip_addresses"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if IP addresses of ec2 instance should be collected",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationElasticSearchElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_nodes"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if IP addresses of ec2 instance should be collected",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationElbElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationEmrElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationIamElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationLambdaElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationRdsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationRedshiftElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationRoute53Elem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationS3Elem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationSnsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
func AwsGovCloudIntegrationSqsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["fetch_tags"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify if tags should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["queue_prefixes"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Specify each name or prefix for the Queues that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
	s["tag_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag key associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	s["tag_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Specify a Tag value associated with the resources that you want to monitor. Filter values are case-sensitive.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}
