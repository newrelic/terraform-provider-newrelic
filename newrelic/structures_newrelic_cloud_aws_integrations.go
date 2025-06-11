package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// Used by the newrelic_cloud_aws_integrations Create & Update functions.
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
		"aws_glue": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsGlue = expandCloudAwsIntegrationGlueInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsGlue = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_kinesis_analytics": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsKinesisanalytics = expandCloudAwsIntegrationKinesisAnalyticsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsKinesisanalytics = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_media_convert": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsMediaconvert = expandCloudAwsIntegrationMediaConvertInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsMediaconvert = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_media_package_vod": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsMediapackagevod = expandCloudAwsIntegrationMediaPackageVodInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsMediapackagevod = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_mq": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsMq = expandCloudAwsIntegrationMqInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsMq = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_msk": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsMsk = expandCloudAwsIntegrationMskInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsMsk = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_neptune": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsNeptune = expandCloudAwsIntegrationNeptuneInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsNeptune = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_qldb": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsQldb = expandCloudAwsIntegrationQldbInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsQldb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_route53resolver": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsRoute53resolver = expandCloudAwsIntegrationRoute53resolverInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsRoute53resolver = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_states": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsStates = expandCloudAwsIntegrationStatesInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsStates = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_transit_gateway": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsTransitgateway = expandCloudAwsIntegrationTransitGatewayInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsTransitgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_waf": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsWaf = expandCloudAwsIntegrationWafInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsWaf = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_wafv2": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsWafv2 = expandCloudAwsIntegrationWafv2Input(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsWafv2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"aws_auto_discovery": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.AwsAutoDiscovery = expandCloudAwsIntegrationAutoDiscoveryInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.AwsAutoDiscovery = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"cloudfront": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Cloudfront = expandCloudAwsIntegrationCloudfrontInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Cloudfront = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"dynamodb": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Dynamodb = expandCloudAwsIntegrationDynamoDBInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Dynamodb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"ec2": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Ec2 = expandCloudAwsIntegrationEc2Input(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Ec2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"ecs": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Ecs = expandCloudAwsIntegrationEcsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Ecs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"efs": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Efs = expandCloudAwsIntegrationEfsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Efs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"elasticbeanstalk": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Elasticbeanstalk = expandCloudAwsIntegrationElasticbeanstalkInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Elasticbeanstalk = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"elasticsearch": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Elasticsearch = expandCloudAwsIntegrationElasticsearchInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Elasticsearch = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"elb": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Elb = expandCloudAwsIntegrationElbInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Elb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"emr": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Emr = expandCloudAwsIntegrationEmrInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Emr = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"iam": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Iam = expandCloudAwsIntegrationIamInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Iam = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"iot": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Iot = expandCloudAwsIntegrationIotInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Iot = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"kinesis": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Kinesis = expandCloudAwsIntegrationKinesisInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Kinesis = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"kinesis_firehose": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.KinesisFirehose = expandCloudAwsIntegrationKinesisFirehoseInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.KinesisFirehose = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"lambda": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Lambda = expandCloudAwsIntegrationLambdaInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Lambda = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"rds": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Rds = expandCloudAwsIntegrationRdsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Rds = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"redshift": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Redshift = expandCloudAwsIntegrationRedshiftInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Redshift = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"route53": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Route53 = expandCloudAwsIntegrationRoute53Input(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Route53 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"ses": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Ses = expandCloudAwsIntegrationSesInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Ses = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
		"sns": {
			enableFunc: func(a []interface{}, id int) {
				cloudAwsIntegration.Sns = expandCloudAwsIntegrationSnsInput(a, id)
			},
			disableFunc: func(id int) {
				cloudDisableAwsIntegration.Sns = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: id}}
			},
		},
	}

	for key, fun := range awsIntegrationMap {
		if v, ok := d.GetOk(key); ok {
			fun.enableFunc(v.([]interface{}), linkedAccountID)
		} else if o, n := d.GetChange(key); len(n.([]interface{})) < len(o.([]interface{})) {
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

// TODO: Reduce the cyclomatic complexity of this func
// nolint: gocyclo
// Used by the newrelic_cloud_aws_integrations Read function.
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
		case *cloud.CloudAwsGlueIntegration:
			_ = d.Set("aws_glue", flattenCloudAwsGlueIntegration(t))
		case *cloud.CloudAwsKinesisanalyticsIntegration:
			_ = d.Set("aws_kinesis_analytics", flattenCloudAwsKinesisAnalyticsIntegration(t))
		case *cloud.CloudAwsMediaconvertIntegration:
			_ = d.Set("aws_media_convert", flattenCloudAwsMediaConvertIntegration(t))
		case *cloud.CloudAwsMediapackagevodIntegration:
			_ = d.Set("aws_media_package_vod", flattenCloudAwsMediaPackageVodIntegration(t))
		case *cloud.CloudAwsMqIntegration:
			_ = d.Set("aws_mq", flattenCloudAwsMqIntegration(t))
		case *cloud.CloudAwsMskIntegration:
			_ = d.Set("aws_msk", flattenCloudAwsMskIntegration(t))
		case *cloud.CloudAwsNeptuneIntegration:
			_ = d.Set("aws_neptune", flattenCloudAwsNeptuneIntegration(t))
		case *cloud.CloudAwsQldbIntegration:
			_ = d.Set("aws_qldb", flattenCloudAwsQldbIntegration(t))
		case *cloud.CloudAwsRoute53resolverIntegration:
			_ = d.Set("aws_route53resolver", flattenCloudAwsRoute53resolverIntegration(t))
		case *cloud.CloudAwsStatesIntegration:
			_ = d.Set("aws_states", flattenCloudAwsStatesIntegration(t))
		case *cloud.CloudAwsTransitgatewayIntegration:
			_ = d.Set("aws_transit_gateway", flattenCloudAwsTransitGatewayIntegration(t))
		case *cloud.CloudAwsWafIntegration:
			_ = d.Set("aws_waf", flattenCloudAwsWafIntegration(t))
		case *cloud.CloudAwsWafv2Integration:
			_ = d.Set("aws_wafv2", flattenCloudAwsWafv2Integration(t))
		case *cloud.CloudAwsAutoDiscoveryIntegration:
			_ = d.Set("aws_auto_discovery", flattenCloudAwsAutoDiscoveryIntegration(t))
		case *cloud.CloudCloudfrontIntegration:
			_ = d.Set("cloudfront", flattenCloudCloudfrontIntegration(t))
		case *cloud.CloudDynamodbIntegration:
			_ = d.Set("dynamodb", flattenCloudDynamodbIntegration(t))
		case *cloud.CloudEc2Integration:
			_ = d.Set("ec2", flattenCloudEc2Integration(t))
		case *cloud.CloudEcsIntegration:
			_ = d.Set("ecs", flattenCloudEcsIntegration(t))
		case *cloud.CloudEfsIntegration:
			_ = d.Set("efs", flattenCloudEfsIntegration(t))
		case *cloud.CloudElasticbeanstalkIntegration:
			_ = d.Set("elasticbeanstalk", flattenCloudElasticbeanstalkIntegration(t))
		case *cloud.CloudElasticsearchIntegration:
			_ = d.Set("elasticsearch", flattenCloudElasticsearchIntegration(t))
		case *cloud.CloudElbIntegration:
			_ = d.Set("elb", flattenCloudElbIntegration(t))
		case *cloud.CloudEmrIntegration:
			_ = d.Set("emr", flattenCloudEmrIntegration(t))
		case *cloud.CloudIamIntegration:
			_ = d.Set("iam", flattenCloudIamIntegration(t))
		case *cloud.CloudIotIntegration:
			_ = d.Set("iot", flattenCloudIotIntegration(t))
		case *cloud.CloudKinesisIntegration:
			_ = d.Set("kinesis", flattenCloudKinesisIntegration(t))
		case *cloud.CloudKinesisFirehoseIntegration:
			_ = d.Set("kinesis_firehose", flattenCloudKinesisfirehoseIntegration(t))
		case *cloud.CloudLambdaIntegration:
			_ = d.Set("lambda", flattenCloudLambdaIntegration(t))
		case *cloud.CloudRdsIntegration:
			_ = d.Set("rds", flattenCloudRdsIntegration(t))
		case *cloud.CloudRedshiftIntegration:
			_ = d.Set("redshift", flattenCloudRedshiftIntegration(t))
		case *cloud.CloudRoute53Integration:
			_ = d.Set("route53", flattenCloudRoute53Integration(t))
		case *cloud.CloudSesIntegration:
			_ = d.Set("ses", flattenCloudSesIntegration(t))
		case *cloud.CloudSnsIntegration:
			_ = d.Set("sns", flattenCloudSnsIntegration(t))
		}
	}
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint: gocyclo
// Used by the newrelic_cloud_aws_integrations delete function.
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

	if _, ok := d.GetOk("aws_glue"); ok {
		cloudDisableAwsIntegration.AwsGlue = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_kinesis_analytics"); ok {
		cloudDisableAwsIntegration.AwsKinesisanalytics = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_media_convert"); ok {
		cloudDisableAwsIntegration.AwsMediaconvert = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_media_package_vod"); ok {
		cloudDisableAwsIntegration.AwsMediapackagevod = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_mq"); ok {
		cloudDisableAwsIntegration.AwsMq = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_msk"); ok {
		cloudDisableAwsIntegration.AwsMsk = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_neptune"); ok {
		cloudDisableAwsIntegration.AwsNeptune = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_qldb"); ok {
		cloudDisableAwsIntegration.AwsQldb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_route53resolver"); ok {
		cloudDisableAwsIntegration.AwsRoute53resolver = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_states"); ok {
		cloudDisableAwsIntegration.AwsStates = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_transit_gateway"); ok {
		cloudDisableAwsIntegration.AwsTransitgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_waf"); ok {
		cloudDisableAwsIntegration.AwsWaf = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_wafv2"); ok {
		cloudDisableAwsIntegration.AwsWafv2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("aws_auto_discovery"); ok {
		cloudDisableAwsIntegration.AwsAutoDiscovery = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("cloudfront"); ok {
		cloudDisableAwsIntegration.Cloudfront = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("dynamodb"); ok {
		cloudDisableAwsIntegration.Dynamodb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("ec2"); ok {
		cloudDisableAwsIntegration.Ec2 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("ecs"); ok {
		cloudDisableAwsIntegration.Ecs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("efs"); ok {
		cloudDisableAwsIntegration.Efs = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("elasticbeanstalk"); ok {
		cloudDisableAwsIntegration.Elasticbeanstalk = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("elasticsearch"); ok {
		cloudDisableAwsIntegration.Elasticsearch = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("elb"); ok {
		cloudDisableAwsIntegration.Elb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("emr"); ok {
		cloudDisableAwsIntegration.Emr = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("iam"); ok {
		cloudDisableAwsIntegration.Iam = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("iot"); ok {
		cloudDisableAwsIntegration.Iot = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("kinesis"); ok {
		cloudDisableAwsIntegration.Kinesis = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("kinesis_firehose"); ok {
		cloudDisableAwsIntegration.KinesisFirehose = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("lambda"); ok {
		cloudDisableAwsIntegration.Lambda = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("rds"); ok {
		cloudDisableAwsIntegration.Rds = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("redshift"); ok {
		cloudDisableAwsIntegration.Redshift = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("route53"); ok {
		cloudDisableAwsIntegration.Route53 = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("ses"); ok {
		cloudDisableAwsIntegration.Ses = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if _, ok := d.GetOk("sns"); ok {
		cloudDisableAwsIntegration.Sns = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	deleteInput := cloud.CloudDisableIntegrationsInput{
		Aws: cloudDisableAwsIntegration,
	}

	return deleteInput
}

type enableDisableAwsIntegration struct {
	enableFunc  func([]interface{}, int)
	disableFunc func(int)
}

// Expanding the Billing input

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

// Expanding the Cloud Trail input

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

// Expanding the Health input

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

// Expanding the Doc DB input

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
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			docDbInput.AwsRegions = regions
		}

		expanded[i] = docDbInput
	}

	return expanded
}

// Expanding the Trusted Advisor input

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

// Expanding the S3 input

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

// Expanding the VPC input

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

// Expanding the XRay input

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

// Expanding the Sqs input

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

// Expanding the Ebs input

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

// Expanding the Alb input

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

// Expanding the Elasticache input

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

// Expanding the api gateway input

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

// Expanding the auto scaling input

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

// Expanding the aws app sync input

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

// Expanding the aws athena input

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

// Expanding the aws cognito input

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

// Expanding the aws connect input

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

// Expanding the aws direct connect input

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

// Expanding the aws fsx input

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

// Expanding the aws glue input

func expandCloudAwsIntegrationGlueInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsGlueIntegrationInput {
	expanded := make([]cloud.CloudAwsGlueIntegrationInput, len(b))

	for i, glue := range b {
		var glueInput cloud.CloudAwsGlueIntegrationInput

		if glue == nil {
			glueInput.LinkedAccountId = linkedAccountID
			expanded[i] = glueInput
			return expanded
		}

		in := glue.(map[string]interface{})

		glueInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			glueInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			glueInput.AwsRegions = regions
		}

		expanded[i] = glueInput
	}

	return expanded
}

// Expanding the aws kinesis analytics input

func expandCloudAwsIntegrationKinesisAnalyticsInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsKinesisanalyticsIntegrationInput {
	expanded := make([]cloud.CloudAwsKinesisanalyticsIntegrationInput, len(b))

	for i, kinesisAnalytics := range b {
		var kinesisAnalyticsInput cloud.CloudAwsKinesisanalyticsIntegrationInput

		if kinesisAnalytics == nil {
			kinesisAnalyticsInput.LinkedAccountId = linkedAccountID
			expanded[i] = kinesisAnalyticsInput
			return expanded
		}

		in := kinesisAnalytics.(map[string]interface{})

		kinesisAnalyticsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			kinesisAnalyticsInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			kinesisAnalyticsInput.AwsRegions = regions
		}

		expanded[i] = kinesisAnalyticsInput
	}

	return expanded
}

// Expanding the aws media convert input

func expandCloudAwsIntegrationMediaConvertInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsMediaconvertIntegrationInput {
	expanded := make([]cloud.CloudAwsMediaconvertIntegrationInput, len(b))

	for i, mediaConvert := range b {
		var mediaConvertInput cloud.CloudAwsMediaconvertIntegrationInput

		if mediaConvert == nil {
			mediaConvertInput.LinkedAccountId = linkedAccountID
			expanded[i] = mediaConvertInput
			return expanded
		}

		in := mediaConvert.(map[string]interface{})

		mediaConvertInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			mediaConvertInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			mediaConvertInput.AwsRegions = regions
		}

		expanded[i] = mediaConvertInput
	}

	return expanded
}

// Expanding the aws media package vod input

func expandCloudAwsIntegrationMediaPackageVodInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsMediapackagevodIntegrationInput {
	expanded := make([]cloud.CloudAwsMediapackagevodIntegrationInput, len(b))

	for i, mediaPackageVod := range b {
		var mediaPackageVodInput cloud.CloudAwsMediapackagevodIntegrationInput

		if mediaPackageVod == nil {
			mediaPackageVodInput.LinkedAccountId = linkedAccountID
			expanded[i] = mediaPackageVodInput
			return expanded
		}

		in := mediaPackageVod.(map[string]interface{})

		mediaPackageVodInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			mediaPackageVodInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			mediaPackageVodInput.AwsRegions = regions
		}

		expanded[i] = mediaPackageVodInput
	}

	return expanded
}

// Expanding the aws mq input

func expandCloudAwsIntegrationMqInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsMqIntegrationInput {
	expanded := make([]cloud.CloudAwsMqIntegrationInput, len(b))

	for i, mq := range b {
		var mqInput cloud.CloudAwsMqIntegrationInput

		if mq == nil {
			mqInput.LinkedAccountId = linkedAccountID
			expanded[i] = mqInput
			return expanded
		}

		in := mq.(map[string]interface{})

		mqInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			mqInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			mqInput.AwsRegions = regions
		}

		expanded[i] = mqInput
	}

	return expanded
}

// Expanding the aws msk input

func expandCloudAwsIntegrationMskInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsMskIntegrationInput {
	expanded := make([]cloud.CloudAwsMskIntegrationInput, len(b))

	for i, msk := range b {
		var mskInput cloud.CloudAwsMskIntegrationInput

		if msk == nil {
			mskInput.LinkedAccountId = linkedAccountID
			expanded[i] = mskInput
			return expanded
		}

		in := msk.(map[string]interface{})

		mskInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			mskInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			mskInput.AwsRegions = regions
		}

		expanded[i] = mskInput
	}

	return expanded
}

// Expanding the aws neptune input

func expandCloudAwsIntegrationNeptuneInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsNeptuneIntegrationInput {
	expanded := make([]cloud.CloudAwsNeptuneIntegrationInput, len(b))

	for i, neptune := range b {
		var neptuneInput cloud.CloudAwsNeptuneIntegrationInput

		if neptune == nil {
			neptuneInput.LinkedAccountId = linkedAccountID
			expanded[i] = neptuneInput
			return expanded
		}

		in := neptune.(map[string]interface{})

		neptuneInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			neptuneInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			neptuneInput.AwsRegions = regions
		}

		expanded[i] = neptuneInput
	}

	return expanded
}

// Expanding the aws qldb input

func expandCloudAwsIntegrationQldbInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsQldbIntegrationInput {
	expanded := make([]cloud.CloudAwsQldbIntegrationInput, len(b))

	for i, qldb := range b {
		var qldbInput cloud.CloudAwsQldbIntegrationInput

		if qldb == nil {
			qldbInput.LinkedAccountId = linkedAccountID
			expanded[i] = qldbInput
			return expanded
		}

		in := qldb.(map[string]interface{})

		qldbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			qldbInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			qldbInput.AwsRegions = regions
		}

		expanded[i] = qldbInput
	}

	return expanded
}

// Expanding the aws route53resolver input

func expandCloudAwsIntegrationRoute53resolverInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsRoute53resolverIntegrationInput {
	expanded := make([]cloud.CloudAwsRoute53resolverIntegrationInput, len(b))

	for i, route53resolver := range b {
		var route53resolverInput cloud.CloudAwsRoute53resolverIntegrationInput

		if route53resolver == nil {
			route53resolverInput.LinkedAccountId = linkedAccountID
			expanded[i] = route53resolverInput
			return expanded
		}

		in := route53resolver.(map[string]interface{})

		route53resolverInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			route53resolverInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			route53resolverInput.AwsRegions = regions
		}

		expanded[i] = route53resolverInput
	}

	return expanded
}

// Expanding the aws states input

func expandCloudAwsIntegrationStatesInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsStatesIntegrationInput {
	expanded := make([]cloud.CloudAwsStatesIntegrationInput, len(b))

	for i, states := range b {
		var statesInput cloud.CloudAwsStatesIntegrationInput

		if states == nil {
			statesInput.LinkedAccountId = linkedAccountID
			expanded[i] = statesInput
			return expanded
		}

		in := states.(map[string]interface{})

		statesInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			statesInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			statesInput.AwsRegions = regions
		}

		expanded[i] = statesInput
	}

	return expanded
}

// Expanding the aws transit gateway input

func expandCloudAwsIntegrationTransitGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsTransitgatewayIntegrationInput {
	expanded := make([]cloud.CloudAwsTransitgatewayIntegrationInput, len(b))

	for i, transitGateway := range b {
		var transitGatewayInput cloud.CloudAwsTransitgatewayIntegrationInput

		if transitGateway == nil {
			transitGatewayInput.LinkedAccountId = linkedAccountID
			expanded[i] = transitGatewayInput
			return expanded
		}

		in := transitGateway.(map[string]interface{})

		transitGatewayInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			transitGatewayInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			transitGatewayInput.AwsRegions = regions
		}

		expanded[i] = transitGatewayInput
	}

	return expanded
}

// Expanding the aws waf input

func expandCloudAwsIntegrationWafInput(b []interface{}, linkedAccountID int) []cloud.CloudAwsWafIntegrationInput {
	expanded := make([]cloud.CloudAwsWafIntegrationInput, len(b))

	for i, waf := range b {
		var wafInput cloud.CloudAwsWafIntegrationInput

		if waf == nil {
			wafInput.LinkedAccountId = linkedAccountID
			expanded[i] = wafInput
			return expanded
		}

		in := waf.(map[string]interface{})

		wafInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			wafInput.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			wafInput.AwsRegions = regions
		}

		expanded[i] = wafInput
	}

	return expanded
}

// Expanding the aws wafv2 input

func expandCloudAwsIntegrationWafv2Input(b []interface{}, linkedAccountID int) []cloud.CloudAwsWafv2IntegrationInput {
	expanded := make([]cloud.CloudAwsWafv2IntegrationInput, len(b))

	for i, wafv2 := range b {
		var wafv2Input cloud.CloudAwsWafv2IntegrationInput

		if wafv2 == nil {
			wafv2Input.LinkedAccountId = linkedAccountID
			expanded[i] = wafv2Input
			return expanded
		}

		in := wafv2.(map[string]interface{})

		wafv2Input.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			wafv2Input.MetricsPollingInterval = m.(int)
		}
		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			wafv2Input.AwsRegions = regions
		}

		expanded[i] = wafv2Input
	}

	return expanded
}

// Expanding the cloudfront input

func expandCloudAwsIntegrationCloudfrontInput(b []interface{}, linkedAccountID int) []cloud.CloudCloudfrontIntegrationInput {
	expanded := make([]cloud.CloudCloudfrontIntegrationInput, len(b))

	for i, cloudfront := range b {
		var cloudfrontInput cloud.CloudCloudfrontIntegrationInput

		if cloudfront == nil {
			cloudfrontInput.LinkedAccountId = linkedAccountID
			expanded[i] = cloudfrontInput
			return expanded
		}

		in := cloudfront.(map[string]interface{})

		cloudfrontInput.LinkedAccountId = linkedAccountID

		if ft, ok := in["fetch_lambdas_at_edge"]; ok {
			cloudfrontInput.FetchLambdasAtEdge = ft.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			cloudfrontInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			cloudfrontInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			cloudfrontInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			cloudfrontInput.TagValue = tv.(string)
		}

		expanded[i] = cloudfrontInput
	}

	return expanded
}

// Expanding the dynamodb input

func expandCloudAwsIntegrationDynamoDBInput(b []interface{}, linkedAccountID int) []cloud.CloudDynamodbIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			dynamodbInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_extended_inventory"]; ok {
			dynamodbInput.FetchExtendedInventory = ft.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			dynamodbInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			dynamodbInput.MetricsPollingInterval = m.(int)
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

// Expanding the ec2 input

func expandCloudAwsIntegrationEc2Input(b []interface{}, linkedAccountID int) []cloud.CloudEc2IntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			ec2Input.AwsRegions = regions
		}

		if ft, ok := in["duplicate_ec2_tags"]; ok {
			ec2Input.DuplicateEc2Tags = ft.(bool)
		}

		if ft, ok := in["fetch_ip_addresses"]; ok {
			ec2Input.FetchIpAddresses = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			ec2Input.MetricsPollingInterval = m.(int)
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

// Expanding the ecs input

func expandCloudAwsIntegrationEcsInput(b []interface{}, linkedAccountID int) []cloud.CloudEcsIntegrationInput {
	expanded := make([]cloud.CloudEcsIntegrationInput, len(b))

	for i, ecs := range b {
		var ecsInput cloud.CloudEcsIntegrationInput

		if ecs == nil {
			ecsInput.LinkedAccountId = linkedAccountID
			expanded[i] = ecsInput
			return expanded
		}

		in := ecs.(map[string]interface{})

		ecsInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			ecsInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			ecsInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			ecsInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			ecsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			ecsInput.TagValue = tv.(string)
		}

		expanded[i] = ecsInput
	}

	return expanded
}

// Expanding the efs input

func expandCloudAwsIntegrationEfsInput(b []interface{}, linkedAccountID int) []cloud.CloudEfsIntegrationInput {
	expanded := make([]cloud.CloudEfsIntegrationInput, len(b))

	for i, efs := range b {
		var efsInput cloud.CloudEfsIntegrationInput

		if efs == nil {
			efsInput.LinkedAccountId = linkedAccountID
			expanded[i] = efsInput
			return expanded
		}

		in := efs.(map[string]interface{})

		efsInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			efsInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			efsInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			efsInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			efsInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			efsInput.TagValue = tv.(string)
		}

		expanded[i] = efsInput
	}

	return expanded
}

// Expanding the elasticbeanstalk input

func expandCloudAwsIntegrationElasticbeanstalkInput(b []interface{}, linkedAccountID int) []cloud.CloudElasticbeanstalkIntegrationInput {
	expanded := make([]cloud.CloudElasticbeanstalkIntegrationInput, len(b))

	for i, elasticbeanstalk := range b {
		var elasticbeanstalkInput cloud.CloudElasticbeanstalkIntegrationInput

		if elasticbeanstalk == nil {
			elasticbeanstalkInput.LinkedAccountId = linkedAccountID
			expanded[i] = elasticbeanstalkInput
			return expanded
		}

		in := elasticbeanstalk.(map[string]interface{})

		elasticbeanstalkInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elasticbeanstalkInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_extended_inventory"]; ok {
			elasticbeanstalkInput.FetchExtendedInventory = ft.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			elasticbeanstalkInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			elasticbeanstalkInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			elasticbeanstalkInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			elasticbeanstalkInput.TagValue = tv.(string)
		}

		expanded[i] = elasticbeanstalkInput
	}

	return expanded
}

