package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
)

func resourceNewRelicAPIAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAPIAccessKeyCreate,
		ReadContext:   resourceNewRelicAPIAccessKeyRead,
		UpdateContext: resourceNewRelicAPIAccessKeyUpdate,
		DeleteContext: resourceNewRelicAPIAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNewrelicAPIAccessKeyImport,
		},
		Schema:        resourceNewRelicAPIAccessKeySchema(),
		CustomizeDiff: validateAPIAccessKeyAttributes,
	}
}

func resourceNewRelicAPIAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	accountID := getAPIAccessKeyAccountID(d)
	keyType := getAPIAccessKeyType(d)

	opts := apiaccess.APIAccessCreateInput{}

	switch keyType {
	case keyTypeIngest:
		ingestKeyOpts := apiaccess.APIAccessCreateIngestKeyInput{}

		ingestKeyOpts.AccountID = accountID
		ingestKeyOpts.IngestType = apiaccess.APIAccessIngestKeyType(getAPIAccessIngestType(d))
		ingestKeyOpts.Name = getAPIAccessKeyName(d)
		ingestKeyOpts.Notes = getAPIAccessKeyNotes(d)

		opts.Ingest = []apiaccess.APIAccessCreateIngestKeyInput{ingestKeyOpts}
	case keyTypeUser:
		userKeyOpts := apiaccess.APIAccessCreateUserKeyInput{}

		userKeyOpts.AccountID = accountID
		userKeyOpts.Name = getAPIAccessKeyName(d)
		userKeyOpts.Notes = getAPIAccessKeyNotes(d)
		userKeyOpts.UserID = getAPIAccessUserID(d)

		opts.User = []apiaccess.APIAccessCreateUserKeyInput{userKeyOpts}
	default:
		err := fmt.Errorf("unknown api access key type: %s", keyType)
		return diag.FromErr(err)
	}

	keys, createErr := client.APIAccess.CreateAPIAccessKeysWithContext(ctx, opts)
	if createErr != nil {
		return diag.FromErr(createErr)
	}

	// Validate to make sure we only created one key.
	if len(keys) != 1 {
		err := fmt.Errorf("expected 1 new key, got %d", len(keys))
		return diag.FromErr(err)
	}

	// Set the resource ID to be a composite of the key ID and the key type in order to lookup the newly created key
	d.SetId(keys[0].ID)

	// Expose the API key only at creation time; subsequent reads will not refresh/overwrite it.
	if err := d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.Key, keys[0].Key); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicAPIAccessKeyRead(ctx, d, meta)
}

func resourceNewRelicAPIAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	key, readErr := client.APIAccess.GetAPIAccessKeyWithContext(ctx, d.Id(), apiaccess.APIAccessKeyType(getAPIAccessKeyType(d)))
	if readErr != nil {
		if strings.Contains(readErr.Error(), "Key not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(readErr)
	}
	if key == nil {
		return diag.FromErr(fmt.Errorf("no New Relic API Access Key found with given id %s", d.Id()))
	}

	var setErr error
	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.AccountID, key.AccountID)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.KeyType, key.Type)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType, key.IngestType)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.UserID, key.UserID)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.Name, key.Name)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.Notes, key.Notes)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Intentionally do not set the "key" on Read to avoid drift
	// We shall set the "key" in the state only if it is not present already,
	// implying that we would do it on an import (as a create would have the "key" set in the state).
	///////////////////////////////////////////////////////////////////////////////////////////////

	// check if the state already has the key value set (it would only be set at creation time or if manually set in the state)
	// this would help enter the nested conditions only if the key value is not already set in the state, i.e. during an "import".
	if _, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.Key); !ok {
		// if the key value returned by NerdGraph is truncated (which happens when the key belongs to a user that is not
		// associated with the API Key used to perform this read operation via the provider), then emit a warning
		// and do not set the key value in the state (as it would be a truncated value and not the actual key value).
		if strings.HasSuffix(key.Key, "...") {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary: fmt.Sprintf("Importing an API Key belonging to a user that is not associated with the API Key " +
						"used to perform this import operation via the provider is not supported, as NerdGraph will only return a " +
						"truncated key value in such a case, for security reasons.\n" +
						"Please see this article for more information: " +
						"https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-user-keys/#query-keys"),
				}}
		}
		// (else) if the key value is not truncated, then set it in the state (this would be the case when the key belongs to
		// the user associated with the API Key used to perform this read operation via the provider).
		setErr = d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.Key, key.Key)
		if setErr != nil {
			return diag.FromErr(setErr)
		}

	}

	return nil
}

func resourceNewRelicAPIAccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	opts := apiaccess.APIAccessUpdateInput{}

	keyType := getAPIAccessKeyType(d)

	// Construct the key type specific update opts.
	switch keyType {
	case keyTypeIngest:
		updateKeyOpts := apiaccess.APIAccessUpdateIngestKeyInput{
			KeyID: d.Id(),
			Name:  getAPIAccessKeyName(d),
			Notes: getAPIAccessKeyNotes(d),
		}
		opts.Ingest = []apiaccess.APIAccessUpdateIngestKeyInput{updateKeyOpts}
	case keyTypeUser:
		updateKeyOpts := apiaccess.APIAccessUpdateUserKeyInput{
			KeyID: d.Id(),
			Name:  getAPIAccessKeyName(d),
			Notes: getAPIAccessKeyNotes(d),
		}
		opts.User = []apiaccess.APIAccessUpdateUserKeyInput{updateKeyOpts}
	default:
		return diag.Errorf("unknown api access key type: %s", keyType)
	}

	keys, updateErr := client.APIAccess.UpdateAPIAccessKeysWithContext(ctx, opts)
	if updateErr != nil {
		return diag.FromErr(updateErr)
	}

	// Validate to make sure we only updated one key.
	if len(keys) != 1 {
		return diag.Errorf("expected 1 key, got %d", len(keys))
	}

	return resourceNewRelicAPIAccessKeyRead(ctx, d, meta)
}

func resourceNewRelicAPIAccessKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	opts := apiaccess.APIAccessDeleteInput{}

	keyType := getAPIAccessKeyType(d)

	// Construct the key type specific delete opts.
	switch keyType {
	case keyTypeIngest:
		opts.IngestKeyIDs = []string{d.Id()}
	case keyTypeUser:
		opts.UserKeyIDs = []string{d.Id()}
	default:
		return diag.Errorf("unknown api access key type: %s", keyType)
	}

	_, deleteErr := client.APIAccess.DeleteAPIAccessKeyWithContext(ctx, opts)
	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")

	return nil
}

func resourceNewrelicAPIAccessKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keyID, keyType, parseErr := parseCompositeID(d.Id())
	if parseErr != nil {
		return nil, parseErr
	}

	d.SetId(keyID)

	setErr := d.Set(ResourceNewRelicAPIAccessKeyAttributeLabels.KeyType, strings.ToUpper(keyType))
	if setErr != nil {
		return nil, setErr
	}

	readDiagnostics := resourceNewRelicAPIAccessKeyRead(ctx, d, meta)
	if readDiagnostics.HasError() || d.Id() == "" {
		return nil, fmt.Errorf("unable to find New Relic API Access Key with given id %s and type %s", keyID, keyType)
	}

	return []*schema.ResourceData{d}, nil
}
