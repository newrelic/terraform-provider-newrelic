package newrelic

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
)

func expandApmConfigValues(d *schema.ResourceData) *apm.AgentApplicationSettingsApmConfigInput {

	apmConfig := apm.AgentApplicationSettingsApmConfigInput{}
	if v, ok := d.GetOk("app_apdex_threshold"); ok {
		apmConfig.ApdexTarget = v.(float64)
	}
	apmConfig.UseServerSideConfig = getBoolPointer(d.Get("use_server_side_config").(bool))
	if apmConfig == (apm.AgentApplicationSettingsApmConfigInput{}) {
		return nil
	}
	return &apmConfig
}

func expandTransactionTracerValues(d *schema.ResourceData) *apm.AgentApplicationSettingsTransactionTracerInput {

	tracer := apm.AgentApplicationSettingsTransactionTracerInput{}

	if _, ok := d.GetOk("transaction_tracer"); ok {
		flag := true
		tracer.Enabled = &flag
	} else {
		flag := false
		tracer.Enabled = &flag
	}
	if v, ok := d.GetOk("transaction_tracer.0.transaction_threshold_type"); ok {
		tracer.TransactionThresholdType = apm.AgentApplicationSettingsThresholdTypeEnum(v.(string))
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
		tracer.ExplainThresholdType = apm.AgentApplicationSettingsThresholdTypeEnum(v.(string))
	}
	if _, ok := d.GetOk("transaction_tracer.0.sql"); ok {
		flag := true
		tracer.LogSql = &flag
	} else {
		flag := false
		tracer.LogSql = &flag
	}
	if v, ok := d.GetOk("transaction_tracer.0.sql.0.record_sql"); ok {
		tracer.RecordSql = apm.AgentApplicationSettingsRecordSqlEnum(v.(string))
	}
	if v, ok := d.GetOk("transaction_tracer.0.stack_trace_threshold_value"); ok {
		tracer.StackTraceThreshold = getFloatPointer(v.(float64))
	}

	if tracer == (apm.AgentApplicationSettingsTransactionTracerInput{}) {
		return nil
	}

	return &tracer
}

func expandErrorCollectorValues(d *schema.ResourceData) *apm.AgentApplicationSettingsErrorCollectorInput {

	collector := apm.AgentApplicationSettingsErrorCollectorInput{}

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
		var httpStatusSlice []apm.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, apm.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.ExpectedErrorCodes = httpStatusSlice
	}
	if v, ok := d.GetOk("error_collector.0.ignored_error_classes"); ok {
		collector.IgnoredErrorClasses = convertInterfaceToStringSlice(v)
	}
	if v, ok := d.GetOk("error_collector.0.ignored_error_codes"); ok {
		strSlice := convertInterfaceToStringSlice(v)
		var httpStatusSlice []apm.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, apm.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.IgnoredErrorCodes = httpStatusSlice
	}

	if collector.Enabled == nil && len(collector.ExpectedErrorClasses) == 0 && len(collector.ExpectedErrorCodes) == 0 && len(collector.IgnoredErrorClasses) == 0 && len(collector.IgnoredErrorCodes) == 0 {
		return nil
	}

	return &collector
}

func expandMaskInputOptions(v interface{}) apm.AgentApplicationSettingsMaskInputOptionsInput {
	values := v.(map[string]interface{})

	inputOptions := apm.AgentApplicationSettingsMaskInputOptionsInput{}
	for key, value := range values {
		switch key {
		case "color":
			inputOptions.Color = getBoolPointer(value.(bool))
		case "date":
			inputOptions.Date = getBoolPointer(value.(bool))
		case "datetimeLocal":
			inputOptions.DatetimeLocal = getBoolPointer(value.(bool))
		case "email":
			inputOptions.Email = getBoolPointer(value.(bool))
		case "month":
			inputOptions.Month = getBoolPointer(value.(bool))
		case "number":
			inputOptions.Number = getBoolPointer(value.(bool))
		case "range":
			inputOptions.Range = getBoolPointer(value.(bool))
		case "search":
			inputOptions.Search = getBoolPointer(value.(bool))
		case "select":
			inputOptions.Select = getBoolPointer(value.(bool))
		case "tel":
			inputOptions.Tel = getBoolPointer(value.(bool))
		case "text":
			inputOptions.Text = getBoolPointer(value.(bool))
		case "textArea":
			inputOptions.TextArea = getBoolPointer(value.(bool))
		case "time":
			inputOptions.Time = getBoolPointer(value.(bool))
		case "url":
			inputOptions.URL = getBoolPointer(value.(bool))
		case "week":
			inputOptions.Week = getBoolPointer(value.(bool))
		}
	}
	return inputOptions
}

