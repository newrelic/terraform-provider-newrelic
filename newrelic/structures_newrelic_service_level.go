package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicelevel"
)

func flattenServiceLevelIndicator(indicator servicelevel.ServiceLevelIndicator, identifier *serviceLevelIdentifier, d *schema.ResourceData, sliGUID string) error {
	_ = d.Set("guid", indicator.EntityGUID)
	_ = d.Set("sli_id", indicator.ID)
	_ = d.Set("name", indicator.Name)
	_ = d.Set("description", indicator.Description)
	_ = d.Set("sli_guid", sliGUID)

	eventsMap := make(map[string]interface{})
	events := make([]interface{}, 1)

	eventsMap["account_id"] = identifier.AccountID

	if indicator.Events.ValidEvents != nil {
		eventsMap["valid_events"] = flattenServiceLevelEventsQuery(indicator.Events.ValidEvents, d, "valid_events")
	}

	if indicator.Events.GoodEvents != nil {
		eventsMap["good_events"] = flattenServiceLevelEventsQuery(indicator.Events.GoodEvents, d, "good_events")
	}

	if indicator.Events.BadEvents != nil {
		eventsMap["bad_events"] = flattenServiceLevelEventsQuery(indicator.Events.BadEvents, d, "bad_events")
	}

	events[0] = eventsMap
	_ = d.Set("events", events)

	objectives := flattenServiceLevelObjectives(indicator.Objectives)
	_ = d.Set("objective", objectives)

	return nil
}

func flattenServiceLevelEventsQuery(eventsQuery *servicelevel.ServiceLevelEventsQuery, d *schema.ResourceData, eventType string) []interface{} {
	eventsQueryMap := make(map[string]interface{})
	eventsQueryOutput := make([]interface{}, 1)

	eventsQueryMap["from"] = eventsQuery.From
	eventsQueryMap["where"] = eventsQuery.Where

	apiReturnedDefaultValue := len(eventsQuery.Select.Attribute) == 0 && eventsQuery.Select.Function == "COUNT"

	if value, ok := d.GetOk(fmt.Sprintf("events.0.%s.0.select", eventType)); ok && len(value.([]interface{})) > 0 || !apiReturnedDefaultValue {
		eventsQueryMap["select"] = flattenServiceLevelEventsQuerySelect(eventsQuery.Select)
	}

	eventsQueryOutput[0] = eventsQueryMap
	return eventsQueryOutput
}

func flattenServiceLevelEventsQuerySelect(selectValue servicelevel.ServiceLevelEventsQuerySelect) []interface{} {
	selectQueryMap := make(map[string]interface{})
	selectOutput := make([]interface{}, 1)

	selectQueryMap["attribute"] = selectValue.Attribute
	selectQueryMap["function"] = selectValue.Function

	selectOutput[0] = selectQueryMap
	return selectOutput
}

func flattenServiceLevelObjectives(objectives []servicelevel.ServiceLevelObjective) []interface{} {
	objectivesOutput := make([]interface{}, len(objectives))

	for i, objective := range objectives {
		objectivesOutput[i] = flattenServiceLevelObjective(objective)
	}

	return objectivesOutput
}

func flattenServiceLevelObjective(objective servicelevel.ServiceLevelObjective) map[string]interface{} {
	objectiveMap := make(map[string]interface{})

	objectiveMap["name"] = objective.Name
	objectiveMap["description"] = objective.Description
	objectiveMap["target"] = objective.Target
	objectiveMap["time_window"] = flattenServiceLevelTimeWindow(objective.TimeWindow)

	return objectiveMap
}

func flattenServiceLevelTimeWindow(timeWindow servicelevel.ServiceLevelObjectiveTimeWindow) []interface{} {
	timeWindowMap := make(map[string]interface{})
	timeWindowOutput := make([]interface{}, 1)

	timeWindowMap["rolling"] = flattenServiceLevelRollingTimeWindow(timeWindow.Rolling)

	timeWindowOutput[0] = timeWindowMap
	return timeWindowOutput
}

func flattenServiceLevelRollingTimeWindow(rolling servicelevel.ServiceLevelObjectiveRollingTimeWindow) []interface{} {
	rollingMap := make(map[string]interface{})
	rollingOutput := make([]interface{}, 1)

	rollingMap["count"] = rolling.Count
	rollingMap["unit"] = rolling.Unit

	rollingOutput[0] = rollingMap
	return rollingOutput
}

