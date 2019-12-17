package config

import (
	"net/http"
	"time"
)

// RegionType represents the members of the Region enumeration.
type RegionType int

const (
	// US represents New Relic's US-based production deployment.
	US = iota

	// EU represents New Relic's EU-based production deployment.
	EU

	// Staging represents New Relic's US-based staging deployment.
	// This is for internal New Relic use only.
	Staging
)

// Region specifies the New Relic environment to target.
var Region = struct {
	US      RegionType
	EU      RegionType
	Staging RegionType
}{
	US:      US,
	EU:      EU,
	Staging: Staging,
}

// Config contains all the configuration data for the API Client.
type ReplacementConfig struct {
	BaseURL       string
	APIKey        string
	Timeout       *time.Duration
	HTTPTransport *http.RoundTripper
	UserAgent     string
	Region        RegionType
}
