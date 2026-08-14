package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

// ── expand helpers (Terraform state → API input) ──────────────────────────────

func expandPathpointFlowInput(d *schema.ResourceData, accountID int) pathpoint.PathPointFlowInput {
	input := pathpoint.PathPointFlowInput{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v, ok := d.GetOk("health_rollup"); ok {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v.(string))
	}
	if v, ok := d.GetOk("refresh_interval"); ok {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v.(string))
	}
	if v, ok := d.GetOk("kpis"); ok {
		input.Kpis = expandPathpointKpiInputs(v.([]interface{}), accountID)
	}
	if v, ok := d.GetOk("stages"); ok {
		input.Stages = expandPathpointStageInputs(v.([]interface{}), accountID)
	}

	return input
}

func expandPathpointFlowUpdateInput(d *schema.ResourceData, version nrtime.EpochMilliseconds, accountID int) pathpoint.PathPointFlowUpdateInput {
	input := pathpoint.PathPointFlowUpdateInput{
		Name:    d.Get("name").(string),
		Version: version,
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v, ok := d.GetOk("health_rollup"); ok {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v.(string))
	}
	if v, ok := d.GetOk("refresh_interval"); ok {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v.(string))
	}
	if v, ok := d.GetOk("kpis"); ok {
		oldRaw, _ := d.GetChange("kpis")
		input.Kpis = expandPathpointKpiUpdateInputsResolved(v.([]interface{}), oldRaw.([]interface{}), accountID)
	}
	if v, ok := d.GetOk("stages"); ok {
		oldRaw, _ := d.GetChange("stages")
		input.Stages = expandPathpointStageUpdateInputsResolved(v.([]interface{}), oldRaw.([]interface{}), accountID)
	}

	return input
}

func expandPathpointKpiInputs(raw []interface{}, accountID int) []pathpoint.PathPointKpiInput {
	inputs := make([]pathpoint.PathPointKpiInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		kpi := pathpoint.PathPointKpiInput{
			Name: m["name"].(string),
		}
		if v, ok := m["description"].(string); ok && v != "" {
			kpi.Description = v
		}
		if v, ok := m["category"].(string); ok && v != "" {
			kpi.Category = v
		}
		if v, ok := m["account_id"].(int); ok && v != 0 {
			kpi.AccountID = v
		} else {
			kpi.AccountID = accountID
		}
		if q, ok := m["query"].([]interface{}); ok && len(q) > 0 {
			kpi.Query = expandPathpointKpiNRQLInput(q[0].(map[string]interface{}))
		}
		inputs = append(inputs, kpi)
	}
	return inputs
}

func expandPathpointKpiUpdateInputsResolved(newRaw, oldRaw []interface{}, accountID int) []pathpoint.PathPointKpiUpdateInput {
	type oldInfo struct {
		id string
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		info := &oldInfo{id: m["id"].(string)}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	inputs := make([]pathpoint.PathPointKpiUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		kpi := pathpoint.PathPointKpiUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			if info.id != "" {
				kpi.ID = info.id
			}
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			if oldOrdered[i].id != "" {
				kpi.ID = oldOrdered[i].id
			}
			claimedByPos[i] = true
		}

		if v, ok := m["description"].(string); ok && v != "" {
			kpi.Description = v
		}
		if v, ok := m["category"].(string); ok && v != "" {
			kpi.Category = v
		}
		if v, ok := m["account_id"].(int); ok && v != 0 {
			kpi.AccountID = v
		} else {
			kpi.AccountID = accountID
		}
		if q, ok := m["query"].([]interface{}); ok && len(q) > 0 {
			kpi.Query = expandPathpointKpiNRQLInput(q[0].(map[string]interface{}))
		}
		inputs = append(inputs, kpi)
	}
	return inputs
}

