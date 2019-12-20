package alerts

import "strconv"

type listInfrastructureConditionsResponse struct {
	InfrastructureConditions []InfrastructureCondition `json:"data,omitempty"`
}

// ListInfrastructureConditions is used to retrieve New Relic Infrastructure alert conditions.
func (a *Alerts) ListInfrastructureConditions(policyID int) ([]InfrastructureCondition, error) {
	res := listInfrastructureConditionsResponse{}
	paramsMap := map[string]string{"policy_id": strconv.Itoa(policyID)}
	_, err := a.client.Get("/alerts/conditions", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.InfrastructureConditions, nil
}

// GetInfrastructureCondition is used to retrieve a specific New Relic Infrastructure alert condition.
func (a *Alerts) GetInfrastructureCondition(conditionID int) ([]InfrastructureCondition, error) {
	return nil, nil
}

// CreateInfrastructureCondition is used to create a New Relic Infrastructure alert condition.
func (a *Alerts) CreateInfrastructureCondition(condition InfrastructureCondition) (*InfrastructureCondition, error) {
	return nil, nil
}

// UpdateInfrastructureCondition is used to update a New Relic Infrastructure alert condition.
func (a *Alerts) UpdateInfrastructureCondition(condition InfrastructureCondition) (*InfrastructureCondition, error) {
	return nil, nil
}

// DeleteInfrastructureCondition is used to delete a New Relic Infrastructure alert condition.
func (a *Alerts) DeleteInfrastructureCondition(conditionID int) (*InfrastructureCondition, error) {
	return nil, nil
}
