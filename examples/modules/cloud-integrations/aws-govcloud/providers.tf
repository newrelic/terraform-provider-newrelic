terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    newrelic = {
      source = "newrelic/newrelic"
      version = ">= 3.56.0"
      # Using this module requires a minimum version of `3.56.0` of the New Relic Terraform Provider.
    }
  }
}
