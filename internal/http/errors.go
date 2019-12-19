package http

import (
	"fmt"
	"strings"
)

// ErrorResponse provides an interface for obtaining
// a single error message from an error response object.
type ErrorResponse interface {
	Error() string
}

// DefaultErrorResponse represents the default error response from New Relic.
type DefaultErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}

// ErrorDetail represents a New Relic response error detail.
type ErrorDetail struct {
	Title    string   `json:"title"`
	Messages []string `json:"messages"`
}

func (e *DefaultErrorResponse) Error() string {
	m := e.ErrorDetail.Title
	if len(e.ErrorDetail.Messages) > 0 {
		m = fmt.Sprintf("%s: %s", m, strings.Join(e.ErrorDetail.Messages, ", "))
	}

	return m
}

// ErrorNotFound is returned when a 404 response is returned
// from New Relic's APIs.
type ErrorNotFound struct{}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("404 not found")
}

// ErrorUnexpectedStatusCode is returned when an unexpected
// status code is returned from New Relic's APIs.
type ErrorUnexpectedStatusCode struct {
	err        string
	statusCode int
}

func (e *ErrorUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("%d response returned: %s", e.statusCode, e.err)
}
