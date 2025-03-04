package newrelic

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/v2/pkg/agentapplications"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandApmConfigValues(d *schema.ResourceData) *agentapplications.AgentApplicationSettingsApmConfigInput {

	apmConfig := agentapplications.AgentApplicationSettingsApmConfigInput{}
	if v, ok := d.GetOk("app_apdex_threshold"); ok {
		apmConfig.ApdexTarget = v.(float64)
	}
	apmConfig.UseServerSideConfig = getBoolPointer(d.Get("use_server_side_config").(bool))
	if apmConfig == (agentapplications.AgentApplicationSettingsApmConfigInput{}) {
		return nil
	}
	return &apmConfig
}

func expandTransactionTracerValues(d *schema.ResourceData) *agentapplications.AgentApplicationSettingsTransactionTracerInput {

	tracer := agentapplications.AgentApplicationSettingsTransactionTracerInput{}

	if _, ok := d.GetOk("transaction_tracer"); ok {
		flag := true
		tracer.Enabled = &flag
	} else {
		flag := false
		tracer.Enabled = &flag
	}
	if v, ok := d.GetOk("transaction_tracer.0.transaction_threshold_type"); ok {
		tracer.TransactionThresholdType = agentapplications.AgentApplicationSettingsThresholdTypeEnum(v.(string))
	}
	if v, ok := d.GetOk("transaction_tracer.0.transaction_threshold_value"); ok {
		tracer.TransactionThresholdValue = getFloatPointer(v.(float64))
	}
	if _, ok := d.GetOk("transaction_tracer.0.explain_query_plans"); ok {
		flag := true
		tracer.ExplainEnabled = &flag
	} else {
		flag := false
		tracer.ExplainEnabled = &flag
	}
	if v, ok := d.GetOk("transaction_tracer.0.explain_query_plans.0.query_plan_threshold_value"); ok {
		tracer.ExplainThresholdValue = getFloatPointer(v.(float64))
	}
	if v, ok := d.GetOk("transaction_tracer.0.explain_query_plans.0.query_plan_threshold_type"); ok {
		tracer.ExplainThresholdType = agentapplications.AgentApplicationSettingsThresholdTypeEnum(v.(string))
	}
	if _, ok := d.GetOk("transaction_tracer.0.sql"); ok {
		flag := true
		tracer.LogSql = &flag
	} else {
		flag := false
		tracer.LogSql = &flag
	}
	if v, ok := d.GetOk("transaction_tracer.0.sql.0.record_sql"); ok {
		tracer.RecordSql = agentapplications.AgentApplicationSettingsRecordSqlEnum(v.(string))
	}
	if v, ok := d.GetOk("transaction_tracer.0.stack_trace_threshold_value"); ok {
		tracer.StackTraceThreshold = getFloatPointer(v.(float64))
	}

	if tracer == (agentapplications.AgentApplicationSettingsTransactionTracerInput{}) {
		return nil
	}

	return &tracer
}

func expandErrorCollectorValues(d *schema.ResourceData) *agentapplications.AgentApplicationSettingsErrorCollectorInput {

	collector := agentapplications.AgentApplicationSettingsErrorCollectorInput{}

	if _, ok := d.GetOk("error_collector"); ok {
		flag := true
		collector.Enabled = &flag
	} else {
		flag := false
		collector.Enabled = &flag
	}
	if v, ok := d.GetOk("error_collector.0.expected_error_classes"); ok {
		collector.ExpectedErrorClasses = convertInterfaceToStringSlice(v)
	}
	if v, ok := d.GetOk("error_collector.0.expected_error_codes"); ok {
		strSlice := convertInterfaceToStringSlice(v)
		var httpStatusSlice []agentapplications.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, agentapplications.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.ExpectedErrorCodes = httpStatusSlice
	}
	if v, ok := d.GetOk("error_collector.0.ignored_error_classes"); ok {
		collector.IgnoredErrorClasses = convertInterfaceToStringSlice(v)
	}
	if v, ok := d.GetOk("error_collector.0.ignored_error_codes"); ok {
		strSlice := convertInterfaceToStringSlice(v)
		var httpStatusSlice []agentapplications.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, agentapplications.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.IgnoredErrorCodes = httpStatusSlice
	}

	if collector.Enabled == nil && len(collector.ExpectedErrorClasses) == 0 && len(collector.ExpectedErrorCodes) == 0 && len(collector.IgnoredErrorClasses) == 0 && len(collector.IgnoredErrorCodes) == 0 {
		return nil
	}

	return &collector
}

