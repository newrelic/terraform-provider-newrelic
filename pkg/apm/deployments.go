package apm

import (
	"fmt"
)

// ListDeployments returns deployments an application by application ID.
func (apm *APM) ListDeployments(id int) ([]*Deployment, error) {
	response := deploymentsResponse{}
	deployments := []*Deployment{}
	nextURL := fmt.Sprintf("/applications/%d/deployments.json", id)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, nil, &response)

		if err != nil {
			return nil, err
		}

		deployments = append(deployments, response.Deployments...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return deployments, nil
}

type deploymentsResponse struct {
	Deployments []*Deployment `json:"deployments,omitempty"`
}
