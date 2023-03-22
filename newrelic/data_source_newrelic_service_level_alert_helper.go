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
	ALERT_TYPE_FIELD              = "alert_type"
	SLI_GUID_FIELD                = "sli_guid"
	SLO_TARGET_FIELD              = "slo_target"
	SLO_PERIOD_FIELD              = "slo_period"
	CUSTOM_CONSUMPTION_FIELD      = "custom_tolerated_budget_consumption"
	CUSTOM_EVALUTION_PERIOD_FIELD = "custom_evaluation_period"
	CONSUMPTION_FIELD             = "tolerated_budget_consumption"
	EVALUATION_PERIOD_FIELD       = "evaluation_period"
	THRESHOLD_FIELD               = "threshold"
	NRQL_FIELD                    = "nrql"

	CUSTOM                        = "custom"
	FAST_BURN                     = "fast_burn"
	FAST_BURN_PERIOD              = 60
	FAST_BURN_CONSUMPTION         = 2
)

func dataSourceNewRelicServiceLevelAlertHelper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicServiceLevelAlertHelperRead,
		Schema: map[string]*schema.Schema{
			ALERT_TYPE_FIELD: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{CUSTOM, FAST_BURN}, true),
			},
			SLI_GUID_FIELD: {
				Type:     schema.TypeString,
				Required: true,
			},
			SLO_TARGET_FIELD: {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			SLO_PERIOD_FIELD: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 7, 28}),
			},
			CUSTOM_CONSUMPTION_FIELD: {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			CUSTOM_EVALUTION_PERIOD_FIELD: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			CONSUMPTION_FIELD: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			EVALUATION_PERIOD_FIELD: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			THRESHOLD_FIELD: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			NRQL_FIELD: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicServiceLevelAlertHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var sliGuid = d.Get(SLI_GUID_FIELD).(string)
	rnd := strconv.Itoa(rand.Int())
	d.SetId(sliGuid + rnd)

	var sloPeriod = d.Get(SLO_PERIOD_FIELD).(int)
	var sloTarget = d.Get(SLO_TARGET_FIELD).(float64)
	var alertType = d.Get(ALERT_TYPE_FIELD).(string)

	_, tOk := d.GetOk(CUSTOM_CONSUMPTION_FIELD)
	_, eOk := d.GetOk(CUSTOM_EVALUTION_PERIOD_FIELD)

	var err error

	switch alertType {
	case FAST_BURN:
		if tOk || eOk {
			return diag.Errorf("For " + FAST_BURN + " alert type do not fill '" + CUSTOM_EVALUTION_PERIOD_FIELD + "' or '" + CUSTOM_CONSUMPTION_FIELD + "'.")
		}

		threshold := calculateThreshold(sloTarget, FAST_BURN_CONSUMPTION, sloPeriod, FAST_BURN_PERIOD)
		err = d.Set(THRESHOLD_FIELD, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(EVALUATION_PERIOD_FIELD, FAST_BURN_PERIOD)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(CONSUMPTION_FIELD, FAST_BURN_CONSUMPTION)
		if err != nil {
			return diag.FromErr(err)
		}
	case CUSTOM:
		if !tOk || !eOk {
			return diag.Errorf("For " + CUSTOM + " alert type the fields '" + CUSTOM_EVALUTION_PERIOD_FIELD + "' and '" + CUSTOM_CONSUMPTION_FIELD + "' are mandatory.")
		}

		toleratedBudgetConsumption := d.Get(CUSTOM_CONSUMPTION_FIELD).(float64)
		evaluationPeriod := d.Get(CUSTOM_EVALUTION_PERIOD_FIELD).(int)
		var threshold = calculateThreshold(sloTarget, toleratedBudgetConsumption, sloPeriod, evaluationPeriod)

		err = d.Set(THRESHOLD_FIELD, threshold)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(EVALUATION_PERIOD_FIELD, evaluationPeriod)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set(CONSUMPTION_FIELD, toleratedBudgetConsumption)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	err = d.Set(NRQL_FIELD, "FROM Metric SELECT 100 - clamp_max(sum(newrelic.sli.good) / sum(newrelic.sli.valid) * 100, 100) as 'SLO compliance'  WHERE sli.guid = '"+sliGuid+"'")
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func calculateThreshold(sloTarget float64, toleratedBudgetConsumption float64, sloPeriod int, evaluationPeriod int) float64 {
	return (100.0 - sloTarget) * ((toleratedBudgetConsumption / 100 * float64(sloPeriod) * 24) / (float64(evaluationPeriod) / 60.0))
}
