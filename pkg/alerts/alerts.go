package alerts

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Alerts is used to communicate with New Relic Alerts.
type Alerts struct {
	client http.NewRelicClient
}

// New is used to create a new Alerts client instance.
func New(config config.Config) Alerts {
	pkg := Alerts{
		client: http.NewClient(config),
	}

	return pkg
}
