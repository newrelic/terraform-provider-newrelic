package apm

import (
	"github.com/newrelic/newrelic-client-go/internal"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type APM struct {
	client internal.NewRelicClient
}

// New is used to create a new APM client instance.
func New(config newrelic.Config) APM {
	internalConfig := config.ToInternal()

	pkg := APM{
		client: internal.NewClient(internalConfig),
	}

	return pkg
}
