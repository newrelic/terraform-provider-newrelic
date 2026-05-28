terraform {
  required_version = ">= 1.2.0"
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "5.46.0"
    }
  }
}

provider "oci" {
  alias        = "home"
  tenancy_ocid = var.tenancy_ocid
  region       = var.region
}