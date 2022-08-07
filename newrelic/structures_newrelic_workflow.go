package newrelic

import (
	"errors"
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

	destinationConfigurations, destinationConfigurationsOk := d.GetOk("destination_configurations")

	if !destinationConfigurationsOk {
		return nil, errors.New("workflow requires a destination configurations attribute")
	}

	if destinationConfigurationsOk {
		var destination map[string]interface{}

		x := destinationConfigurations.([]interface{})

		for _, destinationConfiguration := range x {
			destination = destinationConfiguration.(map[string]interface{})

			if val, err := expandWorkflowDestinationConfiguration(destination); err == nil {
				workflow.DestinationConfigurations = append(workflow.DestinationConfigurations, *val)
			}
		}
	}

	issuesFilter, issuesFilterOk := d.GetOk("issues_filter")

	if !issuesFilterOk {
		return nil, errors.New("workflow requires a issues filter attribute")
	}

	if issuesFilterOk {
		f, err := expandWorkflowsIssuesFilter(issuesFilter.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		workflow.IssuesFilter = f
	}

	enrichments, enrichmentsOk := d.GetOk("enrichments")

	if !enrichmentsOk {
		return nil, errors.New("workflow requires a enrichments attribute")
	}

	if enrichmentsOk {
		e, err := expandWorkflowsEnrichments(enrichments.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		workflow.Enrichments = e
	}

	return &workflow, nil
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

	destinationConfigurations, destinationConfigurationsOk := d.GetOk("destination_configurations")

	if !destinationConfigurationsOk {
		return nil, errors.New("workflow requires a destination configurations attribute")
	}

	if destinationConfigurationsOk {
		var destination map[string]interface{}

		x := destinationConfigurations.([]interface{})

		for _, destinationConfiguration := range x {
			destination = destinationConfiguration.(map[string]interface{})

			if val, err := expandWorkflowDestinationConfiguration(destination); err == nil {
				workflow.DestinationConfigurations = append(workflow.DestinationConfigurations, *val)
			}
		}
	}

	issuesFilter, issuesFilterOk := d.GetOk("issues_filter")

	if !issuesFilterOk {
		return nil, errors.New("workflow requires a issues filter attribute")
	}

	if issuesFilterOk {
		f, err := expandWorkflowsUpdateIssuesFilter(issuesFilter.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		workflow.IssuesFilter = f
	}

	enrichments, enrichmentsOk := d.GetOk("enrichments")

	if !enrichmentsOk {
		return nil, errors.New("workflow requires a enrichments attribute")
	}

	if enrichmentsOk {
		e, err := expandWorkflowsUpdateEnrichments(enrichments.(*schema.Set).List())
		if err != nil {
			return nil, err
		}

		workflow.Enrichments = e
	}

	return &workflow, nil
}

func expandWorkflowDestinationConfiguration(cfg map[string]interface{}) (*workflows.AiWorkflowsDestinationConfigurationInput, error) {
	destinationConfigurationInput := workflows.AiWorkflowsDestinationConfigurationInput{}

	if channelID, ok := cfg["channel_id"]; ok {
		destinationConfigurationInput.ChannelId = channelID.(string)
	}

	return &destinationConfigurationInput, nil
}

func expandWorkflowsIssuesFilter(issuesFilterSet []interface{}) (workflows.AiWorkflowsFilterInput, error) {
	issuesFilter := make([]workflows.AiWorkflowsFilterInput, len(issuesFilterSet))

	for _, issuesFilterConfig := range issuesFilterSet {
		cfg := issuesFilterConfig.(map[string]interface{})

		if name, ok := cfg["name"]; ok {
			issuesFilter[0].Name = name.(string)
		}

		if filterType, ok := cfg["type"]; ok {
			issuesFilter[0].Type = workflows.AiWorkflowsFilterType(filterType.(string))
		}

		if predicates, ok := cfg["predicates"]; ok {
			var predicateInput map[string]interface{}

			x := predicates.([]interface{})

			for _, predicate := range x {
				predicateInput = predicate.(map[string]interface{})

				if val, err := expandWorkflowPredicate(predicateInput); err == nil {
					issuesFilter[0].Predicates = append(issuesFilter[0].Predicates, *val)
				}
			}
		}

		break
	}

	return issuesFilter[0], nil
}

func expandWorkflowsUpdateIssuesFilter(issuesFilterSet []interface{}) (workflows.AiWorkflowsUpdatedFilterInput, error) {
	issuesFilter := make([]workflows.AiWorkflowsUpdatedFilterInput, len(issuesFilterSet))

	for _, issuesFilterConfig := range issuesFilterSet {
		cfg := issuesFilterConfig.(map[string]interface{})

		if id, ok := cfg["filter_id"]; ok {
			issuesFilter[0].ID = id.(string)
		}

		if name, ok := cfg["name"]; ok {
			issuesFilter[0].FilterInput.Name = name.(string)
		}

		if filterType, ok := cfg["type"]; ok {
			issuesFilter[0].FilterInput.Type = workflows.AiWorkflowsFilterType(filterType.(string))
		}

		if predicates, ok := cfg["predicates"]; ok {
			var predicateInput map[string]interface{}

			x := predicates.([]interface{})

			for _, predicate := range x {
				predicateInput = predicate.(map[string]interface{})

				if val, err := expandWorkflowPredicate(predicateInput); err == nil {
					issuesFilter[0].FilterInput.Predicates = append(issuesFilter[0].FilterInput.Predicates, *val)
				}
			}
		}

		break
	}

	return issuesFilter[0], nil
}

func expandWorkflowPredicate(cfg map[string]interface{}) (*workflows.AiWorkflowsPredicateInput, error) {
	predicateInput := workflows.AiWorkflowsPredicateInput{}

	if attribute, ok := cfg["attribute"]; ok {
		predicateInput.Attribute = attribute.(string)
	}

	if operator, ok := cfg["operator"]; ok {
		predicateInput.Operator = workflows.AiWorkflowsOperator(operator.(string))
	}

	if values, ok := cfg["values"]; ok {
		var valuesList []string

		x := values.([]interface{})

		for _, v := range x {
			vInput := v.(string)
			valuesList = append(valuesList, vInput)
		}

		predicateInput.Values = valuesList
	}

	return &predicateInput, nil
}

func expandWorkflowsEnrichments(enrichmentsSet []interface{}) (*workflows.AiWorkflowsEnrichmentsInput, error) {
	enrichments := make([]workflows.AiWorkflowsEnrichmentsInput, len(enrichmentsSet))

	for _, enrichmentsConfig := range enrichmentsSet {
		cfg := enrichmentsConfig.(map[string]interface{})

		if nrqlList, ok := cfg["nrql"]; ok {
			var nrqlInput map[string]interface{}

			x := nrqlList.([]interface{})

			for _, nrql := range x {
				nrqlInput = nrql.(map[string]interface{})

				if val, err := expandWorkflowNrqlInput(nrqlInput); err == nil {
					enrichments[0].NRQL = append(enrichments[0].NRQL, *val)
				}
			}
		}

		break
	}

	return &enrichments[0], nil
}

func expandWorkflowsUpdateEnrichments(enrichmentsSet []interface{}) (*workflows.AiWorkflowsUpdateEnrichmentsInput, error) {
	enrichments := make([]workflows.AiWorkflowsUpdateEnrichmentsInput, len(enrichmentsSet))

	for _, enrichmentsConfig := range enrichmentsSet {
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

		break
	}

	return &enrichments[0], nil
}

func expandWorkflowNrqlInput(cfg map[string]interface{}) (*workflows.AiWorkflowsNRQLEnrichmentInput, error) {
	nrqlInput := workflows.AiWorkflowsNRQLEnrichmentInput{}

	if name, ok := cfg["name"]; ok {
		nrqlInput.Name = name.(string)
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

	if err := d.Set("destination_configurations", destinationConfigurations); err != nil {
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
