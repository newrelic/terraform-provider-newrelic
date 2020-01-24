package entities

import (
	"github.com/newrelic/newrelic-client-go/internal/http"
	"github.com/newrelic/newrelic-client-go/internal/logging"
	"github.com/newrelic/newrelic-client-go/internal/region"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Entities is used to communicate with the New Relic Entities product.
type Entities struct {
	client *http.GraphQLClient
	logger logging.Logger
}

// New returns a new client for interacting with New Relic One entities.
func New(config config.Config) Entities {
	return Entities{
		client: http.NewGraphQLClient(config),
		logger: config.GetLogger(),
	}
}

// SearchEntitiesParams represents a set of search parameters for retrieving New Relic One entities.
type SearchEntitiesParams struct {
	AlertSeverity                 EntityAlertSeverityType `json:"alertSeverity,omitempty"`
	Domain                        EntityDomainType        `json:"domain,omitempty"`
	InfrastructureIntegrationType string                  `json:"infrastructureIntegrationType,omitempty"`
	Name                          string                  `json:"name,omitempty"`
	Reporting                     *bool                   `json:"reporting,omitempty"`
	Tags                          []Tag                   `json:"tags,omitempty"`
	Type                          EntityType              `json:"type,omitempty"`
}

// SearchEntities searches New Relic One entities based on the provided search parameters.
func (e *Entities) SearchEntities(params SearchEntitiesParams) ([]*Entity, error) {
	entities := []*Entity{}
	var nextCursor *string

	for ok := true; ok; ok = nextCursor != nil {
		resp := searchEntitiesResponse{}
		vars := map[string]interface{}{
			"queryBuilder": params,
			"cursor":       nextCursor,
		}

		if err := e.client.Query(searchEntitiesQuery, vars, &resp); err != nil {
			return nil, err
		}

		entities = append(entities, resp.Actor.EntitySearch.Results.Entities...)

		nextCursor = resp.Actor.EntitySearch.Results.NextCursor
	}

	return entities, nil
}

// GetEntities retrieves a set of New Relic One entities by their entity guids.
func (e *Entities) GetEntities(guids []string) ([]*Entity, error) {
	resp := getEntitiesResponse{}
	vars := map[string]interface{}{
		"guids": guids,
	}

	if err := e.client.Query(getEntitiesQuery, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.Actor.Entities) == 0 {
		return nil, errors.NewNotFound("")
	}

	return resp.Actor.Entities, nil
}

// GetEntity retrieve a set of New Relic One entities by their entity guids.
func (e *Entities) GetEntity(guid string) (*Entity, error) {
	resp := getEntityResponse{}
	vars := map[string]interface{}{
		"guid": guid,
	}

	if err := e.client.Query(getEntityQuery, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Actor.Entity == nil {
		return nil, errors.NewNotFound("")
	}

	return resp.Actor.Entity, nil
}

var searchEntitiesQuery = `
    query($queryBuilder: EntitySearchQueryBuilder, $cursor: String) {
        actor {
            entitySearch(queryBuilder: $queryBuilder)  {
                results(cursor: $cursor) {
                    nextCursor
                    entities {
						accountId
						domain
						entityType
						guid
						name
						permalink
						reporting
						type
                    }
                }
            }
        }
    }
`

type searchEntitiesResponse struct {
	Actor struct {
		EntitySearch struct {
			Results struct {
				NextCursor *string
				Entities   []*Entity
			}
		}
	}
}

var getEntitiesQuery = `
    query($guids: [String!]!) {
        actor {
            entities(guids: $guids)  {
				accountId
				domain
				entityType
				guid
				name
				permalink
				reporting
				type
            }
        }
    }
`

type getEntitiesResponse struct {
	Actor struct {
		Entities []*Entity
	}
}

var getEntityQuery = `
    query($guid: String!) {
        actor {
            entity(guid: $guid)  {
				accountId
				domain
				entityType
				guid
				name
				permalink
				reporting
				type
            }
        }
    }
`

type getEntityResponse struct {
	Actor struct {
		Entity *Entity
	}
}

// BaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var BaseURLs = region.NerdGraphBaseURLs
