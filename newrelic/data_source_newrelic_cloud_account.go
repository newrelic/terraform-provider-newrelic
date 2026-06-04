package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

// gcpDmLinkedAccountsQuery fetches linked accounts for GCP Dimensional Metrics
// (provider slug "gcp_v2") without requesting the integrations field.
// GetLinkedAccountsWithContext("gcp_v2") fails because client-go cannot
// deserialize CloudGcpGenericIntegration in the integrations field.
// gcpDmLinkedAccountsQuery uses the actor-level cloud path (not the account-scoped
// path) because only actor.cloud.linkedAccounts accepts the provider argument.
const gcpDmLinkedAccountsQuery = `query($provider: String) {
	actor {
		cloud {
			linkedAccounts(provider: $provider) {
				id
				name
				nrAccountId
			}
		}
	}
}`

type gcpDmLinkedAccountsResp struct {
	Actor struct {
		Cloud struct {
			LinkedAccounts []struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				NrAccountId int    `json:"nrAccountId"`
			} `json:"linkedAccounts"`
		} `json:"cloud"`
	} `json:"actor"`
}

func dataSourceNewRelicCloudAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicCloudAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the New Relic account.",
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cloud provider of the account, e.g. aws, gcp, azure",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cloud account.",
			},
			"is_dimensional_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Set to true when looking up a GCP Dimensional Metrics linked account " +
					"(cloud_provider must be \"gcp\"). Internally uses the gcp_v2 provider slug.",
			},
		},
	}
}

func dataSourceNewRelicCloudAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*ProviderConfig)
	client := cfg.NewClient

	log.Printf("[INFO] Reading New Relic Cloud Accounts")

	name := d.Get("name").(string)
	provider := d.Get("cloud_provider").(string)
	accountID := selectAccountID(cfg, d)
	isDM := d.Get("is_dimensional_metrics").(bool)

	// GCP Dimensional Metrics accounts use the "gcp_v2" provider slug internally.
	// GetLinkedAccountsWithContext("gcp_v2") fails because client-go cannot
	// deserialize CloudGcpGenericIntegration, so we use a raw NerdGraph query.
	if isDM && strings.EqualFold(provider, "gcp") {
		var resp gcpDmLinkedAccountsResp
		if err := client.NerdGraph.QueryWithResponseAndContext(ctx,
			gcpDmLinkedAccountsQuery,
			map[string]interface{}{"provider": "gcp_v2"},
			&resp,
		); err != nil {
			return diag.FromErr(fmt.Errorf("listing GCP Dimensional Metrics linked accounts: %w", err))
		}

		for _, a := range resp.Actor.Cloud.LinkedAccounts {
			if a.NrAccountId == accountID && strings.EqualFold(a.Name, name) {
				d.SetId(strconv.Itoa(a.ID))
				_ = d.Set("account_id", a.NrAccountId)
				_ = d.Set("name", a.Name)
				return nil
			}
		}

		return diag.Errorf("no GCP Dimensional Metrics linked account named %q found for New Relic account %d", name, accountID)
	}

	accounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, provider)
	if err != nil {
		return diag.FromErr(err)
	}

	var account *cloud.CloudLinkedAccount
	for _, a := range *accounts {
		if a.NrAccountId == accountID && strings.EqualFold(a.Name, name) {
			account = &a
			break
		}
	}

	if account == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any account for provider '%s", name, provider))
	}

	d.SetId(strconv.Itoa(account.ID))

	return diag.FromErr(flattenCloudAccount(account, d, accountID))
}

func flattenCloudAccount(account *cloud.CloudLinkedAccount, d *schema.ResourceData, accountID int) error {
	var err error

	err = d.Set("name", account.Name)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	return nil
}
