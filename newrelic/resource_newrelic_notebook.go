package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicNotebook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotebookCreate,
		ReadContext:   resourceNewRelicNotebookRead,
		UpdateContext: resourceNewRelicNotebookUpdate,
		DeleteContext: resourceNewRelicNotebookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceNewRelicNotebookCustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the notebook.",
			},

			// Exactly one of content or content_json must be set; they cannot
			// be used together. Choose content when you want to author the
			// notebook directly in HCL and get field-level plan diffs. Choose
			// content_json when you are working from JSON exported out of the
			// New Relic UI or stored in a file. Both fields store the content
			// in a normalized form (alphabetically sorted keys, consistent
			// indentation), so purely cosmetic formatting changes never show
			// up as a planned update.

			"content": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ExactlyOneOf:     []string{"content", "content_json"},
				Description: "The notebook body, expressed as an HCL object using " +
					"jsonencode({...}). Terraform evaluates the expression at plan " +
					"time, so terraform plan shows changes at the individual field " +
					"level rather than as an opaque JSON diff. Recommended when the " +
					"notebook is authored and maintained entirely in Terraform. " +
					"Mutually exclusive with content_json.",
			},
			"content_json": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ValidateFunc:     validation.StringIsJSON,
				ExactlyOneOf:     []string{"content", "content_json"},
				Description: "The notebook body as a raw JSON string. Intended for " +
					"notebooks exported from the New Relic UI (for example, via the " +
					"Copy JSON option) or loaded from a file using file(). " +
					"Plan output shows a line-level diff of the normalized JSON, " +
					"making it straightforward to identify what changed. " +
					"Mutually exclusive with content.",
			},

			// organization_id is resolved automatically from the authenticated
			// account and stored in state for subsequent API calls. It is not
			// a user-facing argument; notebooks are organization-scoped and the
			// organization is derived from the provider credentials.
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The New Relic organization ID the notebook belongs to. Resolved automatically from the provider credentials.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the notebook in New Relic.",
			},
		},
	}
}

// resourceNewRelicNotebookCustomizeDiff runs at plan time to catch malformed
// JSON early, so users see a clear error before an apply is attempted.
func resourceNewRelicNotebookCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	for _, key := range []string{"content", "content_json"} {
		raw, ok := d.GetOk(key)
		if !ok || raw.(string) == "" {
			continue
		}
		if _, err := normalizeNotebookContent(raw.(string)); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}
	return nil
}

// notebookContentField returns the name of the active content field
// ("content" or "content_json") along with its current value. When neither
// field is set - for example, immediately after terraform import before the
// first plan - it defaults to "content_json" so the imported state is
// immediately usable.
func notebookContentField(d *schema.ResourceData) (field, raw string) {
	if v, ok := d.GetOk("content"); ok && v.(string) != "" {
		return "content", v.(string)
	}
	if v, ok := d.GetOk("content_json"); ok && v.(string) != "" {
		return "content_json", v.(string)
	}
	return "content_json", ""
}

func resourceNewRelicNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	orgID, err := getOrganizationID(ctx, providerConfig, "")
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get("title").(string)
	field, rawContent := notebookContentField(d)

	normalized, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	var contentBody interface{}
	if unmarshalErr := json.Unmarshal([]byte(normalized), &contentBody); unmarshalErr != nil {
		return diag.Errorf("%s re-parse: %s", field, unmarshalErr)
	}

	log.Printf("[INFO] Creating New Relic notebook: %s", title)

	resp, err := client.Notebooks.CreateNotebookWithContext(ctx, orgID, title, contentBody)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.EntityGUID == "" {
		return diag.Errorf("the New Relic API did not return an entity GUID for the new notebook; this is unexpected - please try again or contact support if the issue persists")
	}

	log.Printf("[INFO] New Relic notebook created, GUID: %s", resp.EntityGUID)
	d.SetId(resp.EntityGUID)
	_ = d.Set("guid", resp.EntityGUID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set(field, normalized)

	return resourceNewRelicNotebookRead(ctx, d, meta)
}

func resourceNewRelicNotebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	log.Printf("[INFO] Reading New Relic notebook %s", guid)

	nb, err := client.Notebooks.GetNotebookWithContext(ctx, guid)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Notebook %s not found, removing from state", guid)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	orgID := nb.Scope.ID
	if orgID == "" {
		orgID, _ = d.Get("organization_id").(string)
	}

	rawContent, err := client.Notebooks.GetNotebookContentWithContext(ctx, orgID, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("title", nb.Name)
	_ = d.Set("guid", nb.ID)
	_ = d.Set("organization_id", orgID)

	// Write the fetched content back into whichever field the user declared.
	// On a fresh import, neither field is set yet; default to content_json so
	// the imported state is immediately usable without a plan change.
	field, _ := notebookContentField(d)
	if err := flattenNotebookContent(rawContent, d, field); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	orgID, _ := d.Get("organization_id").(string)
	if orgID == "" {
		var err error
		orgID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	title := d.Get("title").(string)
	field, rawContent := notebookContentField(d)

	normalized, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	var contentBody interface{}
	if unmarshalErr := json.Unmarshal([]byte(normalized), &contentBody); unmarshalErr != nil {
		return diag.Errorf("%s re-parse: %s", field, unmarshalErr)
	}

	titleChanged := d.HasChange("title")
	contentChanged := d.HasChange("content") || d.HasChange("content_json")

	log.Printf("[INFO] Updating New Relic notebook %s (title_changed=%v, content_changed=%v)", guid, titleChanged, contentChanged)

	if titleChanged {
		// Rename is atomic with a content POST - the Blob API has no rename-only path.
		_, err = client.Notebooks.RenameNotebookWithContext(ctx, orgID, guid, title, contentBody)
	} else if contentChanged {
		_, err = client.Notebooks.UpdateNotebookContentWithContext(ctx, orgID, guid, contentBody)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set(field, normalized)

	return resourceNewRelicNotebookRead(ctx, d, meta)
}

func resourceNewRelicNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	orgID, _ := d.Get("organization_id").(string)
	if orgID == "" {
		var err error
		orgID, err = getOrganizationID(ctx, providerConfig, "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting New Relic notebook %s", guid)

	if err := client.Notebooks.DeleteNotebookWithContext(ctx, orgID, guid); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
