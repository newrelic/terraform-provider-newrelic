terraform {
  required_version = ">= 1.2.0"
  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "7.11.0"
    }
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}