terraform {
  required_version = ">= 1.2.0"
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "5.46.0"
    }
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "oci" {
  alias        = "home"
  tenancy_ocid = var.tenancy_ocid
  region       = var.region
  private_key  = var.private_key
  fingerprint  = var.fingerprint
}

provider "newrelic" {
  region     = var.newrelic_provider_region # US or EU
  account_id = var.newrelic_account_id
  api_key    = local.user_api_key
}

