package http

import (
	"strings"
)

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data interface{} `json:"data"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLErrorResponse struct {
	Errors []graphQLError `json:"errors"`
}

func (r *graphQLErrorResponse) Error() string {
	if len(r.Errors) > 0 {
		messages := []string{}
		for _, e := range r.Errors {
			messages = append(messages, e.Message)
		}
		return strings.Join(messages, ", ")
	}

	return ""
}

func (r *graphQLErrorResponse) New() ErrorResponse {
	return &graphQLErrorResponse{}
}
