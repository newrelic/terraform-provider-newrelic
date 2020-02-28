package apm

import (
	"fmt"
)

// Application represents information about a New Relic application.
type Application struct {
	ID             int                       `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Language       string                    `json:"language,omitempty"`
	HealthStatus   string                    `json:"health_status,omitempty"`
	Reporting      bool                      `json:"reporting"`
	LastReportedAt string                    `json:"last_reported_at,omitempty"`
	Summary        ApplicationSummary        `json:"application_summary,omitempty"`
	EndUserSummary ApplicationEndUserSummary `json:"end_user_summary,omitempty"`
	Settings       ApplicationSettings       `json:"settings,omitempty"`
	Links          ApplicationLinks          `json:"links,omitempty"`
}

// ApplicationSummary represents performance information about a New Relic application.
type ApplicationSummary struct {
	ResponseTime            float64 `json:"response_time"`
	Throughput              float64 `json:"throughput"`
	ErrorRate               float64 `json:"error_rate"`
	ApdexTarget             float64 `json:"apdex_target"`
	ApdexScore              float64 `json:"apdex_score"`
	HostCount               int     `json:"host_count"`
	InstanceCount           int     `json:"instance_count"`
	ConcurrentInstanceCount int     `json:"concurrent_instance_count"`
}

// ApplicationEndUserSummary represents performance information about a New Relic application.
type ApplicationEndUserSummary struct {
	ResponseTime float64 `json:"response_time"`
	Throughput   float64 `json:"throughput"`
	ApdexTarget  float64 `json:"apdex_target"`
	ApdexScore   float64 `json:"apdex_score"`
}

// ApplicationSettings represents some of the settings of a New Relic application.
type ApplicationSettings struct {
	AppApdexThreshold        float64 `json:"app_apdex_threshold,omitempty"`
	EndUserApdexThreshold    float64 `json:"end_user_apdex_threshold,omitempty"`
	EnableRealUserMonitoring bool    `json:"enable_real_user_monitoring"`
	UseServerSideConfig      bool    `json:"use_server_side_config"`
}

// ApplicationLinks represents all the links for a New Relic application.
type ApplicationLinks struct {
	ServerIDs     []int `json:"servers,omitempty"`
	HostIDs       []int `json:"application_hosts,omitempty"`
	InstanceIDs   []int `json:"application_instances,omitempty"`
	AlertPolicyID int   `json:"alert_policy"`
}

// ListApplicationsParams represents a set of filters to be
// used when querying New Relic applications.
type ListApplicationsParams struct {
	Name     string `url:"filter[name],omitempty"`
	Host     string `url:"filter[host],omitempty"`
	IDs      []int  `url:"filter[ids],omitempty,comma"`
	Language string `url:"filter[language],omitempty"`
}

// UpdateApplicationParams represents a set of parameters to be
// used when updating New Relic applications.
type UpdateApplicationParams struct {
	Name     string
	Settings ApplicationSettings
}

// ListApplications is used to retrieve New Relic applications.
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]*Application, error) {
	apps := []*Application{}
	nextURL := "/applications.json"

	for nextURL != "" {
		response := applicationsResponse{}
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		apps = append(apps, response.Applications...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return apps, nil
}

// GetApplication is used to retrieve a single New Relic application.
func (apm *APM) GetApplication(applicationID int) (*Application, error) {
	response := applicationResponse{}
	url := fmt.Sprintf("/applications/%d.json", applicationID)

	_, err := apm.client.Get(url, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Application, nil
}

// UpdateApplication is used to update a New Relic application's name and/or settings.
func (apm *APM) UpdateApplication(applicationID int, params UpdateApplicationParams) (*Application, error) {
	response := applicationResponse{}
	url := fmt.Sprintf("/applications/%d.json", applicationID)
	reqBody := updateApplicationRequest{
		Fields: updateApplicationFields(params),
	}

	_, err := apm.client.Put(url, nil, &reqBody, &response)

	if err != nil {
		return nil, err
	}

	return &response.Application, nil
}

// DeleteApplication is used to delete a New Relic application.
// This process will only succeed if the application is no longer reporting data.
func (apm *APM) DeleteApplication(applicationID int) (*Application, error) {
	response := applicationResponse{}
	url := fmt.Sprintf("/applications/%d.json", applicationID)

	_, err := apm.client.Delete(url, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Application, nil
}

type applicationsResponse struct {
	Applications []*Application `json:"applications,omitempty"`
}

type applicationResponse struct {
	Application Application `json:"application,omitempty"`
}

type updateApplicationRequest struct {
	Fields updateApplicationFields `json:"application"`
}

type updateApplicationFields struct {
	Name     string              `json:"name,omitempty"`
	Settings ApplicationSettings `json:"settings,omitempty"`
}
