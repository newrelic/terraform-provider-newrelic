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

variable "exclude_metric_filters" {
  description = "Map of exclusive metric filters. Use the namespace as the key and the list of metric names as the value."
  type        = map(list(string))
  default     = {}
}

variable "include_metric_filters" {
  description = "Map of inclusive metric filters. Use the namespace as the key and the list of metric names as the value."
  type        = map(list(string))
  default     = {}
}

variable "is_primary_region" {
  description = "Determines if certain mutually exclusive resources are created."
  type = bool
  default = true
}

variable "recorder_enabled" {
  description = "Allows enabling or disabling the AWS Config Recorder"
  type = bool
  default = true
}
