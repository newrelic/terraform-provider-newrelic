package dashboards

import (
	"fmt"
	"strconv"
	"time"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Dashboards is used to communicate with the New Relic Dashboards product.
type Dashboards struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new Dashboards client instance.
func New(config config.Config) Dashboards {
	pkg := Dashboards{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// ListDashboardsParams represents a set of filters to be
// used when querying New Relic dashboards.
type ListDashboardsParams struct {
	Category      string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Page          int
	PerPage       int
	Sort          string
	Title         string
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}

// ListDashboards is used to retrieve New Relic dashboards.
func (dashboards *Dashboards) ListDashboards(params *ListDashboardsParams) ([]Dashboard, error) {
	response := dashboardsResponse{}
	d := []Dashboard{}
	nextURL := "/dashboards.json"
	paramsMap := buildListDashboardsParamsMap(params)

	for nextURL != "" {
		resp, err := dashboards.client.Get(nextURL, &paramsMap, &response)

		if err != nil {
			return nil, err
		}

		d = append(d, response.Dashboards...)

		paging := dashboards.pager.Parse(resp)
		nextURL = paging.Next
	}

	return d, nil
}

// GetDashboard is used to retrieve a single New Relic dashboard.
func (dashboards *Dashboards) GetDashboard(dashboardID int) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboardID)

	_, err := dashboards.client.Get(url, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// CreateDashboard is used to create a New Relic dashboard.
func (dashboards *Dashboards) CreateDashboard(dashboard Dashboard) (*Dashboard, error) {
	response := dashboardResponse{}
	reqBody := dashboardRequest{
		Dashboard: dashboard,
	}
	_, err := dashboards.client.Post("/dashboards.json", nil, &reqBody, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// UpdateDashboard is used to update a New Relic dashboard.
func (dashboards *Dashboards) UpdateDashboard(dashboard Dashboard) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboard.ID)
	reqBody := dashboardRequest{
		Dashboard: dashboard,
	}

	_, err := dashboards.client.Put(url, nil, reqBody, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

// DeleteDashboard is used to delete a New Relic dashboard.
func (dashboards *Dashboards) DeleteDashboard(dashboardID int) (*Dashboard, error) {
	response := dashboardResponse{}
	url := fmt.Sprintf("/dashboards/%d.json", dashboardID)

	_, err := dashboards.client.Delete(url, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Dashboard, nil
}

func buildListDashboardsParamsMap(params *ListDashboardsParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.Title != "" {
		paramsMap["filter[title]"] = params.Title
	}

	if params.Category != "" {
		paramsMap["filter[category]"] = params.Category
	}

	if params.CreatedBefore != nil {
		paramsMap["filter[created_before]"] = params.CreatedBefore.Format(time.RFC3339)
	}

	if params.CreatedAfter != nil {
		paramsMap["filter[created_after]"] = params.CreatedAfter.Format(time.RFC3339)
	}

	if params.UpdatedBefore != nil {
		paramsMap["filter[updated_before]"] = params.UpdatedBefore.Format(time.RFC3339)
	}

	if params.UpdatedAfter != nil {
		paramsMap["filter[updated_after]"] = params.UpdatedAfter.Format(time.RFC3339)
	}

	if params.Sort != "" {
		paramsMap["sort"] = params.Sort
	}

	if params.Page > 0 {
		paramsMap["page"] = strconv.Itoa(params.Page)
	}

	if params.PerPage > 0 {
		paramsMap["per_page"] = strconv.Itoa(params.PerPage)
	}

	return paramsMap
}

type dashboardsResponse struct {
	Dashboards []Dashboard `json:"dashboards,omitempty"`
}

type dashboardResponse struct {
	Dashboard Dashboard `json:"dashboard,omitempty"`
}

type dashboardRequest struct {
	Dashboard Dashboard `json:"dashboard"`
}
