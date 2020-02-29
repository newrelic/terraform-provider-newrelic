package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// PluginsCondition represents an alert condition for New Relic Plugins.
type PluginsCondition struct {
	PolicyID          int             `json:"-"`
	ID                int             `json:"id,omitempty"`
	Name              string          `json:"name,omitempty"`
	Enabled           bool            `json:"enabled"`
	Entities          []string        `json:"entities,omitempty"`
	Metric            string          `json:"metric,omitempty"`
	MetricDescription string          `json:"metric_description,omitempty"`
	RunbookURL        string          `json:"runbook_url,omitempty"`
	Terms             []ConditionTerm `json:"terms,omitempty"`
	ValueFunction     string          `json:"value_function,omitempty"`
	Plugin            AlertPlugin     `json:"plugin,omitempty"`
}

// AlertPlugin represents a plugin to use with a Plugin alert condition.
type AlertPlugin struct {
	ID   string `json:"id,omitempty"`
	GUID string `json:"guid,omitempty"`
}

// ListPluginsConditions returns alert conditions for New Relic plugins for a given alert policy.
func (alerts *Alerts) ListPluginsConditions(policyID int) ([]*PluginsCondition, error) {
	conditions := []*PluginsCondition{}
	queryParams := listPluginsConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_plugins_conditions.json"

	for nextURL != "" {
		response := pluginsConditionsResponse{}
		resp, err := alerts.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		for _, c := range response.PluginsConditions {
			c.PolicyID = policyID
		}

		conditions = append(conditions, response.PluginsConditions...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetPluginsCondition gets information about an alert condition for a plugin
// given a policy ID and plugin ID.
func (alerts *Alerts) GetPluginsCondition(policyID int, pluginID int) (*PluginsCondition, error) {
	conditions, err := alerts.ListPluginsConditions(policyID)

	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == pluginID {
			return condition, nil
		}
	}

	return nil, errors.NewNotFoundf("no condition found for policy %d and condition ID %d", policyID, pluginID)
}

// CreatePluginsCondition creates an alert condition for a plugin.
func (alerts *Alerts) CreatePluginsCondition(condition PluginsCondition) (*PluginsCondition, error) {
	reqBody := pluginConditionRequestBody{
		PluginsCondition: condition,
	}
	resp := pluginConditionResponse{}

	u := fmt.Sprintf("/alerts_plugins_conditions/policies/%d.json", condition.PolicyID)
	_, err := alerts.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.PluginsCondition.PolicyID = condition.PolicyID

	return &resp.PluginsCondition, nil
}

// UpdatePluginsCondition updates an alert condition for a plugin.
func (alerts *Alerts) UpdatePluginsCondition(condition PluginsCondition) (*PluginsCondition, error) {
	reqBody := pluginConditionRequestBody{
		PluginsCondition: condition,
	}
	resp := pluginConditionResponse{}

	u := fmt.Sprintf("/alerts_plugins_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.PluginsCondition.PolicyID = condition.PolicyID

	return &resp.PluginsCondition, nil
}

// DeletePluginsCondition deletes a plugin alert condition.
func (alerts *Alerts) DeletePluginsCondition(id int) (*PluginsCondition, error) {
	resp := pluginConditionResponse{}
	u := fmt.Sprintf("/alerts_plugins_conditions/%d.json", id)

	_, err := alerts.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.PluginsCondition, nil
}

type listPluginsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type pluginsConditionsResponse struct {
	PluginsConditions []*PluginsCondition `json:"plugins_conditions,omitempty"`
}

type pluginConditionResponse struct {
	PluginsCondition PluginsCondition `json:"plugins_condition,omitempty"`
}

type pluginConditionRequestBody struct {
	PluginsCondition PluginsCondition `json:"plugins_condition,omitempty"`
}
