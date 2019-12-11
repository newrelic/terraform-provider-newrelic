package client

// ErrorResponse represents an error response from New Relic.
type ErrorResponse struct {
	Detail *ErrorDetail `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e != nil && e.Detail != nil {
		return e.Detail.Title
	}
	return "Unknown error"
}

// ErrorDetail represents the details of an ErrorResponse from New Relic.
type ErrorDetail struct {
	Title string `json:"title,omitempty"`
}