func expandPathpointKpiNRQLInput(m map[string]interface{}) pathpoint.PathPointKpiNRQLInput {
	q := pathpoint.PathPointKpiNRQLInput{
		From: pathpoint.NRQL(m["from"].(string)),
	}
	if v, ok := m["where"].(string); ok && v != "" {
		q.Where = v
	}
	if sel, ok := m["select"].([]interface{}); ok && len(sel) > 0 {
		q.Select = expandPathpointKpiNRQLSelectInput(sel[0].(map[string]interface{}))
	}
	if tw, ok := m["time_window"].([]interface{}); ok && len(tw) > 0 {
		q.TimeWindow = expandPathpointKpiTimeWindowInput(tw[0].(map[string]interface{}))
	}
	return q
}

func expandPathpointKpiNRQLSelectInput(m map[string]interface{}) pathpoint.PathPointKpiNRQLSelectInput {
	s := pathpoint.PathPointKpiNRQLSelectInput{
		AggregationType: pathpoint.PathPointKpiNRQLAggregations(m["aggregation_type"].(string)),
	}
	if v, ok := m["alias"].(string); ok && v != "" {
		s.Alias = v
	}
	if v, ok := m["attribute"].(string); ok && v != "" {
		s.Attribute = v
	}
	if v, ok := m["threshold"].(float64); ok && v != 0 {
		s.Threshold = v
	}
	return s
}

func expandPathpointKpiTimeWindowInput(m map[string]interface{}) pathpoint.PathPointKpiTimeWindowInput {
	tw := pathpoint.PathPointKpiTimeWindowInput{}
	if v, ok := m["custom_range"].(string); ok && v != "" {
		tw.CustomRange = pathpoint.NRQL(v)
	}
	if rr, ok := m["relative_range"].([]interface{}); ok && len(rr) > 0 {
		rrm := rr[0].(map[string]interface{})
		if since, ok := rrm["since"].(string); ok && since != "" {
			rel := &pathpoint.PathPointKpiTimeWindowRelativeRangeInput{
				Since: pathpoint.PathPointKpiTimeDuration(since),
			}
			if v, ok := rrm["compare_against"].(string); ok && v != "" {
				rel.CompareAgainst = pathpoint.PathPointKpiTimeDuration(v)
			}
			tw.RelativeRange = *rel
		}
	}
	return tw
}

func expandPathpointStageInputs(raw []interface{}, accountID int) []pathpoint.PathPointStageInput {
	stages := make([]pathpoint.PathPointStageInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		stage := pathpoint.PathPointStageInput{
			Name: m["name"].(string),
		}
		if v, ok := m["health_rollup"].(string); ok && v != "" {
			stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			stage.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			stage.Link = v
		}
		if rel, ok := m["related"].([]interface{}); ok && len(rel) > 0 {
			stage.Related = expandPathpointRelatedInput(rel[0].(map[string]interface{}))
		}
		if kpis, ok := m["stage_kpis"].([]interface{}); ok {
			stage.StageKpis = expandPathpointKpiInputs(kpis, accountID)
		}
		if levels, ok := m["levels"].([]interface{}); ok {
			stage.Levels = expandPathpointLevelInputs(levels)
		}
		stages = append(stages, stage)
	}
	return stages
}