func expandServiceLevelCreateInput(d *schema.ResourceData) servicelevel.ServiceLevelIndicatorCreateInput {
	createInput := servicelevel.ServiceLevelIndicatorCreateInput{
		Name: d.Get("name").(string),
	}

	if descr, ok := d.GetOk("description"); ok {
		createInput.Description = descr.(string)
	}

	createInput.Events = expandServiceLevelEventsCreateInput(d.Get("events").([]interface{})[0].(map[string]interface{}))

	if objectives, ok := d.GetOk("objective"); ok {
		createInput.Objectives = expandServiceLevelObjectivesCreateInput(objectives.(*schema.Set).List())
	}

	return createInput
}

func expandServiceLevelEventsCreateInput(cfg map[string]interface{}) servicelevel.ServiceLevelEventsCreateInput {

	events := servicelevel.ServiceLevelEventsCreateInput{}

	events.AccountID = cfg["account_id"].(int)

	validEvents := cfg["valid_events"].([]interface{})[0].(map[string]interface{})
	events.ValidEvents = expandServiceLevelEventsQueryCreateInput(validEvents)

	if value, ok := cfg["good_events"].([]interface{}); ok && len(value) > 0 {
		goodEvents := value[0].(map[string]interface{})
		events.GoodEvents = expandServiceLevelEventsQueryCreateInput(goodEvents)
	}

	if value, ok := cfg["bad_events"].([]interface{}); ok && len(value) > 0 {
		badEvents := value[0].(map[string]interface{})
		events.BadEvents = expandServiceLevelEventsQueryCreateInput(badEvents)
	}

	return events
}

func expandServiceLevelEventsQueryCreateInput(cfg map[string]interface{}) *servicelevel.ServiceLevelEventsQueryCreateInput {
	eventsQuery := servicelevel.ServiceLevelEventsQueryCreateInput{}

	eventsQuery.From = servicelevel.NRQL(cfg["from"].(string))

	if where, ok := cfg["where"]; ok {
		eventsQuery.Where = servicelevel.NRQL(where.(string))
	}

	if value, ok := cfg["select"].([]interface{}); ok && len(value) > 0 {
		eventsQuery.Select = expandServiceLevelEventsQuerySelectCreateInput(value[0].(map[string]interface{}))
	}

	return &eventsQuery
}

func expandServiceLevelEventsQuerySelectCreateInput(cfg map[string]interface{}) *servicelevel.ServiceLevelEventsQuerySelectCreateInput {
	selectValue := servicelevel.ServiceLevelEventsQuerySelectCreateInput{}

	if attribute, ok := cfg["attribute"]; ok {
		selectValue.Attribute = attribute.(string)
	}

	selectValue.Function = servicelevel.ServiceLevelEventsQuerySelectFunction(cfg["function"].(string))

	return &selectValue
}

func expandServiceLevelObjectivesCreateInput(cfg []interface{}) []servicelevel.ServiceLevelObjectiveCreateInput {
	if len(cfg) == 0 {
		return []servicelevel.ServiceLevelObjectiveCreateInput{}
	}

	perms := make([]servicelevel.ServiceLevelObjectiveCreateInput, len(cfg))

	for i, rawCfg := range cfg {
		objective := expandServiceLevelObjectiveCreateInput(rawCfg.(map[string]interface{}))
		perms[i] = objective
	}

	return perms
}

func expandServiceLevelObjectiveCreateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveCreateInput {
	objective := servicelevel.ServiceLevelObjectiveCreateInput{}

	if name, ok := cfg["name"]; ok {
		objective.Name = name.(string)
	}

	if descr, ok := cfg["description"]; ok {
		objective.Description = descr.(string)
	}

	objective.Target = cfg["target"].(float64)

	objective.TimeWindow = expandServiceLevelObjectiveTimeWindowCreateInput(cfg["time_window"].([]interface{})[0].(map[string]interface{}))

	return objective
}

func expandServiceLevelObjectiveTimeWindowCreateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveTimeWindowCreateInput {
	timeWindow := servicelevel.ServiceLevelObjectiveTimeWindowCreateInput{}

	timeWindow.Rolling = expandServiceLevelObjectiveRollingTimeWindowCreateInput(cfg["rolling"].([]interface{})[0].(map[string]interface{}))

	return timeWindow
}

func expandServiceLevelObjectiveRollingTimeWindowCreateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveRollingTimeWindowCreateInput {
	rolling := servicelevel.ServiceLevelObjectiveRollingTimeWindowCreateInput{}

	rolling.Count = cfg["count"].(int)
	rolling.Unit = servicelevel.ServiceLevelObjectiveRollingTimeWindowUnit(cfg["unit"].(string))

	return rolling
}

