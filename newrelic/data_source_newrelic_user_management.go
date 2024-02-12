package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicUserRead,
		Schema: map[string]*schema.Schema{
			"authentication_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Authentication Domain the user being queried would belong to.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the user to be queried.",
				AtLeastOneOf: []string{"name", "email_id"},
			},
			"email_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The email ID of the user to be queried.",
				AtLeastOneOf: []string{"name", "email_id"},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the fetched user.",
			},
		},
	}
}

func dataSourceNewRelicUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Fetching Users")

	name, nameOk := d.GetOk("name")
	email, emailOk := d.GetOk("email_id")
	authDomainID, authDomainIDOk := d.GetOk("authentication_domain_id")

	nameQuery := ""
	emailQuery := ""

	if nameOk && name != "" {
		nameQuery = name.(string)
	}

	if emailOk && email != "" {
		emailQuery = email.(string)
	}

	if !authDomainIDOk {
		return diag.FromErr(fmt.Errorf("'authentication_domain_id' is required"))
	}

	authenticationDomainID := authDomainID.(string)
	userFound := false

	resp, err := client.UserManagement.UserManagementGetUsersWithContext(
		ctx,
		[]string{authenticationDomainID},
		[]string{},
		nameQuery,
		emailQuery,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("failed to fetch users"))
	}

	for _, authDomain := range resp.AuthenticationDomains {
		if authDomain.ID == authenticationDomainID {
			for _, u := range authDomain.Users.Users {
				d.SetId(u.ID)
				_ = d.Set("name", u.Name)
				_ = d.Set("email_id", u.Email)
				userFound = true
				return nil
			}
		}
	}

	if !userFound {
		return diag.FromErr(fmt.Errorf("no user found with the specified parameters"))
	}

	return nil
}
