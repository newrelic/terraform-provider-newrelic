variable "newrelic_account_id" {
  type = string
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