package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the secure credential was created.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the secure credential was last updated.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(15 * time.Second),
		},
	}
}

func dataSourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential")

	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = '%s'", d.Id())

	var entity *entities.EntityOutlineInterface

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		entityResults, err := client.Entities.GetEntitySearchByQueryWithContext(ctx, entities.EntitySearchOptions{}, queryString, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if entityResults.Count != 1 {
			return resource.RetryableError(fmt.Errorf("entity not found, or found more than one"))
		}
		for _, e := range entityResults.Results.Entities {
			// Conditional on case sensitive match
			if e.GetName() == d.Id() {
				entity = &e
				break
			}
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	d.SetId(key)

	return flattenSyntheticsSecureCredential(entity, d)
}
