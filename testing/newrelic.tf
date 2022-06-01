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

resource "newrelic_synthetics_monitor" "foo" {
  custom_headers{
    name="Name"
    value="simpleMonitor"
  }
  treat_redirect_as_failure=true
  validation_string="success"
  bypass_head_request=true
  verify_ssl=true
  locations = ["AP_SOUTH_1"]
  name      = "%[1]s"
  frequency = 5
  status    = "ENABLED"
  type      = "SIMPLE"
  tags{
    key="monitor"
    values=["myMonitor"]
  }
  uri       = "https://www.one.newrelic.com"
}