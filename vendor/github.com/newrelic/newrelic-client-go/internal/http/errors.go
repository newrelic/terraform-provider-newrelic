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
