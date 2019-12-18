package apm

import (
	"strconv"
	"strings"
)

// ListApplicationsParams represents a set of filters to be
// used when querying New Relic applications.
type ListApplicationsParams struct {
	Name     *string
	Host     *string
	IDs      []int
	Language *string
}

type listApplicationsResponse struct {
	Applications []Application `json:"applications,omitempty"`
}

// ListApplications is used to retrieve New Relic applications.
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error) {
	response := listApplicationsResponse{}
	apps := []Application{}
	nextURL := apm.client.Config.BaseURL + "/applications.json"
	paramsMap := buildListApplicationsParamsMap(params)

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &paramsMap, &response)

		if err != nil {
			return nil, err
		}

		apps = append(apps, response.Applications...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return apps, nil
}

func buildListApplicationsParamsMap(params *ListApplicationsParams) map[string]string {
	paramsMap := map[string]string{}

	if params != nil {
		if params.Name != nil {
			paramsMap["filter[name]"] = *params.Name
		}

		if params.Host != nil {
			paramsMap["filter[host]"] = *params.Host
		}

		if params.IDs != nil {
			ids := []string{}
			for _, id := range params.IDs {
				ids = append(ids, strconv.Itoa(id))
			}
			paramsMap["filter[ids]"] = strings.Join(ids, ",")
		}

		if params.Language != nil {
			paramsMap["filter[language]"] = *params.Language
		}
	}

	return paramsMap
}
