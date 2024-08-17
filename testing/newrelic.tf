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

resource "newrelic_synthetics_broken_links_monitor" "foo" {
  name                 = "Sample Broken Links Monitor"
  uri                  = "https://www.one.example.com"
  locations_public     = ["AP_SOUTH_1"]
  period               = "EVERY_6_HOURS"
  status               = "ENABLED"
#   runtime_type         = "NODE_API"
#   runtime_type_version = "16.10"
  use_legacy_runtime_unsupported = true
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
