package synthetics

// MonitorScriptLocation represents a New Relic Synthetics monitor script location.
type MonitorScriptLocation struct {
	Name string `json:"name"`
	HMAC string `json:"hmac"`
}

// MonitorScript represents a New Relic Synthetics monitor script.
type MonitorScript struct {
	Text      string                  `json:"scriptText"`
	Locations []MonitorScriptLocation `json:"scriptLocations"`
}
