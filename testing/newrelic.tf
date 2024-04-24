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

resource "newrelic_alert_policy" "policy" {
  name = "my-policy"
}

resource "newrelic_synthetics_monitor" "monitor" {
  locations_public = ["US_WEST_1"]
  name             = "my-monitor"
  period           = "EVERY_10_MINUTES"
  status           = "DISABLED"
  type             = "SIMPLE"
  uri              = "https://www.one.newrelic.com"
}

resource "newrelic_synthetics_multilocation_alert_condition" "example" {
  policy_id = newrelic_alert_policy.policy.id

  name                         = "Example condition"
  runbook_url                  = "https://example.com"
  enabled                      = true
  violation_time_limit_seconds = 3600

  entities = [
    1234
  ]

  critical {
    threshold = 2
  }

  warning {
    threshold = 1
  }
}