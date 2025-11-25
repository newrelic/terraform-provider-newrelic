# Module Input Variables

# Required: Tenancy OCID (needed for data sources)
variable "tenancy_ocid" {
  description = "The OCID of the tenancy"
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