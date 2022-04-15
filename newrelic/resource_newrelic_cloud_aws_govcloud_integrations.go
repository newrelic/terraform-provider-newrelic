package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
	"strconv"
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

//function to add common schema for AwsGov cloud all resources

func AwsGovCloudIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds",
		},
	}
}

//function to add schema for application load balancer

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

//function to add schema for api gateway

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

//function to add schema for autoscaling

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

//function to add schema for aws direct connect

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

//function to add schema for aws states

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

//function to add schema for cloud trail

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

//function to add schema for dynamo DB

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

//function to add schema for elastic beanstalk

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

//function to add schema for elastic compute cloud

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

//function to add schema for elastic search

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

//function to add schema for elastic load balancing

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

//function to add schema for elastic map reduce

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

//function to add schema for identity access management

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

//function to add schema for lambda

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

//function to add schema for relational database

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

//function to add schema for redshift

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

//function to add schema for route53

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

//function to add schema for s3 bucket

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

//function to add schema for simple notification service

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

//function to add schema for simple queue service

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

func resourceNewRelicAwsGovCloudIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	awsGovCloudIntegrationsInput, _ := expandAwsGovCloudIntegrationsInput(d)

	//cloudLinkAccountWithContext func which integrates azure account with Newrelic
	//which returns payload and error

	awsGovCloudIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, awsGovCloudIntegrationsInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(awsGovCloudIntegrationsPayload.Errors) > 0 {
		for _, err := range awsGovCloudIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	if len(awsGovCloudIntegrationsPayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}

	return nil
}

func expandAwsGovCloudIntegrationsInput(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	awsGovCloudIntegration := cloud.CloudAwsGovcloudIntegrationsInput{}
	cloudDisableAwsGovCloudIntegration := cloud.CloudAwsGovcloudDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if v, ok := d.GetOk("alb"); ok {
		awsGovCloudIntegration.Alb = expandAwsGovCloudIntegrationsAlbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("alb"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Alb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	configureInput := cloud.CloudIntegrationsInput{
		AwsGovcloud: awsGovCloudIntegration,
	}

	disableInput := cloud.CloudDisableIntegrationsInput{
		AwsGovcloud: cloudDisableAwsGovCloudIntegration,
	}

	return configureInput, disableInput
}

// Expanding the alb

func expandAwsGovCloudIntegrationsAlbInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsGovCloudAlbIntegrationInput {
	expanded := make([]cloud.CloudAwsGovCloudAlbIntegrationInput, len(b))

	for i, awsGovCloudAlb := range b {
		var awsGovCloudAlbInput cloud.CloudAwsGovCloudAlbIntegrationInput

		if awsGovCloudAlb == nil {
			awsGovCloudAlbInput.LinkedAccountId = linkedAccountID
			expanded[i] = awsGovCloudAlbInput
			return expanded
		}

		in := awsGovCloudAlbInput.(map[string]interface{})

		awsGovCloudAlbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			awsGovCloudAlbInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = awsGovCloudAlbInput
	}

	return expanded
}

func resourceNewRelicAwsGovCloudIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenAwsGovCloudLinkedAccount(d, linkedAccount)

	return nil
}

/// flatten

//nolint: gocyclo
func flattenAwsGovCloudLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("linked_account_id", result.ID)

	for _, i := range result.Integrations {
		switch t := i.(type) {
		case *cloud.CloudAwsGovCloudAlbIntegration:
			_ = d.Set("alb", flattenAwsGovCloudAlbIntegration(t))
		}
	}
}
func flattenAwsGovCloudAlbIntegration(in *cloud.CloudawsGovCloudAlbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

/// update

func resourceNewRelicAwsGovCloudIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	integrateInput, disableInput := expandAwsGovCloudIntegrationsInput(d)

	awsGovCloudDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(awsGovCloudDisablePayload.Errors) > 0 {
		for _, err := range awsGovCloudDisablePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	awsGovCloudIntegrationPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, integrateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(awsGovCloudIntegrationPayload.Errors) > 0 {
		for _, err := range awsGovCloudIntegrationPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

/// Delete
func resourceNewRelicAwsGovCloudIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	deleteInput := expandAwsGovCloudDisableInputs(d)
	awsGovCloudDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(awsGovCloudDisablePayload.Errors) > 0 {
		for _, err := range awsGovCloudDisablePayload.Errors {
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

//nolint: gocyclo
func expandAwsGovCloudDisableInputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	AwsGovCloudDisableInputs := cloud.CloudAzureDisableIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("alb"); ok {
		AwsGovCloudDisableInputs.Alb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		AwsGovcloud: AwsGovCloudDisableInputs,
	}
	return deleteInput
}
