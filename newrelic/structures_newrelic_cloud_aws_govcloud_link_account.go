package newrelic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func expandAwsGovCloudLinkAccountInputForCreate(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsGovCloudLinkAccountInput{}
	if accessKeyID, ok := d.GetOk("access_key_id"); ok {
		awsGovCloud.AccessKeyId = accessKeyID.(string)
	}
	if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
		awsGovCloud.AwsAccountId = awsAccountID.(string)
	}
	if m, ok := d.GetOk("metric_collection_mode"); ok {
		awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	}
	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}
	if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
		awsGovCloud.SecretAccessKey = cloud.SecureValue(secretKeyID.(string))
	}

	createAwsGovCloudLinkAccountInput := cloud.CloudLinkCloudAccountsInput{
		AwsGovcloud: []cloud.CloudAwsGovCloudLinkAccountInput{awsGovCloud},
	}

	return createAwsGovCloudLinkAccountInput
}

func expandAwsGovCloudLinkAccountInputForRead(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("metric_collection_mode", result.MetricCollectionMode)
	_ = d.Set("name", result.Name)
	_ = d.Set("aws_account_id", result.ExternalId)
	_ = d.Set("account_id", result.NrAccountId)
}

func expandAwsGovCloudLinkAccountInputForUpdate(d *schema.ResourceData, linkedAccountID int) cloud.CloudUpdateCloudAccountsInput {
	awsGovCloud := cloud.CloudAwsGovCloudUpdateAccountInput{}
	awsGovCloud.LinkedAccountId = linkedAccountID
	if accessKeyID, ok := d.GetOk("access_key_id"); ok {
		awsGovCloud.AccessKeyId = accessKeyID.(string)
	}
	if awsAccountID, ok := d.GetOk("aws_account_id"); ok {
		awsGovCloud.AwsAccountId = awsAccountID.(string)
	}

	// The update mutation does not support updating the metric collection mode
	//if m, ok := d.GetOk("metric_collection_mode"); ok {
	//	awsGovCloud.MetricCollectionMode = cloud.CloudMetricCollectionMode(strings.ToUpper(m.(string)))
	//}

	if name, ok := d.GetOk("name"); ok {
		awsGovCloud.Name = name.(string)
	}
	if secretKeyID, ok := d.GetOk("secret_access_key"); ok {
		awsGovCloud.SecretAccessKey = cloud.SecureValue(secretKeyID.(string))
	}

	updateAwsGovCloudLinkAccountInput := cloud.CloudUpdateCloudAccountsInput{
		AwsGovcloud: []cloud.CloudAwsGovCloudUpdateAccountInput{awsGovCloud},
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
