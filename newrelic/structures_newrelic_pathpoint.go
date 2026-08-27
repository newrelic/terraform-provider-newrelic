package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

func expandPathpointFlowInput(d *schema.ResourceData) (*pathpoint.PathPointFlowInput, error) {
	input := &pathpoint.PathPointFlowInput{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("health_rollup"); ok {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v.(string))
	}
	if v, ok := d.GetOk("refresh_interval"); ok {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v.(string))
	}
	if v, ok := d.GetOk("kpis"); ok {
		input.Kpis = expandPathpointKpiInputList(v.([]interface{}))
	}
	if v, ok := d.GetOk("stages"); ok {
		input.Stages = expandPathpointStageInputList(v.([]interface{}))
	}

	return input, nil
}

func expandPathpointFlowUpdateInput(d *schema.ResourceData) (*pathpoint.PathPointFlowUpdateInput, error) {
	input := &pathpoint.PathPointFlowUpdateInput{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("health_rollup"); ok {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v.(string))
	}
	if v, ok := d.GetOk("refresh_interval"); ok {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v.(string))
	}
	if v, ok := d.GetOk("flow_id"); ok {
		input.ID = v.(string)
	}
	if v, ok := d.GetOk("kpis"); ok {
		input.Kpis = expandPathpointKpiUpdateInputList(v.([]interface{}))
	}
	if v, ok := d.GetOk("stages"); ok {
		input.Stages = expandPathpointStageUpdateInputList(v.([]interface{}))
	}
	if v, ok := d.GetOk("version"); ok {
		input.Version = pathpoint.EpochMilliseconds(v.(string))
	}

	return input, nil
}

func expandPathpointScopeInput(d *schema.ResourceData, accountID int) pathpoint.PathPointScopeInput {
	return pathpoint.PathPointScopeInput{
		ID:   accountID,
		Type: pathpoint.PathPointScopeTypeTypes.ACCOUNT,
	}
}

func expandPathpointKpiInputList(list []interface{}) []pathpoint.PathPointKpiInput {
	out := make([]pathpoint.PathPointKpiInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointKpiInput(cfg)
	}
	return out
}

func expandPathpointKpiInput(cfg map[string]interface{}) pathpoint.PathPointKpiInput {
	input := pathpoint.PathPointKpiInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["account_id"].(int); ok && v != 0 {
		input.AccountID = v
	}
	if v, ok := cfg["category"].(string); ok && v != "" {
		input.Category = v
	}
	if v, ok := cfg["description"].(string); ok && v != "" {
		input.Description = v
	}
	if v, ok := cfg["query"].([]interface{}); ok && len(v) > 0 {
		input.Query = expandPathpointKpiNRQLInput(v[0].(map[string]interface{}))
	}
	return input
}

func expandPathpointKpiUpdateInputList(list []interface{}) []pathpoint.PathPointKpiUpdateInput {
	out := make([]pathpoint.PathPointKpiUpdateInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointKpiUpdateInput(cfg)
	}
	return out
}

func expandPathpointKpiUpdateInput(cfg map[string]interface{}) pathpoint.PathPointKpiUpdateInput {
	input := pathpoint.PathPointKpiUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["kpi_id"].(string); ok && v != "" {
		input.ID = v
	}
	if v, ok := cfg["account_id"].(int); ok && v != 0 {
		input.AccountID = v
	}
	if v, ok := cfg["category"].(string); ok && v != "" {
		input.Category = v
	}
	if v, ok := cfg["description"].(string); ok && v != "" {
		input.Description = v
	}
	if v, ok := cfg["query"].([]interface{}); ok && len(v) > 0 {
		input.Query = expandPathpointKpiNRQLInput(v[0].(map[string]interface{}))
	}
	return input
}

func expandPathpointKpiNRQLInput(cfg map[string]interface{}) pathpoint.PathPointKpiNRQLInput {
	input := pathpoint.PathPointKpiNRQLInput{
		From: pathpoint.NRQL(cfg["from"].(string)),
	}
	if v, ok := cfg["where"].(string); ok && v != "" {
		input.Where = v
	}
	if v, ok := cfg["select"].([]interface{}); ok && len(v) > 0 {
		sel := expandPathpointKpiNRQLSelectInput(v[0].(map[string]interface{}))
		input.Select = sel
	}
	if v, ok := cfg["time_window"].([]interface{}); ok && len(v) > 0 {
		tw := expandPathpointKpiTimeWindowInput(v[0].(map[string]interface{}))
		input.TimeWindow = &tw
	}
	return input
}

