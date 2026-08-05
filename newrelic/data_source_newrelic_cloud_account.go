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
					"(cloud_provider must be \"gcp\").",
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

	// is_dimensional_metrics only applies to GCP. Data sources do not run
	// CustomizeDiff, so this relationship is validated here.
	if !isDimensionalMetricsProviderValid(provider, isDM) {
		return diag.Errorf("`is_dimensional_metrics` can only be set when `cloud_provider` is \"gcp\"")
	}

	accounts, err := client.Cloud.GetLinkedAccountsWithContext(ctx, provider)
	if err != nil {
		return diag.FromErr(err)
	}

	// GCP linked accounts are returned together regardless of whether they use
	// Dimensional Metrics; HasDimensionalMetrics disambiguates them so a lookup
	// for one kind never matches an account of the other kind.
	account := findCloudLinkedAccount(*accounts, accountID, name, isDM)

	if account == nil {
		if isDM {
			return diag.Errorf("no GCP Dimensional Metrics linked account named %q found for New Relic account %d", name, accountID)
		}
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any account for provider '%s'", name, provider))
	}

	d.SetId(strconv.Itoa(account.ID))

	return diag.FromErr(flattenCloudAccount(account, d))
}

// isDimensionalMetricsProviderValid reports whether is_dimensional_metrics may be
// set for the given cloud_provider value; it can only be true for GCP.
func isDimensionalMetricsProviderValid(provider string, isDM bool) bool {
	return !isDM || strings.EqualFold(provider, "gcp")
}

// findCloudLinkedAccount returns the linked account matching accountID and name
// whose HasDimensionalMetrics flag equals isDM, so a Dimensional Metrics lookup
// never matches a legacy account of the same name (and vice versa). Returns nil
// if there is no such match.
func findCloudLinkedAccount(accounts []cloud.CloudLinkedAccount, accountID int, name string, isDM bool) *cloud.CloudLinkedAccount {
	for _, a := range accounts {
		if a.NrAccountId == accountID && strings.EqualFold(a.Name, name) && a.HasDimensionalMetrics == isDM {
			return &a
		}
	}
	return nil
}

func flattenCloudAccount(account *cloud.CloudLinkedAccount, d *schema.ResourceData) error {
	if err := d.Set("name", account.Name); err != nil {
		return err
	}

	if err := d.Set("account_id", account.NrAccountId); err != nil {
		return err
	}

	return nil
}
