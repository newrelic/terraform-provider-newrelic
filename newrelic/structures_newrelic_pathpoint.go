package newrelic

import (
	"fmt"

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
		kpis, err := expandPathpointKpiInputList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		input.Kpis = kpis
	}
	if v, ok := d.GetOk("stages"); ok {
		stages, err := expandPathpointStageInputList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		input.Stages = stages
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
	if v, ok := d.GetOk("kpis"); ok {
		kpis, err := expandPathpointKpiUpdateInputList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		input.Kpis = kpis
	}
	if v, ok := d.GetOk("stages"); ok {
		stages, err := expandPathpointStageUpdateInputList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		input.Stages = stages
	}

	// version is required for update
	if v, ok := d.GetOk("version"); ok {
		input.Version = pathpoint.EpochMilliseconds(v.(string))
	}

	return input, nil
}

func expandPathpointKpiInputList(list []interface{}) ([]pathpoint.PathPointKpiInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointKpiInput{}, nil
	}
	out := make([]pathpoint.PathPointKpiInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid kpi element at index %d", i)
		}
		kpi, err := expandPathpointKpiInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = kpi
	}
	return out, nil
}

func expandPathpointKpiInput(cfg map[string]interface{}) (pathpoint.PathPointKpiInput, error) {
	kpi := pathpoint.PathPointKpiInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["account_id"]; ok {
		kpi.AccountID = v.(int)
	}
	if v, ok := cfg["category"]; ok {
		kpi.Category = v.(string)
	}
	if v, ok := cfg["description"]; ok {
		kpi.Description = v.(string)
	}
	if v, ok := cfg["query"]; ok {
		queryList := v.([]interface{})
		if len(queryList) > 0 {
			q, err := expandPathpointKpiNRQLInput(queryList[0].(map[string]interface{}))
			if err != nil {
				return kpi, err
			}
			kpi.Query = q
		}
	}
	return kpi, nil
}

func expandPathpointKpiUpdateInputList(list []interface{}) ([]pathpoint.PathPointKpiUpdateInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointKpiUpdateInput{}, nil
	}
	out := make([]pathpoint.PathPointKpiUpdateInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid kpi element at index %d", i)
		}
		kpi, err := expandPathpointKpiUpdateInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = kpi
	}
	return out, nil
}

func expandPathpointKpiUpdateInput(cfg map[string]interface{}) (pathpoint.PathPointKpiUpdateInput, error) {
	kpi := pathpoint.PathPointKpiUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["kpi_id"]; ok {
		kpi.ID = v.(string)
	}
	if v, ok := cfg["account_id"]; ok {
		kpi.AccountID = v.(int)
	}
	if v, ok := cfg["category"]; ok {
		kpi.Category = v.(string)
	}
	if v, ok := cfg["description"]; ok {
		kpi.Description = v.(string)
	}
	if v, ok := cfg["query"]; ok {
		queryList := v.([]interface{})
		if len(queryList) > 0 {
			q, err := expandPathpointKpiNRQLInput(queryList[0].(map[string]interface{}))
			if err != nil {
				return kpi, err
			}
			kpi.Query = q
		}
	}
	return kpi, nil
}

func expandPathpointKpiNRQLInput(cfg map[string]interface{}) (pathpoint.PathPointKpiNRQLInput, error) {
	q := pathpoint.PathPointKpiNRQLInput{
		From: pathpoint.NRQL(cfg["from"].(string)),
	}
	if v, ok := cfg["where"]; ok {
		q.Where = v.(string)
	}
	if v, ok := cfg["select"]; ok {
		selectList := v.([]interface{})
		if len(selectList) > 0 {
			s := expandPathpointKpiNRQLSelectInput(selectList[0].(map[string]interface{}))
			q.Select = s
		}
	}
	if v, ok := cfg["time_window"]; ok {
		twList := v.([]interface{})
		if len(twList) > 0 {
			tw, err := expandPathpointKpiTimeWindowInput(twList[0].(map[string]interface{}))
			if err != nil {
				return q, err
			}
			q.TimeWindow = tw
		}
	}
	return q, nil
}

func expandPathpointKpiNRQLSelectInput(cfg map[string]interface{}) pathpoint.PathPointKpiNRQLSelectInput {
	s := pathpoint.PathPointKpiNRQLSelectInput{
		AggregationType: pathpoint.PathPointKpiNRQLAggregations(cfg["aggregation_type"].(string)),
	}
	if v, ok := cfg["alias"]; ok {
		s.Alias = v.(string)
	}
	if v, ok := cfg["attribute"]; ok {
		s.Attribute = v.(string)
	}
	if v, ok := cfg["threshold"]; ok {
		s.Threshold = v.(float64)
	}
	return s
}

