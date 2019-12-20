package alerts

// SyntheticsCondition represents a New Relic Synthetics alert condition.
type SyntheticsCondition struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Enabled    bool   `json:"enabled"`
	RunbookURL string `json:"runbook_url,omitempty"`
	MonitorID  string `json:"monitor_id,omitempty"`
}
