package alerts

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/infrastructure"
)

// Alerts is used to communicate with New Relic Alerts.
type Alerts struct {
	client      http.NewRelicClient
	infraClient http.NewRelicClient
	pager       http.Pager
}

// New is used to create a new Alerts client instance.
func New(config config.Config) Alerts {
	infraConfig := config

	if config.BaseURL == "" {
		infraConfig.BaseURL = infrastructure.BaseURLs[config.Region]
	}

	infraClient := http.NewClient(infraConfig)
	infraClient.SetErrorValue(&infrastructure.ErrorResponse{})

	pkg := Alerts{
		client:      http.NewClient(config),
		infraClient: infraClient,
		pager:       &http.LinkHeaderPager{},
	}

	return pkg
}
