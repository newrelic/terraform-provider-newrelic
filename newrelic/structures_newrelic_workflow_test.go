//go:build unit
// +build unit

package newrelic

import (
	"reflect"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"

	"github.com/newrelic/newrelic-client-go/v2/pkg/workflows"

	"github.com/stretchr/testify/assert"
)

func TestExpandWorkflow(t *testing.T) {
	nrql := []map[string]interface{}{{
		"name": "enrichment-test-1",
		"configuration": []map[string]interface{}{{
			"query": "SELECT * FROM Log",
		}},
	}}
	enrichments := []workflows.AiWorkflowsEnrichment{{
		Name: "enrichment-test-1",
		Type: workflows.AiWorkflowsEnrichmentTypeTypes.NRQL,
		Configurations: []ai.AiWorkflowsConfiguration{{
			Query: "SELECT * FROM Log",
		}},
	}}

	destinationConfigurations := []workflows.AiWorkflowsDestinationConfiguration{{
		Name:                 "destination-test",
		Type:                 workflows.AiWorkflowsDestinationTypeTypes.WEBHOOK,
		ChannelId:            "300848f9-c713-463c-9036-40b45c4c970f",
		NotificationTriggers: []workflows.AiWorkflowsNotificationTrigger{workflows.AiWorkflowsNotificationTriggerTypes.ACTIVATED},
	}}

	issuesFilter := workflows.AiWorkflowsFilter{
		Name: "issues-filter-test",
		Type: workflows.AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []workflows.AiWorkflowsPredicate{{
			Attribute: "source",
			Operator:  workflows.AiWorkflowsOperatorTypes.EQUAL,
			Values:    []string{"newrelic"},
		}},
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *workflows.AiWorkflowsWorkflow
	}{
		"valid workflow": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"enrichments": []map[string]interface{}{{
					"nrql": nrql,
				}},
				"issues_filter": []map[string]interface{}{{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				}},
				"destination": []map[string]interface{}{{
					"channel_id":            "300848f9-c713-463c-9036-40b45c4c970f",
					"notification_triggers": []string{"ACTIVATED"},
				}},
			},
			Expanded: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				Enrichments:               enrichments,
				DestinationConfigurations: destinationConfigurations,
				IssuesFilter:              issuesFilter,
			},
		},
		"valid workflow without enrichments": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"issues_filter": []map[string]interface{}{{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				}},
				"destination": []map[string]interface{}{{
					"channel_id":            "300848f9-c713-463c-9036-40b45c4c970f",
					"notification_triggers": []string{"ACTIVATED"},
				}},
			},
			Expanded: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				DestinationConfigurations: destinationConfigurations,
				IssuesFilter:              issuesFilter,
			},
		},
		"valid workflow without notification triggers": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"issues_filter": []map[string]interface{}{{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				}},
				"destination": []map[string]interface{}{{
					"channel_id": "300848f9-c713-463c-9036-40b45c4c970f",
				}},
			},
			Expanded: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				DestinationConfigurations: destinationConfigurations,
				IssuesFilter:              issuesFilter,
			},
		},
	}

	r := resourceNewRelicWorkflow()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded, err := expandWorkflow(d)

		if tc.ExpectErr {
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			assert.Nil(t, err)
		}

		if tc.Expanded != nil {
			assert.Equal(t, tc.Expanded.Name, expanded.Name)
		}
	}
}

