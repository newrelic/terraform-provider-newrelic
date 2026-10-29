package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicelevel"
)

func flattenCompositeServiceLevelIndicator(composite servicelevel.ServiceLevelCompositeIndicator, d *schema.ResourceData) error {
	_ = d.Set("account_id", composite.AccountID)
	_ = d.Set("name", composite.Name)
	_ = d.Set("description", composite.Description)
	_ = d.Set("guid", composite.GUID)
	// component is a TypeSet, so Terraform diffs by membership/hash, not
	// position — no reordering needed to match server response against config.
	_ = d.Set("component", flattenCompositeComponents(composite.Components, composite.WeightType))
	_ = d.Set("objective", flattenCompositeObjectives(composite.Objectives))

	return nil
}

func flattenCompositeComponents(components []servicelevel.ServiceLevelComponent, weightType servicelevel.ServiceLevelWeightType) []interface{} {
	out := make([]interface{}, len(components))

	for i, component := range components {
		componentMap := map[string]interface{}{
			"guid": string(component.GUID),
		}
		// Weights are only ever user-specified when weightType is CUSTOM. When
		// it's AVERAGE, the server fills in an equal-distribution weight on
		// every component even though none was configured — flattening that
		// server-computed value back into state would create a permanent diff
		// against the (weight-less) config, so it's intentionally omitted here.
		if weightType == servicelevel.ServiceLevelWeightTypeTypes.CUSTOM {
			componentMap["weight"] = component.Weight
		}
		out[i] = componentMap
	}

	return out
}

func flattenCompositeObjectives(objectives []servicelevel.ServiceLevelCompositeObjective) []interface{} {
	out := make([]interface{}, len(objectives))

	for i, objective := range objectives {
		out[i] = flattenCompositeObjective(objective)
	}

	return out
}

func flattenCompositeObjective(objective servicelevel.ServiceLevelCompositeObjective) map[string]interface{} {
	return map[string]interface{}{
		"name":        objective.Name,
		"description": objective.Description,
		"target":      objective.Target,
		// Composite's objective reuses the exact same TimeWindow/RollingTimeWindow
		// types the plain service_level resource already uses — reusing the
		// existing flatten helpers directly, not duplicating them.
		"time_window": flattenServiceLevelTimeWindow(objective.TimeWindow),
	}
}

func expandCompositeServiceLevelCreateInput(d *schema.ResourceData) servicelevel.ServiceLevelCompositeCreateInput {
	createInput := servicelevel.ServiceLevelCompositeCreateInput{
		AccountID: d.Get("account_id").(int),
		Name:      d.Get("name").(string),
	}

	if descr, ok := d.GetOk("description"); ok {
		createInput.Description = descr.(string)
	}

	createInput.Components = expandCompositeComponentsCreateInput(d.Get("component").(*schema.Set).List())

	if objectives, ok := d.GetOk("objective"); ok {
		createInput.Objectives = expandCompositeObjectivesCreateInput(objectives.(*schema.Set).List())
	}

	return createInput
}

func expandCompositeComponentsCreateInput(cfg []interface{}) []servicelevel.ServiceLevelComponentCreateInput {
	components := make([]servicelevel.ServiceLevelComponentCreateInput, len(cfg))

	for i, rawCfg := range cfg {
		componentCfg := rawCfg.(map[string]interface{})
		components[i] = servicelevel.ServiceLevelComponentCreateInput{
			GUID: common.EntityGUID(componentCfg["guid"].(string)),
		}
		if weight, ok := componentCfg["weight"].(float64); ok {
			components[i].Weight = weight
		}
	}

	return components
}

func expandCompositeObjectivesCreateInput(cfg []interface{}) []servicelevel.ServiceLevelCompositeObjectiveCreateInput {
	if len(cfg) == 0 {
		return []servicelevel.ServiceLevelCompositeObjectiveCreateInput{}
	}

	objectives := make([]servicelevel.ServiceLevelCompositeObjectiveCreateInput, len(cfg))

	for i, rawCfg := range cfg {
		objectives[i] = expandCompositeObjectiveCreateInput(rawCfg.(map[string]interface{}))
	}

	return objectives
}

func expandCompositeObjectiveCreateInput(cfg map[string]interface{}) servicelevel.ServiceLevelCompositeObjectiveCreateInput {
	objective := servicelevel.ServiceLevelCompositeObjectiveCreateInput{}

	if name, ok := cfg["name"]; ok {
		objective.Name = name.(string)
	}

	if descr, ok := cfg["description"]; ok {
		objective.Description = descr.(string)
	}

	objective.Target = cfg["target"].(float64)

	// Reuses the existing helper directly — composite's objective.timeWindow is
	// the exact same shared type the plain service_level resource uses.
	objective.TimeWindow = expandServiceLevelObjectiveTimeWindowCreateInput(cfg["time_window"].([]interface{})[0].(map[string]interface{}))

	return objective
}

func expandCompositeServiceLevelUpdateInput(d *schema.ResourceData) servicelevel.ServiceLevelCompositeUpdateInput {
	updateInput := servicelevel.ServiceLevelCompositeUpdateInput{}

	if name, ok := d.GetOk("name"); ok {
		updateInput.Name = name.(string)
	}

	if descr, ok := d.GetOk("description"); ok {
		updateInput.Description = descr.(string)
	}

	updateInput.Components = expandCompositeComponentsUpdateInput(d.Get("component").(*schema.Set).List())

	if objectives, ok := d.GetOk("objective"); ok {
		updateInput.Objectives = expandCompositeObjectivesUpdateInput(objectives.(*schema.Set).List())
	}

	return updateInput
}

func expandCompositeComponentsUpdateInput(cfg []interface{}) []servicelevel.ServiceLevelComponentUpdateInput {
	components := make([]servicelevel.ServiceLevelComponentUpdateInput, len(cfg))

	for i, rawCfg := range cfg {
		componentCfg := rawCfg.(map[string]interface{})
		components[i] = servicelevel.ServiceLevelComponentUpdateInput{
			GUID: common.EntityGUID(componentCfg["guid"].(string)),
		}
		if weight, ok := componentCfg["weight"].(float64); ok {
			components[i].Weight = weight
		}
	}

	return components
}

func expandCompositeObjectivesUpdateInput(cfg []interface{}) []servicelevel.ServiceLevelCompositeObjectiveUpdateInput {
	if len(cfg) == 0 {
		return []servicelevel.ServiceLevelCompositeObjectiveUpdateInput{}
	}

	objectives := make([]servicelevel.ServiceLevelCompositeObjectiveUpdateInput, len(cfg))

	for i, rawCfg := range cfg {
		objectives[i] = expandCompositeObjectiveUpdateInput(rawCfg.(map[string]interface{}))
	}

	return objectives
}

func expandCompositeObjectiveUpdateInput(cfg map[string]interface{}) servicelevel.ServiceLevelCompositeObjectiveUpdateInput {
	objective := servicelevel.ServiceLevelCompositeObjectiveUpdateInput{}

	if name, ok := cfg["name"]; ok {
		objective.Name = name.(string)
	}

	if descr, ok := cfg["description"]; ok {
		objective.Description = descr.(string)
	}

	objective.Target = cfg["target"].(float64)

	// Reuses the existing helper directly — same shared TimeWindow type, and it
	// generates as a value (not a pointer) on the composite update struct too,
	// so no extra adaptation is needed here.
	objective.TimeWindow = expandServiceLevelObjectiveTimeWindowUpdateInput(cfg["time_window"].([]interface{})[0].(map[string]interface{}))

	return objective
}