func expandMobileSettingValues(d *schema.ResourceData) *apm.AgentApplicationSettingsMobileSettingsInput {
	a := apm.AgentApplicationSettingsMobileSettingsInput{}

	if v, ok := d.GetOk("use_crash_reports"); ok {
		a.UseCrashReports = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("enable_application_exit_info"); ok {
		a.ApplicationExitInfo.Enabled = getBoolPointer(v.(bool))
	}
	//if v, ok := d.GetOk("log_reporting.0.enabled"); ok {
	//	a.NetworkSettings. = getBoolPointer(v.(bool))
	//}
	//if v, ok := d.GetOk("log_reporting.0.level"); ok {
	//	a.UseCrashReports = getBoolPointer(v.(bool))
	//}
	//if v, ok := d.GetOk("log_reporting.0.sampling_rate"); ok {
	//	a.UseCrashReports = getBoolPointer(v.(bool))
	//}

	if a == (apm.AgentApplicationSettingsMobileSettingsInput{}) {
		return nil
	}

	return &a
}

func expandBrowserMonitoringValue(v interface{}) *apm.AgentApplicationSettingsBrowserMonitoringInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})
	browserMonitoring := &apm.AgentApplicationSettingsBrowserMonitoringInput{}

	browserMonitoring.DistributedTracing = expandDistributedTracingValues(values["distributed_tracing"])
	browserMonitoring.Ajax = expandAjaxValues(values["ajax"])

	return browserMonitoring
}

func expandDistributedTracingValues(v interface{}) *apm.AgentApplicationSettingsBrowserDistributedTracingInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})
	distributedTracing := apm.AgentApplicationSettingsBrowserDistributedTracingInput{}

	if v, ok := values["enabled"]; ok {
		distributedTracing.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["cors_enabled"]; ok {
		distributedTracing.CorsEnabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["exclude_newrelic_header"]; ok {
		distributedTracing.ExcludeNewrelicHeader = getBoolPointer(v.(bool))
	}
	if v, ok := values["cors_use_newrelic_header"]; ok {
		distributedTracing.CorsUseNewrelicHeader = getBoolPointer(v.(bool))
	}
	if v, ok := values["cors_use_trace_context_headers"]; ok {
		distributedTracing.CorsUseTracecontextHeaders = getBoolPointer(v.(bool))
	}
	if v, ok := values["allowed_origins"]; ok {
		distributedTracing.AllowedOrigins = v.([]string)
	}

	return &distributedTracing
}

func expandAjaxValues(v interface{}) *apm.AgentApplicationSettingsBrowserAjaxInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}
	values := v.([]interface{})[0].(map[string]interface{})
	ajax := &apm.AgentApplicationSettingsBrowserAjaxInput{}

	if v, ok := values["deny_list"]; ok {
		ajax.DenyList = v.([]string)
	}
	return ajax
}

func expandSessionReplayValues(v interface{}) *apm.AgentApplicationSettingsSessionReplayInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})
	sessionReplay := &apm.AgentApplicationSettingsSessionReplayInput{}

	if v, ok := values["auto_start"]; ok {
		sessionReplay.AutoStart = getBoolPointer(v.(bool))
	}
	if v, ok := values["enabled"]; ok {
		sessionReplay.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["error_sampling_rate"]; ok {
		sessionReplay.ErrorSamplingRate = v.(float64)
	}
	if v, ok := values["sampling_rate"]; ok {
		sessionReplay.SamplingRate = v.(float64)
	}
	if _, ok := values["mask_input_options"]; ok {
		sessionReplay.MaskInputOptions = expandMaskInputOptions(values["mask_input_options"])
	}
	if v, ok := values["mask_all_inputs"]; ok {
		sessionReplay.MaskAllInputs = getBoolPointer(v.(bool))
	}
	if v, ok := values["block_selector"]; ok {
		sessionReplay.BlockSelector = getStringPointer(v.(string))
	}
	if v, ok := values["mask_text_selector"]; ok {
		sessionReplay.MaskTextSelector = getStringPointer(v.(string))
	}

	return sessionReplay
}