// Expanding the elasticsearch input

func expandCloudAwsIntegrationElasticsearchInput(b []interface{}, linkedAccountID int) []cloud.CloudElasticsearchIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elasticsearchInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_nodes"]; ok {
			elasticsearchInput.FetchNodes = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			elasticsearchInput.MetricsPollingInterval = m.(int)
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

// Expanding the elb input

func expandCloudAwsIntegrationElbInput(b []interface{}, linkedAccountID int) []cloud.CloudElbIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			elbInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_extended_inventory"]; ok {
			elbInput.FetchExtendedInventory = ft.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			elbInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			elbInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = elbInput
	}

	return expanded
}

// Expanding the emr input

func expandCloudAwsIntegrationEmrInput(b []interface{}, linkedAccountID int) []cloud.CloudEmrIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			emrInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			emrInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			emrInput.MetricsPollingInterval = m.(int)
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

// Expanding the iam input

func expandCloudAwsIntegrationIamInput(b []interface{}, linkedAccountID int) []cloud.CloudIamIntegrationInput {
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
			iamInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			iamInput.TagValue = tv.(string)
		}

		expanded[i] = iamInput
	}

	return expanded
}

// Expanding the iot input

func expandCloudAwsIntegrationIotInput(b []interface{}, linkedAccountID int) []cloud.CloudIotIntegrationInput {
	expanded := make([]cloud.CloudIotIntegrationInput, len(b))

	for i, iot := range b {
		var iotInput cloud.CloudIotIntegrationInput

		if iot == nil {
			iotInput.LinkedAccountId = linkedAccountID
			expanded[i] = iotInput
			return expanded
		}

		in := iot.(map[string]interface{})

		iotInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			iotInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			iotInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = iotInput
	}

	return expanded
}

