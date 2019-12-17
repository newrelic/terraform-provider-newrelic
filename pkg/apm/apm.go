package apm

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// APM is used to communicate with the New Relic APM product.
type APM struct {
	client http.ReplacementClient
	pager  http.Pager
}

// New is used to create a new APM client instance.
func New(config config.ReplacementConfig) APM {
	pkg := APM{
		client: http.NewReplacementClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
