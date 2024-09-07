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



resource "newrelic_synthetics_script_monitor" "monitor" {
  status = "DISABLED"
  name   = "name1236"
  type   = "SCRIPT_API"
  period = "EVERY_HOUR"
  script = "$browser.get('https://one.newrelic.com')"
#   browsers = ["CHROME", "FIREFOX"]
#   devices = ["MOBILE_PORTRAIT","TABLET_LANDSCAPE"]
  enable_screenshot_on_failure_and_script = false
  locations_public = ["AP_SOUTH_1","US_EAST_1"]
  runtime_type_version = "16.10"
  runtime_type         = "NODE_API"
  script_language = "JAVASCRIPT"
#   device_orientation = "LANDSCAPE"
#   device_type = "MOBILE"
#   use_unsupported_legacy_runtime = true
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}


resource "newrelic_synthetics_script_monitor" "monitor2" {
  status = "DISABLED"
  name   = "name1237"
  type   = "SCRIPT_BROWSER"
  period = "EVERY_HOUR"
  script = "$browser.get('https://one.newrelic.com')"
#     browsers = ["CHROME", "FIREFOX"]
#     devices = ["MOBILE_PORTRAIT","TABLET_LANDSCAPE"]
  enable_screenshot_on_failure_and_script = false
  locations_public = ["AP_SOUTH_1","US_EAST_1"]
  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language = "JAVASCRIPT"
    device_orientation = "LANDSCAPE"
    device_type = "MOBILE"
#     use_unsupported_legacy_runtime = true
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
