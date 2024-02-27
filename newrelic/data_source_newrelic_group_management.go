package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicGroupRead,
		Schema: map[string]*schema.Schema{
			"authentication_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Authentication Domain the group being queried would belong to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group to be queried.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the fetched group.",
			},
			"user_ids": {
				Type:        schema.TypeList,
				Description: "IDs of users which belong to the group.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNewRelicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Fetching Groups")

	name, nameOk := d.GetOk("name")
	authDomainID, authDomainIDOk := d.GetOk("authentication_domain_id")

	if !authDomainIDOk {
		return diag.FromErr(fmt.Errorf("'authentication_domain_id' is required"))
	}
	if !nameOk {
		return diag.FromErr(fmt.Errorf("'name' is required"))
	}

	nameQuery := ""
	var usersInGroup []string

	authenticationDomainID := authDomainID.(string)
	nameQuery = name.(string)

	resp, err := client.UserManagement.UserManagementGetGroupsWithUsersWithContext(
		ctx,
		[]string{authenticationDomainID},
		[]string{},
		nameQuery,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("failed to fetch groups"))
	}

	groupFound := false
	for _, authDomain := range resp.AuthenticationDomains {
		if !groupFound && (authDomain.ID == authenticationDomainID) {
			for _, g := range authDomain.Groups.Groups {
				if !groupFound {
					d.SetId(g.ID)
					_ = d.Set("name", g.DisplayName)
					for _, u := range g.Users.Users {
						usersInGroup = append(usersInGroup, u.ID)
					}
					_ = d.Set("user_ids", usersInGroup)
					groupFound = true
				}
			}
		}
	}

	if !groupFound {
		return diag.FromErr(fmt.Errorf("no group found with the specified parameters"))
	}

	return nil

}