func expandSessionTraceValues(v interface{}) *apm.AgentApplicationSettingsSessionTraceInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})
	sessionTrace := &apm.AgentApplicationSettingsSessionTraceInput{}

	if v, ok := values["error_sampling_rate"]; ok {
		sessionTrace.ErrorSamplingRate = v.(float64)
	}
	if v, ok := values["mode"]; ok {
		sessionTrace.Mode = apm.AgentApplicationSettingsSessionTraceModeInput(v.(string))
	}
	if v, ok := values["sampling_rate"]; ok {
		sessionTrace.SamplingRate = v.(float64)
	}
	if v, ok := values["enabled"]; ok {
		sessionTrace.Enabled = getBoolPointer(v.(bool))
	}

	return sessionTrace
}

func expandApplication(d *schema.ResourceData) *apm.AgentApplicationSettingsUpdateInput {

	a := apm.AgentApplicationSettingsUpdateInput{}

	if v, ok := d.GetOk("name"); ok { // alias
		a.Alias = getStringPointer(v.(string))
	}

	if v, ok := d.GetOk("tracer_type"); ok { // tracer type
		TracerType := &apm.AgentApplicationSettingsTracerTypeInput{}
		TracerType.Value = apm.AgentApplicationSettingsTracer(v.(string))
		a.TracerType = TracerType
	}

	ThreadProfiler := &apm.AgentApplicationSettingsThreadProfilerInput{}
	ThreadProfiler.Enabled = getBoolPointer(d.Get("enable_thread_profiler").(bool))
	a.ThreadProfiler = ThreadProfiler

	// apm settings
	a.ApmConfig = expandApmConfigValues(d) // APM config

	a.TransactionTracer = expandTransactionTracerValues(d) // transaction_tracing

	a.ErrorCollector = expandErrorCollectorValues(d) // error_collector

	//mobile settings
	a.MobileSettings = expandMobileSettingValues(d)

	// browser settings
	a.BrowserMonitoring = expandBrowserMonitoringValue(d.Get("browser_monitoring")) // browser_monitoring

	a.SessionReplay = expandSessionReplayValues(d.Get("session_replay")) // session replay settings

	a.SessionTrace = expandSessionTraceValues(d.Get("session_trace")) // session trace settings

	if v, ok := d.GetOk("end_user_apdex_threshold"); ok { // browser config apdex target
		a.BrowserConfig.ApdexTarget = v.(float64)
	}

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
	}

	return err
}

func flattenErrorCollectorValues(d *schema.ResourceData, in entities.AgentApplicationSettingsErrorCollector) error {
	isImported := d.Get("is_imported").(bool)
	var expectedErrorCodes []string
	var ignoredErrorCodes []string
	errorCollectorValues := map[string]interface{}{
		"expected_error_classes": in.ExpectedErrorClasses,
		"expected_error_codes":   in.ExpectedErrorCodes,
		"ignored_error_classes":  in.IgnoredErrorClasses,
		"ignored_error_codes":    in.IgnoredErrorCodes,
	}

	if _, ok := d.GetOk("expected_error_codes"); ok || !isImported {
		for _, code := range in.ExpectedErrorCodes {
			expectedErrorCodes = append(expectedErrorCodes, string(code))
		}
		errorCollectorValues["expected_error_codes"] = expectedErrorCodes
	}
	if _, ok := d.GetOk("ignored_error_codes"); ok || !isImported {
		for _, code := range in.IgnoredErrorCodes {
			ignoredErrorCodes = append(ignoredErrorCodes, string(code))
		}
		errorCollectorValues["ignored_error_codes"] = ignoredErrorCodes
	}

	err := d.Set("error_collector", []interface{}{errorCollectorValues})
	return err
}

