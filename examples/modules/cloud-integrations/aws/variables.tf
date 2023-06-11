variable "newrelic_account_id" {
  type = string
}

variable "newrelic_account_region" {
  type    = string
  default = "US"

  validation {
    condition     = contains(["US", "EU"], var.newrelic_account_region)
    error_message = "Valid values for region are 'US' or 'EU'."
  }
}

variable "name" {
  type    = string
  default = "production"
}
