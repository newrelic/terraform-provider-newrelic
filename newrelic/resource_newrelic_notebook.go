package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notebooks"
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
			StateContext: resourceNewRelicNotebookImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		CustomizeDiff: customizeNotebookDiff,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the notebook.",
			},

			// Exactly one of content or content_json must be set. Choose content
			// when authoring in HCL (field-level plan diffs). Choose content_json
			// when working from a UI export or file() (line-level JSON diffs).
			// Both fields store the content in normalized form (alphabetically
			// sorted keys, 2-space indent) so cosmetic formatting never diffs.

			"content": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ValidateFunc:     validateNotebookContent,
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
				ValidateFunc:     validateNotebookContent,
				ExactlyOneOf:     []string{"content", "content_json"},
				Description: "The notebook body as a raw JSON string. Intended for " +
					"notebooks exported from the New Relic UI (for example, via the " +
					"Copy JSON option) or loaded from a file using file(). " +
					"Plan output shows a line-level diff of the normalized JSON, " +
					"making it straightforward to identify what changed. " +
					"Mutually exclusive with content.",
			},

			// organization_id is resolved automatically and stored for API calls.
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
			// blob_id tracks the immutable content blob written by the last
			// Terraform-managed write. Because each write creates a new blob,
			// blob_id equality is a reliable signal that content is unchanged,
			// allowing Read to skip a Blob Storage GET in the common no-diff case.
			"blob_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The blob identifier of the current notebook content. Updated after every write.",
			},
		},
	}
}

// notebookContentField returns the name of the active content field
// ("content" or "content_json") along with its current value. When neither
// field is set — for example immediately after terraform import before the
// first plan — it defaults to "content_json".
func notebookContentField(d *schema.ResourceData) (field, raw string) {
	if v, ok := d.GetOk("content"); ok && v.(string) != "" {
		return "content", v.(string)
	}
	if v, ok := d.GetOk("content_json"); ok && v.(string) != "" {
		return "content_json", v.(string)
	}
	return "content_json", ""
}

// notebookOrgID returns the organization ID from state if present, or resolves
// it from provider credentials. Centralises the fallback logic used by Update
// and Delete.
func notebookOrgID(ctx context.Context, guid string, d *schema.ResourceData, pc *ProviderConfig) (string, error) {
	if id, _ := d.Get("organization_id").(string); id != "" {
		return id, nil
	}
	log.Printf("[DEBUG] organization_id not in state for notebook %s, resolving from provider credentials", guid)
	return getOrganizationID(ctx, pc, "")
}

func resourceNewRelicNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pc := meta.(*ProviderConfig)
	client := pc.NewClient

	orgID, err := getOrganizationID(ctx, pc, "")
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get("title").(string)
	field, rawContent := notebookContentField(d)

	normalized, body, err := parseAndNormalizeContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	log.Printf("[INFO] Creating New Relic notebook: %s", title)

	resp, err := client.Notebooks.CreateNotebookWithContext(ctx, orgID, title, body)
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
	_ = d.Set("blob_id", resp.BlobID)
	_ = d.Set(field, normalized)

	// Poll NerdGraph until content.id is populated.
	//
	// The Blob Storage API and NerdGraph are asynchronously consistent: the
	// notebook exists in the Blob API immediately but NerdGraph may take a few
	// seconds to index the entity and populate content.id. Waiting here ensures
	// blob_id in state is NerdGraph-confirmed, making the short-circuit in
	// subsequent Reads reliable right away.
	indexingDeadline := time.Now().Add(d.Timeout(schema.TimeoutCreate))
	backoff := time.Second
	for {
		nb, nbErr := client.Notebooks.GetNotebookWithContext(ctx, resp.EntityGUID)
		if nbErr == nil && nb != nil && nb.Content.ID != "" {
			log.Printf("[DEBUG] Notebook %s indexed in NerdGraph (blob_id: %s)", resp.EntityGUID, nb.Content.ID)
			_ = d.Set("blob_id", nb.Content.ID)
			break
		}
		if ctx.Err() != nil || time.Now().After(indexingDeadline) {
			log.Printf("[WARN] Notebook %s created but NerdGraph indexing timed out", resp.EntityGUID)
			return append(resourceNewRelicNotebookRead(ctx, d, meta), diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Notebook created but NerdGraph indexing timed out",
				Detail: "Notebook was created successfully via the Blob Storage API but did not " +
					"appear in NerdGraph within the create timeout. Run terraform plan again " +
					"to confirm the resource is in the expected state.",
			})
		}
		log.Printf("[DEBUG] Waiting for notebook %s to be indexed in NerdGraph (retry in %s)...", resp.EntityGUID, backoff)
		time.Sleep(backoff)
		if backoff < 15*time.Second {
			backoff *= 2
		}
	}

	return resourceNewRelicNotebookRead(ctx, d, meta)
}

func resourceNewRelicNotebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pc := meta.(*ProviderConfig)
	client := pc.NewClient

	guid := d.Id()
	log.Printf("[INFO] Reading New Relic notebook %s", guid)

	nb, err := client.Notebooks.GetNotebookWithContext(ctx, guid)
	if err != nil {
		if isNotebookNotFoundError(err) {
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

	// Blob-ID short-circuit: if NerdGraph's content.id matches what we have in
	// state, the blob has not changed and we can skip the Blob Storage GET.
	// NerdGraph may return content as null due to propagation lag; fall through
	// to the full fetch in that case.
	storedBlobID, _ := d.Get("blob_id").(string)
	currentBlobID := nb.Content.ID
	if storedBlobID != "" && currentBlobID != "" && currentBlobID == storedBlobID {
		log.Printf("[DEBUG] Notebook %s content unchanged (blob_id %s) - skipping Blob GET", guid, storedBlobID)
		_ = d.Set("title", nb.Name)
		_ = d.Set("guid", nb.ID)
		_ = d.Set("organization_id", orgID)
		return nil
	}

	rawContent, err := client.Notebooks.GetNotebookContentWithContext(ctx, orgID, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("title", nb.Name)
	_ = d.Set("guid", nb.ID)
	_ = d.Set("organization_id", orgID)
	if currentBlobID != "" {
		_ = d.Set("blob_id", currentBlobID)
	}

	field, _ := notebookContentField(d)
	if err := flattenNotebookContent(rawContent, d, field); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pc := meta.(*ProviderConfig)
	client := pc.NewClient

	guid := d.Id()
	orgID, err := notebookOrgID(ctx, guid, d, pc)
	if err != nil {
		return diag.FromErr(err)
	}

	titleChanged := d.HasChange("title")
	contentChanged := d.HasChange("content") || d.HasChange("content_json")

	log.Printf("[INFO] Updating New Relic notebook %s (title_changed=%v, content_changed=%v)", guid, titleChanged, contentChanged)

	if !titleChanged && !contentChanged {
		return resourceNewRelicNotebookRead(ctx, d, meta)
	}

	title := d.Get("title").(string)
	field, rawContent := notebookContentField(d)

	normalized, body, err := parseAndNormalizeContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	// When the user switches modes (content ↔ content_json) without changing
	// the actual JSON, the SDK sees both fields as changed but the server blob
	// is identical. Skip the API write if the normalized content matches what is
	// already stored in state to avoid a redundant blob version.
	if contentChanged && !titleChanged {
		_, prevRaw := func() (string, string) {
			if d.HasChange("content") {
				old, _ := d.GetChange("content")
				if s, _ := old.(string); s != "" {
					return "content", s
				}
			}
			if d.HasChange("content_json") {
				old, _ := d.GetChange("content_json")
				if s, _ := old.(string); s != "" {
					return "content_json", s
				}
			}
			return "", ""
		}()
		if prevRaw != "" {
			prevNorm, _, normErr := parseAndNormalizeContent(prevRaw)
			if normErr == nil && prevNorm == normalized {
				log.Printf("[DEBUG] Notebook %s: mode switch with identical content — skipping Blob API write", guid)
				// Update only the state field name; the server blob is unchanged.
				_ = d.Set(field, normalized)
				return resourceNewRelicNotebookRead(ctx, d, meta)
			}
		}
	}

	var mutResp *notebooks.NotebookMutationResponse
	if titleChanged {
		mutResp, err = client.Notebooks.RenameNotebookWithContext(ctx, orgID, guid, title, body)
	} else {
		mutResp, err = client.Notebooks.UpdateNotebookContentWithContext(ctx, orgID, guid, body)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set(field, normalized)

	// Set blob_id from the mutation response before calling Read. Read may see
	// a stale NerdGraph blob_id due to propagation lag and would otherwise
	// overwrite our fresh value. We re-assert it after Read completes.
	freshBlobID := ""
	if mutResp != nil && mutResp.BlobID != "" {
		freshBlobID = mutResp.BlobID
		_ = d.Set("blob_id", freshBlobID)
	}

	diags := resourceNewRelicNotebookRead(ctx, d, meta)

	// Re-assert the mutation's blob_id in case Read overwrote it with a stale
	// NerdGraph value. This prevents an unnecessary Blob GET on the next plan.
	if freshBlobID != "" {
		_ = d.Set("blob_id", freshBlobID)
	}

	return diags
}

func resourceNewRelicNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pc := meta.(*ProviderConfig)
	client := pc.NewClient

	guid := d.Id()
	orgID, err := notebookOrgID(ctx, guid, d, pc)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting New Relic notebook %s", guid)

	if err := client.Notebooks.DeleteNotebookWithContext(ctx, orgID, guid); err != nil {
		if isNotebookNotFoundError(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

// resourceNewRelicNotebookImportState handles terraform import for notebooks.
// Accepts either a bare GUID or a composite "GUID:mode" ID where mode is
// "content" or "content_json" (default). The mode signals which field Read
// should populate in state to match the user's configuration.
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

	// Set a minimal valid placeholder so notebookContentField detects the desired
	// mode when Read runs immediately after this function. Read overwrites it
	// with the actual normalized content from the Blob API.
	placeholder := `{"type":"declarative","version":1,"content":[]}`
	if err := d.Set(mode, placeholder); err != nil {
		return nil, fmt.Errorf("failed to signal import mode %q: %w", mode, err)
	}

	return []*schema.ResourceData{d}, nil
}
