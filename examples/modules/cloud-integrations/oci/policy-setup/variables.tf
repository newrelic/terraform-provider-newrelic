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
  description = "The Ingest API key for sending metrics to New Relic endpoints"
}

variable "newrelic_user_api_key" {
  type        = string
  sensitive   = true
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
