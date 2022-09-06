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

data "newrelic_synthetics_private_location" "example" {
  name = "My private location"
}
resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private { guid = data.newrelic_synthetics_private_location.example.id }
}