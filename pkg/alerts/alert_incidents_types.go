package alerts

// Incident represents a New Relic alert incident.
type Incident struct {
	ID                 int          `json:"id,omitempty"`
	OpenedAt           int          `json:"opened_at,omitempty"`
	ClosedAt           int          `json:"closed_at,omitempty"`
	IncidentPreference string       `json:"incident_preference,omitempty"`
	Links              IncidentLink `json:"links"`
}

// IncidentLink represents a link between a New Relic alert incident and its violations
type IncidentLink struct {
	Violations []int `json:"violations,omitempty"`
	PolicyID   int   `json:"policy_id"`
}
