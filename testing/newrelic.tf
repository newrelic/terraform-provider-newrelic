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


resource "newrelic_synthetics_monitor" "my_synth" {
  frequency = 5
  locations = ["AP_SOUTH_1"]
  name      = "my_monitor"
  status    = "ENABLED"
  type      = "SIMPLE"
  uri       = "https://www.one.newrelic.com"
  validation_string=true
  verify_ssl=true
  bypass_head_request=true
  treat_redirect_as_failure=true
  runtime_type="CHROME_BROWSER"
  runtime_type_version="100"
  script_language="JAVASCRIPT"
  tags{
    key="monitor"
    values=["myMonitor"]
  }
  enable_screenshot_on_failure_and_script=true
  custom_headers{
    name="Name"
    value="simpleMonitor"
  }
}