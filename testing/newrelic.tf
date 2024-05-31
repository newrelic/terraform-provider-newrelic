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

data "newrelic_entity" "app" {
  name       = "Dummy App Pro Max"
  domain     = "APM"
  type       = "APPLICATION"
}

output "x" {
  value = data.newrelic_entity.app
}

