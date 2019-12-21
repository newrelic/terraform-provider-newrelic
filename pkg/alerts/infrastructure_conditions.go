package alerts

import (
	"fmt"
	"strconv"
)

// ListInfrastructureConditions is used to retrieve New Relic Infrastructure alert conditions.
func (a *Alerts) ListInfrastructureConditions(policyID int) ([]InfrastructureCondition, error) {
	resp := infrastructureConditionsResponse{}
	paramsMap := map[string]string{"policy_id": strconv.Itoa(policyID)}
	_, err := a.infraClient.Get("/alerts/conditions", &paramsMap, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Conditions, nil
}

// GetInfrastructureCondition is used to retrieve a specific New Relic Infrastructure alert condition.
func (a *Alerts) GetInfrastructureCondition(conditionID int) (*InfrastructureCondition, error) {
	resp := infrastructureConditionResponse{}
	url := fmt.Sprintf("/alerts/conditions/%d", conditionID)
	_, err := a.infraClient.Get(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

// CreateInfrastructureCondition is used to create a New Relic Infrastructure alert condition.
func (a *Alerts) CreateInfrastructureCondition(condition InfrastructureCondition) (*InfrastructureCondition, error) {
	resp := infrastructureConditionResponse{}
	reqBody := infrastructureConditionRequest{condition}

	_, err := a.infraClient.Post("/alerts/conditions", nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

// UpdateInfrastructureCondition is used to update a New Relic Infrastructure alert condition.
func (a *Alerts) UpdateInfrastructureCondition(condition InfrastructureCondition) (*InfrastructureCondition, error) {
	resp := infrastructureConditionResponse{}
	reqBody := infrastructureConditionRequest{condition}

	url := fmt.Sprintf("/alerts/conditions/%d", condition.ID)
	_, err := a.infraClient.Put(url, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Condition, nil
}

// DeleteInfrastructureCondition is used to delete a New Relic Infrastructure alert condition.
func (a *Alerts) DeleteInfrastructureCondition(conditionID int) error {
	url := fmt.Sprintf("/alerts/conditions/%d", conditionID)
	_, err := a.infraClient.Delete(url, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

type infrastructureConditionsResponse struct {
	Conditions []InfrastructureCondition `json:"data,omitempty"`
}

type infrastructureConditionResponse struct {
	Condition InfrastructureCondition `json:"data,omitempty"`
}

type infrastructureConditionRequest struct {
	Condition InfrastructureCondition `json:"data,omitempty"`
}