func expandServiceLevelUpdateInput(d *schema.ResourceData) servicelevel.ServiceLevelIndicatorUpdateInput {
	updateInput := servicelevel.ServiceLevelIndicatorUpdateInput{}

	if name, ok := d.GetOk("name"); ok {
		updateInput.Name = name.(string)
	}

	if descr, ok := d.GetOk("description"); ok {
		updateInput.Description = descr.(string)
	}

	if events, ok := d.GetOk("events"); ok {
		updateInput.Events = expandServiceLevelEventsUpdateInput(events.([]interface{})[0].(map[string]interface{}))
	}

	if objectives, ok := d.GetOk("objective"); ok {
		updateInput.Objectives = expandServiceLevelObjectivesUpdateInput(objectives.(*schema.Set).List())
	}

	return updateInput
}

func expandServiceLevelEventsUpdateInput(cfg map[string]interface{}) *servicelevel.ServiceLevelEventsUpdateInput {
	events := servicelevel.ServiceLevelEventsUpdateInput{}

	if value, ok := cfg["valid_events"].([]interface{}); ok && len(value) > 0 {
		validEvents := value[0].(map[string]interface{})
		events.ValidEvents = expandServiceLevelEventsQueryUpdateInput(validEvents)
	}

	if value, ok := cfg["good_events"].([]interface{}); ok && len(value) > 0 {
		goodEvents := value[0].(map[string]interface{})
		events.GoodEvents = expandServiceLevelEventsQueryUpdateInput(goodEvents)
	}

	if value, ok := cfg["bad_events"].([]interface{}); ok && len(value) > 0 {
		badEvents := value[0].(map[string]interface{})
		events.BadEvents = expandServiceLevelEventsQueryUpdateInput(badEvents)
	}

	return &events
}

func expandServiceLevelEventsQueryUpdateInput(cfg map[string]interface{}) *servicelevel.ServiceLevelEventsQueryUpdateInput {
	eventsQuery := servicelevel.ServiceLevelEventsQueryUpdateInput{}

	eventsQuery.From = servicelevel.NRQL(cfg["from"].(string))

	if where, ok := cfg["where"]; ok {
		eventsQuery.Where = servicelevel.NRQL(where.(string))
	}

	if value, ok := cfg["select"].([]interface{}); ok && len(value) > 0 {
		eventsQuery.Select = expandServiceLevelEventsQuerySelectUpdateInput(value[0].(map[string]interface{}))
	}

	return &eventsQuery
}

func expandServiceLevelEventsQuerySelectUpdateInput(cfg map[string]interface{}) *servicelevel.ServiceLevelEventsQuerySelectUpdateInput {
	selectValue := servicelevel.ServiceLevelEventsQuerySelectUpdateInput{}

	if attribute, ok := cfg["attribute"]; ok {
		selectValue.Attribute = attribute.(string)
	}

	selectValue.Function = servicelevel.ServiceLevelEventsQuerySelectFunction(cfg["function"].(string))

	return &selectValue
}

func expandServiceLevelObjectivesUpdateInput(cfg []interface{}) []servicelevel.ServiceLevelObjectiveUpdateInput {
	if len(cfg) == 0 {
		return []servicelevel.ServiceLevelObjectiveUpdateInput{}
	}

	objectives := make([]servicelevel.ServiceLevelObjectiveUpdateInput, len(cfg))

	for i, rawCfg := range cfg {
		objective := expandServiceLevelObjectiveUpdateInput(rawCfg.(map[string]interface{}))
		objectives[i] = objective
	}

	return objectives
}

func expandServiceLevelObjectiveUpdateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveUpdateInput {
	objective := servicelevel.ServiceLevelObjectiveUpdateInput{}

	if name, ok := cfg["name"]; ok {
		objective.Name = name.(string)
	}

	if descr, ok := cfg["description"]; ok {
		objective.Description = descr.(string)
	}

	objective.Target = cfg["target"].(float64)

	objective.TimeWindow = expandServiceLevelObjectiveTimeWindowUpdateInput(cfg["time_window"].([]interface{})[0].(map[string]interface{}))

	return objective
}

func expandServiceLevelObjectiveTimeWindowUpdateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveTimeWindowUpdateInput {
	timeWindow := servicelevel.ServiceLevelObjectiveTimeWindowUpdateInput{}

	timeWindow.Rolling = expandServiceLevelObjectiveRollingTimeWindowUpdateInput(cfg["rolling"].([]interface{})[0].(map[string]interface{}))

	return timeWindow
}

func expandServiceLevelObjectiveRollingTimeWindowUpdateInput(cfg map[string]interface{}) servicelevel.ServiceLevelObjectiveRollingTimeWindowUpdateInput {
	rolling := servicelevel.ServiceLevelObjectiveRollingTimeWindowUpdateInput{}

	rolling.Count = cfg["count"].(int)
	rolling.Unit = servicelevel.ServiceLevelObjectiveRollingTimeWindowUnit(cfg["unit"].(string))

	return rolling
}
