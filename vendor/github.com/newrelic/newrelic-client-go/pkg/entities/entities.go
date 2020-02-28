package entities

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// Entities is used to communicate with the New Relic Entities product.
type Entities struct {
	client *http.GraphQLClient
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One entities.
func New(config config.Config) Entities {
	return Entities{
		client: http.NewGraphQLClient(config),
		logger: config.GetLogger(),
	}
}

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = region.NerdGraphBaseURLs
