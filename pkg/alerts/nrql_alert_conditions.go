package alerts

import (
	"fmt"
	"strconv"
)

// ListNrqlAlertConditions returns NRQL alert conditions for a specified policy.
func (alerts *Alerts) ListNrqlAlertConditions(policyID int) ([]*NrqlCondition, error) {
	response := nrqlConditionsResponse{}
	conditions := []*NrqlCondition{}
	queryParams := map[string]string{
		"policy_id": strconv.Itoa(policyID),
	}

	nextURL := "/alerts_nrql_conditions.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		// TODO: Is this really necessary?
		for _, c := range response.NrqlConditions {
			c.PolicyID = policyID
		}

		conditions = append(conditions, response.NrqlConditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetNrqlAlertCondition gets information about a NRQL alert condition
// for a specified policy ID and condition ID.
func (alerts *Alerts) GetNrqlAlertCondition(policyID int, id int) (*NrqlCondition, error) {
	conditions, err := alerts.ListNrqlAlertConditions(policyID)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, fmt.Errorf("no condition found for policy %d and condition ID %d", policyID, id)
}

type nrqlConditionsResponse struct {
	NrqlConditions []*NrqlCondition `json:"nrql_conditions,omitempty"`
}
