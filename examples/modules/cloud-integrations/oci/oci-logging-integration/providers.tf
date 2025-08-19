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

# Variables
provider "oci" {
  alias                = "home"
  tenancy_ocid         = var.tenancy_ocid
  user_ocid            = data.oci_identity_user.current_user.user_id
  region               = var.region
  private_key_path     = var.private_key_path != "" ? var.private_key_path : null
  fingerprint          = var.fingerprint
}

provider "newrelic" {
  region = "US" # US or EU
  account_id = var.newrelic_account_id
  api_key = var.newrelic_user_api_key
}

data "oci_identity_user" "current_user" {
  user_id = var.current_user_ocid
}