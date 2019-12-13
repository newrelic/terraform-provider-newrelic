package infrastructure

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var baseURLs = map[config.RegionType]string{
	config.Region.US:      "https://infra-api.newrelic.com/v2/alerts/conditions",
	config.Region.EU:      "https://infra-api.eu.newrelic.com/v2/alerts/conditions",
	config.Region.Staging: "https://staging-infra-api.newrelic.com/v2/alerts/conditions",
}

// Infrastructure is used to communicate with the New Relic Infrastructure product.
type Infrastructure struct {
	client http.NewRelicClient
}

// ErrorResponse represents an error response from New Relic Infrastructure.
type ErrorResponse struct {
	Errors []*ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail represents the details of an error response from New Relic Infrastructure.
type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Error surfaces an error message from the Infrastructure error response.
func (e *ErrorResponse) Error() string {
	if e != nil && len(e.Errors) > 0 && e.Errors[0].Detail != "" {
		return e.Errors[0].Detail
	}
	return "Unknown error"
}

// New is used to create a new Infrastructure client instance.
func New(config config.Config) Infrastructure {
	if config.BaseURL == "" {
		config.BaseURL = baseURLs[config.Region]
	}

	c := http.NewClient(config)
	c.Client.SetError(&ErrorResponse{})

	pkg := Infrastructure{
		client: c,
	}

	return pkg
}
