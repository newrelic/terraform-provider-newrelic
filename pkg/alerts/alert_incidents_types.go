package alerts

// AlertIncident represents a New Relic alert incident.
type AlertIncident struct {
	ID                 int               `json:"id,omitempty"`
	OpenedAt           int               `json:"opened_at,omitempty"`
	ClosedAt           int               `json:"closed_at,omitempty"`
	IncidentPreference string            `json:"incident_preference,omitempty"`
	Links              AlertIncidentLink `json:"links"`
}

// AlertIncidentLink represents a link between a New Relic alert incident and its violations
type AlertIncidentLink struct {
	Violations []int `json:"violations,omitempty"`
	PolicyID   int   `json:"policy_id"`
}