func expandPathpointKpiNRQLSelectInput(cfg map[string]interface{}) pathpoint.PathPointKpiNRQLSelectInput {
	input := pathpoint.PathPointKpiNRQLSelectInput{
		AggregationType: pathpoint.PathPointKpiNRQLAggregations(cfg["aggregation_type"].(string)),
	}
	if v, ok := cfg["alias"].(string); ok && v != "" {
		input.Alias = v
	}
	if v, ok := cfg["attribute"].(string); ok && v != "" {
		input.Attribute = v
	}
	if v, ok := cfg["threshold"].(float64); ok && v != 0 {
		input.Threshold = v
	}
	return input
}

func expandPathpointKpiTimeWindowInput(cfg map[string]interface{}) pathpoint.PathPointKpiTimeWindowInput {
	input := pathpoint.PathPointKpiTimeWindowInput{}
	if v, ok := cfg["custom_range"].(string); ok && v != "" {
		input.CustomRange = pathpoint.NRQL(v)
	}
	if v, ok := cfg["relative_range"].([]interface{}); ok && len(v) > 0 {
		rr := expandPathpointKpiTimeWindowRelativeRangeInput(v[0].(map[string]interface{}))
		input.RelativeRange = &rr
	}
	return input
}

func expandPathpointKpiTimeWindowRelativeRangeInput(cfg map[string]interface{}) pathpoint.PathPointKpiTimeWindowRelativeRangeInput {
	input := pathpoint.PathPointKpiTimeWindowRelativeRangeInput{
		Since: pathpoint.PathPointKpiTimeDuration(cfg["since"].(string)),
	}
	if v, ok := cfg["compare_against"].(string); ok && v != "" {
		input.CompareAgainst = pathpoint.PathPointKpiTimeDuration(v)
	}
	return input
}

func expandPathpointStageInputList(list []interface{}) []pathpoint.PathPointStageInput {
	out := make([]pathpoint.PathPointStageInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointStageInput(cfg)
	}
	return out
}

