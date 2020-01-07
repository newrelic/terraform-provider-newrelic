package plugins

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Plugins is used to communicate with the New Relic Plugins product.
type Plugins struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new Plugins client instance.
func New(config config.Config) Plugins {
	pkg := Plugins{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// ListPluginsParams represents a set of query string parameters
// used as filters when querying New Relic plugins.
type ListPluginsParams struct {
	GUID string `url:"filter[guid],omitempty"`
	IDs  []int  `url:"filter[ids],omitempty,comma"`
}

// ListPlugins returns a list of Plugins associated with an account.
func (plugins *Plugins) ListPlugins(params *ListPluginsParams) ([]*Plugin, error) {
	response := pluginsResponse{}
	results := []*Plugin{}
	nextURL := "/plugins.json"

	for nextURL != "" {
		resp, err := plugins.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		results = append(results, response.Plugins...)

		paging := plugins.pager.Parse(resp)
		nextURL = paging.Next
	}

	return results, nil
}

// GetPluginParams represents a set of query string parameters
// to apply to the request.
type GetPluginParams struct {
	Detailed bool `url:"detailed,omitempty"`
}

// GetPlugin returns a plugin for a given account. If the query paramater `detailed=true`
// is provided, the response will contain an additional `details` property that contains
// additional metadata pertaining to the plugin.
func (plugins *Plugins) GetPlugin(id int, params *GetPluginParams) (*Plugin, error) {
	response := pluginResponse{}

	u := fmt.Sprintf("/plugins/%d.json", id)
	_, err := plugins.client.Get(u, &params, &response)

	if err != nil {
		return nil, err
	}

	return &response.Plugin, nil
}

type pluginsResponse struct {
	Plugins []*Plugin `json:"plugins,omitempty"`
}

type pluginResponse struct {
	Plugin Plugin `json:"plugin,omitempty"`
}
