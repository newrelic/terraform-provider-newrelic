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

resource "newrelic_events_to_metrics_rule" "tech-talk" {
  account_id = "2520528"
  name = "tech-talk test"
  description = "test description"
  nrql = "SELECT uniqueCount(account_id) AS `Transaction.account_id` FROM Transaction FACET appName, name"
  enabled = false
}