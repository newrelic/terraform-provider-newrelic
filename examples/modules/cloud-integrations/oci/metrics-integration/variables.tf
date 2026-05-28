variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five.Do not modify."
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created. Do not modify."
}

variable "nr_prefix" {
  type        = string
  description = "The prefix for naming resources in this module."
  default     = "metrics"
}

variable "region" {
  type        = string
  description = "The name of the OCI region where these resources will be deployed."
}

variable "newrelic_endpoint" {
  type        = string
  default     = "US"
  description = "The endpoint to hit for sending the metrics. Varies by region [US|EU]"
  validation {
    condition     = contains(["US", "EU"], var.newrelic_endpoint)
    error_message = "Valid values for var: newrelic_endpoint are (US, EU)."
  }
}

variable "create_vcn" {
  type        = bool
  default     = true
  description = "Variable to create virtual network for the setup. True by default"
}

variable "function_subnet_id" {
  type        = string
  default     = ""
  description = "The OCID of the subnet to be used for the function app. If create_vcn is set to true, that will take precedence"
}

variable "connector_hubs_data" {
  type        = string
  description = "List of maps containing connector hub configuration data."
}

variable "ingest_api_secret_ocid" {
  type        = string
  description = "The OCID of the vault storing the ingest key for secure access."
}

variable "user_api_secret_ocid" {
  type        = string
  description = "The OCID of the vault storing the user key for secure access."
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

variable "image_version" {
  type = string
  description = "The version of the Docker image for the New Relic function for the region."
  default = "latest"
}

variable "image_bucket" {
  type = string
  description = "The name of the bucket where the Docker image for the New Relic function is stored."
  default = "idptojlonu4e"
}

variable "newrelic_account_id" {
  type        = string
  sensitive   = true
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "provider_account_id" {
  type        = string
  sensitive   = true
  description = "The Provider Account ID that has been linked with New Relic"
}
