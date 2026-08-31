package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
)

func expandOrganization(d *schema.ResourceData) (string, *organization.OrganizationNewManagedAccountInput, organization.OrganizationCreateOrganizationInput, *organization.OrganizationSharedAccountInput, error) {
	customerID := d.Get("customer_id").(string)

	orgInput := organization.OrganizationCreateOrganizationInput{
		Name: d.Get("name").(string),
	}

	var newManagedAccount *organization.OrganizationNewManagedAccountInput
	if v, ok := d.GetOk("new_managed_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			newManagedAccount = expandOrganizationNewManagedAccount(items[0].(map[string]interface{}))
		}
	}

	var sharedAccount *organization.OrganizationSharedAccountInput
	if v, ok := d.GetOk("shared_account"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			sharedAccount = expandOrganizationSharedAccount(items[0].(map[string]interface{}))
		}
	}

	return customerID, newManagedAccount, orgInput, sharedAccount, nil
}

func expandOrganizationNewManagedAccount(cfg map[string]interface{}) *organization.OrganizationNewManagedAccountInput {
	input := &organization.OrganizationNewManagedAccountInput{}

	if v, ok := cfg["name"].(string); ok && v != "" {
		input.Name = v
	}

	if v, ok := cfg["region_code"].(string); ok && v != "" {
		input.RegionCode = organization.OrganizationRegionCodeEnum(v)
	}

	return input
}

func expandOrganizationSharedAccount(cfg map[string]interface{}) *organization.OrganizationSharedAccountInput {
	input := &organization.OrganizationSharedAccountInput{}

	if v, ok := cfg["account_id"].(int); ok {
		input.AccountID = v
	}

	if v, ok := cfg["limiting_role_id"].(int); ok {
		input.LimitingRoleId = v
	}

	return input
}

func expandOrganizationUpdate(d *schema.ResourceData) (organization.OrganizationUpdateInput, error) {
	input := organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}

	return input, nil
}

func flattenOrganization(result *organization.OrganizationCreateOrganizationResponse, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	_ = d.Set("job_id", result.JobId)

	return nil
}

func flattenOrganizationInformation(info organization.OrganizationInformation) map[string]interface{} {
	return map[string]interface{}{
		"name": info.Name,
	}
}

func flattenOrganizationNewManagedAccount(account organization.OrganizationNewManagedAccountInput) []interface{} {
	m := map[string]interface{}{
		"name":        account.Name,
		"region_code": string(account.RegionCode),
	}
	return []interface{}{m}
}

func flattenOrganizationSharedAccount(account organization.OrganizationSharedAccountInput) []interface{} {
	m := map[string]interface{}{
		"account_id":       account.AccountID,
		"limiting_role_id": account.LimitingRoleId,
	}
	return []interface{}{m}
}