func expandPathpointKpiTimeWindowInput(cfg map[string]interface{}) (*pathpoint.PathPointKpiTimeWindowInput, error) {
	tw := &pathpoint.PathPointKpiTimeWindowInput{}
	if v, ok := cfg["custom_range"]; ok {
		tw.CustomRange = pathpoint.NRQL(v.(string))
	}
	if v, ok := cfg["relative_range"]; ok {
		rrList := v.([]interface{})
		if len(rrList) > 0 {
			rr := expandPathpointKpiTimeWindowRelativeRangeInput(rrList[0].(map[string]interface{}))
			tw.RelativeRange = rr
		}
	}
	return tw, nil
}

func expandPathpointKpiTimeWindowRelativeRangeInput(cfg map[string]interface{}) *pathpoint.PathPointKpiTimeWindowRelativeRangeInput {
	rr := &pathpoint.PathPointKpiTimeWindowRelativeRangeInput{
		Since: pathpoint.PathPointKpiTimeDuration(cfg["since"].(string)),
	}
	if v, ok := cfg["compare_against"]; ok {
		rr.CompareAgainst = pathpoint.PathPointKpiTimeDuration(v.(string))
	}
	return rr
}

func expandPathpointStageInputList(list []interface{}) ([]pathpoint.PathPointStageInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointStageInput{}, nil
	}
	out := make([]pathpoint.PathPointStageInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid stage element at index %d", i)
		}
		stage, err := expandPathpointStageInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = stage
	}
	return out, nil
}

func expandPathpointStageInput(cfg map[string]interface{}) (pathpoint.PathPointStageInput, error) {
	stage := pathpoint.PathPointStageInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["health_rollup"]; ok {
		stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v.(string))
	}
	if v, ok := cfg["is_excluded"]; ok {
		stage.IsExcluded = v.(bool)
	}
	if v, ok := cfg["link"]; ok {
		stage.Link = v.(string)
	}
	if v, ok := cfg["related"]; ok {
		relList := v.([]interface{})
		if len(relList) > 0 {
			stage.Related = expandPathpointRelatedInput(relList[0].(map[string]interface{}))
		}
	}
	if v, ok := cfg["stage_kpis"]; ok {
		kpis, err := expandPathpointKpiInputList(v.([]interface{}))
		if err != nil {
			return stage, err
		}
		stage.StageKpis = kpis
	}
	if v, ok := cfg["levels"]; ok {
		levels, err := expandPathpointLevelInputList(v.([]interface{}))
		if err != nil {
			return stage, err
		}
		stage.Levels = levels
	}
	return stage, nil
}

func expandPathpointStageUpdateInputList(list []interface{}) ([]pathpoint.PathPointStageUpdateInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointStageUpdateInput{}, nil
	}
	out := make([]pathpoint.PathPointStageUpdateInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid stage element at index %d", i)
		}
		stage, err := expandPathpointStageUpdateInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = stage
	}
	return out, nil
}

func expandPathpointStageUpdateInput(cfg map[string]interface{}) (pathpoint.PathPointStageUpdateInput, error) {
	stage := pathpoint.PathPointStageUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["stage_id"]; ok {
		stage.ID = v.(string)
	}
	if v, ok := cfg["health_rollup"]; ok {
		stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v.(string))
	}
	if v, ok := cfg["is_excluded"]; ok {
		stage.IsExcluded = v.(bool)
	}
	if v, ok := cfg["link"]; ok {
		stage.Link = v.(string)
	}
	if v, ok := cfg["related"]; ok {
		relList := v.([]interface{})
		if len(relList) > 0 {
			stage.Related = expandPathpointRelatedInput(relList[0].(map[string]interface{}))
		}
	}
	if v, ok := cfg["stage_kpis"]; ok {
		kpis, err := expandPathpointKpiUpdateInputList(v.([]interface{}))
		if err != nil {
			return stage, err
		}
		stage.StageKpis = kpis
	}
	if v, ok := cfg["levels"]; ok {
		levels, err := expandPathpointLevelUpdateInputList(v.([]interface{}))
		if err != nil {
			return stage, err
		}
		stage.Levels = levels
	}
	return stage, nil
}

