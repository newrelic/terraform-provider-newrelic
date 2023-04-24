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

resource "newrelic_synthetics_monitor" "monitor_pns" {
    status           = "ENABLED"
    name             = "monitor_pns3"
    period           = "EVERY_30_MINUTES"
    uri              = "https://www.one.newrelic.com"
    type             = "BROWSER"
    locations_public = ["AP_SOUTH_1"]
  
    custom_header {
      name  = "some_name"
      value = "some_value"
    }
  
    enable_screenshot_on_failure_and_script = true
    validation_string                       = "success"
    verify_ssl                              = true
  
  }