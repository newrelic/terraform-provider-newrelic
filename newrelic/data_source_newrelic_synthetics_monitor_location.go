package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func dataSourceNewRelicSyntheticsMonitorLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicSyntheticsMonitorLocationRead,
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Synthetics monitor private location.",
				Deprecated:  "Use `name` attribute instead.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the Synthetics monitor private location.",
				ConflictsWith: []string{"label"},
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading Synthetics monitor locations")

	// Note: The legacy `label` field is the equivalent of `name` in NerdGraph
	label, labelOk := d.GetOk("label")
	name, nameOk := d.GetOk("name")
	if !labelOk && !nameOk {
		return diag.FromErr(errors.New("one of `label` or `name` must be configured"))
	}

	// If `label` was set in the data source config, set `name` to its value.
	if labelOk {
		name = label
	}

	query := fmt.Sprintf("domain = 'SYNTH' AND type = 'PRIVATE_LOCATION' AND name = '%s'", name)
	entitySearch, err := client.Entities.GetEntitySearchByQuery(
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
		if loc.Name == name {
			location = loc
			break
		}
	}

	if location == nil {
		return diag.FromErr(fmt.Errorf("no matches found for private location with name '%s'", name))
	}

	d.SetId(string(location.GUID))

	if labelOk {
		err = d.Set("label", location.Name)
	} else {
		err = d.Set("name", location.Name)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
