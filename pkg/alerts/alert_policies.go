package alerts

import "fmt"

// ListAlertPoliciesParams represents a set of filters to be
// used when querying New Relic alert policies.
type ListAlertPoliciesParams struct {
	Name *string
}

type listAlertsResponse struct {
	AlertPolicies []AlertPolicy `json:"policies,omitempty"`
}

// GetAlertPolicy returns a specific alert policy by ID
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

// ListAlertPolicies returns a list of Alert Policies for a given account
func (alerts *Alerts) ListAlertPolicies(params *ListAlertPoliciesParams) ([]AlertPolicy, error) {
	res := listAlertsResponse{}
	paramsMap := buildListAlertPoliciesParamsMap(params)
	err := alerts.client.Get("/alerts_policies.json", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.AlertPolicies, nil
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
