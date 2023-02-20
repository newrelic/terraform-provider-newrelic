package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNewRelicServiceLevelAlertThreshold() *schema.Resource { return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertThresholdRead,
		Schema: map[string]*schema.Schema{
			"slo_target": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "",
                ValidateFunc: validation.FloatBetween(0, 100),
			},
			"slo_period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "",
                ValidateFunc: validation.IntInSlice([]int{1,7,28}),
			},
			"tolerated_budget_consumption": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "",
                ValidateFunc: validation.FloatBetween(0, 100),
			},
			"evaluation_period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "",
                ValidateFunc: validation.IntAtLeast(1),
			},
			"alert_threshold": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func dataSourceNewRelicServiceLevelAlertThresholdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    d.SetId("serviceLevelAlertThreshold")

	var sloPeriod = d.Get("slo_period").(int)
	var sloTarget = d.Get("slo_target").(float64)
	var toleratedBudgetConsumption = d.Get("tolerated_budget_consumption").(float64)
	var evaluationPeriod = d.Get("evaluation_period").(int)

	var alertThreshold = calculateAlertThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

    err := d.Set("alert_threshold", alertThreshold)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func calculateAlertThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / float64(evaluationPeriod))
}

