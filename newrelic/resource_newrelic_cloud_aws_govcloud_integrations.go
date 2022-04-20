package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
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
				Elem:        AwsGovCloudIntegrationAPIGatewayElem(),
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
		Type:        schema.TypeList,
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

//function to add schema for api gateway

func AwsGovCloudIntegrationAPIGatewayElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
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

//function to add schema for autoscaling

func AwsGovCloudIntegrationAutoScalingElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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
		Type:        schema.TypeList,
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
		Type:        schema.TypeList,
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
		Type:        schema.TypeList,
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
		Type:        schema.TypeList,
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

//function to add schema for elastic beanstalk

func AwsGovCloudIntegrationEbsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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

//function to add schema for elastic compute cloud

func AwsGovCloudIntegrationEc2Elem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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
		Type:        schema.TypeList,
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

//function to add schema for elastic load balancing

func AwsGovCloudIntegrationElbElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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

//function to add schema for lambda

func AwsGovCloudIntegrationLambdaElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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

//function to add schema for relational database

func AwsGovCloudIntegrationRdsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specify each AWS region that includes the resources that you want to monitor",
		Optional:    true,
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

//function to add schema for redshift

func AwsGovCloudIntegrationRedshiftElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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

//function to add schema for route53

func AwsGovCloudIntegrationRoute53Elem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["fetch_extended_inventory"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Determine if extra inventory data be collected or not. May affect total data collection time and contribute to the Cloud provider API rate limit.",
		Optional:    true,
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

//function to add schema for simple notification service

func AwsGovCloudIntegrationSnsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()
	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for simple queue service

func AwsGovCloudIntegrationSqsElem() *schema.Resource {
	s := AwsGovCloudIntegrationSchemaBase()

	s["aws_regions"] = &schema.Schema{
		Type:        schema.TypeList,
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

func resourceNewRelicAwsGovCloudIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	awsGovCloudIntegrationsInput, _ := expandAwsGovCloudIntegrationsInput(d)

	//cloudLinkAccountWithContext func which integrates aws gov cloud account with Newrelic
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

// TODO: Reduce the cyclomatic complexity of this func
//nolint: gocyclo
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
	if v, ok := d.GetOk("api_gateway"); ok {
		awsGovCloudIntegration.APIgateway = expandAwsGovCloudIntegrationsAPIGatewayInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("api_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.APIgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("auto_scaling"); ok {
		awsGovCloudIntegration.Autoscaling = expandAwsGovCloudIntegrationsAutoScalingInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("auto_scaling"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Autoscaling = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("aws_direct_connect"); ok {
		awsGovCloudIntegration.AwsDirectconnect = expandAwsGovCloudIntegrationsAwsDirectConnectInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("aws_direct_connect"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.AwsDirectconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("aws_states"); ok {
		awsGovCloudIntegration.AwsStates = expandAwsGovCloudIntegrationsAwsStatesInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("aws_states"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.AwsStates = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("cloudtrail"); ok {
		awsGovCloudIntegration.Cloudtrail = expandAwsGovCloudIntegrationsCloudtrailInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("cloudtrail"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("dynamo_db"); ok {
		awsGovCloudIntegration.Dynamodb = expandAwsGovCloudIntegrationsDynamodbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("dynamo_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Dynamodb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("ebs"); ok {
		awsGovCloudIntegration.Ebs = expandAwsGovCloudIntegrationsEbsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("ebs"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Ebs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("ec2"); ok {
		awsGovCloudIntegration.Ec2 = expandAwsGovCloudIntegrationsEc2Input(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("ec2"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Ec2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("elastic_search"); ok {
		awsGovCloudIntegration.Elasticsearch = expandAwsGovCloudIntegrationsElasticsearchInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("elastic_search"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Elasticsearch = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("elb"); ok {
		awsGovCloudIntegration.Elb = expandAwsGovCloudIntegrationsElbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("elb"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Elb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("emr"); ok {
		awsGovCloudIntegration.Emr = expandAwsGovCloudIntegrationsEmrInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("emr"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Emr = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("iam"); ok {
		awsGovCloudIntegration.Iam = expandAwsGovCloudIntegrationsIamInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("iam"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Iam = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("lambda"); ok {
		awsGovCloudIntegration.Lambda = expandAwsGovCloudIntegrationsLambdaInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("lambda"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Lambda = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("rds"); ok {
		awsGovCloudIntegration.Rds = expandAwsGovCloudIntegrationsRdsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("rds"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Rds = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("red_shift"); ok {
		awsGovCloudIntegration.Redshift = expandAwsGovCloudIntegrationsRedshiftInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("red_shift"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Redshift = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("route53"); ok {
		awsGovCloudIntegration.Route53 = expandAwsGovCloudIntegrationsRoute53Input(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("route53"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Route53 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("s3"); ok {
		awsGovCloudIntegration.S3 = expandAwsGovCloudIntegrationsS3Input(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("s3"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.S3 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("sns"); ok {
		awsGovCloudIntegration.Sns = expandAwsGovCloudIntegrationsSnsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("sns"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Sns = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("sqs"); ok {
		awsGovCloudIntegration.Sqs = expandAwsGovCloudIntegrationsSqsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("sqs"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAwsGovCloudIntegration.Sqs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	//
	configureInput := cloud.CloudIntegrationsInput{
		AwsGovcloud: awsGovCloudIntegration,
	}

	disableInput := cloud.CloudDisableIntegrationsInput{
		AwsGovcloud: cloudDisableAwsGovCloudIntegration,
	}

	return configureInput, disableInput
}

///pending
// Expanding the alb

func expandAwsGovCloudIntegrationsAlbInput(b []interface{}, linkedAccountID int) []cloud.CloudAlbIntegrationInput {
	expanded := make([]cloud.CloudAlbIntegrationInput, len(b))

	for i, alb := range b {
		var albInput cloud.CloudAlbIntegrationInput

		if alb == nil {
			albInput.LinkedAccountId = linkedAccountID
			expanded[i] = albInput
			return expanded
		}

		in := alb.(map[string]interface{})

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
			albInput.LoadBalancerPrefixes = lb.([]string)
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

func expandAwsGovCloudIntegrationsAPIGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAPIgatewayIntegrationInput {
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

		if ft, ok := in["stage_prefixes"]; ok {
			apiGatewayInput.StagePrefixes = ft.([]string)

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

func expandAwsGovCloudIntegrationsAutoScalingInput(b []interface{}, linkedAccountID int) []cloud.CloudAutoscalingIntegrationInput {
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

func expandAwsGovCloudIntegrationsAwsDirectConnectInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsDirectconnectIntegrationInput {
	expanded := make([]cloud.CloudAwsDirectconnectIntegrationInput, len(b))

	for i, awsDirectConnect := range b {
		var awsDirectConnectInput cloud.CloudAwsDirectconnectIntegrationInput

		if awsDirectConnect == nil {
			awsDirectConnectInput.LinkedAccountId = linkedAccountID
			expanded[i] = awsDirectConnectInput
			return expanded
		}

		in := awsDirectConnect.(map[string]interface{})

		awsDirectConnectInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			awsDirectConnectInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			awsDirectConnectInput.AwsRegions = regions
		}

		expanded[i] = awsDirectConnectInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsAwsStatesInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsStatesIntegrationInput {
	expanded := make([]cloud.CloudAwsStatesIntegrationInput, len(b))

	for i, awsStates := range b {
		var awsStatesInput cloud.CloudAwsStatesIntegrationInput

		if awsStates == nil {
			awsStatesInput.LinkedAccountId = linkedAccountID
			expanded[i] = awsStatesInput
			return expanded
		}

		in := awsStates.(map[string]interface{})

		awsStatesInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			awsStatesInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			awsStatesInput.AwsRegions = regions
		}

		expanded[i] = awsStatesInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsCloudtrailInput(b []interface{}, linkedAccountID int) []cloud.CloudCloudtrailIntegrationInput {
	expanded := make([]cloud.CloudCloudtrailIntegrationInput, len(b))

	for i, cloudtrail := range b {
		var cloudtrailInput cloud.CloudCloudtrailIntegrationInput

		if cloudtrail == nil {
			cloudtrailInput.LinkedAccountId = linkedAccountID
			expanded[i] = cloudtrailInput
			return expanded
		}

		in := cloudtrail.(map[string]interface{})

		cloudtrailInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			cloudtrailInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
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

func expandAwsGovCloudIntegrationsDynamodbInput(b []interface{}, linkedAccountID int) []cloud.CloudDynamodbIntegrationInput {
	expanded := make([]cloud.CloudDynamodbIntegrationInput, len(b))

	for i, dynamodb := range b {
		var dynamodbInput cloud.CloudDynamodbIntegrationInput

		if dynamodb == nil {
			dynamodbInput.LinkedAccountId = linkedAccountID
			expanded[i] = dynamodbInput
			return expanded
		}

		in := dynamodb.(map[string]interface{})

		dynamodbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			dynamodbInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			dynamodbInput.AwsRegions = regions
		}
		if f, ok := in["fetch_extended_inventory"]; ok {
			dynamodbInput.FetchExtendedInventory = f.(bool)
		}
		if ft, ok := in["fetch_tags"]; ok {
			dynamodbInput.FetchTags = ft.(bool)
		}

		if tk, ok := in["tag_key"]; ok {
			dynamodbInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			dynamodbInput.TagValue = tv.(string)
		}

		expanded[i] = dynamodbInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsEbsInput(b []interface{}, linkedAccountID int) []cloud.CloudEbsIntegrationInput {
	expanded := make([]cloud.CloudEbsIntegrationInput, len(b))

	for i, ebs := range b {
		var ebsInput cloud.CloudEbsIntegrationInput

		if ebs == nil {
			ebsInput.LinkedAccountId = linkedAccountID
			expanded[i] = ebsInput
			return expanded
		}

		in := ebs.(map[string]interface{})

		ebsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			ebsInput.MetricsPollingInterval = m.(int)
		}
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

		expanded[i] = ebsInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsEc2Input(b []interface{}, linkedAccountID int) []cloud.CloudEc2IntegrationInput {
	expanded := make([]cloud.CloudEc2IntegrationInput, len(b))

	for i, ec2 := range b {
		var ec2Input cloud.CloudEc2IntegrationInput

		if ec2 == nil {
			ec2Input.LinkedAccountId = linkedAccountID
			expanded[i] = ec2Input
			return expanded
		}

		in := ec2.(map[string]interface{})

		ec2Input.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			ec2Input.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			ec2Input.AwsRegions = regions
		}
		if f, ok := in["fetch_ip_addresses"]; ok {
			ec2Input.FetchIpAddresses = f.(bool)
		}

		if tk, ok := in["tag_key"]; ok {
			ec2Input.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			ec2Input.TagValue = tv.(string)
		}

		expanded[i] = ec2Input
	}

	return expanded
}

func expandAwsGovCloudIntegrationsElasticsearchInput(b []interface{}, linkedAccountID int) []cloud.CloudElasticsearchIntegrationInput {
	expanded := make([]cloud.CloudElasticsearchIntegrationInput, len(b))

	for i, elasticsearch := range b {
		var elasticsearchInput cloud.CloudElasticsearchIntegrationInput

		if elasticsearch == nil {
			elasticsearchInput.LinkedAccountId = linkedAccountID
			expanded[i] = elasticsearchInput
			return expanded
		}

		in := elasticsearch.(map[string]interface{})

		elasticsearchInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			elasticsearchInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elasticsearchInput.AwsRegions = regions
		}
		if f, ok := in["fetch_nodes"]; ok {
			elasticsearchInput.FetchNodes = f.(bool)
		}

		if tk, ok := in["tag_key"]; ok {
			elasticsearchInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			elasticsearchInput.TagValue = tv.(string)
		}
		expanded[i] = elasticsearchInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsElbInput(b []interface{}, linkedAccountID int) []cloud.CloudElbIntegrationInput {
	expanded := make([]cloud.CloudElbIntegrationInput, len(b))

	for i, elb := range b {
		var elbInput cloud.CloudElbIntegrationInput

		if elb == nil {
			elbInput.LinkedAccountId = linkedAccountID
			expanded[i] = elbInput
			return expanded
		}

		in := elb.(map[string]interface{})

		elbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			elbInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elbInput.AwsRegions = regions
		}
		if f, ok := in["fetch_extended_inventory"]; ok {
			elbInput.FetchExtendedInventory = f.(bool)
		}
		if f, ok := in["fetch_tags"]; ok {
			elbInput.FetchTags = f.(bool)
		}
		expanded[i] = elbInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsEmrInput(b []interface{}, linkedAccountID int) []cloud.CloudEmrIntegrationInput {
	expanded := make([]cloud.CloudEmrIntegrationInput, len(b))

	for i, emr := range b {
		var emrInput cloud.CloudEmrIntegrationInput

		if emr == nil {
			emrInput.LinkedAccountId = linkedAccountID
			expanded[i] = emrInput
			return expanded
		}

		in := emr.(map[string]interface{})

		emrInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			emrInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			emrInput.AwsRegions = regions
		}

		if f, ok := in["fetch_tags"]; ok {
			emrInput.FetchTags = f.(bool)
		}
		if tk, ok := in["tag_key"]; ok {
			emrInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			emrInput.TagValue = tv.(string)
		}
		expanded[i] = emrInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsIamInput(b []interface{}, linkedAccountID int) []cloud.CloudIamIntegrationInput {
	expanded := make([]cloud.CloudIamIntegrationInput, len(b))

	for i, iam := range b {
		var iamInput cloud.CloudIamIntegrationInput

		if iam == nil {
			iamInput.LinkedAccountId = linkedAccountID
			expanded[i] = iamInput
			return expanded
		}

		in := iam.(map[string]interface{})

		iamInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			iamInput.MetricsPollingInterval = m.(int)
		}
		if tk, ok := in["tag_key"]; ok {
			iamInput.LinkedAccountId = linkedAccountID
			iamInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			iamInput.TagValue = tv.(string)
		}
		expanded[i] = iamInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsLambdaInput(b []interface{}, linkedAccountID int) []cloud.CloudLambdaIntegrationInput {
	expanded := make([]cloud.CloudLambdaIntegrationInput, len(b))

	for i, lambda := range b {
		var lambdaInput cloud.CloudLambdaIntegrationInput

		if lambda == nil {
			lambdaInput.LinkedAccountId = linkedAccountID
			expanded[i] = lambdaInput
			return expanded
		}

		in := lambda.(map[string]interface{})

		lambdaInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			lambdaInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			lambdaInput.AwsRegions = regions
		}
		if f, ok := in["fetch_tags"]; ok {
			lambdaInput.FetchTags = f.(bool)
		}
		if tk, ok := in["tag_key"]; ok {
			lambdaInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			lambdaInput.TagValue = tv.(string)
		}

		expanded[i] = lambdaInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsRdsInput(b []interface{}, linkedAccountID int) []cloud.CloudRdsIntegrationInput {
	expanded := make([]cloud.CloudRdsIntegrationInput, len(b))

	for i, rds := range b {
		var rdsInput cloud.CloudRdsIntegrationInput

		if rds == nil {
			rdsInput.LinkedAccountId = linkedAccountID
			expanded[i] = rdsInput
			return expanded
		}

		in := rds.(map[string]interface{})

		rdsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			rdsInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			rdsInput.AwsRegions = regions
		}
		if f, ok := in["fetch_tags"]; ok {
			rdsInput.FetchTags = f.(bool)
		}
		if tk, ok := in["tag_key"]; ok {
			rdsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			rdsInput.TagValue = tv.(string)
		}
		expanded[i] = rdsInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsRedshiftInput(b []interface{}, linkedAccountID int) []cloud.CloudRedshiftIntegrationInput {
	expanded := make([]cloud.CloudRedshiftIntegrationInput, len(b))

	for i, redshift := range b {
		var redshiftInput cloud.CloudRedshiftIntegrationInput

		if redshift == nil {
			redshiftInput.LinkedAccountId = linkedAccountID
			expanded[i] = redshiftInput
			return expanded
		}

		in := redshift.(map[string]interface{})

		redshiftInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			redshiftInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			redshiftInput.AwsRegions = regions
		}
		if tk, ok := in["tag_key"]; ok {
			redshiftInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			redshiftInput.TagValue = tv.(string)
		}
		expanded[i] = redshiftInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsRoute53Input(b []interface{}, linkedAccountID int) []cloud.CloudRoute53IntegrationInput {
	expanded := make([]cloud.CloudRoute53IntegrationInput, len(b))

	for i, route53 := range b {
		var route53Input cloud.CloudRoute53IntegrationInput

		if route53 == nil {
			route53Input.LinkedAccountId = linkedAccountID
			expanded[i] = route53Input
			return expanded
		}

		in := route53.(map[string]interface{})

		route53Input.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			route53Input.MetricsPollingInterval = m.(int)
		}
		if f, ok := in["fetch_extended_inventory"]; ok {
			route53Input.FetchExtendedInventory = f.(bool)
		}
		expanded[i] = route53Input
	}

	return expanded
}

func expandAwsGovCloudIntegrationsS3Input(b []interface{}, linkedAccountID int) []cloud.CloudS3IntegrationInput {
	expanded := make([]cloud.CloudS3IntegrationInput, len(b))

	for i, s3 := range b {
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

func expandAwsGovCloudIntegrationsSnsInput(b []interface{}, linkedAccountID int) []cloud.CloudSnsIntegrationInput {
	expanded := make([]cloud.CloudSnsIntegrationInput, len(b))

	for i, sns := range b {
		var snsInput cloud.CloudSnsIntegrationInput

		if sns == nil {
			snsInput.LinkedAccountId = linkedAccountID
			expanded[i] = snsInput
			return expanded
		}

		in := sns.(map[string]interface{})

		snsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			snsInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			snsInput.AwsRegions = regions
		}
		if f, ok := in["fetch_extended_inventory"]; ok {
			snsInput.FetchExtendedInventory = f.(bool)
		}
		expanded[i] = snsInput
	}

	return expanded
}

func expandAwsGovCloudIntegrationsSqsInput(b []interface{}, linkedAccountID int) []cloud.CloudSqsIntegrationInput {
	expanded := make([]cloud.CloudSqsIntegrationInput, len(b))

	for i, sqs := range b {
		var sqsInput cloud.CloudSqsIntegrationInput

		if sqs == nil {
			sqsInput.LinkedAccountId = linkedAccountID
			expanded[i] = sqsInput
			return expanded
		}

		in := sqs.(map[string]interface{})

		sqsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			sqsInput.MetricsPollingInterval = m.(int)
		}
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
			sqsInput.QueuePrefixes = f.([]string)
		}
		if tk, ok := in["tag_key"]; ok {
			sqsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			sqsInput.TagValue = tv.(string)
		}
		expanded[i] = sqsInput
	}

	return expanded
}

// Read

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
		case *cloud.CloudAlbIntegration:
			_ = d.Set("alb", flattenAwsGovCloudAlbIntegration(t))
		case *cloud.CloudAPIgatewayIntegration:
			_ = d.Set("api_gateway", flattenAwsGovCloudAPIGatewayIntegration(t))
		case *cloud.CloudAutoscalingIntegration:
			_ = d.Set("auto_scaling", flattenAwsGovCloudAutoScalingIntegration(t))
		case *cloud.CloudAwsDirectconnectIntegration:
			_ = d.Set("aws_direct_connect", flattenAwsGovCloudDirectconnectIntegration(t))
		case *cloud.CloudAwsStatesIntegration:
			_ = d.Set("aws_states", flattenAwsGovCloudAwsStatesIntegration(t))
		case *cloud.CloudCloudtrailIntegration:
			_ = d.Set("cloudtrail", flattenAwsGovCloudCloudtrailIntegration(t))
		case *cloud.CloudDynamodbIntegration:
			_ = d.Set("dynamo_db", flattenAwsGovCloudDynamodbIntegration(t))
		case *cloud.CloudEbsIntegration:
			_ = d.Set("ebs", flattenAwsGovCloudEbsIntegration(t))
		case *cloud.CloudEc2Integration:
			_ = d.Set("ec2", flattenAwsGovCloudEc2Integration(t))
		case *cloud.CloudElasticsearchIntegration:
			_ = d.Set("elastic_search", flattenAwsGovCloudElasticsearchIntegration(t))
		case *cloud.CloudElbIntegration:
			_ = d.Set("elb", flattenAwsGovCloudElbIntegration(t))
		case *cloud.CloudEmrIntegration:
			_ = d.Set("emr", flattenAwsGovCloudEmrIntegration(t))
		case *cloud.CloudIamIntegration:
			_ = d.Set("iam", flattenAwsGovCloudIamIntegration(t))
		case *cloud.CloudLambdaIntegration:
			_ = d.Set("lambda", flattenAwsGovCloudLambdaIntegration(t))
		case *cloud.CloudRdsIntegration:
			_ = d.Set("rds", flattenAwsGovCloudRdsIntegration(t))
		case *cloud.CloudRedshiftIntegration:
			_ = d.Set("red_shift", flattenAwsGovCloudRedshiftIntegration(t))
		case *cloud.CloudRoute53Integration:
			_ = d.Set("route53", flattenAwsGovCloudRoute53Integration(t))
		case *cloud.CloudS3Integration:
			_ = d.Set("s3", flattenAwsGovCloudS3Integration(t))
		case *cloud.CloudSnsIntegration:
			_ = d.Set("sns", flattenAwsGovCloudSnsIntegration(t))
		case *cloud.CloudSqsIntegration:
			_ = d.Set("sqs", flattenAwsGovCloudSqsIntegration(t))

		}
	}
}

//flatten for alb

func flattenAwsGovCloudAlbIntegration(in *cloud.CloudAlbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["load_balancer_prefixes"] = in.LoadBalancerPrefixes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for api gateway

func flattenAwsGovCloudAPIGatewayIntegration(in *cloud.CloudAPIgatewayIntegration) []interface{} {
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

//flatten for autoscaling

func flattenAwsGovCloudAutoScalingIntegration(in *cloud.CloudAutoscalingIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

//flatten for direct connect

func flattenAwsGovCloudDirectconnectIntegration(in *cloud.CloudAwsDirectconnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

//flatten for aws states

func flattenAwsGovCloudAwsStatesIntegration(in *cloud.CloudAwsStatesIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

//flatten for cloudtrail

func flattenAwsGovCloudCloudtrailIntegration(in *cloud.CloudCloudtrailIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

//flatten for dynamo db

func flattenAwsGovCloudDynamodbIntegration(in *cloud.CloudDynamodbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for ebs

func flattenAwsGovCloudEbsIntegration(in *cloud.CloudEbsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for ec2

func flattenAwsGovCloudEc2Integration(in *cloud.CloudEc2Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_ip_addresses"] = in.FetchIpAddresses
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for elastic search

func flattenAwsGovCloudElasticsearchIntegration(in *cloud.CloudElasticsearchIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_nodes"] = in.FetchNodes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for elb

func flattenAwsGovCloudElbIntegration(in *cloud.CloudElbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags

	flattened[0] = out

	return flattened
}

//flatten for emr

func flattenAwsGovCloudEmrIntegration(in *cloud.CloudEmrIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for iam

func flattenAwsGovCloudIamIntegration(in *cloud.CloudIamIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue
	flattened[0] = out

	return flattened
}

//flatten for lambda

func flattenAwsGovCloudLambdaIntegration(in *cloud.CloudLambdaIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue
	flattened[0] = out

	return flattened
}

//flatten for rds

func flattenAwsGovCloudRdsIntegration(in *cloud.CloudRdsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for redshift

func flattenAwsGovCloudRedshiftIntegration(in *cloud.CloudRedshiftIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

//flatten for route53

func flattenAwsGovCloudRoute53Integration(in *cloud.CloudRoute53Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory

	flattened[0] = out

	return flattened
}

// flatten for s3
func flattenAwsGovCloudS3Integration(in *cloud.CloudS3Integration) []interface{} {
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

// flatten for sns
func flattenAwsGovCloudSnsIntegration(in *cloud.CloudSnsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory

	flattened[0] = out

	return flattened
}

// flatten for sqs

func flattenAwsGovCloudSqsIntegration(in *cloud.CloudSqsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["queue_prefixes"] = in.QueuePrefixes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

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
	awsGovCloudDisableInputs := cloud.CloudAwsGovcloudDisableIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("alb"); ok {
		awsGovCloudDisableInputs.Alb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("api_gateway"); ok {
		awsGovCloudDisableInputs.APIgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("auto_scaling"); ok {
		awsGovCloudDisableInputs.Autoscaling = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("aws_direct_connect"); ok {
		awsGovCloudDisableInputs.AwsDirectconnect = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("aws_states"); ok {
		awsGovCloudDisableInputs.AwsStates = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("cloudtrail"); ok {
		awsGovCloudDisableInputs.Cloudtrail = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("dynamo_db"); ok {
		awsGovCloudDisableInputs.Dynamodb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("ebs"); ok {
		awsGovCloudDisableInputs.Ebs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("ec2"); ok {
		awsGovCloudDisableInputs.Ec2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("elastic_search"); ok {
		awsGovCloudDisableInputs.Elasticsearch = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("elb"); ok {
		awsGovCloudDisableInputs.Elb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("emr"); ok {
		awsGovCloudDisableInputs.Emr = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("iam"); ok {
		awsGovCloudDisableInputs.Iam = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("lambda"); ok {
		awsGovCloudDisableInputs.Lambda = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("rds"); ok {
		awsGovCloudDisableInputs.Rds = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("red_shift"); ok {
		awsGovCloudDisableInputs.Redshift = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("route53"); ok {
		awsGovCloudDisableInputs.Route53 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("sns"); ok {
		awsGovCloudDisableInputs.Sns = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("sqs"); ok {
		awsGovCloudDisableInputs.Sqs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		AwsGovcloud: awsGovCloudDisableInputs,
	}
	return deleteInput
}