func expandPathpointRelatedInput(cfg map[string]interface{}) pathpoint.PathPointRelatedInput {
	related := pathpoint.PathPointRelatedInput{}
	if v, ok := cfg["source"]; ok {
		related.Source = v.(bool)
	}
	if v, ok := cfg["target"]; ok {
		related.Target = v.(bool)
	}
	return related
}

func expandPathpointLevelInputList(list []interface{}) ([]pathpoint.PathPointLevelInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointLevelInput{}, nil
	}
	out := make([]pathpoint.PathPointLevelInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid level element at index %d", i)
		}
		level, err := expandPathpointLevelInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = level
	}
	return out, nil
}

func expandPathpointLevelInput(cfg map[string]interface{}) (pathpoint.PathPointLevelInput, error) {
	level := pathpoint.PathPointLevelInput{}
	if v, ok := cfg["steps"]; ok {
		steps, err := expandPathpointStepInputList(v.([]interface{}))
		if err != nil {
			return level, err
		}
		level.Steps = steps
	}
	return level, nil
}

func expandPathpointLevelUpdateInputList(list []interface{}) ([]pathpoint.PathPointLevelUpdateInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointLevelUpdateInput{}, nil
	}
	out := make([]pathpoint.PathPointLevelUpdateInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid level element at index %d", i)
		}
		level, err := expandPathpointLevelUpdateInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = level
	}
	return out, nil
}

func expandPathpointLevelUpdateInput(cfg map[string]interface{}) (pathpoint.PathPointLevelUpdateInput, error) {
	level := pathpoint.PathPointLevelUpdateInput{}
	if v, ok := cfg["level_id"]; ok {
		level.ID = v.(string)
	}
	if v, ok := cfg["steps"]; ok {
		steps, err := expandPathpointStepUpdateInputList(v.([]interface{}))
		if err != nil {
			return level, err
		}
		level.Steps = steps
	}
	return level, nil
}

func expandPathpointStepInputList(list []interface{}) ([]pathpoint.PathPointStepInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointStepInput{}, nil
	}
	out := make([]pathpoint.PathPointStepInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid step element at index %d", i)
		}
		step, err := expandPathpointStepInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = step
	}
	return out, nil
}

func expandPathpointStepInput(cfg map[string]interface{}) (pathpoint.PathPointStepInput, error) {
	step := pathpoint.PathPointStepInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["is_excluded"]; ok {
		step.IsExcluded = v.(bool)
	}
	if v, ok := cfg["link"]; ok {
		step.Link = v.(string)
	}
	if v, ok := cfg["scoped_accounts"]; ok {
		rawList := v.([]interface{})
		accounts := make([]int, len(rawList))
		for i, a := range rawList {
			accounts[i] = a.(int)
		}
		step.ScopedAccounts = accounts
	}
	if v, ok := cfg["entity_search_query"]; ok {
		eqList := v.([]interface{})
		if len(eqList) > 0 {
			eq := expandPathpointSignalQueryInput(eqList[0].(map[string]interface{}))
			step.EntitySearchQuery = eq
		}
	}
	if v, ok := cfg["config"]; ok {
		cfgList := v.([]interface{})
		if len(cfgList) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfgList[0].(map[string]interface{}))
		}
	}
	if v, ok := cfg["signals"]; ok {
		signals, err := expandPathpointSignalInputList(v.([]interface{}))
		if err != nil {
			return step, err
		}
		step.Signals = signals
	}
	return step, nil
}

func expandPathpointStepUpdateInputList(list []interface{}) ([]pathpoint.PathPointStepUpdateInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointStepUpdateInput{}, nil
	}
	out := make([]pathpoint.PathPointStepUpdateInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid step element at index %d", i)
		}
		step, err := expandPathpointStepUpdateInput(cfg)
		if err != nil {
			return nil, err
		}
		out[i] = step
	}
	return out, nil
}

func expandPathpointStepUpdateInput(cfg map[string]interface{}) (pathpoint.PathPointStepUpdateInput, error) {
	step := pathpoint.PathPointStepUpdateInput{
		Name: cfg["name"].(string),
	}
	if v, ok := cfg["step_id"]; ok {
		step.ID = v.(string)
	}
	if v, ok := cfg["is_excluded"]; ok {
		step.IsExcluded = v.(bool)
	}
	if v, ok := cfg["link"]; ok {
		step.Link = v.(string)
	}
	if v, ok := cfg["scoped_accounts"]; ok {
		rawList := v.([]interface{})
		accounts := make([]int, len(rawList))
		for i, a := range rawList {
			accounts[i] = a.(int)
		}
		step.ScopedAccounts = accounts
	}
	if v, ok := cfg["entity_search_query"]; ok {
		eqList := v.([]interface{})
		if len(eqList) > 0 {
			eq := expandPathpointSignalQueryInput(eqList[0].(map[string]interface{}))
			step.EntitySearchQuery = eq
		}
	}
	if v, ok := cfg["config"]; ok {
		cfgList := v.([]interface{})
		if len(cfgList) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfgList[0].(map[string]interface{}))
		}
	}
	if v, ok := cfg["signals"]; ok {
		signals, err := expandPathpointSignalInputList(v.([]interface{}))
		if err != nil {
			return step, err
		}
		step.Signals = signals
	}
	return step, nil
}

