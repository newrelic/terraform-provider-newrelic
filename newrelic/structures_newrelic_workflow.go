package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workflows"
)

// migrateStateNewRelicWorkflowV0toV1 currently facilitates migrating:
// `workflow_enabled` to `enabled`
// `destination_configuration` to `destination`
// `predicates` to singular
// `configurations` to singular
func migrateStateNewRelicWorkflowV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["enabled"] = rawState["workflow_enabled"]
	rawState["destination"] = rawState["destination_configuration"]
	delete(rawState, "workflow_enabled")
	delete(rawState, "destination_configuration")

	var issueFilter = rawState["issues_filter"]

	if issueFilter != nil {
		rawState["issues_filter"] = migrateWorkflowIssuesFilterV0toV1(issueFilter.([]interface{}))
	}

	var enrichments = rawState["enrichments"]

	if enrichments != nil {
		rawState["enrichments"] = migrateWorkflowEnrichmentsV0toV1(enrichments.([]interface{}))
	}

	return rawState, nil
}

func migrateWorkflowIssuesFilterV0toV1(issuesFilter []interface{}) []map[string]interface{} {
	var input map[string]interface{}

	if len(issuesFilter) == 1 {
		input = issuesFilter[0].(map[string]interface{})

		input["predicate"] = input["predicates"]
		delete(input, "predicates")
	}

	var returnInput []map[string]interface{}
	returnInput = append(returnInput, input)

	return returnInput
}

func migrateWorkflowEnrichmentsV0toV1(enrichments []interface{}) []map[string]interface{} {
	input := map[string]interface{}{}
	var returnInput []map[string]interface{}

	if len(enrichments) != 1 {
		return returnInput
	}

	richments := enrichments[0].(map[string]interface{})
	nrql := richments["nrql"].([]interface{})
	input["nrql"] = migrateWorkflowNrqlsV0toV1(nrql)

	returnInput = append(returnInput, input)
	return returnInput
}

func migrateWorkflowNrqlsV0toV1(nrqls []interface{}) []map[string]interface{} {
	var input []map[string]interface{}

	for i, n := range nrqls {
		var newN = n.(map[string]interface{})
		newN["configuration"] = newN["configurations"]
		input = append(input, newN)
		delete(input[i], "configurations")
	}

	return input
}

func expandWorkflow(d *schema.ResourceData) (*workflows.AiWorkflowsCreateWorkflowInput, error) {
	workflow := workflows.AiWorkflowsCreateWorkflowInput{
		Name:                d.Get("name").(string),
		EnrichmentsEnabled:  d.Get("enrichments_enabled").(bool),
		DestinationsEnabled: d.Get("destinations_enabled").(bool),
		WorkflowEnabled:     d.Get("enabled").(bool),
		MutingRulesHandling: workflows.AiWorkflowsMutingRulesHandling(d.Get("muting_rules_handling").(string)),
	}

	workflow.DestinationConfigurations = expandWorkflowDestinationConfigurations(d.Get("destination").(*schema.Set).List())
	workflow.IssuesFilter = expandWorkflowIssuesFilter(d.Get("issues_filter").(*schema.Set).List())

	enrichments, enrichmentsOk := d.GetOk("enrichments")
	if enrichmentsOk {
		workflow.Enrichments = expandWorkflowEnrichments(enrichments.(*schema.Set).List())
	}

	return &workflow, nil
}

func expandWorkflowEnrichments(enrichments []interface{}) *workflows.AiWorkflowsEnrichmentsInput {
	input := workflows.AiWorkflowsEnrichmentsInput{}

	if len(enrichments) != 1 {
		return &input
	}

	richments := enrichments[0].(map[string]interface{})
	nrql := richments["nrql"].([]interface{})
	input.NRQL = expandWorkflowNrqls(nrql)

	return &input
}

func expandWorkflowNrqls(nrqls []interface{}) []workflows.AiWorkflowsNRQLEnrichmentInput {
	input := []workflows.AiWorkflowsNRQLEnrichmentInput{}

	for _, n := range nrqls {
		input = append(input, expandWorkflowNrql(n.(map[string]interface{})))
	}

	return input
}

func expandWorkflowNrql(nrqlConfig map[string]interface{}) workflows.AiWorkflowsNRQLEnrichmentInput {
	return workflows.AiWorkflowsNRQLEnrichmentInput{
		Name:          nrqlConfig["name"].(string),
		Configuration: expandWorkflowEnrichmentNrqlConfigurations(nrqlConfig["configuration"].([]interface{})),
	}
}

func expandWorkflowEnrichmentNrqlConfigurations(nrqlConfigs []interface{}) []workflows.AiWorkflowsNRQLConfigurationInput {
	input := []workflows.AiWorkflowsNRQLConfigurationInput{}

	for _, n := range nrqlConfigs {
		input = append(input, expandWorkflowConfiguration(n.(map[string]interface{})))
	}

	return input
}

