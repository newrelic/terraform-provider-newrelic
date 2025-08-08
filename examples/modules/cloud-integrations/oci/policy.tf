locals {
  home_region = [
    for region in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions : region.region_name
    if region.is_home_region
  ][0]
  is_home_region = var.region == local.home_region
}


#Resource for the dynamic group
resource "oci_identity_dynamic_group" "nr_serviceconnector_group" {
  count          = local.is_home_region ? 1 : 0
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Dynamic group for service connector"
  matching_rule  = "Any {resource.type = 'serviceconnector', resource.type = 'fnfunc'}"
  name           = var.dynamic_group_name
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
}

#Resource for the policy
resource "oci_identity_policy" "nr_metrics_policy" {
  count          = local.is_home_region ? 1 : 0
  depends_on     = [oci_identity_dynamic_group.nr_serviceconnector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have any connector hub read from monitoring source and write to a target function"
  name           = var.newrelic_metrics_policy
  statements     = [
    "Allow dynamic-group ${var.dynamic_group_name} to read metrics in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-function in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-invocation in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to read secret-bundles in tenancy",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}