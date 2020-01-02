package apm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ListComponentsParams represents a set of filters to be
// used when querying New Relic applications.
type ListComponentsParams struct {
	Name         string
	IDs          []int
	PluginID     int
	HealthStatus bool
}

// ListComponents is used to retrieve the components associated with
// a New Relic account.
func (apm *APM) ListComponents(params *ListComponentsParams) ([]Component, error) {
	response := componentsResponse{}
	c := []Component{}
	nextURL := "/components.json"
	paramsMap := buildListComponentsParamsMap(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &paramsMap, &response)

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
	Name string
}

// ListComponentMetrics is used to retrieve the metrics for a specific New Relic component.
func (apm *APM) ListComponentMetrics(componentID int, params *ListComponentMetricsParams) ([]ComponentMetric, error) {
	m := []ComponentMetric{}
	response := componentMetricsResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics.json", componentID)
	paramsMap := buildListComponentMetricsParamsMap(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &paramsMap, &response)

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
	Names []string

	// Values allows retrieval of specific metric values.
	Values []string

	// From specifies a begin time for the query.
	From *time.Time

	// To specifies an end time for the query.
	To *time.Time

	// Period represents the period of timeslices in seconds.
	Period int

	// Summarize will summarize the data when set to true.
	Summarize bool

	// Raw will return unformatted raw values when set to true.
	Raw bool
}

// GetComponentMetricData is used to retrieve the metric timeslice data for a specific component metric.
func (apm *APM) GetComponentMetricData(componentID int, params *GetComponentMetricDataParams) ([]Metric, error) {
	m := []Metric{}
	response := componentMetricDataResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics/data.json", componentID)
	paramsMap := buildGetComponentMetricDataParamsMap(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &paramsMap, &response)

		if err != nil {
			return nil, err
		}

		m = append(m, response.MetricData.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return m, nil
}

func buildListComponentMetricsParamsMap(params *ListComponentMetricsParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.Name != "" {
		paramsMap["name"] = params.Name
	}

	return paramsMap
}

func buildListComponentsParamsMap(params *ListComponentsParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.Name != "" {
		paramsMap["filter[name]"] = params.Name

		if params.IDs != nil {
			ids := []string{}
			for _, id := range params.IDs {
				ids = append(ids, strconv.Itoa(id))
			}
			paramsMap["filter[ids]"] = strings.Join(ids, ",")
		}
	}

	if params.PluginID != 0 {
		paramsMap["filter[plugin_id]"] = strconv.Itoa(params.PluginID)
	}

	paramsMap["health_status"] = strconv.FormatBool(params.HealthStatus)

	return paramsMap
}

func buildGetComponentMetricDataParamsMap(params *GetComponentMetricDataParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if len(params.Names) > 0 {
		paramsMap["names[]"] = strings.Join(params.Names, ",")
	}

	if len(params.Values) > 0 {
		paramsMap["values[]"] = strings.Join(params.Values, ",")
	}

	if params.From != nil {
		paramsMap["from"] = params.From.Format(time.RFC3339)
	}

	if params.To != nil {
		paramsMap["to"] = params.From.Format(time.RFC3339)
	}

	if params.Period != 0 {
		paramsMap["period"] = strconv.Itoa(params.Period)
	}

	paramsMap["summarize"] = strconv.FormatBool(params.Summarize)
	paramsMap["raw"] = strconv.FormatBool(params.Raw)

	return paramsMap
}

type componentsResponse struct {
	Components []Component `json:"components,omitempty"`
}

type componentResponse struct {
	Component Component `json:"component,omitempty"`
}

type componentMetricsResponse struct {
	Metrics []ComponentMetric `json:"metrics,omitempty"`
}

type componentMetricDataResponse struct {
	MetricData struct {
		Metrics []Metric `json:"metrics,omitempty"`
	} `json:"metric_data,omitempty"`
}