// setting browser values
func setBrowserApplicationValues(d *schema.ResourceData, BrowserSettings entities.AgentApplicationSettingsBrowserBase) error {

	if err := d.Set("end_user_apdex_threshold", BrowserSettings.BrowserConfig.ApdexTarget); err != nil {
		return fmt.Errorf("[DEBUG] Error setting browser config values: %#v", err)
	}

	if err := d.Set("session_replay", flattenSessionReplayValues(BrowserSettings.SessionReplay)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting session replay values: %#v", err)
	}

	if err := d.Set("session_trace", flattenSessionTracerValues(BrowserSettings.SessionTrace)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting session trace values: %#v", err)
	}

	if err := d.Set("browser_monitoring", flattenBrowserMonitoringValues(BrowserSettings.BrowserMonitoring)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting bowser monitoring values: %#v", err)
	}

	return nil
}

func flattenSessionReplayValues(in entities.AgentApplicationSettingsSessionReplay) []map[string]interface{} {
	sessionReplayValues := make([]map[string]interface{}, 1)
	sessionReplayValues[0] = map[string]interface{}{
		"enabled":             in.Enabled,
		"auto_start":          in.AutoStart,
		"error_sampling_rate": in.ErrorSamplingRate,
		"sampling_rate":       in.SamplingRate,
		"mask_input_options": map[string]interface{}{
			"color":          in.MaskInputOptions.Color,
			"date":           in.MaskInputOptions.Date,
			"datetime_local": in.MaskInputOptions.DatetimeLocal,
			"email":          in.MaskInputOptions.Email,
			"month":          in.MaskInputOptions.Month,
			"number":         in.MaskInputOptions.Number,
			"range":          in.MaskInputOptions.Range,
			"search":         in.MaskInputOptions.Search,
			"select":         in.MaskInputOptions.Select,
			"tel":            in.MaskInputOptions.Tel,
			"text":           in.MaskInputOptions.Text,
			"text_area":      in.MaskInputOptions.TextArea,
			"time":           in.MaskInputOptions.Time,
			"url":            in.MaskInputOptions.URL,
			"week":           in.MaskInputOptions.Week,
		},
		"mask_all_inputs":    in.MaskAllInputs,
		"block_selector":     in.BlockSelector,
		"mask_text_selector": in.MaskTextSelector,
	}
	return sessionReplayValues
}

func flattenSessionTracerValues(in entities.AgentApplicationSettingsSessionTrace) []map[string]interface{} {
	sessionTracerValues := make([]map[string]interface{}, 1)
	sessionTracerValues[0] = map[string]interface{}{
		"enabled":             in.Enabled,
		"sampling_rate":       in.SamplingRate,
		"error_sampling_rate": in.ErrorSamplingRate,
		"mode":                in.Mode,
	}
	return sessionTracerValues
}

func flattenBrowserMonitoringValues(in entities.AgentApplicationSettingsBrowserMonitoring) []map[string]interface{} {
	browserMonitoringValues := make([]map[string]interface{}, 1)
	browserMonitoringValues[0] = map[string]interface{}{
		"distributed_tracing": map[string]interface{}{
			"enabled":                        in.DistributedTracing.Enabled,
			"cors_enabled":                   in.DistributedTracing.CorsEnabled,
			"exclude_newrelic_header":        in.DistributedTracing.ExcludeNewrelicHeader,
			"cors_use_newrelic_header":       in.DistributedTracing.CorsUseNewrelicHeader,
			"cors_use_trace_context_headers": in.DistributedTracing.CorsUseTracecontextHeaders,
			"allowed_origins":                in.DistributedTracing.AllowedOrigins,
		},
		"ajax": map[string]interface{}{
			"deny_list": in.Ajax.DenyList,
		},
	}
	return browserMonitoringValues
}

// setting mobile values
func setMobileApplicationValues(d *schema.ResourceData, MobileSettings entities.AgentApplicationSettingsMobileBase) error {

	if err := d.Set("use_crash_reports", MobileSettings.UseCrashReports); err != nil {
		return fmt.Errorf("[DEBUG] Error setting use crash reports: %#v", err)
	}

	if err := d.Set("enable_application_exit_info", MobileSettings.ApplicationExitInfo.Enabled); err != nil {
		return fmt.Errorf("[DEBUG] Error setting application exit info values: %#v", err)
	}

	return nil
}
