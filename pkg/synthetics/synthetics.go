package synthetics

import (
	"github.com/newrelic/newrelic-client-go/internal"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

const (
	syntheticsBaseURL = "https://synthetics.newrelic.com/synthetics/api/v3"
)

type Synthetics struct {
	client internal.NewRelicClient
}

// New is used to create a new Synthetics client instance.
func New(config newrelic.Config) Synthetics {
	internalConfig := config.ToInternal()

	if internalConfig.BaseURL == "" {
		internalConfig.BaseURL = syntheticsBaseURL
	}

	pkg := Synthetics{
		client: internal.NewClient(internalConfig),
	}

	return pkg
}
