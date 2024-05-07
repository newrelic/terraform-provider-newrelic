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

resource "newrelic_cloud_aws_link_account" "foo" {
  arn = "arn:aws:iam::709144918866:role/NewRelicInfrastructure-Integrations-V2"
  metric_collection_mode = "PUSH"
  name = "foo"
}