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
			// Import accepts either a bare GUID or a composite "GUID:mode" ID.
			// mode must be "content" or "content_json" (default: content_json).
			//
			//   terraform import newrelic_notebook.example NjQy...
			//   terraform import newrelic_notebook.example NjQy...:content_json
			//   terraform import newrelic_notebook.example NjQy...:content
			//
			// Specifying :content means the imported state will have the content
			// field populated (matching a config that uses content = jsonencode({...})).
			// Specifying :content_json (or omitting the mode) populates content_json,
			// matching a config that uses content_json = file("...") or a raw JSON string.
			StateContext: resourceNewRelicNotebookImportState,
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

// resourceNewRelicNotebookImportState handles terraform import for notebooks.
//
// It accepts either a bare GUID or a composite "GUID:mode" ID where mode is
// "content" or "content_json". The mode tells the resource which field to
// populate when Read fetches the notebook body from the Blob Storage API:
//
//   - content_json (default) - populates the content_json field, matching a
//     config that uses content_json = file("...") or an inline JSON string.
//   - content - populates the content field, matching a config that uses
//     content = jsonencode({...}).
//
// After importing, run terraform plan to confirm the imported state matches
// your configuration. If the modes differ (e.g., you import without a mode
// but your config uses content), the next plan will surface the difference so
// you can align your config accordingly.
func resourceNewRelicNotebookImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	guid := parts[0]
	if guid == "" {
		return nil, fmt.Errorf("import ID must be a notebook GUID, optionally followed by :content or :content_json")
	}

	mode := "content_json"
	if len(parts) == 2 {
		mode = parts[1]
		if mode != "content" && mode != "content_json" {
			return nil, fmt.Errorf(
				"invalid import mode %q: use %q for HCL jsonencode authoring or %q for raw JSON / file() authoring "+
					"(e.g. terraform import newrelic_notebook.example %s:%s)",
				mode, "content", "content_json", guid, "content_json",
			)
		}
	}

	d.SetId(guid)

	// Set a minimal valid placeholder in the target field so notebookContentField
	// detects the desired mode when Read runs immediately after this function.
	// Read overwrites it with the actual normalized content from the Blob API.
	placeholder := `{"version":"1","blocks":[]}`
	if err := d.Set(mode, placeholder); err != nil {
		return nil, fmt.Errorf("failed to signal import mode %q: %w", mode, err)
	}

	return []*schema.ResourceData{d}, nil
}
