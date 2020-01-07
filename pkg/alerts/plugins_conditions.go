package alerts

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

type listPluginsConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type pluginsConditionsResponse struct {
	PluginsConditions []*PluginCondition `json:"plugins_conditions,omitempty"`
}
