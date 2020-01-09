package alerts

import "github.com/newrelic/newrelic-client-go/internal/serialization"

// Policy represents a New Relic alert policy.
type Policy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference string                   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}
