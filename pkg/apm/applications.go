package apm

import (
	"github.com/newrelic/newrelic-client-go/internal"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type APM struct {
	client internal.NewRelicClient
}

func New(config newrelic.Config) APM {
	internalConfig := config.ToInternal()

	pkg := APM{
		client: internal.NewClient(internalConfig),
	}

	return pkg
}

func (apm *APM) ListApplications() ([]Application, error) {
	res := struct {
		Applications []Application `json:"applications,omitempty"`
	}{}

	err := apm.client.Get("applications.json", &res)

	if err != nil {
		return nil, err
	}

	return res.Applications, nil
}
