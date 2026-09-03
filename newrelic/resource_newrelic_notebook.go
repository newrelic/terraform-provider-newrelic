package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the notebook.",
			},

			// content and content_json are mutually exclusive.
			//
			// content       — write the notebook body as an HCL object literal
			//                 using jsonencode({...}). Terraform evaluates
			//                 jsonencode at plan time so the plan output shows
			//                 object-level field diffs (e.g. "query changed from
			//                 X to Y") rather than whole-blob text diffs.
			//                 Recommended for notebooks authored and maintained
			//                 entirely in Terraform.
			//
			// content_json  — supply raw JSON directly (e.g. file("export.json")
			//                 or a heredoc). Intended for notebooks exported from
			//                 the New Relic UI via the "Copy JSON" feature and
			//                 pasted into Terraform with minimal editing. Diffs
			//                 are shown at the JSON line level, which is still
			//                 precise enough to identify the changed attribute.
			//
			// Both fields accept any valid notebook JSON; both are normalised
			// (sorted keys, consistent indentation) on every write and read so
			// cosmetic reformatting never surfaces as a plan change.

			"content": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The notebook content as an HCL object. Write using " +
					"jsonencode({...}) for structured authoring with granular plan " +
					"diffs. Mutually exclusive with content_json.",
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ExactlyOneOf:     []string{"content", "content_json"},
			},
			"content_json": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The notebook content as a raw JSON string. Use when " +
					"pasting JSON exported from the New Relic UI (e.g. via the " +
					"\"Copy JSON\" feature) or loading from a file with " +
					"file(\"notebook.json\"). Mutually exclusive with content.",
				DiffSuppressFunc: suppressEquivalentNotebookContent,
				ValidateFunc:     validation.StringIsJSON,
				ExactlyOneOf:     []string{"content", "content_json"},
			},

			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic organization ID the notebook belongs to. Defaults to the organization of the authenticated account when omitted.",
			},
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

// resourceNewRelicNotebookCustomizeDiff validates content at plan time so
// malformed JSON is rejected before an apply ever starts.
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

// notebookContentField returns the name of whichever content field is populated
// ("content" or "content_json") and the raw value. Falls back to "content_json"
// when neither is set (e.g. immediately after an import before a plan).
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

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
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
	if err := json.Unmarshal([]byte(normalized), &contentBody); err != nil {
		return diag.Errorf("%s re-parse: %s", field, err)
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
	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
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
	if err := json.Unmarshal([]byte(normalized), &contentBody); err != nil {
		return diag.Errorf("%s re-parse: %s", field, err)
	}

	titleChanged := d.HasChange("title")
	contentChanged := d.HasChange("content") || d.HasChange("content_json")

	log.Printf("[INFO] Updating New Relic notebook %s (title_changed=%v, content_changed=%v)", guid, titleChanged, contentChanged)

	if titleChanged {
		// Rename is atomic with a content POST — the Blob API has no rename-only path.
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
