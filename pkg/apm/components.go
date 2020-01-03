package apm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/http"
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
func (apm *APM) ListComponents(params *ListComponentsParams) ([]*Component, error) {
	response := componentsResponse{}
	c := []*Component{}
	nextURL := "/components.json"
	queryParams := buildListComponentsQueryParams(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &queryParams, &response)

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
func (apm *APM) ListComponentMetrics(componentID int, params *ListComponentMetricsParams) ([]*ComponentMetric, error) {
	m := []*ComponentMetric{}
	response := componentMetricsResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics.json", componentID)
	queryParams := buildListComponentMetricsQueryParams(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &queryParams, &response)

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
func (apm *APM) GetComponentMetricData(componentID int, params *GetComponentMetricDataParams) ([]*Metric, error) {
	m := []*Metric{}
	response := componentMetricDataResponse{}
	nextURL := fmt.Sprintf("/components/%d/metrics/data.json", componentID)
	queryParams := buildGetComponentMetricDataQueryParams(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		m = append(m, response.MetricData.Metrics...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return m, nil
}

func buildListComponentMetricsQueryParams(params *ListComponentMetricsParams) []http.QueryParam {
	queryParams := []http.QueryParam{}

	if params == nil {
		return queryParams
	}

	if params.Name != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "name", Value: params.Name})
	}

	return queryParams
}

func buildListComponentsQueryParams(params *ListComponentsParams) []http.QueryParam {
	queryParams := []http.QueryParam{}

	if params == nil {
		return queryParams
	}

	if params.Name != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[name]", Value: params.Name})

		if params.IDs != nil {
			ids := []string{}
			for _, id := range params.IDs {
				ids = append(ids, strconv.Itoa(id))
			}

			value := strings.Join(ids, ",")
			queryParams = append(queryParams, http.QueryParam{Name: "filter[ids]", Value: value})
		}
	}

	if params.PluginID != 0 {
		value := strconv.Itoa(params.PluginID)
		queryParams = append(queryParams, http.QueryParam{Name: "filter[plugin_id]", Value: value})
	}

	value := strconv.FormatBool(params.HealthStatus)
	queryParams = append(queryParams, http.QueryParam{Name: "health_status", Value: value})

	return queryParams
}

func buildGetComponentMetricDataQueryParams(params *GetComponentMetricDataParams) []http.QueryParam {
	queryParams := []http.QueryParam{}

	if params == nil {
		return queryParams
	}

	if len(params.Names) > 0 {
		for _, name := range params.Names {
			queryParams = append(queryParams, http.QueryParam{Name: "names[]", Value: name})
		}
	}

	if len(params.Values) > 0 {
		for _, value := range params.Values {
			queryParams = append(queryParams, http.QueryParam{Name: "values[]", Value: value})
		}
	}

	if params.From != nil {
		value := params.From.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "from", Value: value})
	}

	if params.To != nil {
		value := params.From.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "to", Value: value})
	}

	if params.Period != 0 {
		value := strconv.Itoa(params.Period)
		queryParams = append(queryParams, http.QueryParam{Name: "period", Value: value})
	}

	value := strconv.FormatBool(params.Summarize)
	queryParams = append(queryParams, http.QueryParam{Name: "summarize", Value: value})

	value = strconv.FormatBool(params.Raw)
	queryParams = append(queryParams, http.QueryParam{Name: "raw", Value: value})

	return queryParams
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
