# 🎯 PRIMARY OUTPUT: NEW RELIC INTEGRATION DETAILS
# Contains all credentials and URLs needed to configure the New Relic OCI integration
output "newrelic_integration_details" {
  description = "Essential details to share with New Relic for integration setup (IAM domain URL, OAuth client credentials)"
  value = {
    # IAM Domain URL (Required by New Relic)
    iam_domain_url               = regex("^(https?://[^:]+)", local.identity_domain_url)[0]
    token_exchange_client_id     = oci_identity_domains_app.token_exchange_app.name
    token_exchange_client_secret = oci_identity_domains_app.token_exchange_app.client_secret
    # NR-562518: trust type chosen by the customer (UPST or RPST). Pass through to the
    # newrelic_cloud_oci_link_account resource as `trust_type`.
    trust_type = var.trust_type
    # Same value the IAM policy scopes on; pass through so NR stamps it on every JWT.
    resource_tag = var.resource_tag
  }
}

# Static value the customer needs to paste into the OCI Identity Propagation Trust as
# `impersonatingResource` (RPST only). NR sends the same string as `res_type` during exchange.
output "impersonating_resource" {
  description = "Static impersonatingResource string. Same for every NR-OCI integration."
  value       = local.impersonating_resource
}