func TestFlattenWorkflow(t *testing.T) {
	enrichments := []workflows.AiWorkflowsEnrichment{{
		Name: "enrichment-test-1",
		Type: workflows.AiWorkflowsEnrichmentTypeTypes.NRQL,
		Configurations: []ai.AiWorkflowsConfiguration{{
			Query: "SELECT * FROM Log",
		}},
	}}

	destinationConfigurations := []workflows.AiWorkflowsDestinationConfiguration{{
		Name:      "destination-test",
		Type:      workflows.AiWorkflowsDestinationTypeTypes.WEBHOOK,
		ChannelId: "300848f9-c713-463c-9036-40b45c4c970f",
	}}

	destinationConfigurationsWithNotificationTriggers := destinationConfigurations
	destinationConfigurationsWithNotificationTriggers[0].NotificationTriggers =
		[]workflows.AiWorkflowsNotificationTrigger{workflows.AiWorkflowsNotificationTriggerTypes.ACTIVATED}

	issuesFilter := workflows.AiWorkflowsFilter{
		Name: "issues-filter-test",
		Type: workflows.AiWorkflowsFilterTypeTypes.FILTER,
		Predicates: []workflows.AiWorkflowsPredicate{{
			Attribute: "source",
			Operator:  workflows.AiWorkflowsOperatorTypes.EQUAL,
			Values:    []string{"newrelic"},
		}},
	}

	guid := workflows.EntityGUID("testworkflowentityguid")
	r := resourceNewRelicWorkflow()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *workflows.AiWorkflowsWorkflow
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"guid":                  "testworkflowentityguid",
				"enrichments": []map[string]interface{}{{
					"name": "enrichment-test-1",
					"type": "NRQL",
					"configuration": []map[string]interface{}{{
						"query": "SELECT * FROM Log",
					}},
				}},
				"issues_filter": map[string]interface{}{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				},
				"destination": []map[string]interface{}{{
					"channel_id":            "300848f9-c713-463c-9036-40b45c4c970f",
					"name":                  "destination-test",
					"type":                  "WEBHOOK",
					"notification_triggers": []workflows.AiWorkflowsNotificationTrigger{workflows.AiWorkflowsNotificationTriggerTypes.ACTIVATED},
				}},
			},
			Flattened: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				Enrichments:               enrichments,
				DestinationConfigurations: destinationConfigurations,
				IssuesFilter:              issuesFilter,
				GUID:                      guid,
			},
		},
		"no_enrichments": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"guid":                  "testworkflowentityguid",
				"issues_filter": map[string]interface{}{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				},
				"destination": []map[string]interface{}{{
					"channel_id":            "300848f9-c713-463c-9036-40b45c4c970f",
					"name":                  "destination-test",
					"type":                  "WEBHOOK",
					"notification_triggers": []workflows.AiWorkflowsNotificationTrigger{workflows.AiWorkflowsNotificationTriggerTypes.ACTIVATED},
				}},
			},
			Flattened: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				Enrichments:               []workflows.AiWorkflowsEnrichment{},
				DestinationConfigurations: destinationConfigurations,
				IssuesFilter:              issuesFilter,
				GUID:                      guid,
			},
		},
		"no_notification_triggers": {
			Data: map[string]interface{}{
				"name":                  "workflow-test",
				"enrichments_enabled":   true,
				"destinations_enabled":  true,
				"enabled":               true,
				"muting_rules_handling": "NOTIFY_ALL_ISSUES",
				"guid":                  "testworkflowentityguid",
				"issues_filter": map[string]interface{}{
					"name": "issues-filter-test",
					"type": "FILTER",
					"predicate": []map[string]interface{}{{
						"attribute": "source",
						"operator":  "EQUAL",
						"values":    []string{"newrelic"},
					}},
				},
				"destination": []map[string]interface{}{{
					"channel_id": "300848f9-c713-463c-9036-40b45c4c970f",
					"name":       "destination-test",
					"type":       "WEBHOOK",
				}},
			},
			Flattened: &workflows.AiWorkflowsWorkflow{
				Name:                      "workflow-test",
				EnrichmentsEnabled:        true,
				DestinationsEnabled:       true,
				WorkflowEnabled:           true,
				MutingRulesHandling:       workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES,
				Enrichments:               []workflows.AiWorkflowsEnrichment{},
				DestinationConfigurations: destinationConfigurationsWithNotificationTriggers,
				IssuesFilter:              issuesFilter,
				GUID:                      guid,
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenWorkflow(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("err: %s", err)
				}

				if k == "issues_filter" {
					testFlattenWorkflowsIssuesFilter(t, v, tc.Flattened.IssuesFilter)
				} else if k == "enrichments" {
					for _, enrichment := range tc.Flattened.Enrichments {
						testFlattenWorkflowsEnrichment(t, v, enrichment)
					}
				} else if k == "destination" {
					for _, configuration := range tc.Flattened.DestinationConfigurations {
						testFlattenWorkflowsDestinationConfiguration(t, v, configuration)
					}
				} else {
					assert.Equal(t, x, v)
				}
			}
		}
	}
}

