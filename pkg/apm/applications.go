package apm

import (
	"github.com/newrelic/newrelic-client-go/internal"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type APM struct {
	client internal.NewRelicClient
}

func New(config newrelic.Config) APM {
	pkg := APM{
		client: internal.NewClient(internal.Config(config)),
	}

	return pkg
}

func (apm *APM) ListApplications() ([]Application, error) {
	res := struct {
		Applications []Application `json:"applications,omitempty"`
	}{}

	_, err := apm.client.Get("applications.json", &res)

	if err != nil {
		return nil, err
	}

	return res.Applications, nil
}
