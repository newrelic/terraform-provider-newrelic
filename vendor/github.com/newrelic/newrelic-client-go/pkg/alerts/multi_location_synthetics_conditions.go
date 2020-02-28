package alerts

import (
	"fmt"
)

// MultiLocationSyntheticsCondition represents a location-based failure condition.
//
// ViolationTimeLimitSeconds must be one of 3600, 7200, 14400, 28800, 43200, 86400.
type MultiLocationSyntheticsCondition struct {
	ID                        int                                    `json:"id,omitempty"`
	Name                      string                                 `json:"name,omitempty"`
	Enabled                   bool                                   `json:"enabled"`
	RunbookURL                string                                 `json:"runbook_url,omitempty"`
	MonitorID                 string                                 `json:"monitor_id,omitempty"`
	Entities                  []string                               `json:"entities,omitempty"`
	Terms                     []MultiLocationSyntheticsConditionTerm `json:"terms,omitempty"`
	ViolationTimeLimitSeconds int                                    `json:"violation_time_limit_seconds,omitempty"`
}

// MultiLocationSyntheticsConditionTerm represents a single term for a location-based failure condition.
//
// Priority must be "warning" or "critical".
// Threshold must be greater than zero.
type MultiLocationSyntheticsConditionTerm struct {
	Priority  string `json:"priority,omitempty"`
	Threshold int    `json:"threshold,omitempty"`
}

// ListMultiLocationSyntheticsConditions returns alert conditions for a specified policy.
func (alerts *Alerts) ListMultiLocationSyntheticsConditions(policyID int) ([]*MultiLocationSyntheticsCondition, error) {
	response := multiLocationSyntheticsConditionListResponse{}
	multiLocationSyntheticsConditions := []*MultiLocationSyntheticsCondition{}
	queryParams := listMultiLocationSyntheticsConditionsParams{
		PolicyID: policyID,
	}

	nextURL := fmt.Sprintf("/alerts_location_failure_conditions/policies/%d.json", policyID)

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		multiLocationSyntheticsConditions = append(multiLocationSyntheticsConditions, response.MultiLocationSyntheticsConditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return multiLocationSyntheticsConditions, nil
}

// CreateMultiLocationSyntheticsCondition creates an alert condition for a specified policy.
func (alerts *Alerts) CreateMultiLocationSyntheticsCondition(condition MultiLocationSyntheticsCondition, policyID int) (*MultiLocationSyntheticsCondition, error) {
	reqBody := multiLocationSyntheticsConditionRequestBody{
		MultiLocationSyntheticsCondition: condition,
	}
	resp := multiLocationSyntheticsConditionCreateResponse{}

	u := fmt.Sprintf("/alerts_location_failure_conditions/policies/%d.json", policyID)
	_, err := alerts.client.Post(u, nil, &reqBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.MultiLocationSyntheticsCondition, nil
}

// UpdateMultiLocationSyntheticsCondition updates an alert condition.
func (alerts *Alerts) UpdateMultiLocationSyntheticsCondition(condition MultiLocationSyntheticsCondition) (*MultiLocationSyntheticsCondition, error) {
	reqBody := multiLocationSyntheticsConditionRequestBody{
		MultiLocationSyntheticsCondition: condition,
	}
	resp := multiLocationSyntheticsConditionCreateResponse{}

	u := fmt.Sprintf("/alerts_location_failure_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.MultiLocationSyntheticsCondition, nil
}

// DeleteMultiLocationSyntheticsCondition delete an alert condition.
func (alerts *Alerts) DeleteMultiLocationSyntheticsCondition(conditionID int) (*MultiLocationSyntheticsCondition, error) {
	resp := multiLocationSyntheticsConditionCreateResponse{}
	u := fmt.Sprintf("/alerts_conditions/%d.json", conditionID)

	_, err := alerts.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.MultiLocationSyntheticsCondition, nil
}

type listMultiLocationSyntheticsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type multiLocationSyntheticsConditionListResponse struct {
	MultiLocationSyntheticsConditions []*MultiLocationSyntheticsCondition `json:"location_failure_conditions,omitempty"`
}

type multiLocationSyntheticsConditionCreateResponse struct {
	MultiLocationSyntheticsCondition MultiLocationSyntheticsCondition `json:"location_failure_condition,omitempty"`
}

type multiLocationSyntheticsConditionRequestBody struct {
	MultiLocationSyntheticsCondition MultiLocationSyntheticsCondition `json:"location_failure_condition,omitempty"`
}
