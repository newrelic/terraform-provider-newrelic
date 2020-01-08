package alerts

import (
	"fmt"
)

// ListAlertConditions returns alert conditions for a specified policy.
func (alerts *Alerts) ListAlertConditions(policyID int) ([]*AlertCondition, error) {
	response := alertConditionsResponse{}
	alertConditions := []*AlertCondition{}
	queryParams := listAlertConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_conditions.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		for _, c := range response.Conditions {
			c.PolicyID = policyID
		}

		alertConditions = append(alertConditions, response.Conditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertConditions, nil
}

// GetAlertCondition gets an alert condition for a specified policy ID and condition ID.
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

// CreateAlertCondition creates an alert condition for a specified policy.
func (alerts *Alerts) CreateAlertCondition(condition AlertCondition) (*AlertCondition, error) {
	reqBody := alertConditionRequestBody{
		Condition: condition,
	}
	resp := alertConditionResponse{}

	u := fmt.Sprintf("/alerts_conditions/policies/%d.json", condition.PolicyID)
	_, err := alerts.client.Post(u, nil, reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.Condition.PolicyID = condition.PolicyID

	return &resp.Condition, nil
}

// UpdateAlertCondition updates an alert condition.
func (alerts *Alerts) UpdateAlertCondition(condition AlertCondition) (*AlertCondition, error) {
	reqBody := alertConditionRequestBody{
		Condition: condition,
	}
	resp := alertConditionResponse{}

	u := fmt.Sprintf("/alerts_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.Condition.PolicyID = condition.PolicyID

	return &resp.Condition, nil
}

// DeleteAlertCondition delete an alert condition.
func (alerts *Alerts) DeleteAlertCondition(id int) (*AlertCondition, error) {
	resp := alertConditionResponse{}
	u := fmt.Sprintf("/alerts_conditions/%d.json", id)

	_, err := alerts.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

type listAlertConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type alertConditionsResponse struct {
	Conditions []*AlertCondition `json:"conditions,omitempty"`
}

type alertConditionResponse struct {
	Condition AlertCondition `json:"condition,omitempty"`
}

type alertConditionRequestBody struct {
	Condition AlertCondition `json:"condition,omitempty"`
}
