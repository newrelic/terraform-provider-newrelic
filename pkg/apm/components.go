package apm

import (
	"fmt"
	"time"
)

// ListComponentsParams represents a set of filters to be
// used when querying New Relic applications.
type ListComponentsParams struct {
	Name         string `url:"filter[name],omitempty"`
	IDs          []int  `url:"filter[ids],omitempty,comma"`
	PluginID     int    `url:"filter[plugin_id],omitempty"`
	HealthStatus bool   `url:"health_status,omitempty"`
}

// ListComponents is used to retrieve the components associated with
// a New Relic account.
func (apm *APM) ListComponents(params *ListComponentsParams) ([]*Component, error) {
	response := componentsResponse{}
	c := []*Component{}
	nextURL := "/components.json"

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		c = append(c, response.Components...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return c, nil

}

// GetComponent is used to retrieve a specific New Relic component.
func (apm *APM) GetComponent(componentID int) (*Component, error) {
	response := componentResponse{}
	url := fmt.Sprintf("/components/%d.json", componentID)

	_, err := apm.client.Get(url, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Component, nil
}

// ListComponentMetricsParams represents a set of parameters to be
// used when querying New Relic component metrics.
type ListComponentMetricsParams struct {
	// Name allows for filtering the returned list of metrics by name.
	Name string `url:"name,omitempty"`
}

// ListComponentMetrics is used to retrieve the metrics for a specific New Relic component.
func (apm *APM) ListComponentMetrics(componentID int, params *ListComponentMetricsParams) ([]*ComponentMetric, error) {
	m := []*ComponentMetric{}
	response := componentMetricsResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics.json", componentID)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		m = append(m, response.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return m, nil
}

// GetComponentMetricDataParams represents a set of parameters to be
// used when querying New Relic component metric data.
type GetComponentMetricDataParams struct {
	// Names allows retrieval of specific metrics by name.
	// At least one metric name is required.
	Names []string `url:"names[],omitempty"`

	// Values allows retrieval of specific metric values.
	Values []string `url:"values[],omitempty"`

	// From specifies a begin time for the query.
	From *time.Time `url:"from,omitempty"`

	// To specifies an end time for the query.
	To *time.Time `url:"to,omitempty"`

	// Period represents the period of timeslices in seconds.
	Period int `url:"period,omitempty"`

	// Summarize will summarize the data when set to true.
	Summarize bool `url:"summarize,omitempty"`

	// Raw will return unformatted raw values when set to true.
	Raw bool `url:"raw,omitempty"`
}

// GetComponentMetricData is used to retrieve the metric timeslice data for a specific component metric.
func (apm *APM) GetComponentMetricData(componentID int, params *GetComponentMetricDataParams) ([]*Metric, error) {
	m := []*Metric{}
	response := componentMetricDataResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics/data.json", componentID)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		m = append(m, response.MetricData.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return m, nil
}

type componentsResponse struct {
	Components []*Component `json:"components,omitempty"`
}

type componentResponse struct {
	Component Component `json:"component,omitempty"`
}

type componentMetricsResponse struct {
	Metrics []*ComponentMetric `json:"metrics,omitempty"`
}

type componentMetricDataResponse struct {
	MetricData struct {
		Metrics []*Metric `json:"metrics,omitempty"`
	} `json:"metric_data,omitempty"`
}
