terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region     = "US"
}

variable "account_id" {
  type = number
  default = 3806526
}



resource "newrelic_teams_organization_settings" "org" {
  discovery_enabled  = true
  discovery_tag_keys = ["team"]
}

resource "newrelic_teams_hierarchy_level" "squad" {
  name = "Squad"
}
