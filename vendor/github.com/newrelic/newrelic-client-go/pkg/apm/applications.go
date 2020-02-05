package apm

import (
	"fmt"
	"time"
)

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

// GetMetricNames is used to retrieve a list of known metrics and their value names for the given resource.
//
// https://rpm.newrelic.com/api/explore/applications/metric_names
func (apm *APM) GetMetricNames(applicationID int, params MetricNamesParams) ([]*MetricName, error) {
	response := metricNamesResponse{}
	metrics := []*MetricName{}
	nextURL := fmt.Sprintf("/applications/%d/metrics.json", applicationID)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		metrics = append(metrics, response.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return metrics, nil
}

// GetMetricData is used to retrieve a list of values for each of the requested metrics.
//
// https://rpm.newrelic.com/api/explore/applications/metric_data
func (apm *APM) GetMetricData(applicationID int, params MetricDataParams) ([]*MetricData, error) {
	response := metricDataResponse{}
	data := []*MetricData{}
	nextURL := fmt.Sprintf("/applications/%d/metrics/data.json", applicationID)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		data = append(data, response.MetricData.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return data, nil
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

type metricNamesResponse struct {
	Metrics []*MetricName
}

type metricDataResponse struct {
	MetricData struct {
		From            *time.Time    `json:"from"`
		To              *time.Time    `json:"to"`
		MetricsNotFound []string      `json:"metrics_not_found"`
		MetricsFound    []string      `json:"metrics_found"`
		Metrics         []*MetricData `json:"metrics"`
	} `json:"metric_data"`
}
