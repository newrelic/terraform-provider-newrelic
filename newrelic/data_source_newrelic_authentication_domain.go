package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicAuthenticationDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAuthenticationDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the authentication domain to be queried.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the fetched authentication domain.",
			},
		},
	}
}

func dataSourceNewRelicAuthenticationDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	matchingAuthenticationDomainID := ""

	log.Printf("[INFO] Fetching Authentication Domains")

	name, nameOk := d.GetOk("name")
	if !nameOk {
		return diag.FromErr(errors.New("`name` is required"))
	}

	resp, err := client.UserManagement.GetAuthenticationDomainsWithContext(ctx, "", []string{})

	if resp == nil {
		return diag.FromErr(fmt.Errorf("failed to fetch authentication domains"))
	}

	if err != nil {
		return diag.FromErr(err)
	}

	for _, authenticationDomain := range resp.AuthenticationDomains {
		if name == authenticationDomain.Name {
			matchingAuthenticationDomainID = authenticationDomain.ID
			break
		}
	}

	if matchingAuthenticationDomainID == "" {
		return diag.FromErr(fmt.Errorf("no authentication domain found with the name %s", name))
	}

	d.SetId(matchingAuthenticationDomainID)

	return nil
}
