provider "newrelic" {
  account_id = 3342796
  api_key    = "NRAK-0F2WNJV13EEQT404XCWHPS2O1HT"
  region     = "US"
}

resource "newrelic_nrql_alert_condition" "test2" {
  name      = "test after changes"
  type      = "static"
  policy_id = 1746518
  aggregation_method = "event_flow"
  aggregation_delay = 120
  value_function = "single_value"

  nrql {
    query = "From Metric Select count(*)"
  }

  critical {
    threshold = 10000000000
    operator = "above"
    threshold_duration = 300
    threshold_occurrences = "ALL"
  }
}
