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
	queryParams := buildListDashboardsQueryParams(params)

	for nextURL != "" {
		resp, err := dashboards.client.Get(nextURL, &queryParams, &response)

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

func buildListDashboardsQueryParams(params *ListDashboardsParams) []http.QueryParam {
	queryParams := []http.QueryParam{}

	if params == nil {
		return queryParams
	}

	if params.Title != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[title]", Value: params.Title})
	}

	if params.Category != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[category]", Value: params.Category})
	}

	if params.CreatedBefore != nil {
		value := params.CreatedBefore.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "filter[created_before]", Value: value})
	}

	if params.CreatedAfter != nil {
		value := params.CreatedAfter.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "filter[created_after]", Value: value})
	}

	if params.UpdatedBefore != nil {
		value := params.UpdatedBefore.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "filter[updated_before]", Value: value})
	}

	if params.UpdatedAfter != nil {
		value := params.UpdatedAfter.Format(time.RFC3339)
		queryParams = append(queryParams, http.QueryParam{Name: "filter[updated_after]", Value: value})
	}

	if params.Sort != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "sort", Value: params.Sort})
	}

	if params.Page > 0 {
		value := strconv.Itoa(params.Page)
		queryParams = append(queryParams, http.QueryParam{Name: "page", Value: value})
	}

	if params.PerPage > 0 {
		value := strconv.Itoa(params.PerPage)
		queryParams = append(queryParams, http.QueryParam{Name: "per_page", Value: value})
	}

	return queryParams
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
