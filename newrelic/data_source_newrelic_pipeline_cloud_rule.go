package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
)

func dataSourceNewRelicDropRulePipelineCloudRuleRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicDropRulePipelineCloudRuleRelationshipRead,
		Schema: map[string]*schema.Schema{
			"drop_rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NRQL Drop Rule ID to match against Pipeline Cloud Rule tags.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Pipeline Cloud Rule entity ID that matches the drop rule.",
			},
		},
	}
}

func dataSourceNewRelicDropRulePipelineCloudRuleRelationshipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// accountID := selectAccountID(providerConfig, d)
	dropRuleID := strings.Split(d.Get("drop_rule_id").(string), ":")[1]

	query := `{
		actor {
			entityManagement {
				entitySearch(query: "type = 'PIPELINE_CLOUD_RULE'") {
					entities {
						id
						tags {
							key
							values
						}
						name
					}
				}
			}
		}
	}`

	variables := map[string]interface{}{}

	result, err := client.NerdGraph.QueryWithContext(ctx, query, variables)
	if err != nil {
		return diag.FromErr(err)
	}

	queryResponse := result.(nerdgraph.QueryResponse)
	queryResponseActor := queryResponse.Actor.(map[string]interface{})
	queryResponseActorEntityManagement := queryResponseActor["entityManagement"]
	queryResponseActorEntityManagementEntitySearch := queryResponseActorEntityManagement.(map[string]interface{})["entitySearch"]
	entitiesMap := queryResponseActorEntityManagementEntitySearch.(map[string]interface{})

	entitiesJSON, entitiesJSONErr := json.Marshal(entitiesMap)
	if entitiesJSONErr != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	var response RawGetPipelineCloudRuleEntitiesResponse
	if err := json.Unmarshal(entitiesJSON, &response); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse NerdGraph response: %w", err))
	}

	var matchedEntity *RawPipelineCloudRuleEntity
	for _, entity := range response.Entities {
		for _, tag := range entity.Tags {
			if tag.Key == "legacyDropRuleId" {
				for _, value := range tag.Values {
					if value == dropRuleID {
						matchedEntity = &entity
						break
					}
				}
				if matchedEntity != nil {
					break
				}
			}
		}
		if matchedEntity != nil {
			break
		}
	}

	if matchedEntity == nil {
		return diag.Errorf("no Pipeline Cloud Rule found matching drop rule ID: %s", dropRuleID)
	}

	d.SetId(matchedEntity.ID)

	return nil
}

type RawGetPipelineCloudRuleEntitiesResponse struct {
	Entities []RawPipelineCloudRuleEntity `json:"entities"`
}

type RawPipelineCloudRuleEntity struct {
	ID   string                    `json:"id"`
	Name string                    `json:"name"`
	Tags []RawPipelineCloudRuleTag `json:"tags"`
}

type RawPipelineCloudRuleTag struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}
