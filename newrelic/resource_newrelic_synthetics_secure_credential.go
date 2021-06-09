package newrelic

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicSyntheticsSecureCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsSecureCredentialCreate,
		ReadContext:   resourceNewRelicSyntheticsSecureCredentialRead,
		UpdateContext: resourceNewRelicSyntheticsSecureCredentialUpdate,
		DeleteContext: resourceNewRelicSyntheticsSecureCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The secure credential's key name. Regardless of the case used in the configuration, the provider will provide an upcased key to the underlying API.",
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The secure credential's value.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The secure credential's description.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The time the secure credential was created.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The time the secure credential was last updated.",
			},
		},
	}
}

func resourceNewRelicSyntheticsSecureCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	sc := expandSyntheticsSecureCredential(d)

	log.Printf("[INFO] Creating New Relic Synthetics secure credential %s", sc.Key)

	sc, err := client.Synthetics.AddSecureCredential(sc.Key, sc.Value, sc.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sc.Key)
	return resourceNewRelicSyntheticsSecureCredentialRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential %s", d.Id())

	sc, err := client.Synthetics.GetSecureCredential(d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenSyntheticsSecureCredential(sc, d))
}

func resourceNewRelicSyntheticsSecureCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating New Relic Synthetics secure credential %s", d.Id())

	sc := expandSyntheticsSecureCredential(d)

	_, err := client.Synthetics.UpdateSecureCredential(sc.Key, sc.Value, sc.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicSyntheticsSecureCredentialRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics secure credential %s", d.Id())

	if err := client.Synthetics.DeleteSecureCredential(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
