# Retrieves the Identity Domain Administrator role
# This role is required for the admin app to create the identity propagation trust
data "oci_identity_domains_app_roles" "app_roles" {
  #Required
  idcs_endpoint = local.identity_domain_url

  attributes = "id,displayName"
  app_role_filter = "displayName eq \"Identity Domain Administrator\""
}

# Admin OAuth Application (for initial setup only)
# This app has elevated privileges to create the identity propagation trust
resource "oci_identity_domains_app" "admin_app" {
  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:oracle:idcs:App"]

  display_name = "${local.resource_prefix}-ida-app-${local.suffix}"
  active       = var.activate_oauth_apps

  based_on_template {
    value = "CustomWebAppTemplateId"
  }

  # OAuth Configuration
  allowed_grants = ["client_credentials"]

  # Enable OAuth Client Configuration
  is_oauth_client = true

  # Client Configuration
  client_type    = "confidential"
  bypass_consent = true

  attribute_sets = ["all"]

  lifecycle {
    ignore_changes = [active]
  }
}

# Token Exchange OAuth Application (for routine use)
# This app is used by New Relic to exchange tokens and impersonate the service user
resource "oci_identity_domains_app" "token_exchange_app" {
  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:oracle:idcs:App"]

  display_name = "${local.resource_prefix}-token-exchange-app-${local.suffix}"
  active       = var.activate_oauth_apps

  based_on_template {
    value = "CustomWebAppTemplateId"
  }

  # OAuth Configuration
  allowed_grants = ["client_credentials"]

  # Enable OAuth Client Configuration
  is_oauth_client = true

  # Client Configuration
  client_type    = "confidential"
  bypass_consent = true

  attribute_sets = ["all"]
}

# Grants the Identity Domain Administrator role to the admin OAuth app
# This allows the admin app to create and manage identity propagation trusts
resource "oci_identity_domains_grant" "admin_app_domain_admin_grant" {
  idcs_endpoint = local.identity_domain_url
  schemas       = ["urn:ietf:params:scim:schemas:oracle:idcs:Grant"]

  grant_mechanism = "ADMINISTRATOR_TO_APP"

  grantee {
    value = oci_identity_domains_app.admin_app.id
    type  = "App"
  }

  entitlement {
    attribute_name  = "appRoles"
    attribute_value = data.oci_identity_domains_app_roles.app_roles.app_roles[0].id  # This may need to be the actual role ID
  }

  app {
    value = "IDCSAppId"
  }

  attribute_sets = ["all"]
}
