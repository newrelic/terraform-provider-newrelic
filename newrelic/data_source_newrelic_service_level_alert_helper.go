package newrelic

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type serviceLevelAlertType string

var serviceLevelAlertTypes = struct {
	custom   serviceLevelAlertType
	fastBurn serviceLevelAlertType
	slowBurn serviceLevelAlertType
}{
	custom:   "custom",
	fastBurn: "fast_burn",
	slowBurn: "slow_burn",
}

func dataSourceNewRelicServiceLevelAlertHelper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertHelperRead,
		Schema: map[string]*schema.Schema{
			"alert_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"custom", "fast_burn", "slow_burn"}, true),
			},
			"sli_guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slo_target": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			"slo_period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 7, 28}),
			},
			"custom_tolerated_budget_consumption": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			"custom_evaluation_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"tolerated_budget_consumption": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"evaluation_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"nrql": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_bad_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceNewRelicServiceLevelAlertHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sliGUID := d.Get("sli_guid")
	rnd := rand.Int()
	d.SetId(fmt.Sprintf("%v%v", sliGUID, rnd))

	_, tOk := d.GetOk("custom_tolerated_budget_consumption")
	_, eOk := d.GetOk("custom_evaluation_period")

	alertType := d.Get("alert_type").(string)
	switch serviceLevelAlertType(alertType) {
	case serviceLevelAlertTypes.fastBurn:
		if tOk || eOk {
			return diag.Errorf("For 'fast_burn' alert type do not fill 'custom_evaluation_period' or 'custom_tolerated_budget_consumption', we use 60 minutes and 2%%.")
		}

		if err := fillData(d, 60, 2); err != nil {
			return diag.FromErr(err)
		}

	case serviceLevelAlertTypes.slowBurn:
		if tOk || eOk {
			return diag.Errorf("For 'slow_burn' alert type do not fill 'custom_evaluation_period' or 'custom_tolerated_budget_consumption', we use 360 minutes and 5%%.")
		}

		if err := fillData(d, 360, 5); err != nil {
			return diag.FromErr(err)
		}

	case serviceLevelAlertTypes.custom:
		if !tOk || !eOk {
			return diag.Errorf("For 'custom' alert type the fields 'custom_evaluation_period' and 'custom_tolerated_budget_consumption' are mandatory.")
		}

		toleratedBudgetConsumption := d.Get("custom_tolerated_budget_consumption").(float64)
		evaluationPeriod := d.Get("custom_evaluation_period").(int)
		if err := fillData(d, evaluationPeriod, toleratedBudgetConsumption); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func fillData(d *schema.ResourceData, evaluationPeriod int, toleratedBudgetConsumption float64) error {
	if err := d.Set("evaluation_period", evaluationPeriod); err != nil {
		return err
	}

	if err := d.Set("tolerated_budget_consumption", toleratedBudgetConsumption); err != nil {
		return err
	}

	var nrql string
	if d.Get("is_bad_events").(bool) {
		nrql = fmt.Sprintf("FROM Metric SELECT 100 - clamp_max((sum(newrelic.sli.valid) - sum(newrelic.sli.bad)) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance' WHERE sli.guid = '%v'", d.Get("sli_guid"))
	} else {
		nrql = fmt.Sprintf("FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance' WHERE sli.guid = '%v'", d.Get("sli_guid"))
	}
	if err := d.Set("nrql", nrql); err != nil {
		return err
	}

	sloPeriod := d.Get("slo_period").(int)
	sloTarget := d.Get("slo_target").(float64)
	threshold := calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)
	if err := d.Set("threshold", threshold); err != nil {
		return err
	}

	return nil
}

func calculateThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / (float64(evaluationPeriod) / 60.0))
}
