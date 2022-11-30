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

func dataSourceNewRelicSyntheticsPrivateLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicSyntheticsPrivateLocationRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Synthetics monitor private location.",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsPrivateLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading Synthetics monitor locations")

	name, nameOk := d.GetOk("name")
	if !nameOk {
		return diag.FromErr(errors.New("`name` is required"))
	}

	query := fmt.Sprintf("domain = 'SYNTH' AND type = 'PRIVATE_LOCATION' AND name = '%s'", name)
	entitySearch, err := client.Entities.GetEntitySearchByQueryWithContext(
		ctx,
		entities.EntitySearchOptions{},
		query,
		[]entities.EntitySearchSortCriteria{},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	privateLocations := entitySearch.Results.Entities
	var location *entities.GenericEntityOutline
	for _, l := range privateLocations {
		loc := l.(*entities.GenericEntityOutline)

		// It's possible to have multiple private locations with the same name.
		// Return the first matching private location.
		if loc.AccountID == accountID && loc.Name == name {
			location = loc
			break
		}
	}

	if location == nil {
		return diag.FromErr(fmt.Errorf("no matches found for private location with name '%s'", name))
	}

	d.SetId(string(location.GUID))

	err = d.Set("name", location.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