// Expanding the kinesis input

func expandCloudAwsIntegrationKinesisInput(b []interface{}, linkedAccountID int) []cloud.CloudKinesisIntegrationInput {
	expanded := make([]cloud.CloudKinesisIntegrationInput, len(b))

	for i, kinesis := range b {
		var kinesisInput cloud.CloudKinesisIntegrationInput

		if kinesis == nil {
			kinesisInput.LinkedAccountId = linkedAccountID
			expanded[i] = kinesisInput
			return expanded
		}

		in := kinesis.(map[string]interface{})

		kinesisInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			kinesisInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_shards"]; ok {
			kinesisInput.FetchShards = ft.(bool)
		}

		if ft, ok := in["fetch_tags"]; ok {
			kinesisInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			kinesisInput.MetricsPollingInterval = m.(int)
		}

		if tk, ok := in["tag_key"]; ok {
			kinesisInput.TagKey = tk.(string)
		}

		if tv, ok := in["tag_value"]; ok {
			kinesisInput.TagValue = tv.(string)
		}

		expanded[i] = kinesisInput
	}

	return expanded
}

// Expanding the kinesisfirehose input

func expandCloudAwsIntegrationKinesisFirehoseInput(b []interface{}, linkedAccountID int) []cloud.CloudKinesisFirehoseIntegrationInput {
	expanded := make([]cloud.CloudKinesisFirehoseIntegrationInput, len(b))

	for i, firehose := range b {
		var firehoseInput cloud.CloudKinesisFirehoseIntegrationInput

		if firehose == nil {
			firehoseInput.LinkedAccountId = linkedAccountID
			expanded[i] = firehoseInput
			return expanded
		}

		in := firehose.(map[string]interface{})

		firehoseInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			firehoseInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			firehoseInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = firehoseInput
	}

	return expanded
}

