package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// Policy represents a New Relic alert policy.
type Policy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference string                   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}

// ListPoliciesParams represents a set of filters to be
// used when querying New Relic alert policies.
type ListPoliciesParams struct {
	Name string `url:"filter[name],omitempty"`
}

// ListPolicies returns a list of Alert Policies for a given account.
func (alerts *Alerts) ListPolicies(params *ListPoliciesParams) ([]Policy, error) {
	alertPolicies := []Policy{}
	nextURL := "/alerts_policies.json"

	for nextURL != "" {
		response := alertPoliciesResponse{}
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

// GetPolicy returns a specific alert policy by ID for a given account.
func (alerts *Alerts) GetPolicy(id int) (*Policy, error) {
	policies, err := alerts.ListPolicies(nil)

	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if policy.ID == id {
			return &policy, nil
		}
	}

	return nil, errors.NewNotFoundf("no alert policy found for id %d", id)
}

// CreatePolicy creates a new alert policy for a given account.
func (alerts *Alerts) CreatePolicy(policy Policy) (*Policy, error) {
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

// UpdatePolicy update an alert policy for a given account.
func (alerts *Alerts) UpdatePolicy(policy Policy) (*Policy, error) {
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

// DeletePolicy deletes an existing alert policy for a given account.
func (alerts *Alerts) DeletePolicy(id int) (*Policy, error) {
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	_, err := alerts.client.Delete(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

type alertPoliciesResponse struct {
	Policies []Policy `json:"policies,omitempty"`
}

type alertPolicyResponse struct {
	Policy Policy `json:"policy,omitempty"`
}

type alertPolicyRequestBody struct {
	Policy Policy `json:"policy"`
}
