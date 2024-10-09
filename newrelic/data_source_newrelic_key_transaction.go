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
			"domain": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				Description: "The Domain of the key transaction in New Relic.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				Description: "The Entity type of the key transaction in New Relic.",
			},
		},
	}
}

func dataSourceNewRelicKeyTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic key transactions")
	// fetch name from TF config
	name, nameOk := d.GetOk("name")
	if !nameOk {
		return diag.FromErr(errors.New("`name` is required"))
	}
	// fetch guid from TF config
	guid := d.Get("guid")
	query := fmt.Sprintf("type = 'KEY_TRANSACTION' AND name = '%s'", name)
	if guid != "" {
		// if guid is provided, add it to the query to filter the results
		query = fmt.Sprintf("%s AND id = '%s'", query, guid)
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
	_ = d.Set("domain", t.Results.Entities[0].GetDomain())
	_ = d.Set("type", t.Results.Entities[0].GetType())
}
