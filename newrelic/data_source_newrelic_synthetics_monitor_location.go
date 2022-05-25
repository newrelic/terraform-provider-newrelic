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

			// The legacy attributes below have been deprecated and removed in NerdGraph.
			"high_security_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Represents if high security mode is enabled for the location. A value of true means that high security mode is enabled, and a value of false means it is disabled.",
				Deprecated:  "The `high_security_mode` field has been deprecated and no longer exists in the API response.",
			},
			"private": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Represents if this location is a private location. A value of true means that the location is private, and a value of false means it is public.",
				Deprecated:  "The `private` field has been deprecated and no longer exists in the API response.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the Synthetics monitor location.",
				Deprecated:  "The `description` field has been deprecated and no longer exists in the API response.",
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
		if loc.Name == label {
			location = loc
			break
		}
	}

	if location == nil {
		return diag.FromErr(fmt.Errorf("the provided `name` or `label` '%s' does not match any Synthetics monitor private location", name))
	}

	d.SetId(location.Name)

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
