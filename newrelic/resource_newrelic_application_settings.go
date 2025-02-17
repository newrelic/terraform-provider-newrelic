package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func resourceNewRelicApplicationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicApplicationSettingsCreate,
		ReadContext:   resourceNewRelicApplicationSettingsRead,
		UpdateContext: resourceNewRelicApplicationSettingsUpdate,
		DeleteContext: resourceNewRelicApplicationSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			applicationSettingCommonSchema(),
			apmApplicationSettingsSchema(),
		),
		CustomizeDiff: validateApplicationSettingsInput,
	}
}

func validateErrorCollector(d *schema.ResourceDiff, errorsList *[]string) {
	if excodes, ok := d.GetOk("error_collector.0.expected_error_codes"); ok {
		for _, excode := range excodes.([]interface{}) {
			excodeStr, ok := excode.(string)
			if !ok {
				*errorsList = append(*errorsList, fmt.Sprintf("expected_error_codes '%v' is not a string", excode))
			} else if !regexp.MustCompile(`^[1-9][0-9]{2}$`).MatchString(excodeStr) { // Validate the expected_error_codes against the regex pattern
				*errorsList = append(*errorsList, fmt.Sprintf("expected_error_codes '%s' must be a valid status code between 100 and 900", excodeStr))
			}
		}
	}
	if excodes, ok := d.GetOk("error_collector.0.ignored_error_codes"); ok {
		for _, excode := range excodes.([]interface{}) {
			excodeStr, ok := excode.(string)
			if !ok {
				*errorsList = append(*errorsList, fmt.Sprintf("ignored_error_codes '%v' is not a string", excode))
			} else if !regexp.MustCompile(`^[1-9][0-9]{2}$`).MatchString(excodeStr) { // Validate the expected_error_codes against the regex pattern
				*errorsList = append(*errorsList, fmt.Sprintf("ignored_error_codes '%s' must be a valid status code between 100 and 900", excodeStr))
			}
		}
	}
}

func validateThresholds(d *schema.ResourceDiff, errorsList *[]string) {
	thresholdType, typeExists := d.GetOk("transaction_tracer.0.explain_query_plans.0.query_plan_threshold_type")
	_, valueExists := d.GetOk("transaction_tracer.0.explain_query_plans.0.query_plan_threshold_value")

	if typeExists && thresholdType == "VALUE" && !valueExists {
		*errorsList = append(*errorsList, "`query_plan_threshold_value` must be set when `query_plan_threshold_type` is 'VALUE'")
	}

	if valueExists && (!typeExists || thresholdType != "VALUE") {
		*errorsList = append(*errorsList, "`query_plan_threshold_type` must be set to 'VALUE' when `query_plan_threshold_value` is provided")
	}

	if typeExists && thresholdType == "APDEX_F" && valueExists {
		*errorsList = append(*errorsList, "`query_plan_threshold_value` should not be set when `query_plan_threshold_type` is 'APDEX_F'")
	}

	transactionThresholdType, transactionTypeExists := d.GetOk("transaction_tracer.0.transaction_threshold_type")
	_, transactionValueExists := d.GetOk("transaction_tracer.0.transaction_threshold_value")

	if transactionTypeExists && transactionThresholdType == "VALUE" && !transactionValueExists {
		*errorsList = append(*errorsList, "`transaction_threshold_value` must be set when `transaction_threshold_type` is 'VALUE'")
	}

	if transactionValueExists && (!transactionTypeExists || transactionThresholdType != "VALUE") {
		*errorsList = append(*errorsList, "`transaction_threshold_type` must be set to 'VALUE' when `transaction_threshold_value` is provided")
	}

	if transactionTypeExists && transactionThresholdType == "APDEX_F" && transactionValueExists {
		*errorsList = append(*errorsList, "`transaction_threshold_value` should not be set when `transaction_threshold_type` is 'APDEX_F'")
	}
}

func validateRealUserMonitoring(d *schema.ResourceDiff) error {
	realUserMonitoring := d.Get("use_server_side_config").(bool)

	attr := map[string]interface{}{
		"transaction_tracer":     d.Get("transaction_tracer"),
		"error_collector":        d.Get("error_collector"),
		"tracer_type":            d.Get("tracer_type"),
		"enable_thread_profiler": d.Get("enable_thread_profiler"),
	}

	for atr := range attr {
		_, blockExists := d.GetOk(atr)
		if blockExists && !realUserMonitoring {
			return fmt.Errorf("use_server_side_config must be set to true when %s is configured", atr)
		}
	}
	return nil
}

// Custom validator that ensures fields are validated correctly
func validateApplicationSettingsInput(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
	var errorsList []string

	if err := validateRealUserMonitoring(d); err != nil {
		return err
	}

	validateThresholds(d, &errorsList)
	validateErrorCollector(d, &errorsList)

	if len(errorsList) == 0 {
		return nil
	}
	errorsString := "the following validation errors have been identified: \n"
	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}
	return errors.New(errorsString)
}

func applicationSettingCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The GUID of the application in New Relic.",
		},
		"is_imported": {
			Type:     schema.TypeBool,
			Computed: true,
			Default:  nil,
		},
	}
}

func apmApplicationSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the application in New Relic.",
		},
		"app_apdex_threshold": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "The response time threshold value for Apdex score calculation.",
		},
		"use_server_side_config": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable or disable server side monitoring.",
		},
		"transaction_tracer": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Configuration for transaction tracing.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"transaction_threshold_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"APDEX_F", "VALUE"}, false),
						Description:  "The type of threshold for transaction tracing, either 'APDEX_F' or 'VALUE'.",
					},
					"transaction_threshold_value": {
						Type:        schema.TypeFloat,
						Optional:    true,
						Description: "The threshold value for transaction tracing when 'transaction_threshold_type' is 'VALUE'.",
					},
					"sql": {
						Type:        schema.TypeList,
						Optional:    true,
						MinItems:    1,
						MaxItems:    1,
						Description: "Configuration for SQL tracing.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"record_sql": {
									Type:         schema.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice([]string{"OBFUSCATED", "OFF", "RAW"}, true),
									Description:  "The level of SQL recording, either 'OBFUSCATED', 'OFF', or 'RAW'.",
								},
							},
						},
					},
					"stack_trace_threshold_value": {
						Type:        schema.TypeFloat,
						Optional:    true,
						Description: "The response time threshold value for capturing stack traces of SQL queries.",
					},
					"explain_query_plans": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Configuration for explain plans of slow SQL queries.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"query_plan_threshold_type": {
									Type:         schema.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice([]string{"APDEX_F", "VALUE"}, false),
									Description:  "The type of threshold for explain plans, either 'APDEX_F' or 'VALUE'.",
								},
								"query_plan_threshold_value": {
									Type:         schema.TypeFloat,
									Optional:     true,
									ValidateFunc: validation.FloatAtLeast(0),
									Description:  "The threshold value for explain plans when 'query_plan_threshold_type' is 'VALUE'.",
								},
							},
						},
					},
				},
			},
		},
		"error_collector": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Configuration for error collection.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expected_error_classes": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "A list of error classes that are expected and should not trigger alerts.",
					},
					"expected_error_codes": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "A list of error codes that are expected and should not trigger alerts.",
					},
					"ignored_error_classes": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "A list of error classes that should be ignored.",
					},
					"ignored_error_codes": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Optional:    true,
						Description: "A list of error codes that should be ignored.",
					},
				},
			},
		},
		"tracer_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"CROSS_APPLICATION_TRACER", "DISTRIBUTED_TRACING", "NONE", "OPT_OUT"}, true),
			Description:  "The type of tracer to use, either 'CROSS_APPLICATION_TRACER', 'DISTRIBUTED_TRACING', 'NONE', or 'OPT_OUT'.",
		},
		"enable_thread_profiler": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable or disable the thread profiler.",
		},
	}
}

func resourceNewRelicApplicationSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return resourceNewRelicApplicationSettingsUpdate(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic application %+v", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no New Relic application found with given guid %s", d.Id()))
	}

	var dig diag.Diagnostics
	switch (*resp).(type) {
	case *entities.ApmApplicationEntity:
		entity := (*resp).(*entities.ApmApplicationEntity)
		d.SetId(string(entity.GUID))
		_ = d.Set("guid", string(entity.GUID))
		dig = diag.FromErr(setAPMApplicationValues(d, entity.ApmSettings))
	default:
		dig = diag.FromErr(fmt.Errorf("problem in retrieving application with GUID %s", d.Id()))
	}
	return dig
}

func resourceNewRelicApplicationSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	updateApplicationParams := expandApplication(d)

	guid := d.Get("guid").(string)
	log.Printf("[INFO] Updating New Relic application %+v with params: %+v", guid, updateApplicationParams)

	agentApplicationSettingResult, err := client.AgentApplications.AgentApplicationSettingsUpdate(common.EntityGUID(guid), *updateApplicationParams)

	if err != nil {
		return diag.FromErr(err)
	}
	if agentApplicationSettingResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while Updating New Relic application"))
	}

	time.Sleep(2 * time.Second)

	d.SetId(string(agentApplicationSettingResult.GUID))
	err = d.Set("is_imported", true)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicApplicationSettingsRead(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// You can only delete applications which have not reported in the past 12 hours
	client := meta.(*ProviderConfig).NewClient
	applicationGUID := common.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Application %s", d.Id())
	result, err := client.AgentApplications.AgentApplicationDelete(applicationGUID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Delete New Relic Application %s response: %+v", d.Id(), result.Success)

	return nil
}
