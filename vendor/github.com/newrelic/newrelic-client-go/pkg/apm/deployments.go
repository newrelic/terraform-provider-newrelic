package apm

import (
	"fmt"
)

// Deployment represents information about a New Relic application deployment.
type Deployment struct {
	Links       *DeploymentLinks `json:"links,omitempty"`
	ID          int              `json:"id,omitempty"`
	Revision    string           `json:"revision"`
	Changelog   string           `json:"changelog,omitempty"`
	Description string           `json:"description,omitempty"`
	User        string           `json:"user,omitempty"`
	Timestamp   string           `json:"timestamp,omitempty"`
}

// DeploymentLinks contain the application ID for the deployment.
type DeploymentLinks struct {
	ApplicationID int `json:"application,omitempty"`
}

// ListDeployments returns deployments for an application.
func (apm *APM) ListDeployments(applicationID int) ([]*Deployment, error) {
	deployments := []*Deployment{}
	nextURL := fmt.Sprintf("/applications/%d/deployments.json", applicationID)

	for nextURL != "" {
		response := deploymentsResponse{}
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

// CreateDeployment creates a deployment marker for an application.
func (apm *APM) CreateDeployment(applicationID int, deployment Deployment) (*Deployment, error) {
	reqBody := deploymentRequestBody{
		Deployment: deployment,
	}
	resp := deploymentResponse{}

	u := fmt.Sprintf("/applications/%d/deployments.json", applicationID)
	_, err := apm.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Deployment, nil
}

// DeleteDeployment deletes a deployment marker for an application.
func (apm *APM) DeleteDeployment(applicationID int, deploymentID int) (*Deployment, error) {
	resp := deploymentResponse{}
	u := fmt.Sprintf("/applications/%d/deployments/%d.json", applicationID, deploymentID)

	_, err := apm.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Deployment, nil
}

type deploymentsResponse struct {
	Deployments []*Deployment `json:"deployments,omitempty"`
}

type deploymentResponse struct {
	Deployment Deployment `json:"deployment,omitempty"`
}

type deploymentRequestBody struct {
	Deployment Deployment `json:"deployment,omitempty"`
}
