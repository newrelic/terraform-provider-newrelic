package entities

// Tag represents a New Relic One entity tag.
type Tag struct {
	Key    string
	Values []string
}

// TagValue represents a New Relic One entity tag and value pair.
type TagValue struct {
	Key   string
	Value string
}

// Entity represents a New Relic One entity.
type Entity struct {
	AccountID  int
	Domain     EntityDomainType
	EntityType EntityType
	GUID       string
	Name       string
	Permalink  string
	Reporting  bool
	Type       string
}

// EntityType represents a New Relic One entity type.
type EntityType string

var (
	// EntityTypes specifies the possible types for a New Relic One entity.
	EntityTypes = struct {
		Application EntityType
		Dashboard   EntityType
		Host        EntityType
		Monitor     EntityType
	}{
		Application: "APPLICATION",
		Dashboard:   "DASHBOARD",
		Host:        "HOST",
		Monitor:     "MONITOR",
	}
)

// EntityDomainType represents a New Relic One entity domain.
type EntityDomainType string

var (
	// EntityDomains specifies the possible domains for a New Relic One entity.
	EntityDomains = struct {
		APM            EntityDomainType
		Browser        EntityDomainType
		Infrastructure EntityDomainType
		Mobile         EntityDomainType
		Synthetics     EntityDomainType
	}{
		APM:            "APM",
		Browser:        "BROWSER",
		Infrastructure: "INFRA",
		Mobile:         "MOBILE",
		Synthetics:     "SYNTH",
	}
)

// EntityAlertSeverityType represents a New Relic One entity alert severity.
type EntityAlertSeverityType string

var (
	// EntityAlertSeverities specifies the possible alert severities for a New Relic One entity.
	EntityAlertSeverities = struct {
		Critical      EntityAlertSeverityType
		NotAlerting   EntityAlertSeverityType
		NotConfigured EntityAlertSeverityType
		Warning       EntityAlertSeverityType
	}{
		Critical:      "APM",
		NotAlerting:   "NOT_ALERTING",
		NotConfigured: "NOT_CONFIGURED",
		Warning:       "WARNING",
	}
)
