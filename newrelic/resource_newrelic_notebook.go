package newrelic

import (
	"context"
	"fmt"
	"log"
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
				Required:         true,
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				Description:      "The notebook body as a JSON string. Accepts a raw JSON string, a file() reference, or a jsonencode({...}) expression.",
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
	rawContent := d.Get("content").(string)

	normalized, body, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("content: %s", err)
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
	_ = d.Set("content", normalized)

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

	if err := flattenNotebookContent(rawContent, d, "content"); err != nil {
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
	contentChanged := d.HasChange("content")

	if !titleChanged && !contentChanged {
		return nil
	}

	log.Printf("[INFO] Updating New Relic notebook %s (title=%v, content=%v)", guid, titleChanged, contentChanged)

	title := d.Get("title").(string)
	normalized, body, err := normalizeNotebookContent(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("content: %s", err)
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

	_ = d.Set("content", normalized)
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

// resourceNewRelicNotebookImportState handles terraform import.
// Accepts the notebook entity GUID as the import ID.
func resourceNewRelicNotebookImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	guid := d.Id()
	if guid == "" {
		return nil, fmt.Errorf("import ID must be the notebook entity GUID")
	}

	d.SetId(guid)
	// Seed content with a valid placeholder so Read populates it with the
	// real content from the API immediately after this function returns.
	if err := d.Set("content", `{"type":"declarative","version":1,"content":[]}`); err != nil {
		return nil, fmt.Errorf("failed to initialise import state: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
