package dashboards

import (
	"fmt"
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

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = config.DefaultBaseURLs

// ListDashboardsParams represents a set of filters to be
// used when querying New Relic dashboards.
type ListDashboardsParams struct {
	Category      string     `url:"filter[category],omitempty"`
	CreatedAfter  *time.Time `url:"filter[created_after],omitempty"`
	CreatedBefore *time.Time `url:"filter[created_before],omitempty"`
	Page          int        `url:"page,omitempty"`
	PerPage       int        `url:"per_page,omitempty"`
	Sort          string     `url:"sort,omitempty"`
	Title         string     `url:"filter[title],omitempty"`
	UpdatedAfter  *time.Time `url:"filter[updated_after],omitempty"`
	UpdatedBefore *time.Time `url:"filter[updated_before],omitempty"`
}

// ListDashboards is used to retrieve New Relic dashboards.
func (dashboards *Dashboards) ListDashboards(params *ListDashboardsParams) ([]*Dashboard, error) {
	response := dashboardsResponse{}
	d := []*Dashboard{}
	nextURL := "/dashboards.json"

	for nextURL != "" {
		resp, err := dashboards.client.Get(nextURL, &params, &response)

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

	_, err := dashboards.client.Put(url, nil, &reqBody, &response)

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

type dashboardsResponse struct {
	Dashboards []*Dashboard `json:"dashboards,omitempty"`
}

type dashboardResponse struct {
	Dashboard Dashboard `json:"dashboard,omitempty"`
}

type dashboardRequest struct {
	Dashboard Dashboard `json:"dashboard"`
}
