terraform {
  required_version = ">= 1.5"

  required_providers {
    newrelic = {
      # Registry source — overridden by dev.tfrc when testing a local build.
      source  = "newrelic/newrelic"
      version = ">= 3.87"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 5.0"
    }
  }
}

provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

provider "google" {
  project = var.gcp_project_id
}
