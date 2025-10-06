locals {
  home_region = [
    for rs in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions :
    rs.region_name if rs.region_key == data.oci_identity_tenancy.current_tenancy.home_region_key
  ][0]

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix = "tf"

  # Names for the network infra
  vcn_name        = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  nat_gateway     = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  service_gateway = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  subnet          = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
}