func expandPathpointStageUpdateInputsResolved(newRaw, oldRaw []interface{}, accountID int) []pathpoint.PathPointStageUpdateInput {
	type oldInfo struct {
		id        string
		levels    []interface{}
		stageKpis []interface{}
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		levels, _ := m["levels"].([]interface{})
		stageKpis, _ := m["stage_kpis"].([]interface{})
		info := &oldInfo{id: m["id"].(string), levels: levels, stageKpis: stageKpis}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	stages := make([]pathpoint.PathPointStageUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		stage := pathpoint.PathPointStageUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		var matched *oldInfo
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			matched = info
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			matched = oldOrdered[i]
			claimedByPos[i] = true
		}

		var oldLevels, oldStageKpis []interface{}
		if matched != nil {
			if matched.id != "" {
				stage.ID = matched.id
			}
			oldLevels = matched.levels
			oldStageKpis = matched.stageKpis
		}

		if v, ok := m["health_rollup"].(string); ok && v != "" {
			stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			stage.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			stage.Link = v
		}
		if rel, ok := m["related"].([]interface{}); ok && len(rel) > 0 {
			stage.Related = expandPathpointRelatedInput(rel[0].(map[string]interface{}))
		}
		if kpis, ok := m["stage_kpis"].([]interface{}); ok {
			stage.StageKpis = expandPathpointKpiUpdateInputsResolved(kpis, oldStageKpis, accountID)
		}
		if levels, ok := m["levels"].([]interface{}); ok {
			stage.Levels = expandPathpointLevelUpdateInputsResolved(levels, oldLevels)
		}
		stages = append(stages, stage)
	}
	return stages
}

func expandPathpointLevelUpdateInputsResolved(newRaw, oldRaw []interface{}) []pathpoint.PathPointLevelUpdateInput {
	type oldInfo struct {
		id        string
		steps     []interface{}
		stepNames map[string]bool
	}
	oldLevels := make([]oldInfo, 0, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		steps, _ := m["steps"].([]interface{})
		names := make(map[string]bool, len(steps))
		for _, s := range steps {
			names[s.(map[string]interface{})["name"].(string)] = true
		}
		oldLevels = append(oldLevels, oldInfo{id: m["id"].(string), steps: steps, stepNames: names})
	}
	claimed := make([]bool, len(oldLevels))

	levels := make([]pathpoint.PathPointLevelUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		level := pathpoint.PathPointLevelUpdateInput{}
		newSteps, _ := m["steps"].([]interface{})

		// Match by step-name overlap first; fall back to position if no overlap found.
		bestIdx, bestOverlap := -1, 0
		for j := range oldLevels {
			if claimed[j] {
				continue
			}
			overlap := 0
			for _, s := range newSteps {
				if oldLevels[j].stepNames[s.(map[string]interface{})["name"].(string)] {
					overlap++
				}
			}
			if overlap > bestOverlap {
				bestOverlap = overlap
				bestIdx = j
			}
		}
		if bestIdx < 0 && i < len(oldLevels) && !claimed[i] {
			bestIdx = i
		}

		var oldSteps []interface{}
		if bestIdx >= 0 {
			claimed[bestIdx] = true
			if oldLevels[bestIdx].id != "" {
				level.ID = oldLevels[bestIdx].id
			}
			oldSteps = oldLevels[bestIdx].steps
		}

		level.Steps = expandPathpointStepUpdateInputsResolved(newSteps, oldSteps)
		levels = append(levels, level)
	}
	return levels
}

func expandPathpointStepUpdateInputsResolved(newRaw, oldRaw []interface{}) []pathpoint.PathPointStepUpdateInput {
	type oldInfo struct {
		id  string
		pos int
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for i, r := range oldRaw {
		m := r.(map[string]interface{})
		info := &oldInfo{id: m["id"].(string), pos: i}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	steps := make([]pathpoint.PathPointStepUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		step := pathpoint.PathPointStepUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			if info.id != "" {
				step.ID = info.id
			}
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			if oldOrdered[i].id != "" {
				step.ID = oldOrdered[i].id
			}
			claimedByPos[i] = true
		}
		if v, ok := m["is_excluded"].(bool); ok {
			step.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			step.Link = v
		}
		if v, ok := m["scoped_accounts"].([]interface{}); ok {
			for _, id := range v {
				step.ScopedAccounts = append(step.ScopedAccounts, id.(int))
			}
		}
		if esq, ok := m["entity_search_query"].([]interface{}); ok && len(esq) > 0 {
			step.EntitySearchQuery = expandPathpointSignalQueryInput(esq[0].(map[string]interface{}))
		}
		if cfg, ok := m["config"].([]interface{}); ok && len(cfg) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfg[0].(map[string]interface{}))
		}
		if sigs, ok := m["signals"].([]interface{}); ok {
			step.Signals = expandPathpointSignalInputs(sigs)
		}
		steps = append(steps, step)
	}
	return steps
}

