package newrelic

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	alertTypeField              = "alert_type"
	sliGUIDField                = "sli_guid"
	sloTargetField              = "slo_target"
	sloPeriodField              = "slo_period"
	customConsumptionField      = "custom_tolerated_budget_consumption"
	customEvaluationPeriodField = "custom_evaluation_period"
	consumptionField            = "tolerated_budget_consumption"
	evaluationPeriodField       = "evaluation_period"
	thresholdField              = "threshold"
	nrqlField                   = "nrql"

	custom              = "custom"
	fastBurn            = "fast_burn"
	fastBurnPeriod      = 60
	fastBurnConsumption = 2
)

func dataSourceNewRelicServiceLevelAlertHelper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertHelperRead,
		Schema: map[string]*schema.Schema{
			alertTypeField: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{custom, fastBurn}, true),
			},
			sliGUIDField: {
				Type:     schema.TypeString,
				Required: true,
			},
			sloTargetField: {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			sloPeriodField: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 7, 28}),
			},
			customConsumptionField: {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			customEvaluationPeriodField: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			consumptionField: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			evaluationPeriodField: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			thresholdField: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			nrqlField: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicServiceLevelAlertHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var sliGUID = d.Get(sliGUIDField).(string)
	rnd := strconv.Itoa(rand.Int())
	d.SetId(sliGUID + rnd)

	var sloPeriod = d.Get(sloPeriodField).(int)
	var sloTarget = d.Get(sloTargetField).(float64)
	var alertType = d.Get(alertTypeField).(string)

	_, tOk := d.GetOk(customConsumptionField)
	_, eOk := d.GetOk(customEvaluationPeriodField)

	switch alertType {
	case fastBurn:
		if tOk || eOk {
			return diag.Errorf("For %s alert type do not fill '%s' or '%s', we use 60 minutes and 2 %.", fastBurn, customEvaluationPeriodField, customConsumptionField)
		}

		threshold := calculateThreshold(sloTarget, fastBurnConsumption, sloPeriod, fastBurnPeriod)
		err := d.Set(thresholdField, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(evaluationPeriodField, fastBurnPeriod)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(consumptionField, fastBurnConsumption)
		if err != nil {
			return diag.FromErr(err)
		}
	case custom:
		if !tOk || !eOk {
			return diag.Errorf("For %s alert type the fields '%s' and '%s' are mandatory.", custom, customEvaluationPeriodField, customConsumptionField)
		}

		toleratedBudgetConsumption := d.Get(customConsumptionField).(float64)
		evaluationPeriod := d.Get(customEvaluationPeriodField).(int)
		var threshold = calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

		err := d.Set(thresholdField, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(evaluationPeriodField, evaluationPeriod)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(consumptionField, toleratedBudgetConsumption)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	err := d.Set(nrqlField, "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = '"+sliGUID+"'")
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func calculateThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / (float64(evaluationPeriod) / 60.0))
}