func expandPathpointSignalQueryInput(cfg map[string]interface{}) *pathpoint.PathPointSignalQueryInput {
	eq := &pathpoint.PathPointSignalQueryInput{
		Query: cfg["query"].(string),
	}
	if v, ok := cfg["is_excluded"]; ok {
		eq.IsExcluded = v.(bool)
	}
	return eq
}

func expandPathpointStepStatusThresholdInput(cfg map[string]interface{}) pathpoint.PathPointStepStatusThresholdInput {
	threshold := pathpoint.PathPointStepStatusThresholdInput{}
	if v, ok := cfg["health_rollup"]; ok {
		threshold.HealthRollup = pathpoint.PathPointStepHealthRollup(v.(string))
	}
	if v, ok := cfg["threshold_type"]; ok {
		threshold.ThresholdType = pathpoint.PathPointThresholdType(v.(string))
	}
	if v, ok := cfg["threshold_value"]; ok {
		threshold.ThresholdValue = v.(int)
	}
	return threshold
}

func expandPathpointSignalInputList(list []interface{}) ([]pathpoint.PathPointSignalInput, error) {
	if len(list) == 0 {
		return []pathpoint.PathPointSignalInput{}, nil
	}
	out := make([]pathpoint.PathPointSignalInput, len(list))
	for i, v := range list {
		cfg, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid signal element at index %d", i)
		}
		out[i] = expandPathpointSignalInput(cfg)
	}
	return out, nil
}

func expandPathpointSignalInput(cfg map[string]interface{}) pathpoint.PathPointSignalInput {
	signal := pathpoint.PathPointSignalInput{
		GUID: pathpoint.EntityGUID(cfg["guid"].(string)),
	}
	if v, ok := cfg["name"]; ok {
		signal.Name = v.(string)
	}
	if v, ok := cfg["type"]; ok {
		signal.Type = pathpoint.PathPointSignalType(v.(string))
	}
	if v, ok := cfg["is_excluded"]; ok {
		signal.IsExcluded = v.(bool)
	}
	return signal
}

// EpochMilliseconds is a type alias helper used for version field conversion.
// pathpoint.EpochMilliseconds is actually nrtime.EpochMilliseconds — use it directly.
func init() {
	// ensure schema package is used
	_ = schema.TypeString
}
func flattenPathpointFlowResult(result *pathpoint.PathPointFlowResult, d *schema.ResourceData) error {
	if result == nil {
		return nil
	}

	_ = d.Set("name", result.Name)
	_ = d.Set("category", result.Category)
	_ = d.Set("description", result.Description)
	_ = d.Set("health_rollup", string(result.HealthRollup))
	_ = d.Set("refresh_interval", string(result.RefreshInterval))
	_ = d.Set("guid", string(result.GUID))
	_ = d.Set("version", result.Version.String())

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
	for i, kpi := range kpis {
		m := map[string]interface{}{
			"kpi_id":      kpi.ID,
			"name":        kpi.Name,
			"account_id":  kpi.AccountID,
			"category":    kpi.Category,
			"description": kpi.Description,
		}
		m["query"] = flattenPathpointKpiNRQL(kpi.Query)
		out[i] = m
	}
	return out
}

func flattenPathpointKpiNRQL(query pathpoint.PathPointKpiNRQL) []interface{} {
	m := map[string]interface{}{
		"from":  string(query.From),
		"where": query.Where,
	}

	m["select"] = flattenPathpointKpiNRQLSelect(query.Select)
	m["time_window"] = flattenPathpointKpiTimeWindow(query.TimeWindow)

	return []interface{}{m}
}

func flattenPathpointKpiNRQLSelect(sel pathpoint.PathPointKpiNRQLSelect) []interface{} {
	m := map[string]interface{}{
		"aggregation_type": string(sel.AggregationType),
		"alias":            sel.Alias,
		"attribute":        sel.Attribute,
		"threshold":        sel.Threshold,
	}
	return []interface{}{m}
}

