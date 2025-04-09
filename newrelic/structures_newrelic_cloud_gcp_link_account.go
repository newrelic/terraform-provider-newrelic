package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// expand function to extract inputs from the schema.
// Here it takes ResourceData as input and returns cloudLinkCloudAccountsInput.
func expandGcpCloudLinkAccountInput(d *schema.ResourceData) cloud.CloudLinkCloudAccountsInput {

	gcpAccount := cloud.CloudGcpLinkAccountInput{}

	if name, ok := d.GetOk("name"); ok {
		gcpAccount.Name = name.(string)
	}

	if projectID, ok := d.GetOk("project_id"); ok {
		gcpAccount.ProjectId = projectID.(string)
	}

	input := cloud.CloudLinkCloudAccountsInput{
		Gcp: []cloud.CloudGcpLinkAccountInput{gcpAccount},
	}

	return input

}

// readGcpLinkedAccount function to store name and ExternalId.
// Using set func to store the output values.
func readGcpLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("name", result.Name)
	_ = d.Set("project_id", result.ExternalId)
	_ = d.Set("disabled", result.Disabled)
}