func expandApplication(d *schema.ResourceData) *agentapplications.AgentApplicationSettingsUpdateInput {

	a := agentapplications.AgentApplicationSettingsUpdateInput{}

	if v, ok := d.GetOk("name"); ok { // alias
		a.Alias = getStringPointer(v.(string))
	}

	slowSQL := &agentapplications.AgentApplicationSettingsSlowSqlInput{}
	slowSQL.Enabled = getBoolPointer(d.Get("enable_slow_sql").(bool))
	a.SlowSql = slowSQL

	TracerType := &agentapplications.AgentApplicationSettingsTracerTypeInput{}
	TracerType.Value = "NONE"
	a.TracerType = TracerType
	if v, ok := d.GetOk("tracer_type"); ok { // tracer type
		TracerType.Value = agentapplications.AgentApplicationSettingsTracer(v.(string))
		a.TracerType = TracerType
	}

	ThreadProfiler := &agentapplications.AgentApplicationSettingsThreadProfilerInput{}
	ThreadProfiler.Enabled = getBoolPointer(d.Get("enable_thread_profiler").(bool))
	a.ThreadProfiler = ThreadProfiler

	// apm settings
	a.ApmConfig = expandApmConfigValues(d) // APM config

	a.TransactionTracer = expandTransactionTracerValues(d) // transaction_tracing

	a.ErrorCollector = expandErrorCollectorValues(d) // error_collector

	return &a
}

// setting APM values
func setAPMApplicationValues(d *schema.ResourceData, ApmSettings entities.AgentApplicationSettingsApmBase) error {
	var err error
	isImported := d.Get("is_imported").(bool)

	if err = setBasicAPMValues(d, ApmSettings, isImported); err != nil {
		return err
	}
	if err = setTransactionTracerValues(d, ApmSettings, isImported); err != nil {
		return err
	}
	if err = setErrorCollectorValues(d, ApmSettings, isImported); err != nil {
		return err
	}

	return nil
}

