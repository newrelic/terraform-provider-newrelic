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
  user_ocid            = var.current_user_ocid
  region               = var.region
}

provider "newrelic" {
  region = "US" # US or EU
  account_id = var.newrelic_account_id
  api_key = var.newrelic_user_api_key
}