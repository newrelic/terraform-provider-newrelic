package alerts

import "strconv"

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

type alertConditionsResponse struct {
	Conditions []AlertCondition `json:"conditions,omitempty"`
}
