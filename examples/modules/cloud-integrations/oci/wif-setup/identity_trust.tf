# Generate and retrieve IDA OAuth Token
# Obtains an OAuth token from the admin app to authenticate API calls for trust setup
data "external" "ida_oauth_token" {
  count = 1
  program = ["sh", "-c", <<-EOT
      OAUTH_TOKEN=$(echo -n "${oci_identity_domains_app.admin_app.name}:${oci_identity_domains_app.admin_app.client_secret}" | base64)
      RESPONSE=$(curl --silent --location '${local.identity_domain_url}/oauth2/v1/token' \
      --header 'Content-Type: application/x-www-form-urlencoded;charset=UTF-8' \
      --header 'Authorization: Basic ${base64encode("${oci_identity_domains_app.admin_app.name}:${oci_identity_domains_app.admin_app.client_secret}")}' \
      --data 'grant_type=client_credentials&scope=urn:opc:idm:__myscopes__')

      ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token // empty')
      if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
        echo {\"access_token\":\"$ACCESS_TOKEN\"}
      else
        ERROR_MESSAGE=$(echo "$RESPONSE" | jq -r '.error_description // .error // "Unknown error occurred"')
        echo {\"error\":\"$ERROR_MESSAGE\"}
      fi
  EOT
  ]

  depends_on = [oci_identity_domains_app.admin_app]
}

locals {
  ida_oauth_token = data.external.ida_oauth_token[0].result.access_token
}

# Creates the Identity Propagation Trust between New Relic and OCI
# This trust enables New Relic to exchange JWT tokens for OCI access tokens
resource "null_resource" "trust_setup" {
  provisioner "local-exec" {
    command = <<EOT
      RESPONSE=$(curl --location '${local.identity_domain_url}/admin/v1/IdentityPropagationTrusts' \
      --header 'Content-Type: application/json' \
      --header 'Authorization: Bearer ${local.ida_oauth_token}' \
      --data '{
        "active": true,
        "allowImpersonation": true,
        "issuer": "${local.newrelic_config.issuer_name}",
        "name": "newrelic-trust-setup",
        "oauthClients": ["${oci_identity_domains_app.token_exchange_app.name}"],
        "publicKeyEndpoint": "${local.newrelic_config.public_jwks_url}",
        "impersonationServiceUsers": [{
          "rule": "sub eq '${local.newrelic_config.subject_name}'",
          "value": "${oci_identity_domains_user.svc_user.id}"
        }],
        "subjectType": "User",
        "type": "JWT",
        "schemas": ["urn:ietf:params:scim:schemas:oracle:idcs:IdentityPropagationTrust"]
      }')
      echo $RESPONSE
    EOT
  }

  depends_on = [oci_identity_domains_app.admin_app]
}
