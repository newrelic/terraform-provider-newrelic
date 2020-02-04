package alerts

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/serialization"
)

// Channel represents a New Relic alert notification channel
type Channel struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	Type          string               `json:"type,omitempty"`
	Configuration ChannelConfiguration `json:"configuration,omitempty"`
	Links         ChannelLinks         `json:"links,omitempty"`
}

// ChannelLinks represent the links between policies and alert channels
type ChannelLinks struct {
	PolicyIDs []int `json:"policy_ids,omitempty"`
}

// ChannelConfiguration represents a Configuration type within Channels
type ChannelConfiguration struct {
	Recipients            string `json:"recipients,omitempty"`
	IncludeJSONAttachment string `json:"include_json_attachment,omitempty"`
	AuthToken             string `json:"auth_token,omitempty"`
	APIKey                string `json:"api_key,omitempty"`
	Teams                 string `json:"teams,omitempty"`
	Tags                  string `json:"tags,omitempty"`
	URL                   string `json:"url,omitempty"`
	Channel               string `json:"channel,omitempty"`
	Key                   string `json:"key,omitempty"`
	RouteKey              string `json:"route_key,omitempty"`
	ServiceKey            string `json:"service_key,omitempty"`
	BaseURL               string `json:"base_url,omitempty"`
	AuthUsername          string `json:"auth_username,omitempty"`
	AuthPassword          string `json:"auth_password,omitempty"`
	PayloadType           string `json:"payload_type,omitempty"`
	Region                string `json:"region,omitempty"`
	UserID                string `json:"user_id,omitempty"`

	// Payload is unmarshaled to type map[string]string
	Payload MapStringInterface `json:"payload,omitempty"`

	// Headers is unmarshaled to type map[string]string
	Headers MapStringInterface `json:"headers,omitempty"`
}

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

// InfrastructureCondition represents a New Relic Infrastructure alert condition.
type InfrastructureCondition struct {
	Comparison          string                            `json:"comparison,omitempty"`
	CreatedAt           *serialization.EpochTime          `json:"created_at_epoch_millis,omitempty"`
	Critical            *InfrastructureConditionThreshold `json:"critical_threshold,omitempty"`
	Enabled             bool                              `json:"enabled"`
	Event               string                            `json:"event_type,omitempty"`
	ID                  int                               `json:"id,omitempty"`
	IntegrationProvider string                            `json:"integration_provider,omitempty"`
	Name                string                            `json:"name,omitempty"`
	PolicyID            int                               `json:"policy_id,omitempty"`
	ProcessWhere        string                            `json:"process_where_clause,omitempty"`
	RunbookURL          string                            `json:"runbook_url,omitempty"`
	Select              string                            `json:"select_value,omitempty"`
	Type                string                            `json:"type,omitempty"`
	UpdatedAt           *serialization.EpochTime          `json:"updated_at_epoch_millis,omitempty"`
	ViolationCloseTimer *int                              `json:"violation_close_timer,omitempty"`
	Warning             *InfrastructureConditionThreshold `json:"warning_threshold,omitempty"`
	Where               string                            `json:"where_clause,omitempty"`
}

// InfrastructureConditionThreshold represents an New Relic Infrastructure alert condition threshold.
type InfrastructureConditionThreshold struct {
	Duration int     `json:"duration_minutes,omitempty"`
	Function string  `json:"time_function,omitempty"`
	Value    float64 `json:"value,omitempty"`
}

// NrqlQuery represents a NRQL query to use with a NRQL alert condition
type NrqlQuery struct {
	Query      string `json:"query,omitempty"`
	SinceValue string `json:"since_value,omitempty"`
}

// NrqlCondition represents a New Relic NRQL Alert condition.
type NrqlCondition struct {
	Terms               []ConditionTerm `json:"terms,omitempty"`
	Nrql                NrqlQuery       `json:"nrql,omitempty"`
	Type                string          `json:"type,omitempty"`
	Name                string          `json:"name,omitempty"`
	RunbookURL          string          `json:"runbook_url,omitempty"`
	ValueFunction       string          `json:"value_function,omitempty"`
	PolicyID            int             `json:"-"`
	ID                  int             `json:"id,omitempty"`
	ViolationCloseTimer int             `json:"violation_time_limit_seconds,omitempty"`
	ExpectedGroups      int             `json:"expected_groups,omitempty"`
	IgnoreOverlap       bool            `json:"ignore_overlap,omitempty"`
	Enabled             bool            `json:"enabled"`
}

// AlertPlugin represents a plugin to use with a Plugin alert condition.
type AlertPlugin struct {
	ID   string `json:"id,omitempty"`
	GUID string `json:"guid,omitempty"`
}

// PluginsCondition represents an alert condition for New Relic Plugins.
type PluginsCondition struct {
	PolicyID          int             `json:"-"`
	ID                int             `json:"id,omitempty"`
	Name              string          `json:"name,omitempty"`
	Enabled           bool            `json:"enabled"`
	Entities          []string        `json:"entities,omitempty"`
	Metric            string          `json:"metric,omitempty"`
	MetricDescription string          `json:"metric_description,omitempty"`
	RunbookURL        string          `json:"runbook_url,omitempty"`
	Terms             []ConditionTerm `json:"terms,omitempty"`
	ValueFunction     string          `json:"value_function,omitempty"`
	Plugin            AlertPlugin     `json:"plugin,omitempty"`
}

// Policy represents a New Relic alert policy.
type Policy struct {
	ID                 int                      `json:"id,omitempty"`
	IncidentPreference string                   `json:"incident_preference,omitempty"`
	Name               string                   `json:"name,omitempty"`
	CreatedAt          *serialization.EpochTime `json:"created_at,omitempty"`
	UpdatedAt          *serialization.EpochTime `json:"updated_at,omitempty"`
}

// PolicyChannels represents an association of alert channels to a specific alert policy.
type PolicyChannels struct {
	ID         int   `json:"id,omitempty"`
	ChannelIDs []int `json:"channel_ids,omitempty"`
}

// SyntheticsCondition represents a New Relic Synthetics alert condition.
type SyntheticsCondition struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Enabled    bool   `json:"enabled"`
	RunbookURL string `json:"runbook_url,omitempty"`
	MonitorID  string `json:"monitor_id,omitempty"`
}

// MapStringInterface is used for custom unmarshaling of
// fields that have potentially dynamic types.
// E.g. when a field can be a string or an object/map
type MapStringInterface map[string]interface{}
type mapStringInterfaceProxy MapStringInterface

// UnmarshalJSON is a custom unmarshal method to guard against
// fields that can have more than one type returned from an API.
func (c *MapStringInterface) UnmarshalJSON(data []byte) error {
	var mapStrInterface mapStringInterfaceProxy

	str := string(data)

	// Check for empty JSON string
	if str == `""` {
		return nil
	}

	// Remove quotes if this is a string representation of JSON
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		s, err := strconv.Unquote(str)
		if err != nil {
			return nil
		}

		data = []byte(s)
	}

	err := json.Unmarshal(data, &mapStrInterface)
	if err != nil {
		return err
	}

	*c = MapStringInterface(mapStrInterface)

	return nil
}
