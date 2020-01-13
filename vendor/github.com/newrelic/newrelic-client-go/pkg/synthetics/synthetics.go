package synthetics

import (
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// BaseURLs represents the base API URLs for the different environments of the Synthetics API.
var BaseURLs = map[region.Region]string{
	region.US:      "https://synthetics.newrelic.com/synthetics/api/v3",
	region.EU:      "https://synthetics.eu.newrelic.com/synthetics/api/v3",
	region.Staging: "https://staging-synthetics.newrelic.com/synthetics/api/v3",
}

// Synthetics is used to communicate with the New Relic Synthetics product.
type Synthetics struct {
	client http.NewRelicClient
	logger logging.Logger
	pager  http.Pager
}

// ErrorResponse represents an error response from New Relic Synthetics.
type ErrorResponse struct {
	Message            string        `json:"error,omitempty"`
	Messages           []ErrorDetail `json:"errors,omitempty"`
	ServerErrorMessage string        `json:"message,omitempty"`
}

// ErrorDetail represents an single error from New Relic Synthetics.
type ErrorDetail struct {
	Message string `json:"error,omitempty"`
}

// Error surfaces an error message from the New Relic Synthetics error response.
func (e *ErrorResponse) Error() string {
	if e.ServerErrorMessage != "" {
		return e.ServerErrorMessage
	}

	if e.Message != "" {
		return e.Message
	}

	if len(e.Messages) > 0 {
		messages := []string{}
		for _, m := range e.Messages {
			messages = append(messages, m.Message)
		}
		return strings.Join(messages, ", ")
	}

	return ""
}

// New is used to create a new Synthetics client instance.
func New(config config.Config) Synthetics {

	if config.BaseURL == "" {
		config.BaseURL = BaseURLs[region.Parse(config.Region)]
	}

	client := http.NewClient(config)
	client.SetErrorValue(&ErrorResponse{})

	pkg := Synthetics{
		client: client,
		logger: config.GetLogger(),
		pager:  &http.LinkHeaderPager{},
	}

	return pkg
}
