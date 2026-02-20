package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func flattenAlertPolicy(policy *alerts.AlertsPolicy, d *schema.ResourceData, accountID int) error {
	var err error

	err = d.Set("name", policy.Name)
	if err != nil {
		return err
	}

	err = d.Set("incident_preference", policy.IncidentPreference)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	return nil
}

func flattenAlertPolicyWithEntityGUID(ctx context.Context, client *newrelic.NewRelic, policy *alerts.AlertsPolicy, d *schema.ResourceData, accountID int) error {
	var err error

	err = d.Set("name", policy.Name)
	if err != nil {
		return err
	}

	err = d.Set("incident_preference", policy.IncidentPreference)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	// Fetch entity GUID using entity search
	entityGUID, err := getAlertPolicyEntityGUID(ctx, client, policy.Name, accountID)
	if err != nil {
		log.Printf("[WARN] Error fetching entity GUID for alert policy %s: %s", policy.Name, err)
		// Don't fail the entire operation if entity GUID fetch fails
		// Just log the error and continue
	} else if entityGUID != "" {
		err = d.Set("entity_guid", entityGUID)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAlertPolicyEntityGUID(ctx context.Context, client *newrelic.NewRelic, policyName string, accountID int) (string, error) {
	// Escape single quotes in policy name for the query
	escapedName := escapeSingleQuote(policyName)

	// Build query to search for alert policy entity
	query := fmt.Sprintf("name = '%s' AND domain = 'AIOPS' AND type = 'POLICY'", escapedName)

	entityResults, err := client.Entities.GetEntitySearchByQueryWithContext(ctx,
		entities.EntitySearchOptions{},
		query,
		[]entities.EntitySearchSortCriteria{},
	)

	if err != nil {
		return "", err
	}

	if entityResults == nil || len(entityResults.Results.Entities) == 0 {
		return "", fmt.Errorf("no entity found for alert policy: %s", policyName)
	}

	// Find the entity matching the account ID
	for _, entity := range entityResults.Results.Entities {
		if entity.GetAccountID() == accountID && entity.GetName() == policyName {
			return string(entity.GetGUID()), nil
		}
	}

	// If no exact match with account ID, return the first result
	if len(entityResults.Results.Entities) > 0 {
		return string(entityResults.Results.Entities[0].GetGUID()), nil
	}

	return "", fmt.Errorf("no entity found for alert policy: %s in account: %d", policyName, accountID)
}
