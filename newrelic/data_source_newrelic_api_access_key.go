package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
)

func dataSourceNewRelicAPIAccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAPIAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID the key belongs to. Defaults to the account ID configured on the provider when not specified.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the key. When specified, the key is fetched directly by its ID instead of searching by other attributes.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("The type of the key, one of %s or %s.", keyTypeIngest, keyTypeUser),
			},
			"ingest_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("The type of the ingest key, one of %s or %s. Only applies when `key_type` is %s.", keyTypeIngestLicense, keyTypeIngestBrowser, keyTypeIngest),
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("The ID of the user that owns the key. Only applies when `key_type` is %s.", keyTypeUser),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the key. Used to narrow down the search when `key_id` is not specified.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Any notes attached to the key.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The value of the key.",
			},
		},
	}
}

func dataSourceNewRelicAPIAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	keyType := strings.ToUpper(d.Get("key_type").(string))
	if keyType != keyTypeIngest && keyType != keyTypeUser {
		return diag.Errorf("the `key_type` attribute must be set to either %s or %s", keyTypeIngest, keyTypeUser)
	}

	var key *apiaccess.APIKey
	var err error

	// When a `key_id` is provided, fetch the key directly by its ID. Otherwise, search for a
	// matching key using the other attributes provided in the configuration.
	if keyID, ok := d.GetOk("key_id"); ok {
		key, err = client.APIAccess.GetAPIAccessKeyWithContext(ctx, keyID.(string), apiaccess.APIAccessKeyType(keyType))
		if err != nil {
			return diag.FromErr(err)
		}
		if key == nil {
			return diag.FromErr(fmt.Errorf("no New Relic API access key found with id %s and type %s", keyID, keyType))
		}
	} else {
		key, err = searchNewRelicAPIAccessKey(ctx, providerConfig, d, keyType)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return setNewRelicAPIAccessKeyDataSourceAttributes(d, key)
}

// searchNewRelicAPIAccessKey searches for a single API access key matching the attributes
// provided in the configuration. It returns an error when no key or more than one key matches.
func searchNewRelicAPIAccessKey(ctx context.Context, providerConfig *ProviderConfig, d *schema.ResourceData, keyType string) (*apiaccess.APIKey, error) {
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	query := apiaccess.APIAccessKeySearchQuery{
		Types: []apiaccess.APIAccessKeyType{apiaccess.APIAccessKeyType(keyType)},
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs: []int{accountID},
		},
	}

	ingestType, ingestTypeOk := d.GetOk("ingest_type")
	if ingestTypeOk {
		query.Scope.IngestTypes = []apiaccess.APIAccessIngestKeyType{apiaccess.APIAccessIngestKeyType(strings.ToUpper(ingestType.(string)))}
	}

	userID, userIDOk := d.GetOk("user_id")
	if userIDOk {
		query.Scope.UserIDs = []int{userID.(int)}
	}

	keys, err := client.APIAccess.SearchAPIAccessKeysWithContext(ctx, query)
	if err != nil {
		return nil, err
	}

	name, nameOk := d.GetOk("name")

	matches := make([]apiaccess.APIKey, 0, len(keys))
	for _, k := range keys {
		// The search is scoped on the server side, but `name` is matched here as the API does
		// not support filtering by name.
		if nameOk && k.Name != name.(string) {
			continue
		}
		matches = append(matches, k)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no New Relic API access key found matching the given criteria; try specifying `key_id`, `name`, `ingest_type` or `user_id` to narrow down the search")
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("found %d New Relic API access keys matching the given criteria; please specify additional attributes such as `key_id`, `name`, `ingest_type` or `user_id` to narrow down the search to a single key", len(matches))
	}

	return &matches[0], nil
}

func setNewRelicAPIAccessKeyDataSourceAttributes(d *schema.ResourceData, key *apiaccess.APIKey) diag.Diagnostics {
	d.SetId(key.ID)

	if err := d.Set("key_id", key.ID); err != nil {
		return diag.FromErr(err)
	}
	if key.AccountID != nil {
		if err := d.Set("account_id", *key.AccountID); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("key_type", key.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ingest_type", key.IngestType); err != nil {
		return diag.FromErr(err)
	}
	if key.UserID != nil {
		if err := d.Set("user_id", *key.UserID); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("name", key.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", key.Notes); err != nil {
		return diag.FromErr(err)
	}

	// NerdGraph returns a truncated key value (suffixed with "...") when the key belongs to a user
	// that is not associated with the API key used by the provider. In that case the actual key
	// value cannot be retrieved, so emit a warning instead of setting a truncated value.
	if strings.HasSuffix(key.Key, "...") {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "The API access key value could not be retrieved as it belongs to a user that is not associated " +
					"with the API key used by the provider. NerdGraph only returns a truncated key value in such a case, for " +
					"security reasons.\n" +
					"Please see this article for more information: " +
					"https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys/#query-keys",
			}}
	}

	if err := d.Set("key", key.Key); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
