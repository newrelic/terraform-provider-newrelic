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
	AlertTypeField              = "alert_type"
	SliGUIDField                = "sli_guid"
	SloTargetField              = "slo_target"
	SloPeriodField              = "slo_period"
	CustomConsumptionField      = "custom_tolerated_budget_consumption"
	CustomEvaluationPeriodField = "custom_evaluation_period"
	ConsumptionField            = "tolerated_budget_consumption"
	EvaluationPeriodField       = "evaluation_period"
	ThresholdField              = "threshold"
	NRQLField                   = "nrql"

	Custom              = "custom"
	FastBurn            = "fast_burn"
	FastBurnPeriod      = 60
	FastBurnConsumption = 2
)

func dataSourceNewRelicServiceLevelAlertHelper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertHelperRead,
		Schema: map[string]*schema.Schema{
			AlertTypeField: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{Custom, FastBurn}, true),
			},
			SliGUIDField: {
				Type:     schema.TypeString,
				Required: true,
			},
			SloTargetField: {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			SloPeriodField: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 7, 28}),
			},
			CustomConsumptionField: {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			CustomEvaluationPeriodField: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			ConsumptionField: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			EvaluationPeriodField: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			ThresholdField: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			NRQLField: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicServiceLevelAlertHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var sliGUID = d.Get(SliGUIDField).(string)
	rnd := strconv.Itoa(rand.Int())
	d.SetId(sliGUID + rnd)

	var sloPeriod = d.Get(SloPeriodField).(int)
	var sloTarget = d.Get(SloTargetField).(float64)
	var alertType = d.Get(AlertTypeField).(string)

	_, tOk := d.GetOk(CustomConsumptionField)
	_, eOk := d.GetOk(CustomEvaluationPeriodField)

	var err error

	switch alertType {
	case FastBurn:
		if tOk || eOk {
			return diag.Errorf("For " + FastBurn + " alert type do not fill '" + CustomEvaluationPeriodField + "' or '" + CustomConsumptionField + "'.")
		}

		threshold := calculateThreshold(sloTarget, FastBurnConsumption, sloPeriod, FastBurnPeriod)
		err = d.Set(ThresholdField, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(EvaluationPeriodField, FastBurnPeriod)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(ConsumptionField, FastBurnConsumption)
		if err != nil {
			return diag.FromErr(err)
		}
	case Custom:
		if !tOk || !eOk {
			return diag.Errorf("For " + Custom + " alert type the fields '" + CustomEvaluationPeriodField + "' and '" + CustomConsumptionField + "' are mandatory.")
		}

		toleratedBudgetConsumption := d.Get(CustomConsumptionField).(float64)
		evaluationPeriod := d.Get(CustomEvaluationPeriodField).(int)
		var threshold = calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

		err = d.Set(ThresholdField, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(EvaluationPeriodField, evaluationPeriod)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(ConsumptionField, toleratedBudgetConsumption)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	err = d.Set(NRQLField, "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = '"+sliGUID+"'")
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func calculateThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / (float64(evaluationPeriod) / 60.0))
}
