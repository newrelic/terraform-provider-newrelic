# Create IAM Group for New Relic Service Users (UPST only).
# RPST uses ephemeral `identityfederateddomainapp` principals — no group needed.
resource "oci_identity_domains_group" "newrelic_service_group" {
  count = local.is_upst ? 1 : 0

  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:core:2.0:Group"]

  display_name = "${local.resource_prefix}-svc-user-group-${local.suffix}"

  attribute_sets = ["all"]

  lifecycle {
    ignore_changes = [schemas]
  }
}

# Create IAM Policy for New Relic.
# UPST: grants the service-user group read access.
# RPST: claim-based — allows any-user where principal type is `identityfederateddomainapp` and
# both `ext_account_id` and `ext_tenancy_id` match. The tenancy match is defense in depth: an
# RPST minted for one tenancy cannot be replayed against another. Customers who want tag-based
# scoping (`ext_resource_tag`) can add their own policy statement; we don't bake it in here so
# the default policy stays minimal.
resource "oci_identity_policy" "newrelic_service_policy" {
  compartment_id = var.compartment_id != "" ? var.compartment_id : var.tenancy_ocid
  name           = "${local.resource_prefix}-svc-user-policy-${local.suffix}"
  description    = "Policy granting New Relic read-only access to OCI resources"

  statements = local.is_upst ? [
    "Allow group '${oci_identity_domains_group.newrelic_service_group[0].display_name}' to read all-resources in tenancy",
    ] : [
    format("allow any-user to read all-resources in tenancy where all { request.principal.type = 'identityfederateddomainapp', request.principal.ext_account_id = '%s', request.principal.ext_tenancy_id = '%s' }", var.newrelic_account_id, var.tenancy_ocid),
  ]

  depends_on = [oci_identity_domains_group.newrelic_service_group]
}
