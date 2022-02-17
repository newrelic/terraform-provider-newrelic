terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {

  region = "US" # US or EU
}

resource "newrelic_alert_policy" "name" {
  name = "name"
  account_id = 2520528
  incident_preference = "PER_POLICY"
}

resource "newrelic_nrql_alert_condition" "cond" {
  name      = "nrql"
  policy_id = newrelic_alert_policy.name.id
  type = "static"
  value_function = "single_value"
  nrql {
    query = "SELECT average(duration) FROM Transaction where appName = 'Your App'"
  }
  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }
}