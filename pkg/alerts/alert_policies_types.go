package alerts

import "github.com/newrelic/newrelic-client-go/internal/serialization"

// AlertPolicy represents a New Relic alert policy.
type AlertPolicy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference string                   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}
