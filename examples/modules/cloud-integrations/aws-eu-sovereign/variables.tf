variable "newrelic_account_id" {
  type = string
}

variable "newrelic_region" {
  description = "New Relic region. Must be 'EU' for EU Sovereign cloud integration."
  type        = string
  default     = "EU"

  validation {
    condition     = var.newrelic_region == "EU"
    error_message = "newrelic_region must be 'EU' for EU Sovereign cloud integration. Other regions are not supported."
  }
}

variable "name" {
  type    = string
  default = "production"
}

variable "metric_collection_mode" {
  description = "How metrics are collected. PUSH (metric streams), PULL (API polling), or BOTH"
  type        = string
  default     = "BOTH"

  validation {
    condition     = contains(["PUSH", "PULL", "BOTH"], var.metric_collection_mode)
    error_message = "metric_collection_mode must be 'PUSH', 'PULL', or 'BOTH'."
  }
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
  default     = "opentelemetry1.0"

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