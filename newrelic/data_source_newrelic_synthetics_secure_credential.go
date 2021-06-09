package newrelic

import (
	"context"
	"log"
	"strings"

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
	}
}

func dataSourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential")

	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	sc, err := client.Synthetics.GetSecureCredential(key)
	if err != nil {
		return diag.Errorf("the key '%s' does not match any New Relic Synthetics secure credential", key)
	}

	d.SetId(key)

	return diag.FromErr(flattenSyntheticsSecureCredential(sc, d))
}
