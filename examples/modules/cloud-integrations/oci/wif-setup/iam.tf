# Create IAM Group for New Relic Service Users
# This group will contain the service user that New Relic will impersonate
resource "oci_identity_domains_group" "newrelic_service_group" {
  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:core:2.0:Group"]

  display_name  = "${local.resource_prefix}-svc-user-group-${local.suffix}"

  attribute_sets = ["all"]

  lifecycle {
    create_before_destroy = true
  }
}

# Create IAM Policy for New Relic Group
# Grants read-only access to all OCI resources for monitoring purposes
resource "oci_identity_policy" "newrelic_service_policy" {
  compartment_id = var.compartment_id != "" ? var.compartment_id : var.tenancy_ocid
  name           = "${local.resource_prefix}-svc-user-policy-${local.suffix}"
  description    = "Policy granting New Relic service users read-only access to all OCI resources"

  statements = [
    "Allow group '${oci_identity_domains_group.newrelic_service_group.display_name}' to read all-resources in tenancy",
  ]

  depends_on = [oci_identity_domains_group.newrelic_service_group]
}