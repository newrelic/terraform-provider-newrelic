terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "Staging" # US or EU
}

resource "newrelic_monitor_downtime" "foo" {
  name = "Random Monitor Downtime 2111-1548-IST"
  monitor_guids = [
    "MXxTWU5USHxNT05JVE9SfDk2ZTBiMWUyLTBhMmEtNGI1MS05OTI0LTgwNjBmNTA0N2ZkMw",
    "MXxTWU5USHxNT05JVE9SfGQ2OGY2YWJjLThlNTgtNDNmMC05ZDczLWMwZWFmMjE5MTMwZA"
  ]
  mode = "MONTHLY"
  start_time = "2023-11-22T10:11:15"
  end_time = "2023-11-29T10:11:15"
  time_zone = "Asia/Kolkata"
  end_repeat {
    on_date = "2023-11-26"
  }
  frequency {
    days_of_month = [2,3,5,7,11,13]
  }
}