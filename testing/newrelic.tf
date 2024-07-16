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

resource "newrelic_browser_application" "foo" {
  name                        = "example-browser-app4"
  distributed_tracing_enabled = false
  cookies_enabled             =  false
}
# resource "newrelic_alert_policy" "gen_servers_alerts_policy" {
#   name = "example-policy1"
#   incident_preference = "PER_POLICY" # PER_POLICY is default
#   account_id = 3957524
# }
# resource "newrelic_infra_alert_condition" "agent_not_reporting_alert" {
#
#   policy_id = newrelic_alert_policy.gen_servers_alerts_policy.id
#   name    = "Host not reporting"
#   description = "Critical alert when the host is not reporting"
#   type    = "infra_host_not_reporting"
#   where    = "(hostname LIKE '%frontend%')"
#
#   critical {
#     duration = 5
#   }
# }