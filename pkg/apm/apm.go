package apm

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// APM is used to communicate with the New Relic APM product.
type APM struct {
	client http.NewRelicClient
	pager  http.Pager
}

// New is used to create a new APM client instance.
func New(config config.Config) APM {
	pkg := APM{
		client: http.NewClient(config),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}

// nolint
func newIntegrationTestClient(t *testing.T) APM {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey: apiKey,
	})
}
