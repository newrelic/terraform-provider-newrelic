package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicCloudGcpV2Account() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicCloudGcpV2AccountRead,
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
				Description: "The name of the GCP v2 linked account to look up.",
			},
		},
	}
}

func dataSourceNewRelicCloudGcpV2AccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	targetName := d.Get("name").(string)

	// GCP v2 linked accounts use "gcp_v2" as the provider slug internally.
	// If the API does not return results for "gcp_v2", fall back to "gcp".
	linkedAccounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, "gcp_v2")
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
		"no GCP v2 linked account named %q found for New Relic account %d",
		targetName, accountID,
	)
}
