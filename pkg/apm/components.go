package apm

import (
	"fmt"
	"strconv"
	"strings"
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

func buildListComponentsParamsMap(params *ListComponentsParams) map[string]string {
	paramsMap := map[string]string{}

	if params != nil {
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
	}

	return paramsMap
}

type componentsResponse struct {
	Components []Component `json:"components,omitempty"`
}

type componentResponse struct {
	Component Component `json:"component,omitempty"`
}
