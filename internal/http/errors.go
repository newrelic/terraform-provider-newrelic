package http

import "fmt"

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
	Title string `json:"title"`
}

func (e *DefaultErrorResponse) Error() string {
	return e.ErrorDetail.Title
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
