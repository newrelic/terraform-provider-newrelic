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
        exit 1
      fi
  EOT
  ]

  depends_on = [oci_identity_domains_app.admin_app]
}

locals {
  ida_oauth_token = data.external.ida_oauth_token[0].result.access_token
}

locals {
  # NR-562518: Build the trust JSON body conditionally. UPST uses User-type subject + service-user
  # impersonation mapping; RPST uses Resource-type subject + impersonatingResource + claim
  # propagation. Same JWKS endpoint either way (NR reuses the existing UPST signing key).
  trust_body_upst = jsonencode({
    active             = true
    allowImpersonation = true
    issuer             = local.newrelic_config.issuer_name
    name               = "newrelic-trust-setup"
    oauthClients       = [oci_identity_domains_app.token_exchange_app.name]
    publicKeyEndpoint  = local.newrelic_config.public_jwks_url
    impersonationServiceUsers = local.is_upst ? [{
      rule  = "sub eq '${local.newrelic_config.subject_name}'"
      value = oci_identity_domains_user.svc_user[0].id
    }] : []
    subjectType = "User"
    type        = "JWT"
    schemas     = ["urn:ietf:params:scim:schemas:oracle:idcs:IdentityPropagationTrust"]
  })
  trust_body_rpst = jsonencode({
    active                = true
    allowImpersonation    = true
    issuer                = local.newrelic_config.rpst_issuer_name
    name                  = "newrelic-rpst-trust-setup"
    oauthClients          = [oci_identity_domains_app.token_exchange_app.name]
    publicKeyEndpoint     = local.newrelic_config.public_jwks_url
    impersonatingResource = local.impersonating_resource
    # NR-562518: propagate every claim NR may send so customer IAM policies can scope on any of
    # them. NR's worker/api-v2 always send `account_id` and `tenancy_id`; `resource_tag` is sent
    # only when configured in `params.oci.custom_claims`. OCI silently skips propagating claims
    # that aren't present in the JWT, so listing all three is safe.
    claimPropagations = ["ext_account_id", "ext_tenancy_id", "ext_resource_tag"]
    subjectType       = "Resource"
    type              = "JWT"
    schemas           = ["urn:ietf:params:scim:schemas:oracle:idcs:IdentityPropagationTrust"]
  })
  trust_body = local.is_upst ? local.trust_body_upst : local.trust_body_rpst
}

# Creates the Identity Propagation Trust between New Relic and OCI.
# UPST: maps NR's `sub eq '<subject>'` to a customer service user (impersonation).
# RPST: validates `res_type` matches `impersonatingResource` and propagates `ext_account_id`
# from the JWT into the RPST so customer IAM policies can reference it.
resource "null_resource" "trust_setup" {
  provisioner "local-exec" {
    command = <<EOT
      sleep 20
      RESPONSE=$(curl --location '${local.identity_domain_url}/admin/v1/IdentityPropagationTrusts' \
      --header 'Content-Type: application/json' \
      --header 'Authorization: Bearer ${local.ida_oauth_token}' \
      --data '${local.trust_body}')
      TRUST_ID=$(echo "$RESPONSE" | jq -r '.id // empty')
      sleep 20
      if [ -n "$TRUST_ID" ] && [ "$TRUST_ID" != "null" ]; then
        echo $RESPONSE
      else
        ERROR_MESSAGE=$(echo "$RESPONSE" | jq -r '.detail // .error // "Unknown error occurred"')
        echo {\"error\":\"$ERROR_MESSAGE\"}
        exit 1
      fi
    EOT
  }

  depends_on = [oci_identity_domains_app.admin_app]
}
