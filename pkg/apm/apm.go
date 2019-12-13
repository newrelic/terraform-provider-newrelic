package apm

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// APM is used to communicate with the New Relic APM product.
type APM struct {
	client http.NewRelicClient
}

// New is used to create a new APM client instance.
func New(config config.Config) APM {
	pkg := APM{
		client: http.NewClient(config),
	}

	return pkg
}
