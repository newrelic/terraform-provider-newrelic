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

resource "newrelic_synthetics_cert_check_monitor" "cert-check-monitor" {
  name                   = "random-cert-check-monitor"
  domain                 = "www.example.com"
  locations_public       = ["AP_EAST_1", "AWS_AP_SOUTH_1", "AWS_AP_EAST_1"]
  certificate_expiration = "10"
  period                 = "EVERY_6_HOURS"
  status                 = "ENABLED"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}

#resource "newrelic_synthetics_script_monitor" "monitor" {
#  status           = "ENABLED"
#  name             = "Script Monitor 2470"
#  type             = "SCRIPT_BROWSER"
#  locations_public = ["AWS_AP_SOUTH_1", "AP_EAST_1", "AWS_US_EAST_2"]
#  period           = "EVERY_HOUR"
#
#  enable_screenshot_on_failure_and_script = false
#
#  script = "$browser.get('https://one.newrelic.com')"
#
#  runtime_type_version = "100"
#  runtime_type         = "CHROME_BROWSER"
#  script_language      = "JAVASCRIPT"
#
#  tag {
#    key    = "some_key"
#    values = ["some_value"]
#  }
#}