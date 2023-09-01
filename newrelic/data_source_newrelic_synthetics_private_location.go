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
			"key": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Optional:    true,
				Description: "The key of the queried private location.",
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

	key := fetchPrivateLocationKey(location.Tags)
	if len(key) == 0 {
		//logs the absence of a key but does not throw an error to prevent halting execution
		log.Printf("[INFO] No keys found corresponding to the queried private location.")
	}
	err = d.Set("key", key)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// fetchPrivateLocationKey extracts the 'key' of the private location by iterating through all tags
// and finding the tag by the name 'key', the 'values' of which would contain the key of the private location
func fetchPrivateLocationKey(tagList []entities.EntityTag) []string {
	var key []string
	for _, tag := range tagList {
		if tag.Key == "key" {
			// tag.Values is a list of strings returned by the API, though it is expected to contain only one string
			key = tag.Values
		}
	}
	return key
}
