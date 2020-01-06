package plugins

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Plugins is used to communicate with the New Relic Plugins product.
type Plugins struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new APM client instance.
func New(config config.Config) Plugins {
	pkg := Plugins{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// ListPluginsParams represents a set of filters to be
// used when querying New Relic key transactions.
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

type pluginsResponse struct {
	Plugins []*Plugin `json:"plugins,omitempty"`
}
