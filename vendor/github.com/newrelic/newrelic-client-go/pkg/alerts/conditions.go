package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Condition represents a New Relic alert condition.
// TODO: custom unmarshal entities to ints?
type Condition struct {
	PolicyID            int                  `json:"-"`
	ID                  int                  `json:"id,omitempty"`
	Type                string               `json:"type,omitempty"`
	Name                string               `json:"name,omitempty"`
	Enabled             bool                 `json:"enabled"`
	Entities            []string             `json:"entities,omitempty"`
	Metric              string               `json:"metric,omitempty"`
	RunbookURL          string               `json:"runbook_url,omitempty"`
	Terms               []ConditionTerm      `json:"terms,omitempty"`
	UserDefined         ConditionUserDefined `json:"user_defined,omitempty"`
	Scope               string               `json:"condition_scope,omitempty"`
	GCMetric            string               `json:"gc_metric,omitempty"`
	ViolationCloseTimer int                  `json:"violation_close_timer,omitempty"`
}

// ConditionUserDefined represents user defined metrics for the New Relic alert condition.
type ConditionUserDefined struct {
	Metric        string `json:"metric,omitempty"`
	ValueFunction string `json:"value_function,omitempty"`
}

// ConditionTerm represents the terms of a New Relic alert condition.
type ConditionTerm struct {
	Duration     int     `json:"duration,string,omitempty"`
	Operator     string  `json:"operator,omitempty"`
	Priority     string  `json:"priority,omitempty"`
	Threshold    float64 `json:"threshold,string"`
	TimeFunction string  `json:"time_function,omitempty"`
}

// ListConditions returns alert conditions for a specified policy.
func (alerts *Alerts) ListConditions(policyID int) ([]*Condition, error) {
	alertConditions := []*Condition{}
	queryParams := listConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_conditions.json"

	for nextURL != "" {
		response := alertConditionsResponse{}
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
