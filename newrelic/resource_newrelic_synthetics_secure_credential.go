package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
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
				ForceNew:    true,
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
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The time the secure credential was last updated.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
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
	_ = d.Set("key", res.Key)
	_ = d.Set("description", res.Description)
	_ = d.Set("last_updated", time.Time(*res.LastUpdate).Format(time.RFC3339))
	_ = d.Set("account_id", accountID)

	return nil
}

func resourceNewRelicSyntheticsSecureCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	queryString := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = '%s' AND accountId = '%d'", d.Id(), accountID)

	var entity *entities.EntityOutlineInterface
	var entityResults *entities.EntitySearch
	var reqErr error
	var diags diag.Diagnostics

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		entityResults, reqErr = client.Entities.GetEntitySearchByQueryWithContext(ctx, entities.EntitySearchOptions{}, queryString, []entities.EntitySearchSortCriteria{})
		if reqErr != nil {
			return resource.NonRetryableError(reqErr)
		}

		// commenting this checks as this check disables us to detect the changes made via UI
		//if entityResults.Count != 1 {
		//	return resource.RetryableError(fmt.Errorf("failed to read secure credential"))
		//}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	if len(diags) > 0 {
		return diags
	}

	if entityResults.Count != 1 {
		d.SetId("")
		return nil
	}

	for _, e := range entityResults.Results.Entities {
		// Conditional on case-sensitive match
		if e.GetName() == d.Id() {
			entity = &e
			break
		}
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

	_ = d.Set("key", res.Key)
	_ = d.Set("description", res.Description)
	_ = d.Set("last_updated", time.Time(*res.LastUpdate).Format(time.RFC3339))
	_ = d.Set("account_id", accountID)

	return nil
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
	switch e := (*sc).(type) {
	case *entities.SecureCredentialEntityOutline:
		_ = d.Set("account_id", e.AccountID)
		_ = d.Set("key", e.GetName())
		_ = d.Set("description", e.Description)
		_ = d.Set("last_updated", time.Time(*e.UpdatedAt).Format(time.RFC3339))
	}

	return nil
}
