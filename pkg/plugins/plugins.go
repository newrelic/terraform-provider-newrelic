package plugins

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Plugins is used to communicate with the New Relic Plugins product.
type Plugins struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new APM client instance.
func New(config config.Config) Plugins {
	pkg := Plugins{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
