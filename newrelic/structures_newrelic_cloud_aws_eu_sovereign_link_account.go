package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// expandAwsEuSovereignLinkAccountInputForCreate expands the schema data into a CloudLinkCloudAccountsInput for create operations
func expandAwsEuSovereignLinkAccountInputForCreate(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {
	input := cloud.CloudLinkCloudAccountsInput{
		AwsEuSovereign: []cloud.CloudAwsEuSovereignLinkAccountInput{
			{
				Arn:                  d.Get("arn").(string),
				MetricCollectionMode: cloud.CloudMetricCollectionMode(d.Get("metric_collection_mode").(string)),
				Name:                 d.Get("name").(string),
			},
		},
	}

	return input
}

// flattenAwsEuSovereignLinkAccountForRead flattens the linked account data into the schema for read operations
func flattenAwsEuSovereignLinkAccountForRead(linkedAccount *cloud.CloudLinkedAccount, d *schema.ResourceData, accountID int) error {
	_ = d.Set("account_id", accountID)
	_ = d.Set("name", linkedAccount.Name)
	_ = d.Set("arn", linkedAccount.AuthLabel)

	if linkedAccount.MetricCollectionMode != "" {
		_ = d.Set("metric_collection_mode", string(linkedAccount.MetricCollectionMode))
	}

	return nil
}

// expandAwsEuSovereignLinkAccountInputForUpdate expands the schema data into a CloudUpdateCloudAccountsInput for update operations
func expandAwsEuSovereignLinkAccountInputForUpdate(d *schema.ResourceData, linkedAccountID int) cloud.CloudUpdateCloudAccountsInput {
	input := cloud.CloudUpdateCloudAccountsInput{
		AwsEuSovereign: []cloud.CloudAwsEuSovereignUpdateAccountInput{
			{
				LinkedAccountId: linkedAccountID,
				Name:            d.Get("name").(string),
			},
		},
	}

	return input
}

// getAwsEuSovereignLinkedAccountIDFromState extracts the linked account ID from the terraform state
func getAwsEuSovereignLinkedAccountIDFromState(d *schema.ResourceData) int {
	linkedAccountID, _ := parseIDs(d.Id(), 1)
	return linkedAccountID[0]
}