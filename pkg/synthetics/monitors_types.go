package synthetics

import "time"

// MonitorOptions represents the options for a New Relic Synthetics monitor.
type MonitorOptions struct {
	ValidationString       string `json:"validationString,omitempty"`
	VerifySSL              bool   `json:"verifySSL,omitempty"`
	BypassHEADRequest      bool   `json:"bypassHEADRequest,omitempty"`
	TreatRedirectAsFailure bool   `json:"treatRedirectAsFailure,omitempty"`
}

// Monitor represents a New Relic Synthetics monitor.
type Monitor struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Frequency    uint           `json:"frequency"`
	URI          string         `json:"uri"`
	Locations    []string       `json:"locations"`
	Status       string         `json:"status"`
	SLAThreshold float64        `json:"slaThreshold"`
	UserID       uint           `json:"userId,omitempty"`
	APIVersion   string         `json:"apiVersion,omitempty"`
	ModifiedAt   time.Time      `json:"modified_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	Options      MonitorOptions `json:"options,omitempty"`
}
