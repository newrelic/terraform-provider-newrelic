package region

import (
	"strings"
)

const (
	// US represents New Relic's US-based production deployment.
	US Name = "US"

	// EU represents New Relic's EU-based production deployment.
	EU Name = "EU"

	// Staging represents New Relic's US-based staging deployment.
	// This is for internal New Relic use only.
	Staging Name = "Staging"
)

var Regions = map[Name]*Region{
	US: {
		name:                  "US",
		infrastructureBaseURL: "https://infra-api.newrelic.com/v2",
		nerdGraphBaseURL:      "https://api.newrelic.com/graphql",
		restBaseURL:           "https://api.newrelic.com/v2",
		syntheticsBaseURL:     "https://synthetics.newrelic.com/synthetics/api",
	},
	EU: {
		name:                  "EU",
		infrastructureBaseURL: "https://infra-api.eu.newrelic.com/v2",
		nerdGraphBaseURL:      "https://api.eu.newrelic.com/graphql",
		restBaseURL:           "https://api.eu.newrelic.com/v2",
		syntheticsBaseURL:     "https://synthetics.eu.newrelic.com/synthetics/api",
	},
	Staging: {
		name:                  "Staging",
		infrastructureBaseURL: "https://staging-infra-api.newrelic.com/v2",
		nerdGraphBaseURL:      "https://staging-api.newrelic.com/graphql",
		restBaseURL:           "https://staging-api.newrelic.com/v2",
		syntheticsBaseURL:     "https://staging-synthetics.newrelic.com/synthetics/api",
	},
}

// Default represents the region returned if nothing was specified
var Default *Region = Regions[US]

// Parse takes a Region string and returns a RegionType
func Parse(r string) *Region {
	var ret Region

	switch strings.ToLower(r) {
	case "us":
		ret = *Regions[US]
	case "eu":
		ret = *Regions[EU]
	case "staging":
		ret = *Regions[Staging]
	default:
		ret = *Default
	}

	return &ret
}
