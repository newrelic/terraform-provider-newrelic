// Package region describes the operational regions defined for New Relic
//
// Regions are geographical locations where the New Relic platform operates
// and this package provides an abstraction layer for handling them within
// the New Relic Client and underlying APIs
package region

import (
	"strings"
)

// Region represents the members of the Region enumeration.
type Region int

const (
	_ = iota // Ignore zero

	// US represents New Relic's US-based production deployment.
	US Region = iota

	// EU represents New Relic's EU-based production deployment.
	EU

	// Staging represents New Relic's US-based staging deployment.
	// This is for internal New Relic use only.
	Staging
)

// Default represents the region returned if nothing was specified
const Default Region = US

// DefaultBaseURLs represents the base API URLs for the different environments of the New Relic REST API V2.
var DefaultBaseURLs = map[Region]string{
	US:      "https://api.newrelic.com/v2",
	EU:      "https://api.eu.newrelic.com/v2",
	Staging: "https://staging-api.newrelic.com/v2",
}

// NerdGraphBaseURLs describes the base URLs for the NerdGraph API.
var NerdGraphBaseURLs = map[Region]string{
	US:      "https://api.newrelic.com/graphql",
	EU:      "https://api.eu.newrelic.com/graphql",
	Staging: "https://staging-api.newrelic.com/graphql",
}

// Parse takes a Region string and returns a RegionType
func Parse(r string) Region {
	switch strings.ToLower(r) {
	case "us":
		return US
	case "eu":
		return EU
	case "staging":
		return Staging
	default:
		return Default
	}
}

// String returns a human readable value for the specified Region
func (r Region) String() string {
	switch r {
	case US:
		return "US"
	case EU:
		return "EU"
	case Staging:
		return "Staging"
	default:
		return "(Unknown)"
	}
}
