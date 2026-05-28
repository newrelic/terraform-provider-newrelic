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

variable "output_format" {
  description = "The output format for the CloudWatch metric stream"
  type        = string
  default     = "opentelemetry0.7"

  validation {
    condition     = contains(["opentelemetry0.7", "opentelemetry1.0"], var.output_format)
    error_message = "The output_format must be either 'opentelemetry0.7' or 'opentelemetry1.0'."
  }
}
variable "enable_config_recorder" {
  description = "Set to true to enable AWS Configuration Recorder."
  type        = bool
  default     = false
}
