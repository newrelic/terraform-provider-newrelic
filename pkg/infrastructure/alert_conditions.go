package infrastructure

import "strconv"

type listAlertConditionsResponse struct {
	AlertConditions []AlertCondition `json:"data,omitempty"`
}

// ListAlertConditions is used to retrieve New Relic Infrastructure alert conditions.
func (i *Infrastructure) ListAlertConditions(policyID int) ([]AlertCondition, error) {
	res := listAlertConditionsResponse{}
	paramsMap := map[string]string{"policy_id": strconv.Itoa(policyID)}
	err := i.client.Get("/alerts/conditions", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.AlertConditions, nil
}
