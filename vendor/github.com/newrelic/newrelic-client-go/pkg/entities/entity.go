package entities

import (
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// Entity represents a New Relic One entity.
type Entity struct {
	AccountID     int
	ApplicationID int
	Domain        EntityDomainType
	EntityType    EntityType
	GUID          string
	Name          string
	Permalink     string
	Reporting     bool
	Type          string
}

// EntityType represents a New Relic One entity type.
type EntityType string

var (
	// EntityTypes specifies the possible types for a New Relic One entity.
	EntityTypes = struct {
		Application EntityType
		Dashboard   EntityType
		Host        EntityType
		Monitor     EntityType
	}{
		Application: "APPLICATION",
		Dashboard:   "DASHBOARD",
		Host:        "HOST",
		Monitor:     "MONITOR",
	}
)

// EntityDomainType represents a New Relic One entity domain.
type EntityDomainType string

var (
	// EntityDomains specifies the possible domains for a New Relic One entity.
	EntityDomains = struct {
		APM            EntityDomainType
		Browser        EntityDomainType
		Infrastructure EntityDomainType
		Mobile         EntityDomainType
		Synthetics     EntityDomainType
	}{
		APM:            "APM",
		Browser:        "BROWSER",
		Infrastructure: "INFRA",
		Mobile:         "MOBILE",
		Synthetics:     "SYNTH",
	}
)

// EntityAlertSeverityType represents a New Relic One entity alert severity.
type EntityAlertSeverityType string

var (
	// EntityAlertSeverities specifies the possible alert severities for a New Relic One entity.
	EntityAlertSeverities = struct {
		Critical      EntityAlertSeverityType
		NotAlerting   EntityAlertSeverityType
		NotConfigured EntityAlertSeverityType
		Warning       EntityAlertSeverityType
	}{
		Critical:      "APM",
		NotAlerting:   "NOT_ALERTING",
		NotConfigured: "NOT_CONFIGURED",
		Warning:       "WARNING",
	}
)

// SearchEntitiesParams represents a set of search parameters for retrieving New Relic One entities.
type SearchEntitiesParams struct {
	AlertSeverity                 EntityAlertSeverityType `json:"alertSeverity,omitempty"`
	Domain                        EntityDomainType        `json:"domain,omitempty"`
	InfrastructureIntegrationType string                  `json:"infrastructureIntegrationType,omitempty"`
	Name                          string                  `json:"name,omitempty"`
	Reporting                     *bool                   `json:"reporting,omitempty"`
	Tags                          *TagValue               `json:"tags,omitempty"`
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

const (
	getEntitiesQuery = `
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
	getEntityQuery = `
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
	searchEntitiesQuery = `
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
          ... on ApmApplicationEntityOutline {
              applicationId
          }
        }
      }
    }
  }
}`
)

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

type getEntitiesResponse struct {
	Actor struct {
		Entities []*Entity
	}
}

type getEntityResponse struct {
	Actor struct {
		Entity *Entity
	}
}
