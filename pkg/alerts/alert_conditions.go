package alerts

import (
	"fmt"
	"strconv"
)

func (alerts *Alerts) ListAlertConditions(policyID int) ([]*AlertCondition, error) {
	response := alertConditionsResponse{}
	alertConditions := []AlertCondition{}
	queryParams := map[string]string{
		"policy_id": strconv.Itoa(policyID),
	}

	nextURL := "/alerts_conditions.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		alertConditions = append(alertConditions, response.Conditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	results := []*AlertCondition{}
	for _, condition := range alertConditions {
		results = append(results, &condition)
	}

	return results, nil
}

func (alerts *Alerts) GetAlertCondition(policyID int, id int) (*AlertCondition, error) {
	conditions, err := alerts.ListAlertConditions(policyID)
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

type alertConditionsResponse struct {
	Conditions []AlertCondition `json:"conditions,omitempty"`
}
