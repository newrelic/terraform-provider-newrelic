# 🎯 PRIMARY OUTPUT: NEW RELIC INTEGRATION DETAILS
# Contains all credentials and URLs needed to configure the New Relic OCI integration
output "newrelic_integration_details" {
  description = "Essential details to share with New Relic for integration setup (IAM domain URL, OAuth client credentials)"
  value = {
    # IAM Domain URL (Required by New Relic)
    iam_domain_url = regex("^(https?://[^:]+)", local.identity_domain_url)[0]
    token_exchange_client_id = oci_identity_domains_app.token_exchange_app.name
    token_exchange_client_secret = oci_identity_domains_app.token_exchange_app.client_secret
  }
}
