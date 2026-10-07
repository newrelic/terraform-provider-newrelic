package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const lensConnectorsNerdGraphQuery = `{
  actor {
    organization {
      lens {
        catalogs {
          items {
            name
            properties {
              key
              value
            }
            connector
            scope {
              id
              type
            }
            type
          }
        }
      }
    }
  }
}`

type lensConnectorCatalogItem struct {
	Name       string         `json:"name"`
	Properties []lensProperty `json:"properties"`
	Connector  string         `json:"connector"`
	Scope      lensScope      `json:"scope"`
	Type       string         `json:"type"`
}

type lensProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type lensScope struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type lensQueryResponse struct {
	Actor struct {
		Organization struct {
			Lens struct {
				Catalogs struct {
					Items []lensConnectorCatalogItem `json:"items"`
				} `json:"catalogs"`
			} `json:"lens"`
		} `json:"organization"`
	} `json:"actor"`
}

func dataSourceNewRelicLensConnectors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicLensConnectorsRead,
		Schema: map[string]*schema.Schema{
			"connectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Lens connector catalog items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the connector catalog item.",
						},
						"connector": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector identifier.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the connector.",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Key/value properties associated with the connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scope": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The scope of the connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the scope.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the scope.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNewRelicLensConnectorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Lens Connectors")

	var resp lensQueryResponse
	if err := client.NerdGraph.QueryWithResponseAndContext(ctx, lensConnectorsNerdGraphQuery, map[string]interface{}{}, &resp); err != nil {
		return diag.FromErr(fmt.Errorf("error fetching Lens Connectors: %w", err))
	}

	items := resp.Actor.Organization.Lens.Catalogs.Items

	if err := d.Set("connectors", flattenLensConnectors(items)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("lens_connectors")

	return nil
}

func flattenLensConnectors(items []lensConnectorCatalogItem) []interface{} {
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"name":       item.Name,
			"connector":  item.Connector,
			"type":       item.Type,
			"properties": flattenLensProperties(item.Properties),
			"scope":      flattenLensScope(item.Scope),
		}
	}
	return result
}

func flattenLensProperties(props []lensProperty) []interface{} {
	result := make([]interface{}, len(props))
	for i, p := range props {
		result[i] = map[string]interface{}{
			"key":   p.Key,
			"value": p.Value,
		}
	}
	return result
}

func flattenLensScope(s lensScope) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":   s.ID,
			"type": s.Type,
		},
	}
}
