package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const lensCreateCatalogMutation = `
mutation($catalog: LensCatalogInput!, $scope: LensScopeInput!) {
  lensCreateCatalog(catalog: $catalog, scope: $scope) {
    connector
    name
    scope {
      id
      type
    }
  }
}`

const lensDeleteCatalogMutation = `
mutation($deleteCatalog: LensDeleteCatalogInput!) {
  lensDeleteCatalog(deleteCatalog: $deleteCatalog)
}`

type lensCreateCatalogResponse struct {
	LensCreateCatalog struct {
		Connector string    `json:"connector"`
		Name      string    `json:"name"`
		Scope     lensScope `json:"scope"`
	} `json:"lensCreateCatalog"`
}

func resourceNewRelicLensConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicLensConnectorCreate,
		ReadContext:   resourceNewRelicLensConnectorRead,
		DeleteContext: resourceNewRelicLensConnectorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Lens connector catalog.",
			},
			"connector": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The connector type enum (e.g. AWSGLUE).",
			},
			"scope": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The scope of the connector.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of the scope (e.g. organization UUID).",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of the scope (e.g. ORGANIZATION).",
						},
					},
				},
			},
			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Key/value properties for the connector.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicLensConnectorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	name := d.Get("name").(string)
	scopeMap := d.Get("scope").([]interface{})[0].(map[string]interface{})

	variables := map[string]interface{}{
		"catalog": map[string]interface{}{
			"connector":  d.Get("connector").(string),
			"name":       name,
			"properties": expandLensProperties(d.Get("properties").([]interface{})),
		},
		"scope": map[string]interface{}{
			"id":   scopeMap["id"].(string),
			"type": scopeMap["type"].(string),
		},
	}

	log.Printf("[INFO] Creating New Relic Lens Connector: %s", name)

	var resp lensCreateCatalogResponse
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, lensCreateCatalogMutation, variables, &resp); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Lens Connector: %w", err))
	}

	created := resp.LensCreateCatalog
	d.SetId(lensConnectorID(created.Name, created.Scope))

	return resourceNewRelicLensConnectorRead(ctx, d, meta)
}

func resourceNewRelicLensConnectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Lens Connector: %s", d.Id())

	name, scopeType, scopeID := parseLensConnectorID(d.Id())

	var resp lensQueryResponse
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, lensConnectorsNerdGraphQuery, map[string]interface{}{}, &resp); err != nil {
		return diag.FromErr(fmt.Errorf("error reading Lens Connectors: %w", err))
	}

	for _, item := range resp.Actor.Organization.Lens.Catalogs.Items {
		if item.Name == name && item.Scope.Type == scopeType && item.Scope.ID == scopeID {
			_ = d.Set("name", item.Name)
			_ = d.Set("connector", item.Connector)
			_ = d.Set("properties", flattenLensProperties(item.Properties))
			_ = d.Set("scope", flattenLensScope(item.Scope))
			return nil
		}
	}

	// Not found — remove from state
	d.SetId("")
	return nil
}

func resourceNewRelicLensConnectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	name, scopeType, scopeID := parseLensConnectorID(d.Id())

	log.Printf("[INFO] Deleting New Relic Lens Connector: %s", name)

	variables := map[string]interface{}{
		"deleteCatalog": map[string]interface{}{
			"name": name,
			"scope": map[string]interface{}{
				"id":   scopeID,
				"type": scopeType,
			},
		},
	}

	var resp interface{}
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, lensDeleteCatalogMutation, variables, &resp); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Lens Connector %q: %w", name, err))
	}

	return nil
}

func expandLensProperties(raw []interface{}) []map[string]interface{} {
	props := make([]map[string]interface{}, len(raw))
	for i, p := range raw {
		m := p.(map[string]interface{})
		props[i] = map[string]interface{}{
			"key":   m["key"].(string),
			"value": m["value"].(string),
		}
	}
	return props
}

// lensConnectorID encodes name:scopeType:scopeID into a single Terraform resource ID.
// SplitN(3) in the parser ensures scopeID may itself contain colons.
func lensConnectorID(name string, scope lensScope) string {
	return strings.Join([]string{name, scope.Type, scope.ID}, ":")
}

func parseLensConnectorID(id string) (name, scopeType, scopeID string) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	}
	return id, "", ""
}
