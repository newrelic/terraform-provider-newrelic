package apm

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/http"
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
func (a *APM) ListDeployments(applicationID int) ([]*Deployment, error) {
	deployments := []*Deployment{}
	nextURL := fmt.Sprintf("/applications/%d/deployments.json", applicationID)

	for nextURL != "" {
		response := deploymentsResponse{}
		req, err := http.NewRequest(a.client, "GET", nextURL, nil, nil, &response)
		if err != nil {
			return nil, err
		}

		req.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

		resp, err := a.client.Do(req)

		if err != nil {
			return nil, err
		}

		deployments = append(deployments, response.Deployments...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return deployments, nil
}

// CreateDeployment creates a deployment marker for an application.
func (a *APM) CreateDeployment(applicationID int, deployment Deployment) (*Deployment, error) {
	reqBody := deploymentRequestBody{
		Deployment: deployment,
	}
	resp := deploymentResponse{}

	u := fmt.Sprintf("/applications/%d/deployments.json", applicationID)
	req, err := http.NewRequest(a.client, "POST", u, nil, &reqBody, &resp)
	if err != nil {
		return nil, err
	}

	req.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	_, err = a.client.Do(req)

	if err != nil {
		return nil, err
	}

	return &resp.Deployment, nil
}

// DeleteDeployment deletes a deployment marker for an application.
func (a *APM) DeleteDeployment(applicationID int, deploymentID int) (*Deployment, error) {
	resp := deploymentResponse{}
	u := fmt.Sprintf("/applications/%d/deployments/%d.json", applicationID, deploymentID)

	req, err := http.NewRequest(a.client, "DELETE", u, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	req.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	_, err = a.client.Do(req)

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
