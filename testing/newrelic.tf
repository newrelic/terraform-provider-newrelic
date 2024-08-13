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

#
# resource "newrelic_synthetics_script_monitor" "monitor" {
#   status           = "DISABLED"
#   name             = "name123"
#   type             = "SCRIPT_BROWSER"
#   locations_public = ["AP_SOUTH_1", "AP_EAST_1"]
#   period           = "EVERY_HOUR"
#   devices         = ["MOBILE_LANDSCAPE", "DESKTOP", "TABLET_LANDSCAPE"]
#   browsers        = ["EDGE", "CHROME", "FIREFOX"]
#   enable_screenshot_on_failure_and_script = false
#
#   script = "$browser.get('https://one.newrelic.com')"
#
#   runtime_type_version = "100"
#   runtime_type         = "CHROME_BROWSER"
#   script_language      = "JAVASCRIPT"
#
#   tag {
#     key    = "some_key"
#     values = ["some_value"]
#   }
# }

resource "newrelic_synthetics_step_monitor" "foo" {
  name                                    = "step_name123"
  enable_screenshot_on_failure_and_script = true
  locations_public                        = ["US_EAST_1", "US_EAST_2"]
  period                                  = "EVERY_6_HOURS"
  status                                  = "DISABLED"
  runtime_type                            = "CHROME_BROWSER"
  runtime_type_version                    = "100"
  devices         = ["MOBILE_LANDSCAPE", "DESKTOP"]
  browsers        = ["EDGE", "CHROME"]
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://www.newrelic.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}