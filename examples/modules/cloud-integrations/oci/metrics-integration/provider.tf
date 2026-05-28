terraform {
  required_version = ">= 1.2.0"
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "5.46.0"
    }
    external = {
      source  = "hashicorp/external"
      version = "2.3.5"
    }
  }
}


# --- Home Provider Configurations ---
provider "oci" {
  alias        = "home_provider"
  tenancy_ocid = var.tenancy_ocid
  region       = local.home_region
  private_key  = var.private_key
  fingerprint  = var.fingerprint
}

provider "oci" {
  alias        = "current_region"
  tenancy_ocid = var.tenancy_ocid
  region       = var.region
  private_key  = var.private_key
  fingerprint  = var.fingerprint
}