// Expanding the lambda input

func expandCloudAwsIntegrationLambdaInput(b []interface{}, linkedAccountID int) []cloud.CloudLambdaIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			lambdaInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			lambdaInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			lambdaInput.MetricsPollingInterval = m.(int)
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

// Expanding the rds input

func expandCloudAwsIntegrationRdsInput(b []interface{}, linkedAccountID int) []cloud.CloudRdsIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			rdsInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_tags"]; ok {
			rdsInput.FetchTags = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			rdsInput.MetricsPollingInterval = m.(int)
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

// Expanding the redshift input

func expandCloudAwsIntegrationRedshiftInput(b []interface{}, linkedAccountID int) []cloud.CloudRedshiftIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			redshiftInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			redshiftInput.MetricsPollingInterval = m.(int)
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

// Expanding the route53 input

func expandCloudAwsIntegrationRoute53Input(b []interface{}, linkedAccountID int) []cloud.CloudRoute53IntegrationInput {
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

		if ft, ok := in["fetch_extended_inventory"]; ok {
			route53Input.FetchExtendedInventory = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			route53Input.MetricsPollingInterval = m.(int)
		}

		expanded[i] = route53Input
	}

	return expanded
}

// Expanding the ses input

func expandCloudAwsIntegrationSesInput(b []interface{}, linkedAccountID int) []cloud.CloudSesIntegrationInput {
	expanded := make([]cloud.CloudSesIntegrationInput, len(b))

	for i, ses := range b {
		var sesInput cloud.CloudSesIntegrationInput

		if ses == nil {
			sesInput.LinkedAccountId = linkedAccountID
			expanded[i] = sesInput
			return expanded
		}

		in := ses.(map[string]interface{})

		sesInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			sesInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			sesInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = sesInput
	}

	return expanded
}

