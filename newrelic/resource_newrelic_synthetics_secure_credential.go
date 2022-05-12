package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	sc := expandSyntheticsSecureCredential(d)

	log.Printf("[INFO] Creating New Relic Synthetics secure credential %s", sc.Key)

	res, err := client.Synthetics.SyntheticsCreateSecureCredentialWithContext(ctx, accountID, sc.Description, sc.Key, synthetics.SecureValue(sc.Value))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sc.Key)

	flattenSyntheticsSecureCredential(res, d)

	return nil
}

func resourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient

	queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED'AND name = %s", d.Id())

	entityResults, err := client.Entities.GetEntitySearchWithContext(ctx, entities.EntitySearchOptions{}, queryString, entities.EntitySearchQueryBuilder{}, []entities.EntitySearchSortCriteria{})
	if err != nil {
		return diag.FromErr(err)
	}

	var entity *entities.EntityOutlineInterface
	for _, e := range entityResults.Results.Entities {
		// Conditional on case sensitive match
		if e.GetName() == d.Id() {
			entity = &e
			break
		}
	}
	
	return nil
}

func resourceNewRelicSyntheticsSecureCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Updating New Relic Synthetics secure credential %s", d.Id())

	sc := expandSyntheticsSecureCredential(d)

	res, err := client.Synthetics.SyntheticsUpdateSecureCredentialWithContext(ctx, accountID, sc.Description, sc.Key, synthetics.SecureValue(sc.Value))
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSyntheticsSecureCredential(res, d)

	return nil
}

func resourceNewRelicSyntheticsSecureCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Deleting New Relic Synthetics secure credential %s", d.Id())

	res, err := client.Synthetics.SyntheticsDeleteSecureCredentialWithContext(ctx, accountID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		d.SetId("")
	}

	return nil
}

func expandSyntheticsSecureCredential(d *schema.ResourceData) *synthetics.SecureCredential {
	key := d.Get("key").(string)
	key = strings.ToUpper(key)

	sc := synthetics.SecureCredential{
		Key:         key,
		Value:       d.Get("value").(string),
		Description: d.Get("description").(string),
	}

	return &sc
}

func flattenSyntheticsSecureCredential(sc *synthetics.SyntheticsSecureCredentialMutationResult, d *schema.ResourceData) error {
	_ = d.Set("key", sc.Key)
	_ = d.Set("description", sc.Description)

	createdAt := &sc.CreatedAt
	_ = d.Set("created_at", createdAt)

	lastUpdated := &sc.LastUpdate
	_ = d.Set("last_updated", lastUpdated)

	return nil
}