func expandPathpointLevelInputs(raw []interface{}) []pathpoint.PathPointLevelInput {
	levels := make([]pathpoint.PathPointLevelInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		level := pathpoint.PathPointLevelInput{}
		if steps, ok := m["steps"].([]interface{}); ok {
			level.Steps = expandPathpointStepInputs(steps)
		}
		levels = append(levels, level)
	}
	return levels
}

func expandPathpointStepInputs(raw []interface{}) []pathpoint.PathPointStepInput {
	steps := make([]pathpoint.PathPointStepInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		step := pathpoint.PathPointStepInput{
			Name: m["name"].(string),
		}
		if v, ok := m["is_excluded"].(bool); ok {
			step.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			step.Link = v
		}
		if v, ok := m["scoped_accounts"].([]interface{}); ok {
			for _, id := range v {
				step.ScopedAccounts = append(step.ScopedAccounts, id.(int))
			}
		}
		if esq, ok := m["entity_search_query"].([]interface{}); ok && len(esq) > 0 {
			step.EntitySearchQuery = expandPathpointSignalQueryInput(esq[0].(map[string]interface{}))
		}
		if cfg, ok := m["config"].([]interface{}); ok && len(cfg) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfg[0].(map[string]interface{}))
		}
		if sigs, ok := m["signals"].([]interface{}); ok {
			step.Signals = expandPathpointSignalInputs(sigs)
		}
		steps = append(steps, step)
	}
	return steps
}

func expandPathpointSignalQueryInput(m map[string]interface{}) pathpoint.PathPointSignalQueryInput {
	q := pathpoint.PathPointSignalQueryInput{
		Query: m["query"].(string),
	}
	if v, ok := m["is_excluded"].(bool); ok {
		q.IsExcluded = v
	}
	return q
}

func expandPathpointStepStatusThresholdInput(m map[string]interface{}) pathpoint.PathPointStepStatusThresholdInput {
	cfg := pathpoint.PathPointStepStatusThresholdInput{}
	if v, ok := m["health_rollup"].(string); ok && v != "" {
		cfg.HealthRollup = pathpoint.PathPointStepHealthRollup(v)
	}
	if v, ok := m["threshold_type"].(string); ok && v != "" {
		cfg.ThresholdType = pathpoint.PathPointThresholdType(v)
	}
	if v, ok := m["threshold_value"].(int); ok {
		cfg.ThresholdValue = v
	}
	return cfg
}

