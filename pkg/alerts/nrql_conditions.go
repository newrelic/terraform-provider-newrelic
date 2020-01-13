package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// ListNrqlConditions returns NRQL alert conditions for a specified policy.
func (alerts *Alerts) ListNrqlConditions(policyID int) ([]*NrqlCondition, error) {
	response := nrqlConditionsResponse{}
	conditions := []*NrqlCondition{}
	queryParams := listNrqlConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_nrql_conditions.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		for _, c := range response.NrqlConditions {
			c.PolicyID = policyID
		}

		conditions = append(conditions, response.NrqlConditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetNrqlCondition gets information about a NRQL alert condition
// for a specified policy ID and condition ID.
func (alerts *Alerts) GetNrqlCondition(policyID int, id int) (*NrqlCondition, error) {
	conditions, err := alerts.ListNrqlConditions(policyID)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, &errors.ErrorNotFound{
		Message: fmt.Sprintf("no condition found for policy %d and condition ID %d", policyID, id),
	}
}

// CreateNrqlCondition creates a NRQL alert condition.
func (alerts *Alerts) CreateNrqlCondition(condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/policies/%d.json", condition.PolicyID)
	_, err := alerts.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.NrqlCondition.PolicyID = condition.PolicyID

	return &resp.NrqlCondition, nil
}

// UpdateNrqlCondition updates a NRQL alert condition.
func (alerts *Alerts) UpdateNrqlCondition(condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.NrqlCondition.PolicyID = condition.PolicyID

	return &resp.NrqlCondition, nil
}

// DeleteNrqlCondition deletes a NRQL alert condition.
func (alerts *Alerts) DeleteNrqlCondition(id int) (*NrqlCondition, error) {
	resp := nrqlConditionResponse{}
	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", id)

	_, err := alerts.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

type listNrqlConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type nrqlConditionsResponse struct {
	NrqlConditions []*NrqlCondition `json:"nrql_conditions,omitempty"`
}

type nrqlConditionResponse struct {
	NrqlCondition NrqlCondition `json:"nrql_condition,omitempty"`
}

type nrqlConditionRequestBody struct {
	NrqlCondition NrqlCondition `json:"nrql_condition,omitempty"`
}
