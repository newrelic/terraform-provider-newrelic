package apm

import (
	client "github.com/newrelic/newrelic-client-go/internal"
)

type APM struct {
	client client.NewRelicClient
}

func New(config client.Config) APM {
	pkg := APM{
		client: client.NewClient(config),
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