func setBasicAPMValues(d *schema.ResourceData, ApmSettings entities.AgentApplicationSettingsApmBase, isImported bool) error {
	var err error
	if _, ok := d.GetOk("name"); ok || !isImported {
		if err = d.Set("name", ApmSettings.Alias); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("app_apdex_threshold"); ok || !isImported {
		if err = d.Set("app_apdex_threshold", ApmSettings.ApmConfig.ApdexTarget); err != nil {
			return fmt.Errorf("[DEBUG] Error setting app_apdex_threshold value : %#v", err)
		}
	}
	if _, ok := d.GetOk("use_server_side_config"); ok || !isImported {
		if err = d.Set("use_server_side_config", ApmSettings.ApmConfig.UseServerSideConfig); err != nil {
			return fmt.Errorf("[DEBUG] Error setting use_server_side_config value : %#v", err)
		}
	}
	if _, ok := d.GetOk("enable_slow_sql"); ok || !isImported {
		if err = d.Set("enable_slow_sql", ApmSettings.SlowSql.Enabled); err != nil {
			return fmt.Errorf("[DEBUG] Error setting enable_slow_sql value : %#v", err)
		}
	}
	if _, ok := d.GetOk("enable_thread_profiler"); ok || !isImported {
		if err = d.Set("enable_thread_profiler", ApmSettings.ThreadProfiler.Enabled); err != nil {
			return fmt.Errorf("[DEBUG] Error setting thread profiler value : %#v", err)
		}
	}
	if _, ok := d.GetOk("tracer_type"); ok || !isImported {
		if err = d.Set("tracer_type", ApmSettings.TracerType); err != nil {
			return fmt.Errorf("[DEBUG] Error setting tracer type value : %#v", err)
		}
	}
	return nil
}

func setTransactionTracerValues(d *schema.ResourceData, ApmSettings entities.AgentApplicationSettingsApmBase, isImported bool) error {
	if _, ok := d.GetOk("transaction_tracer"); ok || !isImported {
		if err := flattenTransactionTracingValues(d, ApmSettings.TransactionTracer); err != nil {
			return fmt.Errorf("[DEBUG] Error setting transaction tracer values : %#v", err)
		}
	}
	return nil
}

func setErrorCollectorValues(d *schema.ResourceData, ApmSettings entities.AgentApplicationSettingsApmBase, isImported bool) error {
	if _, ok := d.GetOk("error_collector"); ok || !isImported {
		if err := flattenErrorCollectorValues(d, ApmSettings.ErrorCollector); err != nil {
			return fmt.Errorf("[DEBUG] Error setting error collector values : %#v", err)
		}
	}
	return nil
}

func flattenTransactionTracingValues(d *schema.ResourceData, in entities.AgentApplicationSettingsTransactionTracer) error {
	isImported := d.Get("is_imported").(bool)
	var err error
	transactionTracingValues := map[string]interface{}{
		"stack_trace_threshold_value": 0,
		"transaction_threshold_value": 0,
		"explain_query_plans":         make([]interface{}, 0),
		"sql":                         make([]interface{}, 0),
	}

	_, ok := d.GetOk("transaction_tracer")
	if (in.Enabled && ok) || !isImported {
		transactionTracingValues["stack_trace_threshold_value"] = in.StackTraceThreshold
		transactionTracingValues["transaction_threshold_type"] = in.TransactionThresholdType
		transactionTracingValues["transaction_threshold_value"] = in.TransactionThresholdValue

		if in.ExplainEnabled || !isImported {
			explainQueryPlans := map[string]interface{}{}
			explainQueryPlans["query_plan_threshold_value"] = in.ExplainThresholdValue
			explainQueryPlans["query_plan_threshold_type"] = in.ExplainThresholdType
			transactionTracingValues["explain_query_plans"] = []interface{}{explainQueryPlans}
		}

		if in.LogSql || !isImported {
			logSQLValues := map[string]interface{}{}
			logSQLValues["record_sql"] = in.RecordSql
			transactionTracingValues["sql"] = []interface{}{logSQLValues}
		}
		err = d.Set("transaction_tracer", []interface{}{transactionTracingValues})
	} else {
		err = d.Set("transaction_tracer", nil)
	}

	return err
}

func flattenErrorCollectorValues(d *schema.ResourceData, in entities.AgentApplicationSettingsErrorCollector) error {
	isImported := d.Get("is_imported").(bool)
	var err error
	errorCollectorValues := map[string]interface{}{
		"expected_error_classes": in.ExpectedErrorClasses,
		"expected_error_codes":   in.ExpectedErrorCodes,
		"ignored_error_classes":  in.IgnoredErrorClasses,
		"ignored_error_codes":    in.IgnoredErrorCodes,
	}

	_, ok := d.GetOk("error_collector")
	if (in.Enabled && ok) || !isImported {
		var expectedErrorClasses []string
		var expectedErrorCodes []string
		var ignoredErrorClasses []string
		var ignoredErrorCodes []string
		if _, ok := d.GetOk("error_collector.0.expected_error_classes"); ok || !isImported {
			expectedErrorClasses = append(expectedErrorClasses, in.ExpectedErrorClasses...)
			errorCollectorValues["expected_error_classes"] = expectedErrorClasses
		} else {
			errorCollectorValues["expected_error_classes"] = []string{}
		}
		if _, ok := d.GetOk("error_collector.0.expected_error_codes"); ok || !isImported {
			for _, code := range in.ExpectedErrorCodes {
				expectedErrorCodes = append(expectedErrorCodes, string(code))
			}
			errorCollectorValues["expected_error_codes"] = expectedErrorCodes
		} else {
			errorCollectorValues["expected_error_codes"] = []string{}
		}
		if _, ok := d.GetOk("error_collector.0.ignored_error_classes"); ok || !isImported {
			ignoredErrorClasses = append(ignoredErrorClasses, in.IgnoredErrorClasses...)
			errorCollectorValues["ignored_error_classes"] = ignoredErrorClasses
		} else {
			errorCollectorValues["ignored_error_classes"] = []string{}
		}
		if _, ok := d.GetOk("error_collector.0.ignored_error_codes"); ok || !isImported {
			for _, code := range in.IgnoredErrorCodes {
				ignoredErrorCodes = append(ignoredErrorCodes, string(code))
			}
			errorCollectorValues["ignored_error_codes"] = ignoredErrorCodes
		} else {
			errorCollectorValues["ignored_error_classes"] = []string{}
		}

		err = d.Set("error_collector", []interface{}{errorCollectorValues})
	} else {
		err = d.Set("error_collector", nil)
	}
	return err
}
