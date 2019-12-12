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

func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error) {
	res := listApplicationsResponse{}
	paramsMap := buildListApplicationsParamsMap(params)
	err := apm.client.Get("applications.json", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.Applications, nil
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
