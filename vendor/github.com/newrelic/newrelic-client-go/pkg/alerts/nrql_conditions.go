package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// NrqlCondition represents a New Relic NRQL Alert condition.
type NrqlCondition struct {
	Terms               []ConditionTerm   `json:"terms,omitempty"`
	Nrql                NrqlQuery         `json:"nrql,omitempty"`
	Type                string            `json:"type,omitempty"`
	Name                string            `json:"name,omitempty"`
	RunbookURL          string            `json:"runbook_url,omitempty"`
	ValueFunction       ValueFunctionType `json:"value_function,omitempty"`
	ID                  int               `json:"id,omitempty"`
	ViolationCloseTimer int               `json:"violation_time_limit_seconds,omitempty"`
	ExpectedGroups      int               `json:"expected_groups,omitempty"`
	IgnoreOverlap       bool              `json:"ignore_overlap,omitempty"`
	Enabled             bool              `json:"enabled"`
}

// NrqlQuery represents a NRQL query to use with a NRQL alert condition
type NrqlQuery struct {
	Query      string `json:"query,omitempty"`
	SinceValue string `json:"since_value,omitempty"`
}

// ListNrqlConditions returns NRQL alert conditions for a specified policy.
func (a *Alerts) ListNrqlConditions(policyID int) ([]*NrqlCondition, error) {
	conditions := []*NrqlCondition{}
	queryParams := listNrqlConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_nrql_conditions.json"

	for nextURL != "" {
		response := nrqlConditionsResponse{}
		resp, err := a.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		conditions = append(conditions, response.NrqlConditions...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetNrqlCondition gets information about a NRQL alert condition
// for a specified policy ID and condition ID.
func (a *Alerts) GetNrqlCondition(policyID int, id int) (*NrqlCondition, error) {
	conditions, err := a.ListNrqlConditions(policyID)
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

// CreateNrqlCondition creates a NRQL alert condition.
func (a *Alerts) CreateNrqlCondition(policyID int, condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/policies/%d.json", policyID)
	_, err := a.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

// UpdateNrqlCondition updates a NRQL alert condition.
func (a *Alerts) UpdateNrqlCondition(condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", condition.ID)
	_, err := a.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

// DeleteNrqlCondition deletes a NRQL alert condition.
func (a *Alerts) DeleteNrqlCondition(id int) (*NrqlCondition, error) {
	resp := nrqlConditionResponse{}
	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", id)

	_, err := a.client.Delete(u, nil, &resp)

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
