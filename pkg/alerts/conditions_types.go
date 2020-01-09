package alerts

// Condition represents a New Relic alert condition.
// TODO: custom unmarshal entities to ints?
type Condition struct {
	PolicyID            int                  `json:"-"`
	ID                  int                  `json:"id,omitempty"`
	Type                string               `json:"type,omitempty"`
	Name                string               `json:"name,omitempty"`
	Enabled             bool                 `json:"enabled"`
	Entities            []string             `json:"entities,omitempty"`
	Metric              string               `json:"metric,omitempty"`
	RunbookURL          string               `json:"runbook_url,omitempty"`
	Terms               []ConditionTerm      `json:"terms,omitempty"`
	UserDefined         ConditionUserDefined `json:"user_defined,omitempty"`
	Scope               string               `json:"condition_scope,omitempty"`
	GCMetric            string               `json:"gc_metric,omitempty"`
	ViolationCloseTimer int                  `json:"violation_close_timer,omitempty"`
}

// ConditionUserDefined represents user defined metrics for the New Relic alert condition.
type ConditionUserDefined struct {
	Metric        string `json:"metric,omitempty"`
	ValueFunction string `json:"value_function,omitempty"`
}

// ConditionTerm represents the terms of a New Relic alert condition.
type ConditionTerm struct {
	Duration     int     `json:"duration,string,omitempty"`
	Operator     string  `json:"operator,omitempty"`
	Priority     string  `json:"priority,omitempty"`
	Threshold    float64 `json:"threshold,string"`
	TimeFunction string  `json:"time_function,omitempty"`
}
