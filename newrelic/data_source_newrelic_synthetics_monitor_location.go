package newrelic

import (
	"context"
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
				Required:    true,
				Description: "The name of the Synthetics monitor location.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Synthetics monitor location.",
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
	label := d.Get("label").(string)

	query := fmt.Sprintf("domain = 'SYNTH' AND type = 'PRIVATE_LOCATION' AND name = '%s'", label)
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
		ll := l.(*entities.GenericEntityOutline)

		// It's possible to have multiple private locations with the same name.
		// Return the first matching private location.
		if ll.Name == label {
			location = ll
			break
		}
	}

	if location == nil {
		return diag.FromErr(fmt.Errorf("the provided label '%s' does not match any Synthetics monitor location tags", label))
	}

	d.SetId(location.Name)
	_ = d.Set("name", location.Name)
	_ = d.Set("label", location.Name)

	// THESE FIELDS NO LONGER EXIST IN THE NEW NERDGRAPH RESPONSE
	// _ = d.Set("high_security_mode", location.HighSecurityMode)
	// _ = d.Set("private", location.Private)
	// _ = d.Set("description", location.Description)

	return nil
}
