package apm

import (
	client "github.com/newrelic/newrelic-client-go/internal"
)

type APM struct{
	client client.NewRelicClient
}

func New(config client.Config) APM {
	pkg := APM{
		client: client.NewClient(config),
	}

	return pkg
}

func (apm *APM) ListApplications() {

}
