package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// ListConditions returns alert conditions for a specified policy.
func (alerts *Alerts) ListConditions(policyID int) ([]*Condition, error) {
	response := alertConditionsResponse{}
	alertConditions := []*Condition{}
	queryParams := listConditionsParams{
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

// GetCondition gets an alert condition for a specified policy ID and condition ID.
func (alerts *Alerts) GetCondition(policyID int, id int) (*Condition, error) {
	conditions, err := alerts.ListConditions(policyID)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, errors.NewNotFoundf("no condition found for policy %d and condition ID %d", policyID, id)
}

// CreateCondition creates an alert condition for a specified policy.
func (alerts *Alerts) CreateCondition(condition Condition) (*Condition, error) {
	reqBody := alertConditionRequestBody{
		Condition: condition,
	}
	resp := alertConditionResponse{}

	u := fmt.Sprintf("/alerts_conditions/policies/%d.json", condition.PolicyID)
	_, err := alerts.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.Condition.PolicyID = condition.PolicyID

	return &resp.Condition, nil
}

// UpdateCondition updates an alert condition.
func (alerts *Alerts) UpdateCondition(condition Condition) (*Condition, error) {
	reqBody := alertConditionRequestBody{
		Condition: condition,
	}
	resp := alertConditionResponse{}

	u := fmt.Sprintf("/alerts_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.Condition.PolicyID = condition.PolicyID

	return &resp.Condition, nil
}

// DeleteCondition delete an alert condition.
func (alerts *Alerts) DeleteCondition(id int) (*Condition, error) {
	resp := alertConditionResponse{}
	u := fmt.Sprintf("/alerts_conditions/%d.json", id)

	_, err := alerts.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

type listConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type alertConditionsResponse struct {
	Conditions []*Condition `json:"conditions,omitempty"`
}

type alertConditionResponse struct {
	Condition Condition `json:"condition,omitempty"`
}

type alertConditionRequestBody struct {
	Condition Condition `json:"condition,omitempty"`
}