func expandWorkflowConfiguration(cfg map[string]interface{}) workflows.AiWorkflowsNRQLConfigurationInput {
	return workflows.AiWorkflowsNRQLConfigurationInput{
		Query: cfg["query"].(string),
	}
}

func expandWorkflowUpdateEnrichments(enrichments []interface{}) workflows.AiWorkflowsUpdateEnrichmentsInput {
	if len(enrichments) != 1 {
		return workflows.AiWorkflowsUpdateEnrichmentsInput{}
	}

	richments := enrichments[0].(map[string]interface{})
	nrql := richments["nrql"].([]interface{})
	return workflows.AiWorkflowsUpdateEnrichmentsInput{
		NRQL: expandWorkflowUpdateNrqls(nrql),
	}
}

func expandWorkflowUpdateNrqls(nrqls []interface{}) []workflows.AiWorkflowsNRQLUpdateEnrichmentInput {
	input := []workflows.AiWorkflowsNRQLUpdateEnrichmentInput{}

	for _, n := range nrqls {
		input = append(input, expandWorkflowUpdateNrql(n.(map[string]interface{})))
	}

	return input
}

func expandWorkflowUpdateNrql(nrqlConfig map[string]interface{}) workflows.AiWorkflowsNRQLUpdateEnrichmentInput {
	return workflows.AiWorkflowsNRQLUpdateEnrichmentInput{
		Name:          nrqlConfig["name"].(string),
		Configuration: expandWorkflowEnrichmentNrqlConfigurations(nrqlConfig["configuration"].([]interface{})),
	}
}

func expandWorkflowDestinationConfigurations(destinationConfigurations []interface{}) []workflows.AiWorkflowsDestinationConfigurationInput {
	input := []workflows.AiWorkflowsDestinationConfigurationInput{}

	for _, d := range destinationConfigurations {
		input = append(input, expandWorkflowDestinationConfiguration(d.(map[string]interface{})))
	}

	return input
}

func expandWorkflowDestinationConfiguration(cfg map[string]interface{}) workflows.AiWorkflowsDestinationConfigurationInput {
	destinationConfigurationInput := workflows.AiWorkflowsDestinationConfigurationInput{}
	var notificationTriggersInput []workflows.AiWorkflowsNotificationTrigger

	if channelID, ok := cfg["channel_id"]; ok {
		destinationConfigurationInput.ChannelId = channelID.(string)
	}

	if notificationTriggers, ok := cfg["notification_triggers"]; ok {
		for _, p := range notificationTriggers.([]interface{}) {
			notificationTriggersInput = append(notificationTriggersInput, workflows.AiWorkflowsNotificationTrigger(p.(string)))
		}
		destinationConfigurationInput.NotificationTriggers = notificationTriggersInput
	}

	return destinationConfigurationInput
}

func expandWorkflowIssuesFilter(issuesFilter []interface{}) workflows.AiWorkflowsFilterInput {
	if len(issuesFilter) == 1 {
		filter := issuesFilter[0].(map[string]interface{})
		predicates := []workflows.AiWorkflowsPredicateInput{}

		if p, ok := filter["predicate"]; ok {
			predicates = expandWorkflowIssuePredicates(p.([]interface{}))
		}

		return workflows.AiWorkflowsFilterInput{
			Name:       filter["name"].(string),
			Type:       workflows.AiWorkflowsFilterType(filter["type"].(string)),
			Predicates: predicates,
		}
	}

	return workflows.AiWorkflowsFilterInput{}
}

func expandWorkflowIssuePredicates(predicates []interface{}) []workflows.AiWorkflowsPredicateInput {
	input := []workflows.AiWorkflowsPredicateInput{}

	for _, p := range predicates {
		input = append(input, expandWorkflowIssuePredicate(p.(map[string]interface{})))
	}

	return input
}

func expandWorkflowIssuePredicate(predicate map[string]interface{}) workflows.AiWorkflowsPredicateInput {
	var valuesList []string

	for _, v := range predicate["values"].([]interface{}) {
		vInput := v.(string)
		valuesList = append(valuesList, vInput)
	}

	return workflows.AiWorkflowsPredicateInput{
		Attribute: predicate["attribute"].(string),
		Operator:  workflows.AiWorkflowsOperator(predicate["operator"].(string)),
		Values:    valuesList,
	}
}

