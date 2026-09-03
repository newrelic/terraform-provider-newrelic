package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the notebook.",
			},
			// content is the full notebook body, sent verbatim to the Blob
			// Storage API. No Go struct expansion is performed; the provider
			// normalises (canonical JSON: sorted keys, consistent indentation)
			// on every write so Terraform's line-level diff shows only the
			// lines that actually changed.
			//
			// Tip: use jsonencode({...}) in your config to write the content
			// as HCL object syntax. Terraform evaluates jsonencode at plan
			// time, which yields the most granular diff output.
			"content": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The notebook content as a JSON string. Use jsonencode({...}) for inline HCL-object syntax. The provider stores a normalised (sorted-key, indented) copy so that terraform plan shows only semantically meaningful changes.",
				DiffSuppressFunc: suppressEquivalentNotebookContent,
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic organization ID the notebook belongs to. Defaults to the organization of the authenticated account when omitted.",
			},
			// account_id reflects the numeric account scope reported by NerdGraph.
			// It is informational; the Blob Storage API uses organization_id.
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID. Defaults to the account configured in the provider.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the notebook in New Relic.",
			},
		},
	}
}

// resourceNewRelicNotebookCustomizeDiff runs at terraform plan time. It
// validates that content is well-formed JSON before the user ever reaches
// apply, mirroring the pattern used by other resources in this provider.
func resourceNewRelicNotebookCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	raw, ok := d.GetOk("content")
	if !ok {
		return nil
	}
	rawStr := raw.(string)
	if rawStr == "" {
		return nil
	}

	normalized, err := normalizeNotebookContent(rawStr)
	if err != nil {
		return fmt.Errorf("content: %w", err)
	}

	// Suppress cosmetic-only diffs: if the normalised form of the new value
	// matches the normalised value already in state, clear the planned change.
	if !d.HasChange("content") {
		return nil
	}
	oldRaw, _ := d.GetChange("content")
	normOld, _ := normalizeNotebookContent(oldRaw.(string))
	if string(normalized) == string(normOld) {
		_ = d.Clear("content")
	}
	return nil
}

func resourceNewRelicNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get("title").(string)
	rawContent := d.Get("content").(string)

	normalized, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("content: %s", err)
	}

	var contentBody interface{}
	if err := json.Unmarshal([]byte(normalized), &contentBody); err != nil {
		return diag.Errorf("content re-parse: %s", err)
	}

	log.Printf("[INFO] Creating New Relic notebook: %s", title)

	resp, err := client.Notebooks.CreateNotebookWithContext(ctx, orgID, title, contentBody)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.EntityGUID == "" {
		return diag.Errorf("notebook create: no entityGuid returned")
	}

	log.Printf("[INFO] New Relic notebook created, GUID: %s", resp.EntityGUID)
	d.SetId(resp.EntityGUID)

	_ = d.Set("guid", resp.EntityGUID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("content", string(normalized))

	return resourceNewRelicNotebookRead(ctx, d, meta)
}

func resourceNewRelicNotebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	log.Printf("[INFO] Reading New Relic notebook %s", guid)

	// Fetch NerdGraph metadata (title, scope/org, version).
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

	// Fetch content blob from Blob Storage.
	rawContent, err := client.Notebooks.GetNotebookContentWithContext(ctx, orgID, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("title", nb.Name)
	_ = d.Set("guid", nb.ID)
	_ = d.Set("organization_id", orgID)

	if err := flattenNotebookContent(rawContent, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get("title").(string)
	rawContent := d.Get("content").(string)

	normalized, err := normalizeNotebookContent(rawContent)
	if err != nil {
		return diag.Errorf("content: %s", err)
	}

	var contentBody interface{}
	if err := json.Unmarshal([]byte(normalized), &contentBody); err != nil {
		return diag.Errorf("content re-parse: %s", err)
	}

	titleChanged := d.HasChange("title")
	contentChanged := d.HasChange("content")

	log.Printf("[INFO] Updating New Relic notebook %s (title_changed=%v, content_changed=%v)", guid, titleChanged, contentChanged)

	if titleChanged {
		// RenameNotebook sends both the new title and the content in one
		// atomic call. The Blob Storage API has no rename-only path.
		_, err = client.Notebooks.RenameNotebookWithContext(ctx, orgID, guid, title, contentBody)
	} else if contentChanged {
		_, err = client.Notebooks.UpdateNotebookContentWithContext(ctx, orgID, guid, contentBody)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("content", string(normalized))

	return resourceNewRelicNotebookRead(ctx, d, meta)
}

func resourceNewRelicNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()
	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
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
