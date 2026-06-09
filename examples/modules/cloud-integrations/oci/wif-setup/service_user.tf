# Creates a service user for New Relic Workload Identity Federation (UPST only).
# UPST impersonates this user; RPST uses ephemeral `identityfederateddomainapp` and skips this.
resource "oci_identity_domains_user" "svc_user" {
  count = local.is_upst ? 1 : 0

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

# Adds the service user to the New Relic service group (UPST only).
resource "oci_identity_user_group_membership" "svc_user_group_membership" {
  count = local.is_upst ? 1 : 0

  #Required
  group_id = oci_identity_domains_group.newrelic_service_group[0].ocid
  user_id  = oci_identity_domains_user.svc_user[0].ocid
}