// Expanding the sns input

func expandCloudAwsIntegrationSnsInput(b []interface{}, linkedAccountID int) []cloud.CloudSnsIntegrationInput {
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

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			snsInput.AwsRegions = regions
		}

		if ft, ok := in["fetch_extended_inventory"]; ok {
			snsInput.FetchExtendedInventory = ft.(bool)
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			snsInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = snsInput
	}

	return expanded
}

// Expanding the Aws Auto Discovery input

func expandCloudAwsIntegrationAutoDiscoveryInput(b []interface{}, linkedAccountID int) []cloud.CloudAutoDiscoveryIntegrationInput {
	expanded := make([]cloud.CloudAutoDiscoveryIntegrationInput, len(b))

	for i, iot := range b {
		var autoDiscoveryInput cloud.CloudAutoDiscoveryIntegrationInput

		if iot == nil {
			autoDiscoveryInput.LinkedAccountId = linkedAccountID
			expanded[i] = autoDiscoveryInput
			return expanded
		}

		in := iot.(map[string]interface{})

		autoDiscoveryInput.LinkedAccountId = linkedAccountID

		if a, ok := in["aws_regions"]; ok {
			awsRegions := a.([]interface{})
			var regions []string

			for _, region := range awsRegions {
				regions = append(regions, region.(string))
			}
			autoDiscoveryInput.AwsRegions = regions
		}

		if m, ok := in["metrics_polling_interval"]; ok {
			autoDiscoveryInput.MetricsPollingInterval = m.(int)
		}
		expanded[i] = autoDiscoveryInput
	}

	return expanded
}

