# Module Input Variables

variable "tenancy_ocid" {
  description = "The OCID of the tenancy"
  type        = string
}

variable "region" {
  description = "The OCI region to create the repository in (e.g. us-ashburn-1). OCIR repository names are unique tenancy+region-wide, not per compartment -- run this module once per region you deploy the logs/metrics integration into, not once per tenancy."
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

variable "compartment_ocid" {
  description = "The OCID of the compartment to create the repository in. Use the compartment OCID already created as part of Policy Setup module."
  type        = string
}

variable "newrelic_logging_identifier" {
  description = "A unique label or name identifier matching the one passed to the logs-integration module. Leave it blank if not needed."
  type        = string
  default     = "logs"
}
