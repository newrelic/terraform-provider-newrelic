# OCI Tenancy and User Variables
variable "tenancy_ocid" {
  type        = string
  description = "OCI tenant OCID, more details can be found at https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five"
}

variable "compartment_ocid" {
  type        = string
  description = "The OCID of the compartment where resources will be created."
}

# Subnet OCID for Function Application
variable "subnet_id" {
  type        = string
  description = "Subnet OCID for the function application"
}

# Log Forwarding Function Custom Variables
variable "nr_region" {
  type        = string
  description = "New Relic Region to forward Logs to. Valid values are 'US' or 'EU'."
  default = "US"
}

variable "newrelic_account_id" {
  type        = string
  description = "The New Relic account ID for sending logs to New Relic endpoints"
}

variable "log_group_id" {
  type        = string
  description = "log group OCID to send logs to New Relic."
}

variable "log_id" {
  type        = string
  description = "log OCID to send logs to New Relic."
}

variable "debug_enabled" {
  type        = string
  description = "Enable debug mode."
  default     = "FALSE"
}