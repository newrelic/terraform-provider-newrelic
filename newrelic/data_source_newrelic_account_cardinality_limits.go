package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicAccountCardinalityLimits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicAccountCardinalityLimitsRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID to query limits for. Defaults to the account ID set in the provider configuration.",
			},
			"limits": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cardinality limits for the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique name of the limit.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The current limit value.",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unit for the limit value (e.g. COUNT).",
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The category of the limit (e.g. INGEST).",
						},
					},
				},
			},
		},
	}
}

func dataSourceNewRelicAccountCardinalityLimitsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	limits, err := client.DataManagement.GetLimitsWithContext(ctx, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", accountID))

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if limits != nil {
		flatLimits := make([]map[string]interface{}, 0, len(*limits))
		for _, l := range *limits {
			flatLimits = append(flatLimits, map[string]interface{}{
				"name":     l.Name,
				"value":    l.Value,
				"unit":     string(l.Unit),
				"category": string(l.Category),
			})
		}
		if err := d.Set("limits", flatLimits); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
