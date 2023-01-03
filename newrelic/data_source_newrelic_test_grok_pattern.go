package newrelic

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicTestGrokPattern() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicTestGrokRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The account id associated with the test grok.",
			},
			"grok": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Grok pattern to test.",
			},
			"log_lines": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The log lines to test the Grok pattern against.",
			},
			"test_grok": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Test a Grok pattern against a list of log lines.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matched": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the Grok pattern matched.",
						},
						"log_line": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The log line that was tested against.",
						},
						"attributes": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Any attributes that were extracted.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute name.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A string representation of the extracted value (which might not be a String).",
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

func dataSourceNewRelicTestGrokRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Test a Grok pattern against a list of log lines.")
	lines := d.Get("log_lines").(*schema.Set).List()

	if len(lines) == 0 {
		return diag.FromErr(errors.New("`log_lines` is required"))
	}
	perms := make([]string, len(lines))

	for i, line := range lines {
		perms[i] = line.(string)
	}

	res, err := client.Logconfigurations.GetTestGrokWithContext(ctx, accountID, d.Get("grok").(string),
		perms)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	_ = d.Set("test_grok", flattenTestGrokResponse(res))

	return nil

}

func flattenTestGrokResponse(res *[]logconfigurations.LogConfigurationsGrokTestResult) interface{} {
	out := make([]interface{}, len(*res))
	for i, e := range *res {
		m := make(map[string]interface{})
		m["matched"] = e.Matched
		m["log_line"] = e.LogLine
		m["attributes"] = flattenTestGrokAttributes(e.Attributes)
		out[i] = m
	}
	return out
}

func flattenTestGrokAttributes(in []logconfigurations.LogConfigurationsGrokTestExtractedAttribute) interface{} {
	out := make([]interface{}, len(in))

	for i, e := range in {
		m := make(map[string]interface{})
		m["name"] = e.Name
		m["value"] = e.Value

		out[i] = m
	}
	return out
}
