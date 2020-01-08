package alerts

import (
	"fmt"
)

// ListAlertPoliciesParams represents a set of filters to be
// used when querying New Relic alert policies.
type ListAlertPoliciesParams struct {
	Name string `url:"filter[name],omitempty"`
}

// ListAlertPolicies returns a list of Alert Policies for a given account.
func (alerts *Alerts) ListAlertPolicies(params *ListAlertPoliciesParams) ([]AlertPolicy, error) {
	response := alertPoliciesResponse{}
	alertPolicies := []AlertPolicy{}
	nextURL := "/alerts_policies.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		alertPolicies = append(alertPolicies, response.Policies...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertPolicies, nil
}

// GetAlertPolicy returns a specific alert policy by ID for a given account.
func (alerts *Alerts) GetAlertPolicy(id int) (*AlertPolicy, error) {
	policies, err := alerts.ListAlertPolicies(nil)

	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if policy.ID == id {
			return &policy, nil
		}
	}

	return nil, fmt.Errorf("no alert policy found for id %d", id)
}

// CreateAlertPolicy creates a new alert policy for a given account.
func (alerts *Alerts) CreateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error) {
	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}

	_, err := alerts.client.Post("/alerts_policies.json", nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// UpdateAlertPolicy update an alert policy for a given account.
func (alerts *Alerts) UpdateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error) {
	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}

	url := fmt.Sprintf("/alerts_policies/%d.json", policy.ID)

	_, err := alerts.client.Put(url, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeleteAlertPolicy deletes an existing alert policy for a given account.
func (alerts *Alerts) DeleteAlertPolicy(id int) (*AlertPolicy, error) {
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	_, err := alerts.client.Delete(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

type alertPoliciesResponse struct {
	Policies []AlertPolicy `json:"policies,omitempty"`
}

type alertPolicyResponse struct {
	Policy AlertPolicy `json:"policy,omitempty"`
}

type alertPolicyRequestBody struct {
	Policy AlertPolicy `json:"policy"`
}
