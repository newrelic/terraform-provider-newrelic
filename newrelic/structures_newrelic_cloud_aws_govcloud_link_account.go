package newrelic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func expandAwsGovCloudLinkAccountInputForCreate(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	// NOTE: The AwsGovCloudLinkAccountInput datatype is no longer supported to facilitate linking an AWS GovCloud
	// account to New Relic; AwsLinkAccountInput is intended to be used instead, since a link for AWS/AWS GovCloud
	// both can now be facilitated via the "aws" field in the CloudLinkCloudAccountsInput datatype, with the same
	// authentication mechanism, i.e. an ARN.

	awsGovCloud := cloud.CloudAwsLinkAccountInput{}

	// NOTE: The following arguments are no longer supported, as the establishment of a connection
	// with New Relic from AWS GovCloud is no longer supported with these credentials (an ARN is needed
	// to facilitate a working connection.

	//if accessKeyID, ok := d.GetOk("access_key_id"); ok {
	//	awsGovCloud.AccessKeyId = accessKeyID.(string)
	//}
	//if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
	//	awsGovCloud.AwsAccountId = awsAccountID.(string)
	//}
	//if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
	//	awsGovCloud.SecretAccessKey = cloud.SecureValue(secretKeyID.(string))
	//}

	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}
	if m, ok := d.GetOk("metric_collection_mode"); ok {
		awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	}
	if arn, ok := d.GetOk("arn"); ok {
		awsGovCloud.Arn = arn.(string)
	}

	createAwsGovCloudLinkAccountInput := cloud.CloudLinkCloudAccountsInput{
		Aws: []cloud.CloudAwsLinkAccountInput{awsGovCloud},
	}

	return createAwsGovCloudLinkAccountInput
}

func expandAwsGovCloudLinkAccountInputForRead(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("name", result.Name)
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("arn", result.AuthLabel)
}

func expandAwsGovCloudLinkAccountInputForUpdate(d *schema.ResourceData, linkedAccountID int) cloud.CloudUpdateCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsUpdateAccountInput{}
	awsGovCloud.LinkedAccountId = linkedAccountID

	// NOTE: The following arguments are no longer supported; see `expandAwsGovCloudLinkAccountInputForCreate` to know why

	//if accessKeyID, ok := d.GetOk("access_key_id"); ok {
	//	awsGovCloud.AccessKeyId = accessKeyID.(string)
	//}
	//if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
	//	awsGovCloud.AwsAccountId = awsAccountID.(string)
	//}
	//if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
	//	awsGovCloud.SecretAccessKey = cloud.SecureValue(secretKeyID.(string))
	//}

	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}

	// The update mutation does not support updating the metric collection mode
	// This is also why a 'ForceNew' constraint has been applied on this argument in the schema

	//if m, ok := d.GetOk("metric_collection_mode"); ok {
	//	awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	//}

	if arn, ok := d.GetOk("arn"); ok {
		awsGovCloud.Arn = arn.(string)
	}

	updateAwsGovCloudLinkAccountInput := cloud.CloudUpdateCloudAccountsInput{
		Aws: []cloud.CloudAwsUpdateAccountInput{awsGovCloud},
	}

	return updateAwsGovCloudLinkAccountInput
}

func getLinkedAccountIDFromState(d *schema.ResourceData) (int, error) {
	linkedAccountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return 0, fmt.Errorf("error converting linked account ID to int: %s", err)
	}
	return linkedAccountID, nil
}
