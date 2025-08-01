locals {
  home_region = [
    for region in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions : region.region_name
    if region.is_home_region
  ][0]
  is_home_region = var.region == local.home_region
  policy_not_exists = length([
    for policy in data.oci_identity_policies.existing_policies.policies : policy.name
    if policy.name == var.newrelic_metrics_policy
  ]) == 0
  dynamic_group_not_exists = length([
    for dg in data.oci_identity_dynamic_groups.existing_dynamic_groups.dynamic_groups : dg.name
    if dg.name == var.dynamic_group_name
  ]) == 0  
}


#Resource for the dynamic group
resource "oci_identity_dynamic_group" "nr_serviceconnector_group" {
  count          = local.is_home_region && local.dynamic_group_not_exists ? 1 : 0
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Dynamic group for service connector"
  matching_rule  = "Any {resource.type = 'serviceconnector', resource.type = 'fnfunc'}"
  name           = var.dynamic_group_name
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
}

#Resource for the policy
resource "oci_identity_policy" "nr_metrics_policy" {
  count          = local.is_home_region && local.policy_not_exists ? 1 : 0
  depends_on     = [oci_identity_dynamic_group.nr_serviceconnector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have any connector hub read from monitoring source and write to a target function"
  name           = var.newrelic_metrics_policy
  statements     = [
    "Allow service keymanagementservice to manage vaults in tenancy",
    "Allow service keymanagementservice to manage secret-bundles in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to read metrics in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-function in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-invocation in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to manage stream-family in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to manage repos in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to read secret-bundles in tenancy",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}