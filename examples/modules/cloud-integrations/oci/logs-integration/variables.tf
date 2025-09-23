# OCI variables
variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created. Use the compartment OCID already created as part of Policy Setup module."
}

variable "region" {
  type        = string
  description = "The name of the OCI region where these resources will be deployed."
}

variable "newrelic_logging_identifier" {
  type        = string
  description = "A unique label or name identifier for all resources in this deployment. Leave it blank if not needed."
  default = "logs"
}

# VCN variables
variable "create_vcn" {
  type        = bool
  default     = true
  description = "Variable to create virtual network for the setup. True by default"
}

variable "function_subnet_id" {
  type        = string
  description = "The OCID of the subnet to be used for the function app. If create_vcn is set to true, that will take precedence"
}

# Connector Hub Variables
variable "connector_hub_details" {
  type        = string
  description = "JSON formatted string for the service connectors to be created. See README for details and example"
  default     = null
}

variable "batch_size_in_kbs" {
  type = number
  description = "The maximum size of the batch of events to process. Maximum is 6000 KB."
  default = 6000
}

variable "batch_time_in_sec" {
  type = number
  description = "The maximum amount of time to wait before processing a batch of events. Maximum is 300 seconds."
  default = 60
}

# New Relic Function variables
variable "image_version" {
  type = string
  description = "The version of the Docker image for the New Relic function for the region."
  default = "latest"
}

variable "debug_enabled" {
  type        = string
  default     = "FALSE"
  description = "Enable debug mode."
}

variable "new_relic_region" {
  type        = string
  default     = "US"
  description = "New Relic Region. US or EU"
}

variable "secret_ocid" {
  type        = string
  description = "OCI Vault Secret OCID that contains the New Relic License Key. Use the secret OCID already created as part of Policy Setup module."
}