func expandPathpointStageInput(cfg map[string]interface{}) pathpoint.PathPointStageInput {
	input := pathpoint.PathPointStageInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["health_rollup"].(string); ok && v != "" {
		input.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	if v, ok := cfg["link"].(string); ok && v != "" {
		input.Link = v
	}
	if v, ok := cfg["related"].([]interface{}); ok && len(v) > 0 {
		input.Related = expandPathpointRelatedInput(v[0].(map[string]interface{}))
	}
	if v, ok := cfg["stage_kpis"].([]interface{}); ok && len(v) > 0 {
		input.StageKpis = expandPathpointKpiInputList(v)
	}
	if v, ok := cfg["levels"].([]interface{}); ok && len(v) > 0 {
		input.Levels = expandPathpointLevelInputList(v)
	}
	return input
}

func expandPathpointStageUpdateInputList(list []interface{}) []pathpoint.PathPointStageUpdateInput {
	out := make([]pathpoint.PathPointStageUpdateInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointStageUpdateInput(cfg)
	}
	return out
}

func expandPathpointStageUpdateInput(cfg map[string]interface{}) pathpoint.PathPointStageUpdateInput {
	input := pathpoint.PathPointStageUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["stage_id"].(string); ok && v != "" {
		input.ID = v
	}
	if v, ok := cfg["health_rollup"].(string); ok && v != "" {
		input.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	if v, ok := cfg["link"].(string); ok && v != "" {
		input.Link = v
	}
	if v, ok := cfg["related"].([]interface{}); ok && len(v) > 0 {
		input.Related = expandPathpointRelatedInput(v[0].(map[string]interface{}))
	}
	if v, ok := cfg["stage_kpis"].([]interface{}); ok && len(v) > 0 {
		input.StageKpis = expandPathpointKpiUpdateInputList(v)
	}
	if v, ok := cfg["levels"].([]interface{}); ok && len(v) > 0 {
		input.Levels = expandPathpointLevelUpdateInputList(v)
	}
	return input
}

func expandPathpointRelatedInput(cfg map[string]interface{}) pathpoint.PathPointRelatedInput {
	input := pathpoint.PathPointRelatedInput{}
	if v, ok := cfg["source"].(bool); ok {
		input.Source = v
	}
	if v, ok := cfg["target"].(bool); ok {
		input.Target = v
	}
	return input
}

func expandPathpointLevelInputList(list []interface{}) []pathpoint.PathPointLevelInput {
	out := make([]pathpoint.PathPointLevelInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointLevelInput(cfg)
	}
	return out
}

func expandPathpointLevelInput(cfg map[string]interface{}) pathpoint.PathPointLevelInput {
	input := pathpoint.PathPointLevelInput{}
	if v, ok := cfg["steps"].([]interface{}); ok && len(v) > 0 {
		input.Steps = expandPathpointStepInputList(v)
	}
	return input
}

func expandPathpointLevelUpdateInputList(list []interface{}) []pathpoint.PathPointLevelUpdateInput {
	out := make([]pathpoint.PathPointLevelUpdateInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointLevelUpdateInput(cfg)
	}
	return out
}

func expandPathpointLevelUpdateInput(cfg map[string]interface{}) pathpoint.PathPointLevelUpdateInput {
	input := pathpoint.PathPointLevelUpdateInput{}
	if v, ok := cfg["level_id"].(string); ok && v != "" {
		input.ID = v
	}
	if v, ok := cfg["steps"].([]interface{}); ok && len(v) > 0 {
		input.Steps = expandPathpointStepUpdateInputList(v)
	}
	return input
}

func expandPathpointStepInputList(list []interface{}) []pathpoint.PathPointStepInput {
	out := make([]pathpoint.PathPointStepInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointStepInput(cfg)
	}
	return out
}

func expandPathpointStepInput(cfg map[string]interface{}) pathpoint.PathPointStepInput {
	input := pathpoint.PathPointStepInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	if v, ok := cfg["link"].(string); ok && v != "" {
		input.Link = v
	}
	if v, ok := cfg["scoped_accounts"].([]interface{}); ok && len(v) > 0 {
		accounts := make([]int, len(v))
		for i, a := range v {
			accounts[i] = a.(int)
		}
		input.ScopedAccounts = accounts
	}
	if v, ok := cfg["config"].([]interface{}); ok && len(v) > 0 {
		input.Config = expandPathpointStepStatusThresholdInput(v[0].(map[string]interface{}))
	}
	if v, ok := cfg["entity_search_query"].([]interface{}); ok && len(v) > 0 {
		esq := expandPathpointSignalQueryInput(v[0].(map[string]interface{}))
		input.EntitySearchQuery = &esq
	}
	if v, ok := cfg["signals"].([]interface{}); ok && len(v) > 0 {
		input.Signals = expandPathpointSignalInputList(v)
	}
	return input
}

func expandPathpointStepUpdateInputList(list []interface{}) []pathpoint.PathPointStepUpdateInput {
	out := make([]pathpoint.PathPointStepUpdateInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointStepUpdateInput(cfg)
	}
	return out
}

func expandPathpointStepUpdateInput(cfg map[string]interface{}) pathpoint.PathPointStepUpdateInput {
	input := pathpoint.PathPointStepUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["step_id"].(string); ok && v != "" {
		input.ID = v
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	if v, ok := cfg["link"].(string); ok && v != "" {
		input.Link = v
	}
	if v, ok := cfg["scoped_accounts"].([]interface{}); ok && len(v) > 0 {
		accounts := make([]int, len(v))
		for i, a := range v {
			accounts[i] = a.(int)
		}
		input.ScopedAccounts = accounts
	}
	if v, ok := cfg["config"].([]interface{}); ok && len(v) > 0 {
		input.Config = expandPathpointStepStatusThresholdInput(v[0].(map[string]interface{}))
	}
	if v, ok := cfg["entity_search_query"].([]interface{}); ok && len(v) > 0 {
		esq := expandPathpointSignalQueryInput(v[0].(map[string]interface{}))
		input.EntitySearchQuery = &esq
	}
	if v, ok := cfg["signals"].([]interface{}); ok && len(v) > 0 {
		input.Signals = expandPathpointSignalInputList(v)
	}
	return input
}

func expandPathpointStepStatusThresholdInput(cfg map[string]interface{}) pathpoint.PathPointStepStatusThresholdInput {
	input := pathpoint.PathPointStepStatusThresholdInput{}
	if v, ok := cfg["health_rollup"].(string); ok && v != "" {
		input.HealthRollup = pathpoint.PathPointStepHealthRollup(v)
	}
	if v, ok := cfg["threshold_type"].(string); ok && v != "" {
		input.ThresholdType = pathpoint.PathPointThresholdType(v)
	}
	if v, ok := cfg["threshold_value"].(int); ok && v != 0 {
		input.ThresholdValue = v
	}
	return input
}

func expandPathpointSignalQueryInput(cfg map[string]interface{}) pathpoint.PathPointSignalQueryInput {
	input := pathpoint.PathPointSignalQueryInput{
		Query: cfg["query"].(string),
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	return input
}

func expandPathpointSignalInputList(list []interface{}) []pathpoint.PathPointSignalInput {
	out := make([]pathpoint.PathPointSignalInput, len(list))
	for i, v := range list {
		cfg := v.(map[string]interface{})
		out[i] = expandPathpointSignalInput(cfg)
	}
	return out
}

func expandPathpointSignalInput(cfg map[string]interface{}) pathpoint.PathPointSignalInput {
	input := pathpoint.PathPointSignalInput{
		GUID: pathpoint.EntityGUID(cfg["guid"].(string)),
	}
	if v, ok := cfg["name"].(string); ok && v != "" {
		input.Name = v
	}
	if v, ok := cfg["type"].(string); ok && v != "" {
		input.Type = pathpoint.PathPointSignalType(v)
	}
	if v, ok := cfg["is_excluded"].(bool); ok {
		input.IsExcluded = v
	}
	return input
}
func flattenPathpointFlow(result *pathpoint.PathPointFlowResult, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	_ = d.Set("name", result.Name)
	_ = d.Set("category", result.Category)
	_ = d.Set("description", result.Description)
	_ = d.Set("health_rollup", string(result.HealthRollup))
	_ = d.Set("refresh_interval", string(result.RefreshInterval))
	_ = d.Set("guid", string(result.GUID))
	_ = d.Set("flow_id", result.ID)
	_ = d.Set("version", result.Version.String())
	_ = d.Set("health_status", string(result.HealthStatus))
	_ = d.Set("message", result.Message)

	excludedKpis := make([]interface{}, len(result.ExcludedKpis))
	for i, v := range result.ExcludedKpis {
		excludedKpis[i] = v
	}
	_ = d.Set("excluded_kpis", excludedKpis)

	if err := d.Set("kpis", flattenPathpointKpis(result.Kpis)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `kpis`: %v", err)
	}

	if err := d.Set("stages", flattenPathpointStages(result.Stages.Items)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `stages`: %v", err)
	}

	return nil
}

func flattenPathpointKpis(kpis []pathpoint.PathPointKpi) []interface{} {
	out := make([]interface{}, len(kpis))
	for i, k := range kpis {
		out[i] = flattenPathpointKpi(k)
	}
	return out
}

func flattenPathpointKpi(k pathpoint.PathPointKpi) map[string]interface{} {
	m := map[string]interface{}{
		"kpi_id":      k.ID,
		"account_id":  k.AccountID,
		"name":        k.Name,
		"category":    k.Category,
		"description": k.Description,
	}

	queryList := flattenPathpointKpiNRQL(k.Query)
	m["query"] = queryList

	return m
}

func flattenPathpointKpiNRQL(q pathpoint.PathPointKpiNRQL) []interface{} {
	m := map[string]interface{}{
		"from":  string(q.From),
		"where": q.Where,
	}

	selectList := flattenPathpointKpiNRQLSelect(q.Select)
	m["select"] = selectList

	timeWindowList := flattenPathpointKpiTimeWindow(q.TimeWindow)
	m["time_window"] = timeWindowList

	return []interface{}{m}
}

func flattenPathpointKpiNRQLSelect(s pathpoint.PathPointKpiNRQLSelect) []interface{} {
	m := map[string]interface{}{
		"aggregation_type": string(s.AggregationType),
		"alias":            s.Alias,
		"attribute":        s.Attribute,
		"threshold":        s.Threshold,
	}
	return []interface{}{m}
}

func flattenPathpointKpiTimeWindow(tw pathpoint.PathPointKpiTimeWindow) []interface{} {
	if tw == (pathpoint.PathPointKpiTimeWindow{}) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"custom_range": string(tw.CustomRange),
	}

	relativeRangeList := flattenPathpointKpiTimeWindowRelativeRange(tw.RelativeRange)
	m["relative_range"] = relativeRangeList

	return []interface{}{m}
}

func flattenPathpointKpiTimeWindowRelativeRange(rr pathpoint.PathPointKpiTimeWindowRelativeRange) []interface{} {
	if rr == (pathpoint.PathPointKpiTimeWindowRelativeRange{}) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"since":           string(rr.Since),
		"compare_against": string(rr.CompareAgainst),
	}
	return []interface{}{m}
}

func flattenPathpointStages(stages []pathpoint.PathPointStage) []interface{} {
	out := make([]interface{}, len(stages))
	for i, s := range stages {
		out[i] = flattenPathpointStage(s)
	}
	return out
}

func flattenPathpointStage(s pathpoint.PathPointStage) map[string]interface{} {
	m := map[string]interface{}{
		"stage_id":      s.ID,
		"name":          s.Name,
		"health_rollup": string(s.HealthRollup),
		"is_excluded":   s.IsExcluded,
		"link":          s.Link,
		"health_status": string(s.HealthStatus),
	}

	relatedList := flattenPathpointRelated(s.Related)
	m["related"] = relatedList

	m["stage_kpis"] = flattenPathpointKpis(s.StageKpis)

	m["levels"] = flattenPathpointLevels(s.Levels.Items)

	return m
}

func flattenPathpointRelated(r pathpoint.PathPointRelated) []interface{} {
	m := map[string]interface{}{
		"source": r.Source,
		"target": r.Target,
	}
	return []interface{}{m}
}

func flattenPathpointLevels(levels []pathpoint.PathPointLevel) []interface{} {
	out := make([]interface{}, len(levels))
	for i, l := range levels {
		out[i] = flattenPathpointLevel(l)
	}
	return out
}

func flattenPathpointLevel(l pathpoint.PathPointLevel) map[string]interface{} {
	m := map[string]interface{}{
		"level_id":      l.ID,
		"health_status": string(l.HealthStatus),
	}

	m["steps"] = flattenPathpointSteps(l.Steps.Items)

	return m
}

func flattenPathpointSteps(steps []pathpoint.PathPointStep) []interface{} {
	out := make([]interface{}, len(steps))
	for i, s := range steps {
		out[i] = flattenPathpointStep(s)
	}
	return out
}

func flattenPathpointStep(s pathpoint.PathPointStep) map[string]interface{} {
	m := map[string]interface{}{
		"step_id":       s.ID,
		"name":          s.Name,
		"is_excluded":   s.IsExcluded,
		"link":          s.Link,
		"health_status": string(s.HealthStatus),
	}

	scopedAccounts := make([]interface{}, len(s.ScopedAccounts))
	for i, v := range s.ScopedAccounts {
		scopedAccounts[i] = v
	}
	m["scoped_accounts"] = scopedAccounts

	m["config"] = flattenPathpointStepConfig(s.Config)
	m["entity_search_query"] = flattenPathpointSignalQuery(s.EntitySearchQuery)
	m["signals"] = flattenPathpointSignals(s.Signals)

	return m
}

func flattenPathpointStepConfig(c pathpoint.PathPointStepStatusThreshold) []interface{} {
	if c == (pathpoint.PathPointStepStatusThreshold{}) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"health_rollup":   string(c.HealthRollup),
		"threshold_type":  string(c.ThresholdType),
		"threshold_value": c.ThresholdValue,
	}
	return []interface{}{m}
}

func flattenPathpointSignalQuery(q pathpoint.PathPointSignalQuery) []interface{} {
	if q == (pathpoint.PathPointSignalQuery{}) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"query":       q.Query,
		"is_excluded": q.IsExcluded,
	}
	return []interface{}{m}
}

func flattenPathpointSignals(signals []pathpoint.PathPointSignal) []interface{} {
	out := make([]interface{}, len(signals))
	for i, s := range signals {
		out[i] = flattenPathpointSignal(s)
	}
	return out
}

func flattenPathpointSignal(s pathpoint.PathPointSignal) map[string]interface{} {
	return map[string]interface{}{
		"guid":        string(s.GUID),
		"name":        s.Name,
		"type":        string(s.Type),
		"is_excluded": s.IsExcluded,
	}
}