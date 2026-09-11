package newrelic

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicNotebook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicNotebookRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique entity identifier (GUID) of the notebook.",
			},
			"fetch_content": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "When true, the full notebook body is fetched from the Blob Storage API " +
					"and stored in the `content` attribute. When false (default), only NerdGraph " +
					"metadata (title, organization_id, blob_id) is retrieved, which is faster " +
					"and avoids an extra API call.",
			},

			// Computed attributes populated on every read.
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The title of the notebook.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The New Relic organization ID the notebook belongs to.",
			},
			"blob_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The blob identifier of the current notebook content.",
			},

			// Populated only when fetch_content = true.
			"content": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The notebook body as a normalized JSON string. Only populated when " +
					"`fetch_content` is `true`; empty string otherwise.",
			},
		},
	}
}

func dataSourceNewRelicNotebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Get("guid").(string)
	fetchContent := d.Get("fetch_content").(bool)

	log.Printf("[INFO] Reading New Relic notebook data source: guid=%s fetch_content=%v", guid, fetchContent)

	nb, err := client.Notebooks.GetNotebookWithContext(ctx, guid)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return diag.Errorf("notebook %s not found", guid)
		}
		return diag.FromErr(err)
	}

	d.SetId(guid)
	_ = d.Set("title", nb.Name)
	_ = d.Set("blob_id", nb.Content.ID)

	orgID := nb.Scope.ID
	_ = d.Set("organization_id", orgID)

	if fetchContent {
		if orgID == "" {
			// Fall back to provider-level org resolution when scope is absent.
			var resolveErr error
			orgID, resolveErr = getOrganizationID(ctx, providerConfig, "")
			if resolveErr != nil {
				return diag.FromErr(resolveErr)
			}
		}
		rawContent, contentErr := client.Notebooks.GetNotebookContentWithContext(ctx, orgID, guid)
		if contentErr != nil {
			return diag.FromErr(contentErr)
		}
		normalized, normErr := normalizeNotebookContent(string(rawContent))
		if normErr != nil {
			return diag.FromErr(normErr)
		}
		_ = d.Set("content", normalized)
	} else {
		_ = d.Set("content", "")
	}

	return nil
}
