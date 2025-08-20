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

variable "new_relic_tenancy_ocid" {
  description = "The OCID of the New Relic tenancy to which access will be granted."
  type        = string
}

variable "new_relic_group_ocid" {
  description = "The OCID of the New Relic group that will be granted access."
  type        = string
}

variable "newrelic_account_id" {
  type        = string
  description = "The New Relic account ID for sending metrics to New Relic endpoints"
}

variable "newrelic_ingest_api_key" {
  type        = string
  description = "The Ingest API key for sending logs to New Relic endpoints"
}

variable "newrelic_user_api_key" {
  type        = string
  description = "The User API key for Linking the OCI Account to the New Relic account"
}