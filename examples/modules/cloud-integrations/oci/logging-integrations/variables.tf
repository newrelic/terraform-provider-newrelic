variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created. Do not modify."
}

variable "newrelic_logging_prefix" {
  type        = string
  description = "The prefix for naming all the logging resources in this module."
}

variable "region" {
  type        = string
  description = "The name of the OCI region where these resources will be deployed."
}

variable "create_vcn" {
  type        = bool
  description = "Variable to create virtual network for the setup. True by default"
}

variable "function_subnet_id" {
  type        = string
  description = "The OCID of the subnet to be used for the function app. If create_vcn is set to true, that will take precedence"
}

variable "payload_link" {
  type        = string
  description = "The link to the payload for the connector hubs."
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
  description = "OCI Vault Secret OCID that contains the New Relic License Key."
}