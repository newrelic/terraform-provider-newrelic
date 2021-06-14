package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/apiaccess"
)

var (
	keyTypeIngest        = string(apiaccess.APIAccessKeyTypeTypes.INGEST)
	keyTypeIngestBrowser = string(apiaccess.APIAccessIngestKeyTypeTypes.BROWSER)
	keyTypeIngestLicense = string(apiaccess.APIAccessIngestKeyTypeTypes.LICENSE)
	keyTypeUser          = string(apiaccess.APIAccessKeyTypeTypes.USER)
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

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"key_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{keyTypeIngest, keyTypeUser}, false),
			},

			"ingest_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"user_id"},
				ValidateFunc:  validation.StringInSlice([]string{keyTypeIngestBrowser, keyTypeIngestLicense}, false),
			},

			"user_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"ingest_type"},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceNewrelicAPIAccessKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keyID, keyType, parseErr := parseCompositeID(d.Id())
	if parseErr != nil {
		return nil, parseErr
	}

	d.SetId(keyID)

	setErr := d.Set("key_type", strings.ToUpper(keyType))
	if setErr != nil {
		return nil, setErr
	}

	diag := resourceNewRelicAPIAccessKeyRead(ctx, d, meta)
	if diag.HasError() {
		return nil, errors.New("error reading after import")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceNewRelicAPIAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	// Define initial keys to create an API access key.
	opts := apiaccess.APIAccessCreateInput{}

	// Get the account id. This is required regardless of key type.
	var accountID int
	if v, ok := d.GetOk("account_id"); ok {
		accountID = v.(int)
		log.Printf("[DEBUG] new api access account_id: %d", accountID)
	}

	// Get the key type.
	keyType := getAPIAccessKeyType(d)

	// Validate to make sure the following:
	// - If key_type is set to INGEST, `ingest_type` must be set to a value.
	// - If key_type is set to USER, `user_id` must be set to a value.
	switch keyType {
	case keyTypeIngest:
		ingestKeyOpts := apiaccess.APIAccessCreateIngestKeyInput{}

		if v, ok := d.GetOk("ingest_type"); ok {
			ingestKeyOpts.IngestType = apiaccess.APIAccessIngestKeyType(v.(string))
			log.Printf("[DEBUG] new api access ingest_type: %s", ingestKeyOpts.IngestType)
		} else {
			return diag.Errorf("[ERROR] you must define the ingest_type attribute when creating an INGEST key")
		}

		ingestKeyOpts.AccountID = accountID
		ingestKeyOpts.Name = getAPIAccessKeyName(d)
		ingestKeyOpts.Notes = getAPIAccessKeyNotes(d)
		opts.Ingest = []apiaccess.APIAccessCreateIngestKeyInput{ingestKeyOpts}
	case keyTypeUser:
		userKeyOpts := apiaccess.APIAccessCreateUserKeyInput{}

		if v, ok := d.GetOk("user_id"); ok {
			userKeyOpts.UserID = v.(int)
			log.Printf("[DEBUG] new api access user_id: %d", userKeyOpts.UserID)
		} else {
			return diag.Errorf("[ERROR] you must define the user_id attribute when creating an USER key")
		}

		userKeyOpts.AccountID = accountID
		userKeyOpts.Name = getAPIAccessKeyName(d)
		userKeyOpts.Notes = getAPIAccessKeyNotes(d)
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

	return resourceNewRelicAPIAccessKeyRead(ctx, d, meta)
}

func resourceNewRelicAPIAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	key, readErr := client.APIAccess.GetAPIAccessKeyWithContext(ctx, d.Id(), apiaccess.APIAccessKeyType(getAPIAccessKeyType(d)))
	if readErr != nil {
		return diag.FromErr(readErr)
	}

	var setErr error
	setErr = d.Set("account_id", key.AccountID)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("key_type", key.Type)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("ingest_type", key.IngestType)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("user_id", key.UserID)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("name", key.Name)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("notes", key.Notes)
	if setErr != nil {
		return diag.FromErr(setErr)
	}

	setErr = d.Set("key", key.Key)
	if setErr != nil {
		return diag.FromErr(setErr)
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
		updateKeyOpts := apiaccess.APIAccessUpdateIngestKeyInput{}
		updateKeyOpts.KeyID = d.Id()

		// Check if the following attributes have changed. If there are changes, get the new attribute values.
		if d.HasChange("name") {
			updateKeyOpts.Name = d.Get("name").(string)
		}

		if d.HasChange("notes") {
			updateKeyOpts.Notes = d.Get("notes").(string)
		}

		opts.Ingest = []apiaccess.APIAccessUpdateIngestKeyInput{updateKeyOpts}
	case keyTypeUser:
		updateKeyOpts := apiaccess.APIAccessUpdateUserKeyInput{}
		updateKeyOpts.KeyID = d.Id()

		// Check if the following attributes have changed. If there are changes, get the new attribute values.
		if d.HasChange("name") {
			updateKeyOpts.Name = d.Get("name").(string)
		}

		if d.HasChange("notes") {
			updateKeyOpts.Notes = d.Get("notes").(string)
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

func getAPIAccessKeyType(d *schema.ResourceData) string {
	keyType := ""
	if v, ok := d.GetOk("key_type"); ok {
		keyType = v.(string)
	}

	log.Printf("[DEBUG] api access key_type: %s", keyType)

	return keyType
}

func getAPIAccessKeyNotes(d *schema.ResourceData) string {
	notes := ""
	if v, ok := d.GetOk("notes"); ok {
		notes = v.(string)
	}

	log.Printf("[DEBUG] api access key notes: %s", notes)

	return notes
}

func getAPIAccessKeyName(d *schema.ResourceData) string {
	name := ""
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	log.Printf("[DEBUG] api access key name: %s", name)

	return name
}
