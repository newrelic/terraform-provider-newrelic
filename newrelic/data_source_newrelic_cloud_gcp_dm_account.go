package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// gcpDmListLinkedAccountsQuery fetches only id/name/nrAccountId for all gcp_v2
// linked accounts. Requesting the integrations field causes a client-go
// deserialization error on CloudGcpGenericIntegration, so we omit it entirely.
const gcpDmListLinkedAccountsQuery = `query($accountId: Int!) {
	actor {
		account(id: $accountId) {
			cloud {
				linkedAccounts {
					id
					name
					nrAccountId
				}
			}
		}
	}
}`

type gcpDmListLinkedAccountsResp struct {
	Actor struct {
		Account struct {
			Cloud struct {
				LinkedAccounts []struct {
					ID          int    `json:"id"`
					Name        string `json:"name"`
					NrAccountId int    `json:"nrAccountId"`
				} `json:"linkedAccounts"`
			} `json:"cloud"`
		} `json:"account"`
	} `json:"actor"`
}

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

	var resp gcpDmListLinkedAccountsResp
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx,
		gcpDmListLinkedAccountsQuery,
		map[string]interface{}{"accountId": accountID},
		&resp,
	); err != nil {
		return diag.FromErr(fmt.Errorf("listing GCP Dimensional Metrics linked accounts: %w", err))
	}

	for _, account := range resp.Actor.Account.Cloud.LinkedAccounts {
		if strings.EqualFold(account.Name, targetName) {
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
