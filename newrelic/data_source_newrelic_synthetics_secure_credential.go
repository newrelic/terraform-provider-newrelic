package newrelic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func dataSourceNewRelicSyntheticsSecureCredential() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicSyntheticsSecureCredentialRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicSyntheticsSecureCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential")

	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	sc, err := client.Synthetics.GetSecureCredential(key)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return err
		}
	}

	d.SetId(key)

	return flattenSyntheticsSecureCredential(sc, d)
}
