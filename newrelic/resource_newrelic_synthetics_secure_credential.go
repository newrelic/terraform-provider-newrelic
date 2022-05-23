package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
	"log"
	"strings"
	"time"

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
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to create the secure credential.",
			},
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
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Second),
		},
	}
}

func resourceNewRelicSyntheticsSecureCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	sc := expandSyntheticsSecureCredential(d)

	log.Printf("[INFO] Creating New Relic Synthetics secure credential %s", sc.Key)

	var diags diag.Diagnostics

	res, err := client.Synthetics.SyntheticsCreateSecureCredentialWithContext(ctx, accountID, sc.Description, sc.Key, synthetics.SecureValue(sc.Value))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	d.SetId(res.Key)

	_ = d.Set("created_at", time.Time(res.CreatedAt).Format(time.RFC3339))

	return resourceNewRelicSyntheticsSecureCredentialRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient

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

	return flattenSyntheticsSecureCredential(entity, d)
}

func resourceNewRelicSyntheticsSecureCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Updating New Relic Synthetics secure credential %s", d.Id())

	sc := expandSyntheticsSecureCredential(d)

	var diags diag.Diagnostics

	res, err := client.Synthetics.SyntheticsUpdateSecureCredentialWithContext(ctx, accountID, sc.Description, sc.Key, synthetics.SecureValue(sc.Value))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Errors) > 0 {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	return resourceNewRelicSyntheticsSecureCredentialRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Deleting New Relic Synthetics secure credential %s", d.Id())

	var diags diag.Diagnostics

	res, _ := client.Synthetics.SyntheticsDeleteSecureCredentialWithContext(ctx, accountID, d.Id())

	if res != nil {
		for _, err := range res.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Description,
			})
		}
	}

	if len(diags) > 0 {
		return diags
	}

	d.SetId("")

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

func flattenSyntheticsSecureCredential(sc *entities.EntityOutlineInterface, d *schema.ResourceData) diag.Diagnostics {
	_ = d.Set("key", (*sc).GetName())

	switch e := (*sc).(type) {
	case *entities.SecureCredentialEntityOutline:
		_ = d.Set("description", e.Description)
		_ = d.Set("last_updated", time.Time(*e.UpdatedAt).Format(time.RFC3339))
	}

	return nil
}
