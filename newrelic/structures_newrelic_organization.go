package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
)

func expandOrganizationCreate(d *schema.ResourceData) (organization.OrganizationCreateOrganizationInput, *organization.OrganizationNewManagedAccountInput, *organization.OrganizationSharedAccountInput) {
	input := organization.OrganizationCreateOrganizationInput{
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

	return input, newManagedAccount, sharedAccount
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

func expandOrganizationUpdate(d *schema.ResourceData) organization.OrganizationUpdateInput {
	return organization.OrganizationUpdateInput{
		Name: d.Get("name").(string),
	}
}
func flattenOrganization(result *organization.Organization, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	_ = d.Set("name", result.Name)

	return nil
}

func flattenOrganizationNewManagedAccount(input *organization.OrganizationNewManagedAccountInput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"name":        input.Name,
		"region_code": string(input.RegionCode),
	}

	return []interface{}{m}
}

func flattenOrganizationSharedAccount(input *organization.OrganizationSharedAccountInput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"account_id":       input.AccountID,
		"limiting_role_id": input.LimitingRoleId,
	}

	return []interface{}{m}
}