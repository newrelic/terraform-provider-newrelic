package http

import (
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// GraphQLClient represents a graphQL HTTP client.
type GraphQLClient struct {
	Client *NewRelicClient
	config config.Config
	logger logging.Logger
}

// NewGraphQLClient returns a new instance of GraphQLClient.
func NewGraphQLClient(cfg config.Config) *GraphQLClient {
	if cfg.NerdGraphBaseURL == "" {
		cfg.NerdGraphBaseURL = region.NerdGraphBaseURLs[region.Parse(cfg.Region)]
	}

	cfg.BaseURL = cfg.NerdGraphBaseURL

	c := NewClient(cfg)
	c.SetErrorValue(&graphQLErrorResponse{})

	return &GraphQLClient{
		Client: &c,
		config: cfg,
		logger: cfg.GetLogger(),
	}
}

// Query runs a graphQL query.
func (g *GraphQLClient) Query(query string, vars map[string]interface{}, respBody interface{}) error {
	req := graphQLRequest{
		Query:     query,
		Variables: vars,
	}

	resp := graphQLResponse{
		Data: respBody,
	}

	if _, err := g.Client.Post(g.config.BaseURL, nil, &req, &resp); err != nil {
		return err
	}

	return nil
}

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