// flatten for Billing integration

func flattenCloudAwsBillingIntegration(in *cloud.CloudBillingIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Cloud Trail integration

func flattenCloudAwsCloudTrailIntegration(in *cloud.CloudCloudtrailIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Health integration

func flattenCloudAwsHealthIntegration(in *cloud.CloudHealthIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Trusted Advisor integration

func flattenCloudAwsTrustedAdvisorIntegration(in *cloud.CloudTrustedadvisorIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Vpc integration

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

// flatten for XRay integration

func flattenCloudAwsXRayIntegration(in *cloud.CloudAwsXrayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Sqs integration

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

// flatten for Ebs integration

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

// flatten for Alb integration

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

// flatten for Elasticache integration

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

// flatten for S3 integration

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

// flatten for Doc DB integration

func flattenCloudAwsDocDBIntegration(in *cloud.CloudAwsDocdbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for api gateway integration

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

// flatten for autoscaling integration

func flattenCloudAwsAutoScalingIntegration(in *cloud.CloudAutoscalingIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for appsync integration

func flattenCloudAwsAppsyncIntegration(in *cloud.CloudAwsAppsyncIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for athena integration

func flattenCloudAwsAthenaIntegration(in *cloud.CloudAwsAthenaIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for cognito integration

func flattenCloudAwsCognitoIntegration(in *cloud.CloudAwsCognitoIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for connect integration

func flattenCloudAwsConnectIntegration(in *cloud.CloudAwsConnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for direct connect integration

func flattenCloudAwsDirectconnectIntegration(in *cloud.CloudAwsDirectconnectIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for fsx integration

func flattenCloudAwsFsxIntegration(in *cloud.CloudAwsFsxIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for glue integration

func flattenCloudAwsGlueIntegration(in *cloud.CloudAwsGlueIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for kinesis analytics integration

func flattenCloudAwsKinesisAnalyticsIntegration(in *cloud.CloudAwsKinesisanalyticsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for media convert integration

func flattenCloudAwsMediaConvertIntegration(in *cloud.CloudAwsMediaconvertIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for media package vod integration

func flattenCloudAwsMediaPackageVodIntegration(in *cloud.CloudAwsMediapackagevodIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for mq integration

func flattenCloudAwsMqIntegration(in *cloud.CloudAwsMqIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for msk integration

func flattenCloudAwsMskIntegration(in *cloud.CloudAwsMskIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for neptune integration

func flattenCloudAwsNeptuneIntegration(in *cloud.CloudAwsNeptuneIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for qldb integration

func flattenCloudAwsQldbIntegration(in *cloud.CloudAwsQldbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for route53resolver integration

func flattenCloudAwsRoute53resolverIntegration(in *cloud.CloudAwsRoute53resolverIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for aws states integration

func flattenCloudAwsStatesIntegration(in *cloud.CloudAwsStatesIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for transit gateway integration

func flattenCloudAwsTransitGatewayIntegration(in *cloud.CloudAwsTransitgatewayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for waf integration

func flattenCloudAwsWafIntegration(in *cloud.CloudAwsWafIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for wafv2 integration

func flattenCloudAwsWafv2Integration(in *cloud.CloudAwsWafv2Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["aws_regions"] = in.AwsRegions

	flattened[0] = out

	return flattened
}

// flatten for cloudfront integration

func flattenCloudCloudfrontIntegration(in *cloud.CloudCloudfrontIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_lambdas_at_edge"] = in.FetchLambdasAtEdge
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for dynamodb integration

func flattenCloudDynamodbIntegration(in *cloud.CloudDynamodbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for ec2 integration

func flattenCloudEc2Integration(in *cloud.CloudEc2Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["duplicate_ec2_tags"] = in.DuplicateEc2Tags
	out["fetch_ip_addresses"] = in.FetchIpAddresses
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for ecs integration

func flattenCloudEcsIntegration(in *cloud.CloudEcsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for efs integration

func flattenCloudEfsIntegration(in *cloud.CloudEfsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for elasticbeanstalk integration

func flattenCloudElasticbeanstalkIntegration(in *cloud.CloudElasticbeanstalkIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for elasticsearch integration

func flattenCloudElasticsearchIntegration(in *cloud.CloudElasticsearchIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_nodes"] = in.FetchNodes
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for elb integration

func flattenCloudElbIntegration(in *cloud.CloudElbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory
	out["fetch_tags"] = in.FetchTags

	flattened[0] = out

	return flattened
}

// flatten for emr integration

func flattenCloudEmrIntegration(in *cloud.CloudEmrIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for iam integration

func flattenCloudIamIntegration(in *cloud.CloudIamIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for iot integration

func flattenCloudIotIntegration(in *cloud.CloudIotIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for kinesis integration

func flattenCloudKinesisIntegration(in *cloud.CloudKinesisIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_shards"] = in.FetchShards
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for kinesis firehose integration

func flattenCloudKinesisfirehoseIntegration(in *cloud.CloudKinesisFirehoseIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for lambda integration

func flattenCloudLambdaIntegration(in *cloud.CloudLambdaIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for rds integration

func flattenCloudRdsIntegration(in *cloud.CloudRdsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_tags"] = in.FetchTags
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for redshift integration

func flattenCloudRedshiftIntegration(in *cloud.CloudRedshiftIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tag_key"] = in.TagKey
	out["tag_value"] = in.TagValue

	flattened[0] = out

	return flattened
}

// flatten for route53 integration

func flattenCloudRoute53Integration(in *cloud.CloudRoute53Integration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory

	flattened[0] = out

	return flattened
}

// flatten for ses integration

func flattenCloudSesIntegration(in *cloud.CloudSesIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for sns integration

func flattenCloudSnsIntegration(in *cloud.CloudSnsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["fetch_extended_inventory"] = in.FetchExtendedInventory

	flattened[0] = out

	return flattened
}

// flatten for Auto Discovery integration

func flattenCloudAwsAutoDiscoveryIntegration(in *cloud.CloudAutoDiscoveryIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["aws_regions"] = in.AwsRegions
	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}
