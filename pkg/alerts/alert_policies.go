package alerts

type listAlertsResponse struct {
	AlertPolicies []AlertPolicy `json:"policies,omitempty"`
}

// GetAlertPolicy returns a specific alert policy by ID
func (a *Alerts) GetAlertPolicy(id int) (*AlertPolicy, error) {
	return nil, nil
}

// ListAlertPolicies returns a list of Alert Policies for a given account
func (alerts *Alerts) ListAlertPolicies() ([]AlertPolicy, error) {
	res := listAlertsResponse{}
	err := alerts.client.Get("/alerts_policies.json", nil, &res)

	if err != nil {
		return nil, err
	}

	return res.AlertPolicies, nil
}
