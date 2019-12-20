package alerts

import "strconv"

type listInfrastructureConditionsResponse struct {
	InfrastructureConditions []InfrastructureCondition `json:"data,omitempty"`
}

// ListInfrastructureConditions is used to retrieve New Relic Infrastructure alert conditions.
func (i *Alerts) ListInfrastructureConditions(policyID int) ([]InfrastructureCondition, error) {
	res := listInfrastructureConditionsResponse{}
	paramsMap := map[string]string{"policy_id": strconv.Itoa(policyID)}
	_, err := i.client.Get("/alerts/conditions", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.InfrastructureConditions, nil
}
