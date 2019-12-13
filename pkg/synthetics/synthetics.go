package synthetics

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var baseURLs = map[config.Region]string{
	config.Production: "https://synthetics.newrelic.com/synthetics/api/v3",
	config.EU:         "https://synthetics.eu.newrelic.com/synthetics/api/v3",
	config.Staging:    "https://staging-synthetics.newrelic.com/synthetics/api/v3",
}

type Synthetics struct {
	client http.NewRelicClient
}

// New is used to create a new Synthetics client instance.
func New(config config.Config) Synthetics {

	if config.BaseURL == "" {
		config.BaseURL = baseURLs[config.Region]
	}

	pkg := Synthetics{
		client: http.NewClient(config),
	}

	return pkg
}
