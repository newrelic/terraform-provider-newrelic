# Module Input Variables

# Required: Tenancy OCID (needed for data sources)
variable "tenancy_ocid" {
  description = "The OCID of the tenancy"
  type        = string
}

# Required: User OCID (needed for OCI API key authentication)
variable "user_ocid" {
  description = "The OCID of the OCI user whose API key is used to manage the resources this module creates"
  type        = string
}

# Identity Domain Configuration
variable "identity_domain_name" {
  description = "Name of the identity domain to use"
  type        = string
  default     = "Default"
}

# New Relic Configuration
variable "newrelic_region" {
  description = "New Relic region (US or EU)"
  type        = string
  default     = "US"
  validation {
    condition     = contains(["US", "EU"], var.newrelic_region)
    error_message = "New Relic region must be either 'US' or 'EU'."
  }
}

# Resource Naming
variable "resource_prefix" {
  description = "Prefix for all created resources"
  type        = string
  default     = "newrelic"
}

# Optional: Trust configuration name
variable "trust_name" {
  description = "Name for the identity propagation trust"
  type        = string
  default     = "newrelic-wif-trust"
}

# Group and Policy Configuration
variable "compartment_id" {
  description = "Compartment ID where to create the IAM policy (defaults to tenancy root)"
  type        = string
  default     = ""
}

# Service User Configuration
variable "service_user_description" {
  description = "Description for the New Relic service user"
  type        = string
  default     = "Service user for New Relic Workload Identity Federation"
}

# OAuth Apps Configuration
variable "activate_oauth_apps" {
  description = "Whether to activate or deactivate OAuth applications"
  type        = bool
  default     = true
}

variable "home_region" {
  description = "The OCI home region where identity domain resources will be created (e.g., us-ashburn-1, us-phoenix-1)"
  type        = string
}

variable "fingerprint" {
  description = "The fingerprint of the API key used for OCI authentication"
  type        = string
}

variable "private_key" {
  description = "The private key content for OCI API authentication (PEM format)"
  type        = string
  sensitive   = true
}

# NR-562518: OCI WIF trust type. UPST (default) impersonates a service user; RPST uses
# claim-based ephemeral principals via `identityfederateddomainapp`. Backward compatible —
# existing customers default to UPST and see no change.
variable "trust_type" {
  description = "OCI WIF trust type. UPST (default) or RPST."
  type        = string
  default     = "UPST"
  validation {
    condition     = contains(["UPST", "RPST"], var.trust_type)
    error_message = "trust_type must be either 'UPST' or 'RPST'."
  }
}

# NR-562518: NR account id used for the propagated `account_id` claim (becomes
# `ext_account_id` on the RPST and is referenced in IAM policies as
# `request.principal.ext_account_id`). Required when trust_type = RPST.
variable "newrelic_account_id" {
  description = "New Relic account id. Used as the `account_id` JWT claim value for RPST."
  type        = string
  default     = ""
}

# NR-562518: optional value propagated as `ext_resource_tag` on the RPST. When set, the IAM
# policy additionally scopes by this tag so customers can tie an integration to a specific tag
# value (e.g. environment=prod). Ignored for UPST.
variable "resource_tag" {
  description = "Optional value propagated as ext_resource_tag claim on the RPST for tag-based IAM scoping. Ignored for UPST."
  type        = string
  default     = ""
}