func expandWorkflowUpdate(d *schema.ResourceData) (*workflows.AiWorkflowsUpdateWorkflowInput, error) {
	name := d.Get("name").(string)
	enrichmentsEnabled := d.Get("enrichments_enabled").(bool)
	destinationsEnabled := d.Get("destinations_enabled").(bool)
	workflowEnabled := d.Get("enabled").(bool)
	configurations := expandWorkflowDestinationConfigurations(d.Get("destination").(*schema.Set).List())
	filter := expandWorkflowUpdateIssuesFilter(d.Get("issues_filter").(*schema.Set).List())
	enrichments := getAndExpandWorkflowUpdateEnrichments(d)

	workflow := workflows.AiWorkflowsUpdateWorkflowInput{
		ID:                        d.Get("workflow_id").(string),
		Name:                      &name,
		EnrichmentsEnabled:        &enrichmentsEnabled,
		DestinationsEnabled:       &destinationsEnabled,
		WorkflowEnabled:           &workflowEnabled,
		MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandling(d.Get("muting_rules_handling").(string)),
		DestinationConfigurations: &configurations,
		IssuesFilter:              &filter,
		Enrichments:               &enrichments,
	}

	return &workflow, nil
}

func getAndExpandWorkflowUpdateEnrichments(d *schema.ResourceData) workflows.AiWorkflowsUpdateEnrichmentsInput {
	enrichmentsData, enrichmentsOk := d.GetOk("enrichments")
	if !enrichmentsOk {
		return workflows.AiWorkflowsUpdateEnrichmentsInput{
			NRQL: []workflows.AiWorkflowsNRQLUpdateEnrichmentInput{},
		}
	}

	return expandWorkflowUpdateEnrichments(enrichmentsData.(*schema.Set).List())
}

func expandWorkflowUpdateIssuesFilter(issuesFilter []interface{}) workflows.AiWorkflowsUpdatedFilterInput {
	if len(issuesFilter) == 1 {
		filter := issuesFilter[0].(map[string]interface{})

		return workflows.AiWorkflowsUpdatedFilterInput{
			ID:          filter["filter_id"].(string),
			FilterInput: expandWorkflowIssuesFilter(issuesFilter),
		}
	}

	return workflows.AiWorkflowsUpdatedFilterInput{}
}

func flattenWorkflow(workflow *workflows.AiWorkflowsWorkflow, d *schema.ResourceData) error {
	if workflow == nil {
		return nil
	}

	var err error

	if err = d.Set("name", workflow.Name); err != nil {
		return err
	}

	if err = d.Set("account_id", workflow.AccountID); err != nil {
		return err
	}

	if err = d.Set("workflow_id", workflow.ID); err != nil {
		return err
	}

	if err = d.Set("last_run", workflow.LastRun); err != nil {
		return err
	}

	if err = d.Set("enrichments_enabled", workflow.EnrichmentsEnabled); err != nil {
		return err
	}

	if err = d.Set("destinations_enabled", workflow.DestinationsEnabled); err != nil {
		return err
	}

	if err = d.Set("enabled", workflow.WorkflowEnabled); err != nil {
		return err
	}

	if err = d.Set("muting_rules_handling", workflow.MutingRulesHandling); err != nil {
		return err
	}

	destinationConfigurations, destinationConfigurationsErr := flattenWorkflowDestinationConfigurations(d, &workflow.DestinationConfigurations)
	if destinationConfigurationsErr != nil {
		return destinationConfigurationsErr
	}

	if err := d.Set("destination", destinationConfigurations); err != nil {
		return err
	}

	issuesFilter, issuesFilterErr := flattenWorkflowIssuesFilter(&workflow.IssuesFilter)
	if issuesFilterErr != nil {
		return issuesFilterErr
	}

	if err := d.Set("issues_filter", issuesFilter); err != nil {
		return err
	}

	enrichments, enrichmentsErr := flattenWorkflowEnrichments(&workflow.Enrichments)
	if enrichmentsErr != nil {
		return fmt.Errorf("[DEBUG] Error setting workflows enrichments: %#v", enrichmentsErr)
	}

	if len(enrichments) > 0 {
		if err := d.Set("enrichments", enrichments); err != nil {
			return err
		}
	}

	if err := d.Set("guid", workflow.GUID); err != nil {
		return err
	}

	return nil
}

func flattenWorkflowDestinationConfigurations(d *schema.ResourceData, configurations *[]workflows.AiWorkflowsDestinationConfiguration) ([]map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}

	currentStates := d.Get("destination").(*schema.Set).List()

	var destinationConfigurations []map[string]interface{}
	for _, config := range *configurations {
		var currentState map[string]interface{}
		for _, rawState := range currentStates {
			state := rawState.(map[string]interface{})
			if state["channel_id"] == config.ChannelId {
				currentState = state
			}
		}
		flattened, err := flattenWorkflowDestinationConfiguration(&config, currentState)
		if err != nil {
			return nil, err
		}
		destinationConfigurations = append(destinationConfigurations, flattened)
	}

	return destinationConfigurations, nil
}

