package http

import "fmt"

type ErrorResponse interface {
	Error() string
}

// DefaultErrorResponse represents the default error response from New Relic.
type DefaultErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Title string `json:"title"`
}

func (e *DefaultErrorResponse) Error() string {
	return e.ErrorDetail.Title
}

type ErrorNotFound struct{}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("404 not found")
}

type ErrorUnexpectedStatusCode struct {
	err        string
	statusCode int
}

func (e *ErrorUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("%d response returned: %s", e.statusCode, e.err)
}
