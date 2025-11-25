package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func dataSourceNewRelicFleetRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicFleetRoleRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the fleet-scoped role.",
				Optional:    true,
			},
			"scope": {
				Type:        schema.TypeString,
				Description: "The scope of the role (that is intended to be 'fleet').",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the fleet-scoped role.",
				Optional:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(customeradministration.MultiTenantAuthorizationRoleTypeEnumTypes.CUSTOM),
						string(customeradministration.MultiTenantAuthorizationRoleTypeEnumTypes.STANDARD),
					}, false,
				),
			},
		},
	}
}

func dataSourceNewRelicFleetRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	roleType := d.Get("type").(string)

	organization, err := client.Organization.GetOrganization()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch organization information: %v", err))
	}

	organizationID := organization.ID

	filter := customeradministration.MultiTenantAuthorizationRoleFilterInputExpression{}
	filter.OrganizationId = &customeradministration.MultiTenantAuthorizationRoleOrganizationIdInputFilter{
		Eq: organizationID,
	}
	filter.ScopeV2 = &customeradministration.MultiTenantAuthorizationRoleScopeV2InputFilter{
		Eq: "fleet",
	}

	if name == "" && roleType == "" {
		filter.Name = &customeradministration.MultiTenantAuthorizationRoleNameInputFilter{
			Eq: "Fleet Manager",
		}
		filter.Type = &customeradministration.MultiTenantAuthorizationRoleTypeInputFilter{
			Eq: customeradministration.MultiTenantAuthorizationRoleTypeEnumTypes.STANDARD,
		}
	} else {
		if name != "" {
			filter.Name = &customeradministration.MultiTenantAuthorizationRoleNameInputFilter{
				Eq: name,
			}
		}
		if roleType != "" {
			filter.Type = &customeradministration.MultiTenantAuthorizationRoleTypeInputFilter{
				Eq: customeradministration.MultiTenantAuthorizationRoleTypeEnum(roleType),
			}
		}
	}

	roles, err := client.CustomerAdministration.GetRoles(
		"",
		filter,
		[]customeradministration.MultiTenantAuthorizationRoleSortInput{},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(roles.Items) == 0 {
		return diag.Errorf("no fleet role found with the given criteria")
	}

	role := roles.Items[0]

	d.SetId(strconv.Itoa(role.ID))
	if err := d.Set("name", role.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope", role.Scope); err != nil {
		return diag.FromErr(err)
	}

	if role.Type == "" || !strings.HasPrefix(role.Type, "Role::V2::") {
		// do nothing, do not set the role.Type to "type" the state
		// this should not happen though; just adding this to prevent errors
	}
	if err := d.Set("type", strings.ToUpper(strings.TrimPrefix(role.Type, "Role::V2::"))); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
