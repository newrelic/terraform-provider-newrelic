package synthetics

import (
	"github.com/newrelic/newrelic-client-go/internal"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var baseURLs = map[newrelic.Environment]string{
	newrelic.Production: "https://synthetics.newrelic.com/synthetics/api/v3",
	newrelic.EU:         "https://synthetics.eu.newrelic.com/synthetics/api/v3",
	newrelic.Staging:    "https://staging-synthetics.newrelic.com/synthetics/api/v3",
}

type Synthetics struct {
	client internal.NewRelicClient
}

// New is used to create a new Synthetics client instance.
func New(config newrelic.Config) Synthetics {
	internalConfig := config.ToInternal()

	if internalConfig.BaseURL == "" {
		internalConfig.BaseURL = baseURLs[config.Environment]
	}

	pkg := Synthetics{
		client: internal.NewClient(internalConfig),
	}

	return pkg
}
