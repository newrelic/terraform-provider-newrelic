package alerts

import "fmt"

// ListPluginsConditions returns alert conditions for New Relic plugins for an account.
func (alerts *Alerts) ListPluginsConditions(policyID int) ([]*PluginCondition, error) {
	response := pluginsConditionsResponse{}
	conditions := []*PluginCondition{}
	queryParams := listPluginsConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_plugins_conditions.json"

	for nextURL != "" {
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

// GetPluginCondition gets information about an alert condition for a plugin
// given a policy ID and plugin ID.
func (alerts *Alerts) GetPluginCondition(policyID int, id int) (*PluginCondition, error) {
	conditions, err := alerts.ListPluginsConditions(policyID)

	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, fmt.Errorf("no condition found for policy %d and condition ID %d", policyID, id)
}

// UpdatePluginCondition updates an alert condition for a plugin.
func (alerts *Alerts) UpdatePluginCondition(condition PluginCondition) (*PluginCondition, error) {
	reqBody := pluginConditionRequestBody{
		PluginCondition: condition,
	}
	resp := pluginConditionResponse{}

	u := fmt.Sprintf("/alerts_plugins_conditions/%d.json", condition.ID)
	_, err := alerts.client.Put(u, nil, reqBody, &resp)

	if err != nil {
		return nil, err
	}

	resp.PluginCondition.PolicyID = condition.PolicyID

	return &resp.PluginCondition, nil
}

type listPluginsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type pluginsConditionsResponse struct {
	PluginsConditions []*PluginCondition `json:"plugins_conditions,omitempty"`
}

type pluginConditionResponse struct {
	PluginCondition PluginCondition `json:"plugins_condition,omitempty"`
}

type pluginConditionRequestBody struct {
	PluginCondition PluginCondition `json:"plugins_condition,omitempty"`
}
