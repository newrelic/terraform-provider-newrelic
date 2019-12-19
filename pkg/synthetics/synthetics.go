package synthetics

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var baseURLs = map[config.RegionType]string{
	config.Region.US:      "https://synthetics.newrelic.com/synthetics/api/v3",
	config.Region.EU:      "https://synthetics.eu.newrelic.com/synthetics/api/v3",
	config.Region.Staging: "https://staging-synthetics.newrelic.com/synthetics/api/v3",
}

// Synthetics is used to communicate with the New Relic Synthetics product.
type Synthetics struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new Synthetics client instance.
func New(config config.Config) Synthetics {

	if config.BaseURL == "" {
		config.BaseURL = baseURLs[config.Region]
	}

	pkg := Synthetics{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
