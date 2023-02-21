package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNewRelicServiceLevelAlertHelper() *schema.Resource { 
    return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertHelperRead,
		Schema: map[string]*schema.Schema{
			"sli_guid": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"slo_target": {
				Type:        schema.TypeFloat,
				Required:    true,
                ValidateFunc: validation.FloatBetween(0, 100),
			},
			"slo_period": {
				Type:        schema.TypeInt,
				Required:    true,
                ValidateFunc: validation.IntInSlice([]int{1,7,28}),
			},
			"custom_tolerated_budget_consumption": {
				Type:        schema.TypeFloat,
                Optional:    true,
                ValidateFunc: validation.FloatBetween(0, 100),
			},
            "custom_evaluation_period": {
                Type: schema.TypeInt,
                Optional:    true,
                ValidateFunc: validation.IntAtLeast(1),
            },
            "custom_threshold": {
                Type: schema.TypeFloat,
                Computed: true,
            },
            "fast_burn_threshold": {
                Type: schema.TypeFloat,
                Computed: true,
            },
            "fast_burn_evaluation_period": {
                Type: schema.TypeInt,
                Computed: true,
            },
            "slow_burn_threshold": {
                Type: schema.TypeFloat,
                Computed: true,
            },
            "slow_burn_evaluation_period": {
                Type: schema.TypeInt,
                Computed: true,
            },
		},
	}
}

func dataSourceNewRelicServiceLevelAlertHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    d.SetId(d.Get("sli_guid").(string))

	var sloPeriod = d.Get("slo_period").(int)
	var sloTarget = d.Get("slo_target").(float64)
    _, tOk := d.GetOk("custom_tolerated_budget_consumption")
    _, eOk := d.GetOk("custom_evaluation_period")

    var err error
    if tOk && eOk {
        toleratedBudgetConsumption := d.Get("custom_tolerated_budget_consumption").(float64)
        evaluationPeriod := d.Get("custom_evaluation_period").(int)
        var customAlertThreshold = calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

        err = d.Set("custom_threshold", customAlertThreshold)

        if err != nil {
            return diag.FromErr(err)
        }
    }

    fastBurnThreshold := calculateThreshold(sloTarget, 2, sloPeriod, 60)
    err = d.Set("fast_burn_threshold", fastBurnThreshold)
    if err != nil {
        return diag.FromErr(err)
    }
    err = d.Set("fast_burn_evaluation_period", 60)
    if err != nil {
        return diag.FromErr(err)
    }

    slowBurnAlertThreshold := calculateThreshold(sloTarget, 2, sloPeriod, 600)
    err = d.Set("slow_burn_threshold", slowBurnAlertThreshold)
    if err != nil {
        return diag.FromErr(err)
    }
    err = d.Set("slow_burn_evaluation_period", 600)
    if err != nil {
        return diag.FromErr(err)
    }
	return nil
}

func calculateThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / (float64(evaluationPeriod) / 60.0))
}

