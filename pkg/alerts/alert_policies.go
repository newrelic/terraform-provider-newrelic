package alerts

import (
	"fmt"
)

// ListAlertPoliciesParams represents a set of filters to be
// used when querying New Relic alert policies.
type ListAlertPoliciesParams struct {
	Name *string
}

// ListAlertPolicies returns a list of Alert Policies for a given account.
func (alerts *Alerts) ListAlertPolicies(params *ListAlertPoliciesParams) ([]AlertPolicy, error) {
	respBody := alertPoliciesResponse{}
	alertPolicies := []AlertPolicy{}
	nextURL := "/alerts_policies.json"
	paramsMap := buildListAlertPoliciesParamsMap(params)

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, &paramsMap, &respBody)

		if err != nil {
			return nil, err
		}

		alertPolicies = append(alertPolicies, respBody.Policies...)

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
	reqBody := createAlertPolicyRequestBody{
		Policy: policy,
	}
	resBody := alertPolicyResponse{}

	_, err := alerts.client.Post("/alerts_policies.json", nil, &reqBody, &resBody)

	if err != nil {
		return nil, err
	}

	return &resBody.Policy, nil
}

// UpdateAlertPolicy update an alert policy for a given account.
func (alerts *Alerts) UpdateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error) {
	reqBody := updateAlertPolicyRequestBody{
		Policy: policy,
	}
	resBody := alertPolicyResponse{}

	url := fmt.Sprintf("/alerts_policies/%d.json", policy.ID)

	_, err := alerts.client.Put(url, nil, reqBody, &resBody)

	if err != nil {
		return nil, err
	}

	return &resBody.Policy, nil
}

// DeleteAlertPolicy deletes an existing alert policy for a given account.
func (alerts *Alerts) DeleteAlertPolicy(id int) (*AlertPolicy, error) {
	respBody := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	_, err := alerts.client.Delete(url, nil, &respBody)

	if err != nil {
		return nil, err
	}

	return &respBody.Policy, nil
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

type alertPoliciesResponse struct {
	Policies []AlertPolicy `json:"policies,omitempty"`
}

type alertPolicyResponse struct {
	Policy AlertPolicy `json:"policy,omitempty"`
}

type createAlertPolicyRequestBody struct {
	Policy AlertPolicy `json:"policy,omitempty"`
}

type updateAlertPolicyRequestBody struct {
	Policy AlertPolicy `json:"policy"`
}
