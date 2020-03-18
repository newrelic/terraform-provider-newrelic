package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// IncidentPreferenceType specifies rollup settings for alert policies.
type IncidentPreferenceType string

var (
	// IncidentPreferenceTypes specifies the possible incident preferenece types for an alert policy.
	IncidentPreferenceTypes = struct {
		PerPolicy             IncidentPreferenceType
		PerCondition          IncidentPreferenceType
		PerConditionAndTarget IncidentPreferenceType
	}{
		PerPolicy:             "PER_POLICY",
		PerCondition:          "PER_CONDITION",
		PerConditionAndTarget: "PER_CONDITION_AND_TARGET",
	}
)

// Policy represents a New Relic alert policy.
type Policy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference IncidentPreferenceType   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}

// QueryPolicy is similar to a Policy, but the resulting NerdGraph objects are
// string IDs in the JSON response.
type QueryPolicy struct {
	ID                 int                    `json:"id,string"`
	IncidentPreference IncidentPreferenceType `json:"incidentPreference"`
	Name               string                 `json:"name"`
}

type QueryPolicyInput struct {
	IncidentPreference IncidentPreferenceType `json:"incidentPreference"`
	Name               string                 `json:"name"`
}

type QueryPolicyCreateInput struct {
	QueryPolicyInput
}

type QueryPolicyUpdateInput struct {
	QueryPolicyInput
}

// ListPoliciesParams represents a set of filters to be used when querying New
// Relic alert policies.
type ListPoliciesParams struct {
	Name string `url:"filter[name],omitempty"`
}

// ListPolicies returns a list of Alert Policies for a given account.
func (a *Alerts) ListPolicies(params *ListPoliciesParams) ([]Policy, error) {
	alertPolicies := []Policy{}

	nextURL := "/alerts_policies.json"

	for nextURL != "" {
		response := alertPoliciesResponse{}
		resp, err := a.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		alertPolicies = append(alertPolicies, response.Policies...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertPolicies, nil
}

// GetPolicy returns a specific alert policy by ID for a given account.
func (a *Alerts) GetPolicy(id int) (*Policy, error) {
	policies, err := a.ListPolicies(nil)

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
func (a *Alerts) CreatePolicy(policy Policy) (*Policy, error) {
	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}

	_, err := a.client.Post("/alerts_policies.json", nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// UpdatePolicy update an alert policy for a given account.
func (a *Alerts) UpdatePolicy(policy Policy) (*Policy, error) {

	reqBody := alertPolicyRequestBody{
		Policy: policy,
	}
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", policy.ID)

	_, err := a.client.Put(url, nil, &reqBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeletePolicy deletes an existing alert policy for a given account.
func (a *Alerts) DeletePolicy(id int) (*Policy, error) {
	resp := alertPolicyResponse{}
	url := fmt.Sprintf("/alerts_policies/%d.json", id)

	_, err := a.client.Delete(url, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

func (a *Alerts) CreatePolicyMutation(accountID int, policy QueryPolicyCreateInput) (*QueryPolicy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policy":    policy,
	}

	resp := alertQueryPolicyCreateResponse{}

	if err := a.client.Query(alertsPolicyCreatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.QueryPolicy, nil
}

func (a *Alerts) UpdatePolicyMutation(accountID int, policyID int, policy QueryPolicyUpdateInput) (*QueryPolicy, error) {
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  policyID,
		"policy":    policy,
	}

	resp := alertQueryPolicyUpdateResponse{}

	if err := a.client.Query(alertsPolicyUpdatePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.QueryPolicy, nil
}

// QueryPolicy queries NerdGraph for a policy matching the given account ID and
// policy ID.
func (a *Alerts) QueryPolicy(accountID, id int) (*QueryPolicy, error) {
	resp := alertQueryPolicyResponse{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := a.client.Query(alertPolicyQueryPolicy, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.Alerts.Policy, nil
}

// DeletePolicyMutation is the NerdGraph mutation to delete a policy given the
// account ID and the policy ID.
func (a *Alerts) DeletePolicyMutation(accountID, id int) (*QueryPolicy, error) {
	policy := &QueryPolicy{}

	resp := alertQueryPolicyDeleteRespose{}
	vars := map[string]interface{}{
		"accountID": accountID,
		"policyID":  id,
	}

	if err := a.client.Query(alertPolicyDeletePolicy, vars, &resp); err != nil {
		return nil, err
	}

	return policy, nil
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

type alertQueryPolicyCreateResponse struct {
	QueryPolicy QueryPolicy `json:"alertsPolicyCreate"`
}

type alertQueryPolicyUpdateResponse struct {
	QueryPolicy QueryPolicy `json:"alertsPolicyCreate"`
}

type alertQueryPolicyResponse struct {
	Actor struct {
		Account struct {
			Alerts struct {
				Policy QueryPolicy `json:"policy"`
			} `json:"alerts"`
		} `json:"account"`
	} `json:"actor"`
}

type alertQueryPolicyDeleteRespose struct {
	AlertsPolicyDelete struct {
		ID int `json:"id,string"`
	} `json:"alertsPolicyDelete"`
}

const (
	graphqlAlertPolicyFields = `
						id
						name
						incidentPreference
	`
	alertPolicyQueryPolicy = `query($accountID: Int!, $policyID: ID!) {
		actor {
			account(id: $accountID) {
				alerts {
					policy(id: $policyID) {` + graphqlAlertPolicyFields + `
					}
				}
			}
		}
	}`

	alertsPolicyCreatePolicy = `mutation CreatePolicy($accountID: Int!, $policy: AlertsPolicyInput!){
		alertsPolicyCreate(accountId: $accountID, policy: $policy) {` + graphqlAlertPolicyFields + `
		} }`

	alertsPolicyUpdatePolicy = `mutation UpdatePolicy($accountID: Int!, $policyID: ID!, $policy: AlertsPolicyUpdateInput!){
			alertsPolicyUpdate(accountId: $accountID, id: $policyID, policy: $policy) {` + graphqlAlertPolicyFields + `
			}
		}`

	alertPolicyDeletePolicy = `mutation DeletePolicy($accountID: Int!, $policyID: ID!){
		alertsPolicyDelete(accountId: $accountID, id: $policyID) {
			id
		} }`
)
