terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

resource "newrelic_synthetics_script_monitor" "appliance" {
  count = 10
  name   = "Test GH Synth Script Monitor ${1000+count.index}"
  type   = "SCRIPT_API"
  period = "EVERY_5_MINUTES"
  status = "ENABLED"
  locations_public = [
    "AP_EAST_1", // Hong Kong
    "AP_SOUTH_1",     // Mumbai
    "AP_SOUTHEAST_1", // Singapore
    "AP_NORTHEAST_2", // Seoul
    "AP_NORTHEAST_1", // Tokyo
    "AP_SOUTHEAST_2", // Sydney
    "US_WEST_1",      // San Francisco
    "US_WEST_2", // Portland
    "US_EAST_2",      // Columbus
    "US_EAST_1",      // Washington
    "CA_CENTRAL_1",   // Montreal
    "SA_EAST_1",      // São Paulo
    "EU_WEST_1",      // Dublin
    "EU_WEST_2",      // London
    "EU_WEST_3",      // Paris
    "EU_CENTRAL_1",   // Frankfurt
    "EU_NORTH_1",     // Stockholm
    "EU_SOUTH_1",     // Milan
    "ME_SOUTH_1",     // Manama (Bahrain)
    "AF_SOUTH_1",     // Cape Town (South Africa)
  ]

  script = "console.log('it works!')"
  script_language      = "JAVASCRIPT"
  runtime_type         = "NODE_API"
  runtime_type_version = "16.10"
}
