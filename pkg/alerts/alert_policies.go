package alerts

import (
	"fmt"
)

// ListAlertPoliciesParams represents a set of filters to be
// used when querying New Relic alert policies.
type ListAlertPoliciesParams struct {
	Name *string
}

type listAlertsResponse struct {
	AlertPolicies []AlertPolicy `json:"policies,omitempty"`
}

type createAlertPolicyRequestBody struct {
	Policy AlertPolicy `json:"policy"`
}

type createAlertPolicyResponse struct {
	Policy AlertPolicy `json:"policy,omitempty"`
}

type updateAlertPolicyRequestBody struct {
	Policy AlertPolicy `json:"policy"`
}

type updateAlertPolicyResponse struct {
	Policy AlertPolicy `json:"policy,omitempty"`
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

// ListAlertPolicies returns a list of Alert Policies for a given account.
func (alerts *Alerts) ListAlertPolicies(params *ListAlertPoliciesParams) ([]AlertPolicy, error) {
	res := listAlertsResponse{}
	paramsMap := buildListAlertPoliciesParamsMap(params)
	err := alerts.client.Get("/alerts_policies.json", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.AlertPolicies, nil
}

// CreateAlertPolicy creates a new alert policy for a given account.
func (alerts *Alerts) CreateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error) {
	reqBody := createAlertPolicyRequestBody{
		Policy: policy,
	}
	resp := createAlertPolicyResponse{}

	err := alerts.client.Post("/alerts_policies.json", reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// UpdateAlertPolicy update an alert policy for a given account.
func (alerts *Alerts) UpdateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error) {
	reqBody := updateAlertPolicyRequestBody{
		Policy: policy,
	}
	resp := updateAlertPolicyResponse{}

	url := fmt.Sprintf("/alerts_policies/%d.json", policy.ID)

	err := alerts.client.Put(url, reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeleteAlertPolicy deletes an existing alert policy for a given account.
func (alerts *Alerts) DeleteAlertPolicy(id int) error {
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	err := alerts.client.Delete(url)

	return err
}

func buildListAlertPoliciesParamsMap(params *ListAlertPoliciesParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.Name != nil {
		paramsMap["filter[name]"] = *params.Name
	}

	return paramsMap
}
