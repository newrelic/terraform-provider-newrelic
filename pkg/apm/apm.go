package apm

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// APM is used to communicate with the New Relic APM product.
type APM struct {
	client http.NewRelicClient
	logger logging.Logger
	pager  http.Pager
}

// New is used to create a new APM client instance.
func New(config config.Config) APM {
	pkg := APM{
		client: http.NewClient(config),
		logger: config.GetLogger(),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = config.DefaultBaseURLs
