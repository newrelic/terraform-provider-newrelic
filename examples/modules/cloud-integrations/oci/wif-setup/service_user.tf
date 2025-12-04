# Creates a service user for New Relic Workload Identity Federation
# This user will be impersonated by New Relic to access OCI resources
resource "oci_identity_domains_user" "svc_user" {
  #Required
  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:core:2.0:User"]
  user_name     = "${local.resource_prefix}-wif-svc-user-${local.suffix}"
  urnietfparamsscimschemasoracleidcsextensionuser_user {
    service_user = true
  }
  lifecycle {
    ignore_changes = [schemas]
  }
}

# Adds the service user to the New Relic service group
# This membership grants the service user the permissions defined in the IAM policy
resource "oci_identity_user_group_membership" "svc_user_group_membership" {
  #Required
  group_id = oci_identity_domains_group.newrelic_service_group.ocid
  user_id  = oci_identity_domains_user.svc_user.ocid
}
