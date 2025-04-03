package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func dataSourceNewRelicKeyTransaction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicKeyTransactionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key transaction in New Relic.",
			},
			"guid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "GUID of the key transaction in New Relic.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the New Relic account the key transaction would need to belong to. Uses the account_id in the provider{} block by default, if not specified.",
				Optional:    true,
				Computed:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Domain of the key transaction in New Relic.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Entity type of the key transaction in New Relic.",
			},
		},
	}
}

func dataSourceNewRelicKeyTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic key transactions")
	// fetch name from TF config
	name, nameOk := d.GetOk("name")
	if !nameOk {
		return diag.FromErr(errors.New("`name` is required"))
	}
	// fetch guid from TF config
	guid := d.Get("guid")
	query := ""

	if guid != "" {
		// if guid is provided, irrespective of other arguments added, add it to the query to filter the results, while also adding
		// the KEY_TRANSACTION type filter to make sure the provided GUID is that of a real key transaction
		query = fmt.Sprintf("type = 'KEY_TRANSACTION' AND id = '%s'", guid)
	} else {
		// if the GUID is not given, use all the other arguments specified, i.e. name, which is required, and accountId, though optional,
		// shall be picked from account_id in the configuration, or the account_id in the provider's configuration if not specified
		query = fmt.Sprintf("type = 'KEY_TRANSACTION' AND name = '%s' AND accountId = %d", name, accountID)
	}

	keyTransactionsFound, err := client.Entities.GetEntitySearchByQueryWithContext(ctx, entities.EntitySearchOptions{}, query, []entities.EntitySearchSortCriteria{})
	if err != nil {
		return diag.FromErr(err)
	}

	if keyTransactionsFound == nil || len(keyTransactionsFound.Results.Entities) == 0 {
		return diag.FromErr(fmt.Errorf("no key transaction that matches the specified parameters is found in New Relic"))
	}

	flattenKeyTransaction(keyTransactionsFound, d)

	return nil
}

func flattenKeyTransaction(t *entities.EntitySearch, d *schema.ResourceData) {
	// iterate over the tags to get the key transaction id
	for _, tag := range t.Results.Entities[0].GetTags() {
		if tag.Key == "keyTransactionId" {
			d.SetId(tag.Values[0])
			break
		}
	}
	_ = d.Set("guid", t.Results.Entities[0].GetGUID())
	_ = d.Set("name", t.Results.Entities[0].GetName())
	_ = d.Set("account_id", t.Results.Entities[0].GetAccountID())
	_ = d.Set("domain", t.Results.Entities[0].GetDomain())
	_ = d.Set("type", t.Results.Entities[0].GetType())
}
