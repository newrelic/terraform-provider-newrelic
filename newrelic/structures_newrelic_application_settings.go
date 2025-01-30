package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
)

func expandApplication(d *schema.ResourceData) *apm.Application {
	a := apm.Application{
		Name: d.Get("name").(string),
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] expanding application, id %s", err)
	} else {
		a.ID = id
	}

	a.Settings.ApmConfig.ApdexTarget = d.Get("app_apdex_threshold").(float64)
	a.Settings.ApmConfig.ApdexTarget = d.Get("end_user_apdex_threshold").(float64)
	a.Settings.ApmConfig.UseServerSideConfig = d.Get("enable_real_user_monitoring").(bool)

	return &a
}

func flattenApplication(a *apm.Application, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(a.ID))
	var err error

	err = d.Set("name", a.Name)
	if err != nil {
		return err
	}

	err = d.Set("app_apdex_threshold", a.Settings.ApmConfig.ApdexTarget)
	if err != nil {
		return err
	}

	err = d.Set("end_user_apdex_threshold", a.Settings.ApmConfig.ApdexTarget)
	if err != nil {
		return err
	}

	err = d.Set("enable_real_user_monitoring", a.Settings.ApmConfig.UseServerSideConfig)
	if err != nil {
		return err
	}

	return nil
}

func getBoolPointer(value bool) *bool {
	return &value
}

func getFloatPointer(value float64) *float64 {
	return &value
}

func getStringPointer(value string) *string {
	return &value
}

func expandApmConfigValues(v interface{}) *apm.AgentApplicationSettingsApmConfigInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})

	perms := &apm.AgentApplicationSettingsApmConfigInput{}

	if v, ok := values["apdex_target"]; ok {
		perms.ApdexTarget = v.(float64)
	}
	if v, ok := values["enable_server_side_config"]; ok {
		perms.UseServerSideConfig = getBoolPointer(v.(bool))
	}

	return perms
}

