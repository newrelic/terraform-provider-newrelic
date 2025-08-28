variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "compartment_ocid" {
  description = "The OCID of the compartment where resources will be created."
  type        = string
}

variable "current_user_ocid" {
  type        = string
  description = "The OCID of the current user executing the terraform script. Do not modify."
}

variable "region" {
  description = "The home region where the vault and policies will be created."
  type        = string
  default     = "us-ashburn-1"
}

variable "newrelic_account_id" {
  type        = string
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "newrelic_logging_prefix" {
  type        = string
  description = "The prefix for naming all the logging resources in this module."
}

variable "newrelic_region" {
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

variable "debug_enabled" {
  type        = string
  default     = "FALSE"
  description = "Enable debug mode."
}

variable "newrelic_user_api_key" {
  type        = string
  description = "New Relic user api key for account linking call"
}

variable "payload_link" {
  type        = string
  description = "The link to the payload for the connector hubs."
}