func expandPathpointSignalInputs(raw []interface{}) []pathpoint.PathPointSignalInput {
	signals := make([]pathpoint.PathPointSignalInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		sig := pathpoint.PathPointSignalInput{
			GUID: pathpoint.EntityGUID(m["guid"].(string)),
		}
		if v, ok := m["name"].(string); ok && v != "" {
			sig.Name = v
		}
		if v, ok := m["type"].(string); ok && v != "" {
			sig.Type = pathpoint.PathPointSignalType(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			sig.IsExcluded = v
		}
		signals = append(signals, sig)
	}
	return signals
}

func expandPathpointRelatedInput(m map[string]interface{}) pathpoint.PathPointRelatedInput {
	return pathpoint.PathPointRelatedInput{
		Source: m["source"].(bool),
		Target: m["target"].(bool),
	}
}

// ── flatten helpers (API response → Terraform state) ─────────────────────────

func flattenPathpointKpis(kpis []pathpoint.PathPointKpi) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(kpis))
	for _, k := range kpis {
		m := map[string]interface{}{
			"id":          k.ID,
			"name":        k.Name,
			"description": k.Description,
			"category":    k.Category,
			"account_id":  k.AccountID,
			"query":       flattenPathpointKpiNRQL(k.Query),
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointKpiNRQL(q pathpoint.PathPointKpiNRQL) []map[string]interface{} {
	m := map[string]interface{}{
		"from":  string(q.From),
		"where": q.Where,
		"select": []map[string]interface{}{
			{
				"aggregation_type": string(q.Select.AggregationType),
				"alias":            q.Select.Alias,
				"attribute":        q.Select.Attribute,
				"threshold":        q.Select.Threshold,
			},
		},
	}
	if tw := flattenPathpointKpiTimeWindow(q.TimeWindow); len(tw) > 0 {
		m["time_window"] = tw
	}
	return []map[string]interface{}{m}
}

func flattenPathpointKpiTimeWindow(tw pathpoint.PathPointKpiTimeWindow) []map[string]interface{} {
	if tw.CustomRange == "" && tw.RelativeRange.Since == "" {
		return nil
	}
	m := map[string]interface{}{
		"custom_range": string(tw.CustomRange),
	}
	if tw.RelativeRange.Since != "" {
		rr := map[string]interface{}{
			"since":           string(tw.RelativeRange.Since),
			"compare_against": string(tw.RelativeRange.CompareAgainst),
		}
		m["relative_range"] = []map[string]interface{}{rr}
	}
	return []map[string]interface{}{m}
}

func flattenPathpointStages(stages []pathpoint.PathPointStage) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(stages))
	for _, s := range stages {
		m := map[string]interface{}{
			"id":            s.ID,
			"name":          s.Name,
			"health_rollup": string(s.HealthRollup),
			"is_excluded":   s.IsExcluded,
			"link":          s.Link,
			"stage_kpis":    flattenPathpointKpis(s.StageKpis),
			"levels":        flattenPathpointLevels(s.Levels.Items),
		}
		if s.Related.Source || s.Related.Target {
			m["related"] = []map[string]interface{}{
				{
					"source": s.Related.Source,
					"target": s.Related.Target,
				},
			}
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointLevels(levels []pathpoint.PathPointLevel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(levels))
	for _, l := range levels {
		m := map[string]interface{}{
			"id":    l.ID,
			"steps": flattenPathpointSteps(l.Steps.Items),
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointSteps(steps []pathpoint.PathPointStep) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(steps))
	for _, s := range steps {
		scopedAccounts := make([]interface{}, 0, len(s.ScopedAccounts))
		for _, id := range s.ScopedAccounts {
			scopedAccounts = append(scopedAccounts, id)
		}
		m := map[string]interface{}{
			"id":              s.ID,
			"name":            s.Name,
			"is_excluded":     s.IsExcluded,
			"link":            s.Link,
			"scoped_accounts": scopedAccounts,
			"signals":         flattenPathpointSignals(s.Signals),
		}
		if s.EntitySearchQuery.Query != "" {
			m["entity_search_query"] = []map[string]interface{}{
				{
					"query":       s.EntitySearchQuery.Query,
					"is_excluded": s.EntitySearchQuery.IsExcluded,
				},
			}
		}
		if s.Config.HealthRollup != "" || s.Config.ThresholdType != "" || s.Config.ThresholdValue != 0 {
			m["config"] = []map[string]interface{}{
				{
					"health_rollup":   string(s.Config.HealthRollup),
					"threshold_type":  string(s.Config.ThresholdType),
					"threshold_value": s.Config.ThresholdValue,
				},
			}
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointSignals(signals []pathpoint.PathPointSignal) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(signals))
	for _, s := range signals {
		result = append(result, map[string]interface{}{
			"guid":        string(s.GUID),
			"name":        s.Name,
			"type":        string(s.Type),
			"is_excluded": s.IsExcluded,
		})
	}
	return result
}
