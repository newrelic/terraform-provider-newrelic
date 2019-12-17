package http

import "fmt"

type ErrorResponse interface {
	Error() string
}

type DefaultErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Title string `json:"title"`
}

func (e *DefaultErrorResponse) Error() string {
	return e.ErrorDetail.Title
}

// RestyErrorResponse represents an error response from New Relic.
type RestyErrorResponse struct {
	Detail *ErrorDetail `json:"error,omitempty"`
}

func (e *RestyErrorResponse) Error() string {
	if e != nil && e.Detail != nil {
		return e.Detail.Title
	}
	return "Unknown error"
}

// ErrorDetail represents the details of an ErrorResponse from New Relic.
type RestyErrorDetail struct {
	Title string `json:"title,omitempty"`
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
