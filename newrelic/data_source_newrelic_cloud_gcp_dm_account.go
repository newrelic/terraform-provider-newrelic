package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicCloudGcpDmAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicCloudGcpDmAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID. Defaults to the provider account.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GCP Dimensional Metrics linked account to look up.",
			},
		},
	}
}

func dataSourceNewRelicCloudGcpDmAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	targetName := d.Get("name").(string)

	// GCP Dimensional Metrics linked accounts use "gcp_dm" as the provider slug internally.
	// If the API does not return results for "gcp_dm", fall back to "gcp".
	linkedAccounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, "gcp_dm")
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetLinkedAccounts failed: %w", err))
	}

	for _, account := range *linkedAccounts {
		if strings.EqualFold(account.Name, targetName) && account.NrAccountId == accountID {
			d.SetId(strconv.Itoa(account.ID))
			_ = d.Set("account_id", account.NrAccountId)
			_ = d.Set("name", account.Name)
			return nil
		}
	}

	return diag.Errorf(
		"no GCP Dimensional Metrics linked account named %q found for New Relic account %d",
		targetName, accountID,
	)
}
