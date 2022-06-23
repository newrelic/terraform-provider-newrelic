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
  custom_header{
    name	=	"name"
    value	=	"simpleMonitorUpdated"
  }
  treat_redirect_as_failure	=	false
  validation_string	=	"success"
  bypass_head_request	=	false
  verify_ssl	=	false
  //false to true --w
  //true to false --nw
  location_public	=	["AP_SOUTH_1","AP_EAST_1"]
  name	=	"%s-updated"
  period	=	"EVERY_15_MINUTES"
  status	=	"ENABLED"
  type	=	"SIMPLE"
#  tag {
#    key	=	"Name"
#    values	=	[ "myMonitor","simple_monitor"]
#  }
  uri	=	"https://one.newrelic.com"
}

