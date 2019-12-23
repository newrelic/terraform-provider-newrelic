package alerts

// AlertCondition represents a New Relic alert condition.
// TODO: custom unmarshal entities to ints?
type AlertCondition struct {
	PolicyID            int                       `json:"-"`
	ID                  int                       `json:"id,omitempty"`
	Type                string                    `json:"type,omitempty"`
	Name                string                    `json:"name,omitempty"`
	Enabled             bool                      `json:"enabled"`
	Entities            []string                  `json:"entities,omitempty"`
	Metric              string                    `json:"metric,omitempty"`
	RunbookURL          string                    `json:"runbook_url,omitempty"`
	Terms               []AlertConditionTerm      `json:"terms,omitempty"`
	UserDefined         AlertConditionUserDefined `json:"user_defined,omitempty"`
	Scope               string                    `json:"condition_scope,omitempty"`
	GCMetric            string                    `json:"gc_metric,omitempty"`
	ViolationCloseTimer int                       `json:"violation_close_timer,omitempty"`
}

// AlertConditionUserDefined represents user defined metrics for the New Relic alert condition.
type AlertConditionUserDefined struct {
	Metric        string `json:"metric,omitempty"`
	ValueFunction string `json:"value_function,omitempty"`
}

// AlertConditionTerm represents the terms of a New Relic alert condition.
type AlertConditionTerm struct {
	Duration     int     `json:"duration,string,omitempty"`
	Operator     string  `json:"operator,omitempty"`
	Priority     string  `json:"priority,omitempty"`
	Threshold    float64 `json:"threshold,string"`
	TimeFunction string  `json:"time_function,omitempty"`
}
