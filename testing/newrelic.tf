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
  name   = "name1236891625"
  type   = "SCRIPT_BROWSER"
  period = "EVERY_HOUR"
  script = "$browser.get('https://one.newrelic.com')"

  enable_screenshot_on_failure_and_script = false
  locations_public = ["AP_SOUTH_1"]
#   browsers = ["CHROME"]
#   devices = [ "MOBILE_PORTRAIT", "MOBILE_LANDSCAPE"]
  runtime_type_version = ""
  runtime_type         = ""
  script_language      = "JAVASCRIPT"
  use_unsupported_legacy_runtime = true
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}

