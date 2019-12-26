package alerts

import "strconv"

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

type nrqlConditionsResponse struct {
	NrqlConditions []*NrqlCondition `json:"nrql_conditions,omitempty"`
}
