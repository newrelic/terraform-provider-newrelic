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
	"github.com/newrelic/newrelic-client-go/v2/pkg/notebooks"
)

func resourceNewRelicNotebook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNotebookCreate,
		ReadContext:   resourceNewRelicNotebookRead,
		UpdateContext: resourceNewRelicNotebookUpdate,
		DeleteContext: resourceNewRelicNotebookDelete,
		Importer: &schema.ResourceImporter{
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
			"content": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ValidateFunc:     validateNotebookContent,
				ExactlyOneOf:     []string{"content", "content_json"},
				Description:      "The notebook body as a jsonencode({...}) expression. Mutually exclusive with content_json.",
			},
			"content_json": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ValidateFunc:     validateNotebookContent,
				ExactlyOneOf:     []string{"content", "content_json"},
				Description:      "The notebook body as a raw JSON string. Mutually exclusive with content.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The New Relic organization ID the notebook belongs to.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity GUID of the notebook.",
			},
			"blob_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The blob identifier of the current notebook content version.",
			},
		},
	}
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

	normalized, body, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	log.Printf("[INFO] Creating New Relic notebook: %s", title)

	resp, err := client.Notebooks.CreateNotebookWithContext(ctx, orgID, title, body)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.EntityGUID == "" {
		return diag.Errorf("the New Relic API did not return an entity GUID for the new notebook")
	}

	log.Printf("[INFO] New Relic notebook created, GUID: %s", resp.EntityGUID)
	d.SetId(resp.EntityGUID)
	_ = d.Set("guid", resp.EntityGUID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("blob_id", resp.BlobID)
	_ = d.Set(field, normalized)

	// NerdGraph indexes the notebook entity asynchronously after the Blob Storage
	// write. RetryContext polls until content.id is populated, ensuring blob_id
	// in state is NerdGraph-confirmed before the first plan runs.
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		nb, nbErr := client.Notebooks.GetNotebookWithContext(ctx, resp.EntityGUID)
		if nbErr != nil {
			return resource.NonRetryableError(fmt.Errorf("error reading notebook after create: %w", nbErr))
		}
		if nb == nil || nb.Content.ID == "" {
			return resource.RetryableError(fmt.Errorf("notebook %s not yet indexed in NerdGraph", resp.EntityGUID))
		}
		log.Printf("[DEBUG] Notebook %s indexed in NerdGraph (blob_id: %s)", resp.EntityGUID, nb.Content.ID)
		_ = d.Set("blob_id", nb.Content.ID)
		return nil
	})
	if retryErr != nil {
		return append(diag.FromErr(retryErr), diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Notebook created but NerdGraph indexing timed out",
			Detail:   "The notebook was written to Blob Storage but did not appear in NerdGraph within the create timeout. Run terraform plan to confirm state.",
		})
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

	// Blob-ID short-circuit: skip the Blob Storage GET when NerdGraph confirms
	// the content is unchanged. Falls through when NerdGraph lags and returns "".
	storedBlobID, _ := d.Get("blob_id").(string)
	currentBlobID := nb.Content.ID
	if storedBlobID != "" && currentBlobID != "" && currentBlobID == storedBlobID {
		log.Printf("[DEBUG] Notebook %s content unchanged (blob_id %s) — skipping Blob GET", guid, storedBlobID)
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

	if !titleChanged && !contentChanged {
		return nil
	}

	log.Printf("[INFO] Updating New Relic notebook %s (title=%v, content=%v)", guid, titleChanged, contentChanged)

	title := d.Get("title").(string)
	field, rawContent := notebookContentField(d)

	normalized, body, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("%s: %s", field, err)
	}

	// Skip the Blob API write when the user switches between content and
	// content_json modes without changing the actual JSON value.
	if contentChanged && !titleChanged {
		prevRaw := previousContentRaw(d)
		if prevRaw != "" {
			prevNorm, _, normErr := normalizeNotebookContent(prevRaw)
			if normErr == nil && prevNorm == normalized {
				log.Printf("[DEBUG] Notebook %s: mode switch with identical content — skipping Blob API write", guid)
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
	if mutResp != nil && mutResp.BlobID != "" {
		_ = d.Set("blob_id", mutResp.BlobID)
	}

	return resourceNewRelicNotebookRead(ctx, d, meta)
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

// resourceNewRelicNotebookImportState handles terraform import. Accepts either
// a bare GUID or "GUID:mode" where mode is "content" or "content_json" (default).
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
	// Seed the mode field with a valid placeholder so notebookContentField picks
	// up the intended mode when Read runs immediately after this function.
	placeholder := `{"type":"declarative","version":1,"content":[]}`
	if err := d.Set(mode, placeholder); err != nil {
		return nil, fmt.Errorf("failed to signal import mode %q: %w", mode, err)
	}

	return []*schema.ResourceData{d}, nil
}

// previousContentRaw returns the old raw value from the content field that changed,
// used by Update to detect a mode-switch with semantically identical JSON.
func previousContentRaw(d *schema.ResourceData) string {
	if d.HasChange("content") {
		if old, _ := d.GetChange("content"); old != nil {
			if s, _ := old.(string); s != "" {
				return s
			}
		}
	}
	if d.HasChange("content_json") {
		if old, _ := d.GetChange("content_json"); old != nil {
			if s, _ := old.(string); s != "" {
				return s
			}
		}
	}
	return ""
}
