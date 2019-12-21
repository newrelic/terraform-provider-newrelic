package infrastructure

import (
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// BaseURLs describes the base URLs for the Infrastructure Alerts API.
var BaseURLs = map[config.RegionType]string{
	config.Region.US:      "https://infra-api.newrelic.com/v2",
	config.Region.EU:      "https://infra-api.eu.newrelic.com/v2",
	config.Region.Staging: "https://staging-infra-api.newrelic.com/v2",
}

// ErrorResponse represents an error response from New Relic Infrastructure.
type ErrorResponse struct {
	Errors  []*ErrorDetail `json:"errors,omitempty"`
	Message string         `json:"description,omitempty"`
}

// ErrorDetail represents the details of an error response from New Relic Infrastructure.
type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Error surfaces an error message from the Infrastructure error response.
func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if len(e.Errors) > 0 && e.Errors[0].Detail != "" {
		return e.Errors[0].Detail
	}
	return "Unknown error"
}
