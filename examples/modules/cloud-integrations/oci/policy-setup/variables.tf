variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "nr_prefix" {
  type        = string
  description = "The prefix for naming resources in this module."
  default     = "newrelic"
}

variable "region" {
  type        = string
  description = "OCI Region as documented at https://docs.cloud.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm"
}

variable "newrelic_ingest_api_key" {
  type        = string
  sensitive   = true
  default     = ""
  description = "The Ingest API key for sending metrics to New Relic endpoints"
}

variable "newrelic_user_api_key" {
  type        = string
  sensitive   = true
  default     = ""
  description = "The User API key for Linking the OCI Account to the New Relic account"
}

variable "newrelic_account_id" {
  type        = string
  sensitive   = true
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "newrelic_provider_region" {
  type        = string
  description = "The region where the New Relic provider resources will be created."
}

variable "instrumentation_type" {
  type        = string
  description = "Specifies which policy types to create. Valid values: 'METRICS', 'LOGS', or a comma-separated list such as 'METRICS,LOGS'."
}

variable "client_id" {
  type        = string
  sensitive   = true
  description = "Client ID for API access"
}

variable "client_secret" {
  type        = string
  sensitive   = true
  description = "Client Secret for API access"
}

variable "oci_domain_url" {
  type        = string
  description = "OCI domain URL"
}

variable "private_key" {
  type        = string
  sensitive   = true
  description = "The private key content for OCI API authentication (alternative to private_key_path). Use this if you want to pass the key content directly instead of a file path."
  default     = ""
}

variable "fingerprint" {
  type        = string
  description = "The fingerprint of the public key. Get this from OCI Console -> User Settings -> API Keys"
}

variable "user_key_secret_ocid" {
  type        = string
  default     = ""
  description = "The OCID of the secret containing the New Relic User API key"
}

variable "ingest_key_secret_ocid" {
    type        = string
    default     = ""
    description = "The OCID of the secret containing the New Relic Ingest License API key"
}

# NR-562518: trust_type + resource_tag are propagated to the newrelic_cloud_oci_link_account
# resource that this module creates internally. UPST (default) preserves existing behavior;
# RPST customers wire these inputs from the wif-setup module's outputs so the linked account
# is registered with the matching trust shape.
variable "trust_type" {
  type        = string
  default     = "UPST"
  description = "OCI WIF trust type. UPST (default, service-user impersonation) or RPST (claim-based ephemeral principal). Must match what the wif-setup module created."
  validation {
    condition     = contains(["UPST", "RPST"], var.trust_type)
    error_message = "trust_type must be either 'UPST' or 'RPST'."
  }
}

variable "resource_tag" {
  type        = string
  default     = ""
  description = "Optional value propagated as the ext_resource_tag claim on the RPST. Customers use this for tag-based resource scoping in their IAM policies. Ignored for UPST."
}
