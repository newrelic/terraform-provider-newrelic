package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/organization"
)

func expandOrganization(d *schema.ResourceData) (*organization.OrganizationCreateOrganizationInput, error) {
	input := organization.OrganizationCreateOrganizationInput{
		Name: d.Get("name").(string),
	}
	return &input, nil
}

func expandOrganizationUpdate(d *schema.ResourceData) (*organization.OrganizationUpdateInput, error) {
	input := organization.OrganizationUpdateInput{}
	if v, ok := d.GetOk("name"); ok {
		input.Name = v.(string)
	}
	return &input, nil
}

func expandOrganizationNewManagedAccount(d *schema.ResourceData) *organization.OrganizationNewManagedAccountInput {
	v, ok := d.GetOk("new_managed_account")
	if !ok {
		return nil
	}
	list := v.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	cfg := list[0].(map[string]interface{})
	input := &organization.OrganizationNewManagedAccountInput{}
	if name, ok := cfg["name"].(string); ok && name != "" {
		input.Name = name
	}
	if regionCode, ok := cfg["region_code"].(string); ok && regionCode != "" {
		input.RegionCode = organization.OrganizationRegionCodeEnum(regionCode)
	}
	return input
}

func expandOrganizationSharedAccount(d *schema.ResourceData) *organization.OrganizationSharedAccountInput {
	v, ok := d.GetOk("shared_account")
	if !ok {
		return nil
	}
	list := v.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	cfg := list[0].(map[string]interface{})
	input := &organization.OrganizationSharedAccountInput{
		AccountID: cfg["account_id"].(int),
	}
	if limitingRoleId, ok := cfg["limiting_role_id"].(int); ok {
		input.LimitingRoleId = limitingRoleId
	}
	return input
}
func flattenOrganization(result *organization.Organization, d *schema.ResourceData) error {
	_ = d.Set("name", result.Name)
	_ = d.Set("organization_id", result.ID)
	return nil
}

func flattenOrganizationUpdateResponse(result *organization.OrganizationUpdateResponse, d *schema.ResourceData) error {
	if result.OrganizationInformation.ID != "" {
		_ = d.Set("organization_id", result.OrganizationInformation.ID)
	}
	if result.OrganizationInformation.Name != "" {
		_ = d.Set("name", result.OrganizationInformation.Name)
	}
	return nil
}