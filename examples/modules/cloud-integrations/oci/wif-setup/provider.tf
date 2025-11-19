terraform {
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = ">= 5.0.0"
    }
  }
}

# Configure the OCI Provider
provider "oci" {
  tenancy_ocid = var.tenancy_ocid
  region       = var.region
  fingerprint  = var.fingerprint
  private_key  = var.private_key
}
