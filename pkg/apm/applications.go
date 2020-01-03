package apm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/http"
)

// ListApplicationsParams represents a set of filters to be
// used when querying New Relic applications.
type ListApplicationsParams struct {
	Name     string
	Host     string
	IDs      []int
	Language string
}

// ListApplications is used to retrieve New Relic applications.
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error) {
	response := applicationsResponse{}
	apps := []Application{}
	nextURL := "/applications.json"
	queryParams := buildListApplicationsQueryParams(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &queryParams, &response)

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

// UpdateApplicationParams represents a set of parameters to be
// used when updating New Relic applications.
type UpdateApplicationParams struct {
	Name     string
	Settings ApplicationSettings
}

// UpdateApplication is used to update a New Relic application's name and/or settings.
func (apm *APM) UpdateApplication(applicationID int, params UpdateApplicationParams) (*Application, error) {
	response := applicationResponse{}
	url := fmt.Sprintf("/applications/%d.json", applicationID)
	reqBody := updateApplicationRequest{
		Fields: updateApplicationFields(params),
	}

	_, err := apm.client.Put(url, nil, reqBody, &response)

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

func buildListApplicationsQueryParams(params *ListApplicationsParams) []http.QueryParam {
	queryParams := []http.QueryParam{}

	if params == nil {
		return queryParams
	}

	if params.Name != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[name]", Value: params.Name})
	}

	if params.Host != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[host]", Value: params.Host})
	}

	if params.IDs != nil {
		ids := []string{}
		for _, id := range params.IDs {
			ids = append(ids, strconv.Itoa(id))
		}

		values := strings.Join(ids, ",")
		queryParams = append(queryParams, http.QueryParam{Name: "filter[ids]", Value: values})
	}

	if params.Language != "" {
		queryParams = append(queryParams, http.QueryParam{Name: "filter[language]", Value: params.Language})
	}

	return queryParams
}

type applicationsResponse struct {
	Applications []Application `json:"applications,omitempty"`
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
