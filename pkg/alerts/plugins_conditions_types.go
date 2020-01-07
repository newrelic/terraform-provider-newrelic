package alerts

// AlertPlugin represents a plugin to use with a Plugin alert condition.
type AlertPlugin struct {
	ID   string `json:"id,omitempty"`
	GUID string `json:"guid,omitempty"`
}

// PluginCondition represents an alert condition for New Relic Plugins.
type PluginCondition struct {
	PolicyID          int                  `json:"-"`
	ID                int                  `json:"id,omitempty"`
	Name              string               `json:"name,omitempty"`
	Enabled           bool                 `json:"enabled"`
	Entities          []string             `json:"entities,omitempty"`
	Metric            string               `json:"metric,omitempty"`
	MetricDescription string               `json:"metric_description,omitempty"`
	RunbookURL        string               `json:"runbook_url,omitempty"`
	Terms             []AlertConditionTerm `json:"terms,omitempty"`
	ValueFunction     string               `json:"value_function,omitempty"`
	Plugin            AlertPlugin          `json:"plugin,omitempty"`
}