func expandTransactionTracerValues(v interface{}) *apm.AgentApplicationSettingsTransactionTracerInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})

	tracer := &apm.AgentApplicationSettingsTransactionTracerInput{}

	if v, ok := values["enabled"]; ok {
		tracer.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["explain_enabled"]; ok {
		tracer.ExplainEnabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["explain_threshold_value"]; ok {
		tracer.ExplainThresholdValue = getFloatPointer(v.(float64))
	}
	if v, ok := values["explain_threshold_type"]; ok {
		tracer.ExplainThresholdType = apm.AgentApplicationSettingsThresholdTypeEnum(v.(string))
	}
	if v, ok := values["log_sql"]; ok {
		tracer.LogSql = getBoolPointer(v.(bool))
	}
	if v, ok := values["record_sql"]; ok {
		tracer.RecordSql = apm.AgentApplicationSettingsRecordSqlEnum(v.(string))
	}
	if v, ok := values["stack_trace_threshold_value"]; ok {
		tracer.StackTraceThreshold = getFloatPointer(v.(float64))
	}
	if v, ok := values["transaction_threshold_type"]; ok {
		tracer.TransactionThresholdType = apm.AgentApplicationSettingsThresholdTypeEnum(v.(string))
	}
	if v, ok := values["transaction_threshold_value"]; ok {
		tracer.TransactionThresholdValue = getFloatPointer(v.(float64))
	}

	return tracer
}

func expandErrorCollectorValues(v interface{}) *apm.AgentApplicationSettingsErrorCollectorInput {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	values := v.([]interface{})[0].(map[string]interface{})

	collector := &apm.AgentApplicationSettingsErrorCollectorInput{}

	if v, ok := values["enabled"]; ok {
		collector.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := values["expected_error_classes"]; ok {
		collector.ExpectedErrorClasses = v.([]string)
	}
	if v, ok := values["expected_error_codes"]; ok {
		strSlice := v.([]string)
		var httpStatusSlice []apm.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, apm.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.ExpectedErrorCodes = httpStatusSlice
	}
	if v, ok := values["ignored_error_classes"]; ok {
		collector.IgnoredErrorClasses = v.([]string)
	}
	if v, ok := values["ignored_error_codes"]; ok {
		strSlice := v.([]string)
		var httpStatusSlice []apm.AgentApplicationSettingsErrorCollectorHttpStatus
		for _, str := range strSlice {
			httpStatusSlice = append(httpStatusSlice, apm.AgentApplicationSettingsErrorCollectorHttpStatus(str))
		}
		collector.IgnoredErrorCodes = httpStatusSlice
	}

	return collector
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

func expandApplicationCopy(d *schema.ResourceData) *apm.AgentApplicationSettingsUpdateInput {

	a := apm.AgentApplicationSettingsUpdateInput{}

	if v, ok := d.GetOk("alias"); ok {
		a.Alias = getStringPointer(v.(string))
	}

	if v, ok := d.GetOk("tracer_type"); ok {
		a.TracerType = &apm.AgentApplicationSettingsTracerTypeInput{Value: apm.AgentApplicationSettingsTracer(v.(string))}
	}

	if v, ok := d.GetOk("thread_profiler_enabled"); ok {
		if a.ThreadProfiler == nil {
			a.ThreadProfiler = &apm.AgentApplicationSettingsThreadProfilerInput{}
		}
		a.ThreadProfiler.Enabled = getBoolPointer(v.(bool))
	}
	// apm settings
	a.ApmConfig = expandApmConfigValues(d.Get("apm_config"))
	// transaction_tracing
	a.TransactionTracer = expandTransactionTracerValues(d.Get("transaction_tracing"))
	// error_collector
	a.ErrorCollector = expandErrorCollectorValues(d.Get("error_collector"))

	// mobile settings
	if v, ok := d.GetOk("use_crash_reports"); ok {
		a.MobileSettings.UseCrashReports = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("application_exit_info.0.enabled"); ok {
		a.MobileSettings.ApplicationExitInfo.Enabled = getBoolPointer(v.(bool))
	}

	// browser settings

	// browser_monitoring
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.enabled"); ok {
		a.BrowserMonitoring.DistributedTracing.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.cors_enabled"); ok {
		a.BrowserMonitoring.DistributedTracing.CorsEnabled = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.exclude_newrelic_header"); ok {
		a.BrowserMonitoring.DistributedTracing.ExcludeNewrelicHeader = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.cors_use_newrelic_header"); ok {
		a.BrowserMonitoring.DistributedTracing.CorsUseNewrelicHeader = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.cors_use_tracecontext_headers"); ok {
		a.BrowserMonitoring.DistributedTracing.CorsUseTracecontextHeaders = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("browser_monitoring.0.distributed_tracing.0.allowed_origins"); ok {
		a.BrowserMonitoring.DistributedTracing.AllowedOrigins = v.([]string)
	}
	if v, ok := d.GetOk("browser_monitoring.0.ajax.0.deny_list"); ok {
		a.BrowserMonitoring.Ajax.DenyList = v.([]string)
	}
	// session replay settings
	if v, ok := d.GetOk("session_replay.0.auto_start"); ok {
		a.SessionReplay.AutoStart = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("session_replay.0.enabled"); ok {
		a.SessionReplay.Enabled = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("session_replay.0.error_sampling_rate"); ok {
		a.SessionReplay.ErrorSamplingRate = v.(float64)
	}
	if v, ok := d.GetOk("session_replay.0.sampling_rate"); ok {
		a.SessionReplay.SamplingRate = v.(float64)
	}
	if _, ok := d.GetOk("session_replay.0.mask_input_options"); ok {
		a.SessionReplay.MaskInputOptions = expandMaskInputOptions(d.Get("session_replay.0.mask_input_options"))
	}

	if v, ok := d.GetOk("session_replay.0.mask_all_inputs"); ok {
		a.SessionReplay.MaskAllInputs = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("session_replay.0.block_selector"); ok {
		a.SessionReplay.BlockSelector = getStringPointer(v.(string))
	}
	if v, ok := d.GetOk("session_replay.0.mask_text_selector"); ok {
		a.SessionReplay.MaskTextSelector = getStringPointer(v.(string))
	}

	// session trace settings
	if v, ok := d.GetOk("session_trace.0.error_sampling_rate"); ok {
		a.SessionTrace.ErrorSamplingRate = v.(float64)
	}
	if v, ok := d.GetOk("session_trace.0.mode"); ok {
		a.SessionTrace.Mode = apm.AgentApplicationSettingsSessionTraceModeInput(v.(string))
	}
	if v, ok := d.GetOk("session_trace.0.sampling_rate"); ok {
		a.SessionTrace.SamplingRate = v.(float64)
	}
	if v, ok := d.GetOk("session_trace.0.enabled"); ok {
		a.SessionTrace.Enabled = getBoolPointer(v.(bool))
	}

	// Mobile settings
	if v, ok := d.GetOk("use_crash_reports"); ok {
		a.MobileSettings.UseCrashReports = getBoolPointer(v.(bool))
	}
	if v, ok := d.GetOk("application_exit_info.0.enabled"); ok {
		a.MobileSettings.ApplicationExitInfo.Enabled = getBoolPointer(v.(bool))
	}
	return &a
}

// setting APM values
func setAPMApplicationValues(d *schema.ResourceData, ApmSettings entities.AgentApplicationSettingsApmBase) error {

	var err error
	err = d.Set("alias", ApmSettings.Alias)
	if err != nil {
		return err
	}

	if _, isok := d.GetOk("apm_config"); isok {
		if err := d.Set("apm_config", flattenApmConfigValues(ApmSettings.ApmConfig)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting apm config values: %#v", err)
		}
	}

	if _, isok := d.GetOk("thread_profiler_enabled"); isok {
		if err = d.Set("thread_profiler_enabled", ApmSettings.ThreadProfiler.Enabled); err != nil {
			return err
		}
	}

	if _, isok := d.GetOk("tracer_type"); isok {
		err = d.Set("tracer_type", ApmSettings.TracerType)
		if err != nil {
			return err
		}
	}

	if _, isok := d.GetOk("transaction_tracing"); isok {
		if err := d.Set("transaction_tracing", flattenTransactionTracingValues(ApmSettings.TransactionTracer)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting transaction tracer values: %#v", err)
		}
	}

	if _, isok := d.GetOk("error_collector"); isok {
		if err := d.Set("error_collector", flattenErrorCollectorValues(ApmSettings.ErrorCollector)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting error collector values: %#v", err)
		}
	}

	return nil
}

func flattenApmConfigValues(in entities.AgentApplicationSettingsApmConfig) []interface{} {
	apmConfigValues := map[string]interface{}{
		"apdex_target":              in.ApdexTarget,
		"enable_server_side_config": in.UseServerSideConfig,
	}
	return []interface{}{apmConfigValues}
}

func flattenTransactionTracingValues(in entities.AgentApplicationSettingsTransactionTracer) []interface{} {
	transactionTracingValues := map[string]interface{}{
		"enabled":                     in.Enabled,
		"log_sql":                     in.LogSql,
		"record_sql":                  in.RecordSql,
		"explain_enabled":             in.ExplainEnabled,
		"explain_threshold_value":     in.ExplainThresholdValue,
		"explain_threshold_type":      in.ExplainThresholdType,
		"stack_trace_threshold_value": in.StackTraceThreshold,
		"transaction_threshold_type":  in.TransactionThresholdType,
		"transaction_threshold_value": in.TransactionThresholdValue,
	}
	return []interface{}{transactionTracingValues}
}

func flattenErrorCollectorValues(in entities.AgentApplicationSettingsErrorCollector) []interface{} {
	var expectedErrorCodes []string
	var ignoredErrorCodes []string
	errorCollectorValues := map[string]interface{}{
		"enabled":                in.Enabled,
		"expected_error_classes": in.ExpectedErrorClasses,
		"expected_error_codes":   in.ExpectedErrorCodes,
		"ignored_error_classes":  in.IgnoredErrorClasses,
		"ignored_error_codes":    in.IgnoredErrorCodes,
	}

	for _, code := range in.ExpectedErrorCodes {
		expectedErrorCodes = append(expectedErrorCodes, string(code))
	}
	errorCollectorValues["expected_error_codes"] = expectedErrorCodes

	for _, code := range in.IgnoredErrorCodes {
		ignoredErrorCodes = append(ignoredErrorCodes, string(code))
	}
	errorCollectorValues["ignored_error_codes"] = ignoredErrorCodes

	return []interface{}{errorCollectorValues}
}

// setting browser values
func setBrowserApplicationValues(d *schema.ResourceData, BrowserSettings entities.AgentApplicationSettingsBrowserBase) error {

	if err := d.Set("browser_config", flattenBrowserConfigValues(BrowserSettings.BrowserConfig)); err != nil {
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

func flattenBrowserConfigValues(in entities.AgentApplicationSettingsBrowserConfig) []map[string]interface{} {
	apmConfigValues := make([]map[string]interface{}, 1)
	apmConfigValues[0] = map[string]interface{}{
		"apdex_target": in.ApdexTarget,
	}
	return apmConfigValues
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

	if err := d.Set("application_exit_info", flattenApplicationExitInfoValues(MobileSettings.ApplicationExitInfo)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting application exit info values: %#v", err)
	}

	return nil
}

func flattenApplicationExitInfoValues(in entities.AgentApplicationSettingsApplicationExitInfo) []map[string]interface{} {
	applicationExitInfoValues := make([]map[string]interface{}, 1)
	applicationExitInfoValues[0] = map[string]interface{}{
		"enabled": in.Enabled,
	}
	return applicationExitInfoValues
}
