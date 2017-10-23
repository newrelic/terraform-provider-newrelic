package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/api"
)

func dataSourceNewRelicKeyTransaction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicKeyTransactionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNewRelicKeyTransactionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*newrelic.Client)

	log.Printf("[INFO] Reading New Relic key transactions")

	transactions, err := client.ListKeyTransactions()
	if err != nil {
		return err
	}

	var transaction *newrelic.KeyTransaction
	name := d.Get("name").(string)

	for _, t := range transactions {
		if t.Name == name {
			transaction = &t
			break
		}
	}

	if transaction == nil {
		return fmt.Errorf("The name '%s' does not match any New Relic key transaction.", name)
	}

	d.SetId(strconv.Itoa(transaction.ID))
	d.Set("name", transaction.Name)

	return nil
}
