package alerts

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/infrastructure"
)

// Alerts is used to communicate with New Relic Alerts.
type Alerts struct {
	client      http.Client
	infraClient http.Client
	logger      logging.Logger
	pager       http.Pager
}

// New is used to create a new Alerts client instance.
func New(config config.Config) Alerts {
	infraConfig := config

	if infraConfig.InfrastructureBaseURL == "" {
		infraConfig.InfrastructureBaseURL = infrastructure.BaseURLs[region.Parse(string(config.Region))]
	}

	infraConfig.BaseURL = infraConfig.InfrastructureBaseURL

	infraClient := http.NewClient(infraConfig)
	infraClient.SetErrorValue(&infrastructure.ErrorResponse{})

	client := http.NewClient(config)
	client.SetAuthStrategy(&http.PersonalAPIKeyCapableV2Authorizer{})

	pkg := Alerts{
		client:      client,
		infraClient: infraClient,
		logger:      config.GetLogger(),
		pager:       &http.LinkHeaderPager{},
	}

	return pkg
}

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = region.DefaultBaseURLs
