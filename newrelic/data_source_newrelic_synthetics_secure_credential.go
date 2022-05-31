package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func dataSourceNewRelicSyntheticsSecureCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicSyntheticsSecureCredentialRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The secure credential's key name. Regardless of the case used in the configuration, the provider will provide an upcased key to the underlying API.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secure credential's description.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the secure credential was last updated.",
			},
		},
	}
}

func dataSourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential")

	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = '%s'", key)

	entityResults, err := client.Entities.GetEntitySearchByQueryWithContext(ctx, entities.EntitySearchOptions{}, queryString, []entities.EntitySearchSortCriteria{})
	if err != nil {
		return diag.FromErr(err)
	}

	var entity *entities.EntityOutlineInterface
	for _, e := range entityResults.Results.Entities {
		// Conditional on case sensitive match
		if e.GetName() == key {
			entity = &e
			break
		}
	}

	d.SetId(key)

	return flattenSyntheticsSecureCredential(entity, d)
}
