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

resource "newrelic_synthetics_script_monitor" "bar" {
  enable_screenshot_on_failure_and_script	=	true
  locations_public	=	["AP_SOUTH_1"]
  name	=	"Script_name"
  period	=	"EVERY_15_MINUTES"
  runtime_type_version	=	"100"
  runtime_type	=	"CHROME_BROWSER"
  script_language	=	"JAVASCRIPT"
  status	=	"ENABLED"
  type	=	"SCRIPT_BROWSER"
  script	=	"$browser.get('https://one.newrelic.com')"
  tag {
    key	= "Name"
    values	= ["scriptedMonitor","simpler"]
  }
}

