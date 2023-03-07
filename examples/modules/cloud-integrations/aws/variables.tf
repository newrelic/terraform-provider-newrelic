variable "account_id" {
  type = string
}

variable "account_name" {
  type    = string
  default = "production"
}

variable "region" {
  type    = string
  default = "US"

  validation {
    condition     = contains(["US", "EU"], var.region)
    error_message = "Valid values for region are 'US' or 'EU'."
  }
}