func TestWorkflowStateUpgradeV0(t *testing.T) {
	expected := testWorkflowStateDataV1()
	actual, err := migrateStateNewRelicWorkflowV0toV1(nil, testWorkflowStateDataV0(), nil)

	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func testWorkflowStateDataV0() map[string]any {
	return map[string]interface{}{
		"workflow_enabled": true,
		"destination_configuration": []map[string]interface{}{{
			"channel_id": "300848f9-c713-463c-9036-40b45c4c970f",
		}},
	}
}

func testWorkflowStateDataV1() map[string]any {
	v0 := testWorkflowStateDataV0()

	return map[string]interface{}{
		"enabled":     v0["workflow_enabled"],
		"destination": v0["destination_configuration"],
	}
}

func testFlattenWorkflowsIssuesFilter(t *testing.T, v interface{}, issuesFilter workflows.AiWorkflowsFilter) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "type":
			assert.Equal(t, cv, string(issuesFilter.Type))
		case "name":
			assert.Equal(t, cv, issuesFilter.Name)
		case "predicate":
			for _, predicate := range issuesFilter.Predicates {
				testFlattenWorkflowsIssuesFilterPredicate(t, v, predicate)
			}
		}
	}
}

func testFlattenWorkflowsIssuesFilterPredicate(t *testing.T, v interface{}, predicate workflows.AiWorkflowsPredicate) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "attribute":
			assert.Equal(t, cv, predicate.Attribute)
		case "operator":
			assert.Equal(t, cv, string(predicate.Operator))
		case "values":
			assert.Equal(t, cv, predicate.Values)
		}
	}
}

func testFlattenWorkflowsDestinationConfiguration(t *testing.T, v interface{}, configuration workflows.AiWorkflowsDestinationConfiguration) {
	for _, v1 := range v.([]map[string]interface{}) {
		for ck, cv := range v1 {
			switch ck {
			case "channel_id":
				assert.Equal(t, cv, configuration.ChannelId)
			case "name":
				assert.Equal(t, cv, configuration.Name)
			case "type":
				assert.Equal(t, cv, string(configuration.Type))
			case "notification_triggers":
				assert.Equal(t, cv, configuration.NotificationTriggers)
			}
		}
	}
}

func testFlattenWorkflowsEnrichment(t *testing.T, v interface{}, enrichment workflows.AiWorkflowsEnrichment) {
	for _, v1 := range v.([]map[string]interface{}) {
		for ck, cv := range v1 {
			switch ck {
			case "configuration":
				for _, configuration := range enrichment.Configurations {
					testFlattenWorkflowsEnrichmentConfiguration(t, cv, configuration)
				}
			case "name":
				assert.Equal(t, cv, enrichment.Name)
			case "type":
				assert.Equal(t, cv, string(enrichment.Type))
			}
		}
	}
}

func testFlattenWorkflowsEnrichmentConfiguration(t *testing.T, v interface{}, configuration ai.AiWorkflowsConfiguration) {
	for _, v1 := range v.([]map[string]interface{}) {
		for ck, cv := range v1 {
			switch ck {
			case "query":
				assert.Equal(t, cv, configuration.Query)
			}
		}
	}
}
