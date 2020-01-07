package errors

import (
	"fmt"
)

// ErrorNotFound is returned when a 404 response is returned
// from New Relic's APIs.
type ErrorNotFound struct{}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("404 not found")
}

// ErrorUnexpectedStatusCode is returned when an unexpected
// status code is returned from New Relic's APIs.
type ErrorUnexpectedStatusCode struct {
	Err        string
	StatusCode int
}

func (e *ErrorUnexpectedStatusCode) Error() string {

	msg := fmt.Sprintf("%d response returned", e.StatusCode)

	if e.Err != "" {
		msg += fmt.Sprintf(": %s", e.Err)
	}

	return msg
}
