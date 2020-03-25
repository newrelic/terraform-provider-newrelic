// Package nerdgraph provides a programmatic API for interacting with NerdGraph, New Relic One's GraphQL API.
package nerdgraph

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

// NerdGraph is used to communicate with the New Relic's GraphQL API, NerdGraph.
type NerdGraph struct {
	client http.Client
	logger logging.Logger
}

// QueryResponse represents the top-level GraphQL response object returned
// from a NerdGraph query request.
type QueryResponse struct {
	Actor          interface{} `json:"actor,omitempty"`
	Docs           interface{} `json:"docs,omitempty"`
	RequestContext interface{} `json:"requestContext,omitempty"`
}

// New returns a new GraphQL client for interacting with New Relic's GraphQL API, NerdGraph.
func New(config config.Config) NerdGraph {
	return NerdGraph{
		client: http.NewClient(config),
		logger: config.GetLogger(),
	}
}

// Query facilitates making a NerdGraph request with a raw GraphQL query. Variables may be provided
// in the form of a map. The response's data structure will vary based on the query provided.
func (n *NerdGraph) Query(query string, variables map[string]interface{}) (interface{}, error) {
	respBody := QueryResponse{}

	if err := n.client.Query(query, variables, &respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func (n *NerdGraph) QuerySchema() (*Schema, error) {

	schemaResponse := allTypesResponse{}
	vars := map[string]interface{}{}
	err := n.client.Query(allTypes, vars, &schemaResponse)
	if err != nil {
		return nil, err
	}

	return &schemaResponse.Schema, nil
}

// AccountReference represents the NerdGraph schema for a New Relic account.
type AccountReference struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

const (
	// https://github.com/graphql/graphql-js/blob/master/src/utilities/getIntrospectionQuery.js#L35
	allTypes = ` query IntrospectionQuery {
      __schema {
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          locations
          args {
            ...InputValue
          }
        }
      }
    }
    fragment FullType on __Type {
      kind
      name
      description
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }
    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
    }
    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
	`
)

// Wheee... :)
type SchemaType struct {
	InputFields []SchemaInputValue `json:"inputFields"`
	Kind        string             `json:"kind"`
	Name        string             `json:"name"`
	// Description string             `json:"description"`
	Fields []struct {
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Args        []SchemaInputValue `json:"args"`
		Type        SchemaTypeRef      `json:"type"`
	} `json:"fields"`
	Interfaces    []SchemaTypeRef `json:"interfaces"`
	PossibleTypes []SchemaTypeRef `json:"possibleTypes"`
	EnumValues    []struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		IsDeprecated      bool   `json:"isDeprecated"`
		DeprecationReason string `json:"deprecationReason"`
	} `json:"enumValues"`
}

type SchemaInputValue struct {
	DefaultValue interface{}   `json:"defaultValue"`
	Description  string        `json:"description"`
	Name         string        `json:"name"`
	Type         SchemaTypeRef `json:"type"`
}

type SchemaTypeRef struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	OfType struct {
		Name   string `json:"name"`
		Kind   string `json:"kind"`
		OfType struct {
			Name   string `json:"name"`
			Kind   string `json:"kind"`
			OfType struct {
				Name   string `json:"name"`
				Kind   string `json:"kind"`
				OfType struct {
					Name   string `json:"name"`
					Kind   string `json:"kind"`
					OfType struct {
						Name   string `json:"name"`
						Kind   string `json:"kind"`
						OfType struct {
							Name   string `json:"name"`
							Kind   string `json:"kind"`
							OfType struct {
								Name string `json:"name"`
								Kind string `json:"kind"`
							} `json:"ofType"`
						} `json:"ofType"`
					} `json:"ofType"`
				} `json:"ofType"`
			} `json:"ofType"`
		} `json:"ofType"`
	} `json:"ofType"`
}

type Schema struct {
	// TODO Implement the rest of the schema if needed.

	Types []*SchemaType `json:"types"`
}

type allTypesResponse struct {
	Schema Schema `json:"__schema"`
}
