package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/workflows"
)

func expandWorkflow(d *schema.ResourceData) (*workflows.AiWorkflowsCreateWorkflowInput, error) {
	workflow := workflows.AiWorkflowsCreateWorkflowInput{
		Name:                d.Get("name").(string),
		EnrichmentsEnabled:  d.Get("enrichments_enabled").(bool),
		DestinationsEnabled: d.Get("destinations_enabled").(bool),
		WorkflowEnabled:     d.Get("workflow_enabled").(bool),
		MutingRulesHandling: workflows.AiWorkflowsMutingRulesHandling(d.Get("muting_rules_handling").(string)),
	}

	workflow.DestinationConfigurations = expandWorkflowDestinationConfigurations(d.Get("destination_configuration").(*schema.Set).List())
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
		Configuration: expandWorkflowEnrichmentNrqlConfigurations(nrqlConfig["configurations"].([]interface{})),
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

func expandWorkflowDestinationConfigurations(destinationConfigurations []interface{}) []workflows.AiWorkflowsDestinationConfigurationInput {
	input := []workflows.AiWorkflowsDestinationConfigurationInput{}

	for _, d := range destinationConfigurations {
		input = append(input, expandWorkflowDestinationConfiguration(d.(map[string]interface{})))
	}

	return input
}

func expandWorkflowDestinationConfiguration(cfg map[string]interface{}) workflows.AiWorkflowsDestinationConfigurationInput {
	destinationConfigurationInput := workflows.AiWorkflowsDestinationConfigurationInput{}

	if channelID, ok := cfg["channel_id"]; ok {
		destinationConfigurationInput.ChannelId = channelID.(string)
	}

	return destinationConfigurationInput
}

func expandWorkflowIssuesFilter(issuesFilter []interface{}) workflows.AiWorkflowsFilterInput {
	if len(issuesFilter) == 1 {
		filter := issuesFilter[0].(map[string]interface{})
		predicates := []workflows.AiWorkflowsPredicateInput{}

		if p, ok := filter["predicates"]; ok {
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
	workflow := workflows.AiWorkflowsUpdateWorkflowInput{
		ID:                  d.Get("workflow_id").(string),
		Name:                d.Get("name").(string),
		EnrichmentsEnabled:  d.Get("enrichments_enabled").(bool),
		DestinationsEnabled: d.Get("destinations_enabled").(bool),
		WorkflowEnabled:     d.Get("workflow_enabled").(bool),
		MutingRulesHandling: workflows.AiWorkflowsMutingRulesHandling(d.Get("muting_rules_handling").(string)),
	}

	workflow.DestinationConfigurations = expandWorkflowDestinationConfigurations(d.Get("destination_configuration").(*schema.Set).List())
	workflow.IssuesFilter = expandWorkflowUpdateIssuesFilter(d.Get("issues_filter").(*schema.Set).List())

	enrichments, enrichmentsOk := d.GetOk("enrichments")
	if enrichmentsOk {
		e, err := expandWorkflowsUpdateEnrichments(enrichments.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		workflow.Enrichments = e
	}

	return &workflow, nil
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

func expandWorkflowsUpdateEnrichments(enrichmentsSet []interface{}) (*workflows.AiWorkflowsUpdateEnrichmentsInput, error) {
	enrichments := make([]workflows.AiWorkflowsUpdateEnrichmentsInput, len(enrichmentsSet))

	enrichmentsConfig := enrichmentsSet[0]
	cfg := enrichmentsConfig.(map[string]interface{})

	if nrqlList, ok := cfg["nrql"]; ok {
		var nrqlInput map[string]interface{}

		x := nrqlList.([]interface{})

		for _, nrql := range x {
			nrqlInput = nrql.(map[string]interface{})

			if val, err := expandWorkflowUpdateNrqlInput(nrqlInput); err == nil {
				enrichments[0].NRQL = append(enrichments[0].NRQL, *val)
			}
		}
	}

	return &enrichments[0], nil
}

func expandWorkflowUpdateNrqlInput(cfg map[string]interface{}) (*workflows.AiWorkflowsNRQLUpdateEnrichmentInput, error) {
	nrqlInput := workflows.AiWorkflowsNRQLUpdateEnrichmentInput{}

	if name, ok := cfg["name"]; ok {
		nrqlInput.Name = name.(string)
	}

	if id, ok := cfg["enrichment_id"]; ok {
		nrqlInput.ID = id.(string)
	}

	if configurationList, ok := cfg["configurations"]; ok {
		var configurationsInput map[string]interface{}

		x := configurationList.([]interface{})

		for _, configuration := range x {
			configurationsInput = configuration.(map[string]interface{})

			if val, err := expandWorkflowConfigurationInput(configurationsInput); err == nil {
				nrqlInput.Configuration = append(nrqlInput.Configuration, *val)
			}
		}
	}

	return &nrqlInput, nil
}

func expandWorkflowConfigurationInput(cfg map[string]interface{}) (*workflows.AiWorkflowsNRQLConfigurationInput, error) {
	configurationInput := workflows.AiWorkflowsNRQLConfigurationInput{}

	if query, ok := cfg["query"]; ok {
		configurationInput.Query = query.(string)
	}

	return &configurationInput, nil
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

	if err = d.Set("workflow_enabled", workflow.WorkflowEnabled); err != nil {
		return err
	}

	if err = d.Set("muting_rules_handling", workflow.MutingRulesHandling); err != nil {
		return err
	}

	destinationConfigurations, destinationConfigurationsErr := flattenWorkflowDestinationConfigurations(&workflow.DestinationConfigurations)
	if destinationConfigurationsErr != nil {
		return destinationConfigurationsErr
	}

	if err := d.Set("destination_configuration", destinationConfigurations); err != nil {
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

	if err := d.Set("enrichments", enrichments); err != nil {
		return err
	}

	return nil
}

func flattenWorkflowDestinationConfigurations(d *[]workflows.AiWorkflowsDestinationConfiguration) ([]map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}

	var destinationConfigurations []map[string]interface{}

	for _, destinationConfiguration := range *d {
		if val, err := flattenWorkflowDestinationConfiguration(&destinationConfiguration); err == nil {
			destinationConfigurations = append(destinationConfigurations, val)
		}
	}

	return destinationConfigurations, nil
}

func flattenWorkflowDestinationConfiguration(d *workflows.AiWorkflowsDestinationConfiguration) (map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}

	destinationConfigurationResult := make(map[string]interface{})

	destinationConfigurationResult["channel_id"] = d.ChannelId
	destinationConfigurationResult["name"] = d.Name
	destinationConfigurationResult["type"] = d.Type

	return destinationConfigurationResult, nil
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
	issuesFilterResult["predicates"] = predicates

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

func flattenWorkflowEnrichments(e *[]workflows.AiWorkflowsEnrichment) (interface{}, error) {
	if e == nil {
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
	enrichmentResult["configurations"] = configuration

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
