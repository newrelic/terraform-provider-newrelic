package infrastructure

import "github.com/newrelic/newrelic-client-go/internal/serialization"

// AlertCondition represents a New Relic Infrastructure alert condition.
type AlertCondition struct {
	Comparison          string                  `json:"comparison,omitempty"`
	CreatedAt           serialization.EpochTime `json:"created_at_epoch_millis,omitempty"`
	Critical            *Threshold              `json:"critical_threshold,omitempty"`
	Enabled             bool                    `json:"enabled"`
	Event               string                  `json:"event_type,omitempty"`
	ID                  int                     `json:"id,omitempty"`
	IntegrationProvider string                  `json:"integration_provider,omitempty"`
	Name                string                  `json:"name,omitempty"`
	PolicyID            int                     `json:"policy_id,omitempty"`
	ProcessWhere        string                  `json:"process_where_clause,omitempty"`
	RunbookURL          string                  `json:"runbook_url,omitempty"`
	Select              string                  `json:"select_value,omitempty"`
	Type                string                  `json:"type,omitempty"`
	UpdatedAt           serialization.EpochTime `json:"updated_at_epoch_millis,omitempty"`
	ViolationCloseTimer *int                    `json:"violation_close_timer,omitempty"`
	Warning             *Threshold              `json:"warning_threshold,omitempty"`
	Where               string                  `json:"where_clause,omitempty"`
}

// Threshold represents an New Relic Infrastructure alert condition threshold.
type Threshold struct {
	Duration int    `json:"duration_minutes,omitempty"`
	Function string `json:"time_function,omitempty"`
	Value    int    `json:"value,omitempty"`
}
