package alerts

import (
	"fmt"
)

// ListSyntheticsConditions returns a list of Synthetics alert conditions for a given policy.
func (alerts *Alerts) ListSyntheticsConditions(policyID int) ([]SyntheticsCondition, error) {
	response := syntheticsConditionsResponse{}
	conditions := []SyntheticsCondition{}
	nextURL := fmt.Sprintf("/alerts_synthetics_conditions.json")
	queryParams := listSyntheticsConditionsParams{
		PolicyID: policyID,
	}

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		conditions = append(conditions, response.Conditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// CreateSyntheticsCondition creates a new Synthetics alert condition.
func (alerts *Alerts) CreateSyntheticsCondition(policyID int, condition SyntheticsCondition) (*SyntheticsCondition, error) {
	resp := syntheticsConditionResponse{}
	reqBody := syntheticsConditionRequest{condition}
	url := fmt.Sprintf("/alerts_synthetics_conditions/policies/%d.json", policyID)
	_, err := alerts.client.Post(url, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

// UpdateSyntheticsCondition updates an existing Synthetics alert condition.
func (alerts *Alerts) UpdateSyntheticsCondition(condition SyntheticsCondition) (*SyntheticsCondition, error) {
	resp := syntheticsConditionResponse{}
	reqBody := syntheticsConditionRequest{condition}
	url := fmt.Sprintf("/alerts_synthetics_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(url, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

// DeleteSyntheticsCondition deletes a Synthetics alert condition.
func (alerts *Alerts) DeleteSyntheticsCondition(conditionID int) (*SyntheticsCondition, error) {
	resp := syntheticsConditionResponse{}
	url := fmt.Sprintf("/alerts_synthetics_conditions/%d.json", conditionID)
	_, err := alerts.client.Delete(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

type listSyntheticsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type syntheticsConditionsResponse struct {
	Conditions []SyntheticsCondition `json:"synthetics_conditions,omitempty"`
}

type syntheticsConditionResponse struct {
	Condition SyntheticsCondition `json:"synthetics_condition,omitempty"`
}

type syntheticsConditionRequest struct {
	Condition SyntheticsCondition `json:"synthetics_condition,omitempty"`
}