func flattenWorkflowDestinationConfiguration(d *workflows.AiWorkflowsDestinationConfiguration, currentState map[string]interface{}) (map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}

	destinationConfigurationResult := make(map[string]interface{})

	destinationConfigurationResult["channel_id"] = d.ChannelId
	destinationConfigurationResult["name"] = d.Name
	destinationConfigurationResult["type"] = d.Type

	if currentState == nil || currentState["notification_triggers"] == nil {
		destinationConfigurationResult["notification_triggers"] = d.NotificationTriggers
	} else {
		destinationConfigurationResult["notification_triggers"] = normaliseTriggerList(
			d.NotificationTriggers,
			currentState["notification_triggers"].([]interface{}),
		)
	}

	return destinationConfigurationResult, nil
}

// preserve the state order to avoid state drift as trigger order is not preserved
func normaliseTriggerList(triggers []workflows.AiWorkflowsNotificationTrigger, currentState []interface{}) []workflows.AiWorkflowsNotificationTrigger {
	var result []workflows.AiWorkflowsNotificationTrigger
	for _, trigger := range currentState {
		triggerStillExists := false
		for _, actualTrigger := range triggers {
			if string(actualTrigger) == trigger {
				triggerStillExists = true
			}
		}

		if triggerStillExists {
			result = append(result, workflows.AiWorkflowsNotificationTrigger(trigger.(string)))
		}
	}

	for _, actualTrigger := range triggers {
		triggerExistedBefore := false
		for _, trigger := range currentState {
			if string(actualTrigger) == trigger {
				triggerExistedBefore = true
			}
		}

		if !triggerExistedBefore {
			result = append(result, actualTrigger)
		}
	}

	return result
}

func flattenWorkflowIssuesFilter(f *workflows.AiWorkflowsFilter) (interface{}, error) {
	if f == nil {
		return nil, nil
	}

	issuesFilter := make([]interface{}, 1)

	issuesFilterResult := make(map[string]interface{})
	issuesFilterResult["filter_id"] = f.ID
	issuesFilterResult["name"] = f.Name
	issuesFilterResult["type"] = f.Type

	predicates, predicatesErr := flattenWorkflowPredicates(f.Predicates)
	if predicatesErr != nil {
		return nil, predicatesErr
	}
	issuesFilterResult["predicate"] = predicates

	issuesFilter[0] = issuesFilterResult

	return issuesFilter, nil
}

func flattenWorkflowPredicates(p []workflows.AiWorkflowsPredicate) ([]map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	var predicates []map[string]interface{}

	for _, predicate := range p {
		if val, err := flattenWorkflowPredicate(&predicate); err == nil {
			predicates = append(predicates, val)
		}
	}

	return predicates, nil
}

func flattenWorkflowPredicate(p *workflows.AiWorkflowsPredicate) (map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	predicateResult := make(map[string]interface{})

	predicateResult["attribute"] = p.Attribute
	predicateResult["operator"] = p.Operator
	predicateResult["values"] = p.Values

	return predicateResult, nil
}

func flattenWorkflowEnrichments(e *[]workflows.AiWorkflowsEnrichment) ([]interface{}, error) {
	if e == nil || len(*e) == 0 {
		return nil, nil
	}

	nrql := make([]map[string]interface{}, len(*e))

	for i, enrichment := range *e {
		if val, err := flattenWorkflowEnrichment(&enrichment); err == nil {
			nrql[i] = val
		}
	}

	enrichmentsResult := make(map[string]interface{})
	enrichmentsResult["nrql"] = nrql

	enrichments := make([]interface{}, 1)
	enrichments[0] = enrichmentsResult

	return enrichments, nil
}

func flattenWorkflowEnrichment(e *workflows.AiWorkflowsEnrichment) (map[string]interface{}, error) {
	if e == nil {
		return nil, nil
	}

	enrichmentResult := make(map[string]interface{})

	enrichmentResult["enrichment_id"] = e.ID
	enrichmentResult["account_id"] = e.AccountID
	enrichmentResult["name"] = e.Name
	enrichmentResult["type"] = e.Type

	configuration, configurationErr := flattenWorkflowEnrichmentConfigurations(&e.Configurations)
	if configurationErr != nil {
		return nil, configurationErr
	}
	enrichmentResult["configuration"] = configuration

	return enrichmentResult, nil
}

func flattenWorkflowEnrichmentConfigurations(c *[]ai.AiWorkflowsConfiguration) ([]map[string]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	var configurations []map[string]interface{}

	for _, configuration := range *c {
		if val, err := flattenWorkflowEnrichmentConfiguration(&configuration); err == nil {
			configurations = append(configurations, val)
		}
	}

	return configurations, nil
}

func flattenWorkflowEnrichmentConfiguration(c *ai.AiWorkflowsConfiguration) (map[string]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	configurationResult := make(map[string]interface{})

	configurationResult["query"] = c.Query

	return configurationResult, nil
}
