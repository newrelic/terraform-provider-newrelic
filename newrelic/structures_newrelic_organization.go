package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
)

func expandOrganization(d *schema.ResourceData) (string, *organization.OrganizationNewManagedAccountInput, organization.OrganizationCreateOrganizationInput, *organization.OrganizationSharedAccountInput, error) {
	customerID := d.Get("customer_id").(string)

	org := organization.OrganizationCreateOrganizationInput{
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

	return customerID, newManagedAccount, org, sharedAccount, nil
}

func expandOrganizationNewManagedAccount(cfg map[string]interface{}) *organization.OrganizationNewManagedAccountInput {
	input := &organization.OrganizationNewManagedAccountInput{}

	if v, ok := cfg["name"]; ok {
		input.Name = v.(string)
	}

	if v, ok := cfg["region_code"]; ok {
		input.RegionCode = organization.OrganizationRegionCodeEnum(v.(string))
	}

	return input
}

func expandOrganizationSharedAccount(cfg map[string]interface{}) *organization.OrganizationSharedAccountInput {
	input := &organization.OrganizationSharedAccountInput{}

	if v, ok := cfg["account_id"]; ok {
		input.AccountID = v.(int)
	}

	if v, ok := cfg["limiting_role_id"]; ok {
		input.LimitingRoleId = v.(int)
	}

	return input
}

func expandOrganizationUpdate(d *schema.ResourceData) (organization.OrganizationUpdateInput, error) {
	input := organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}

	return input, nil
}

func flattenOrganization(result *organization.Organization, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	_ = d.Set("name", result.Name)

	return nil
}