func flattenPathpointKpiTimeWindow(tw pathpoint.PathPointKpiTimeWindow) []interface{} {
	if tw.CustomRange == "" && tw.RelativeRange == (pathpoint.PathPointKpiTimeWindowRelativeRange{}) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"custom_range": string(tw.CustomRange),
	}

	if tw.RelativeRange != (pathpoint.PathPointKpiTimeWindowRelativeRange{}) {
		rr := map[string]interface{}{
			"since":           string(tw.RelativeRange.Since),
			"compare_against": string(tw.RelativeRange.CompareAgainst),
		}
		m["relative_range"] = []interface{}{rr}
	} else {
		m["relative_range"] = []interface{}{}
	}

	return []interface{}{m}
}

func flattenPathpointStages(stages []pathpoint.PathPointStage) []interface{} {
	out := make([]interface{}, len(stages))
	for i, stage := range stages {
		m := map[string]interface{}{
			"stage_id":     stage.ID,
			"name":         stage.Name,
			"health_rollup": string(stage.HealthRollup),
			"is_excluded":  stage.IsExcluded,
			"link":         stage.Link,
		}

		m["related"] = flattenPathpointRelated(stage.Related)
		m["stage_kpis"] = flattenPathpointStageKpis(stage.StageKpis)
		m["levels"] = flattenPathpointLevels(stage.Levels.Items)

		out[i] = m
	}
	return out
}

func flattenPathpointRelated(related pathpoint.PathPointRelated) []interface{} {
	m := map[string]interface{}{
		"source": related.Source,
		"target": related.Target,
	}
	return []interface{}{m}
}

func flattenPathpointStageKpis(kpis []pathpoint.PathPointKpi) []interface{} {
	out := make([]interface{}, len(kpis))
	for i, kpi := range kpis {
		m := map[string]interface{}{
			"kpi_id":      kpi.ID,
			"name":        kpi.Name,
			"account_id":  kpi.AccountID,
			"category":    kpi.Category,
			"description": kpi.Description,
		}
		m["query"] = flattenPathpointKpiNRQL(kpi.Query)
		out[i] = m
	}
	return out
}

func flattenPathpointLevels(levels []pathpoint.PathPointLevel) []interface{} {
	out := make([]interface{}, len(levels))
	for i, level := range levels {
		m := map[string]interface{}{
			"level_id": level.ID,
		}
		m["steps"] = flattenPathpointSteps(level.Steps.Items)
		out[i] = m
	}
	return out
}

func flattenPathpointSteps(steps []pathpoint.PathPointStep) []interface{} {
	out := make([]interface{}, len(steps))
	for i, step := range steps {
		m := map[string]interface{}{
			"step_id":     step.ID,
			"name":        step.Name,
			"is_excluded": step.IsExcluded,
			"link":        step.Link,
		}

		scopedAccounts := make([]interface{}, len(step.ScopedAccounts))
		for j, v := range step.ScopedAccounts {
			scopedAccounts[j] = v
		}
		m["scoped_accounts"] = scopedAccounts

		m["entity_search_query"] = flattenPathpointEntitySearchQuery(step.EntitySearchQuery)
		m["config"] = flattenPathpointStepConfig(step.Config)
		m["signals"] = flattenPathpointSignals(step.Signals)

		out[i] = m
	}
	return out
}

func flattenPathpointEntitySearchQuery(q pathpoint.PathPointSignalQuery) []interface{} {
	if q.Query == "" && !q.IsExcluded {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"query":       q.Query,
		"is_excluded": q.IsExcluded,
	}
	return []interface{}{m}
}

func flattenPathpointStepConfig(config pathpoint.PathPointStepStatusThreshold) []interface{} {
	if config.HealthRollup == "" && config.ThresholdType == "" && config.ThresholdValue == 0 {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"health_rollup":   string(config.HealthRollup),
		"threshold_type":  string(config.ThresholdType),
		"threshold_value": config.ThresholdValue,
	}
	return []interface{}{m}
}

func flattenPathpointSignals(signals []pathpoint.PathPointSignal) []interface{} {
	out := make([]interface{}, len(signals))
	for i, sig := range signals {
		m := map[string]interface{}{
			"guid":        string(sig.GUID),
			"name":        sig.Name,
			"type":        string(sig.Type),
			"is_excluded": sig.IsExcluded,
		}
		out[i] = m
	}
	return out
}