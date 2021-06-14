package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
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
		},
	}
}

func dataSourceNewRelicKeyTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic key transactions")

	name := d.Get("name").(string)

	params := apm.ListKeyTransactionsParams{
		Name: name,
	}

	transactions, err := client.APM.ListKeyTransactionsWithContext(ctx, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	var transaction *apm.KeyTransaction

	for _, t := range transactions {
		if t.Name == name {
			transaction = t
			break
		}
	}

	if transaction == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic key transaction", name))
	}

	flattenKeyTransaction(transaction, d)

	return nil
}

func flattenKeyTransaction(t *apm.KeyTransaction, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(t.ID))
	_ = d.Set("name", t.Name)
}
