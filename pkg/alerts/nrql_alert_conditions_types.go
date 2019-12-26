package alerts

// NrqlQuery represents a NRQL query to use with a NRQL alert condition
type NrqlQuery struct {
	Query      string `json:"query,omitempty"`
	SinceValue string `json:"since_value,omitempty"`
}

// NrqlCondition represents a New Relic NRQL Alert condition.
type NrqlCondition struct {
	Terms               []AlertConditionTerm `json:"terms,omitempty"`
	Nrql                NrqlQuery            `json:"nrql,omitempty"`
	Type                string               `json:"type,omitempty"`
	Name                string               `json:"name,omitempty"`
	RunbookURL          string               `json:"runbook_url,omitempty"`
	ValueFunction       string               `json:"value_function,omitempty"`
	PolicyID            int                  `json:"-"`
	ID                  int                  `json:"id,omitempty"`
	ViolationCloseTimer int                  `json:"violation_time_limit_seconds,omitempty"`
	ExpectedGroups      int                  `json:"expected_groups,omitempty"`
	IgnoreOverlap       bool                 `json:"ignore_overlap,omitempty"`
	Enabled             bool                 `json:"enabled"`
}
