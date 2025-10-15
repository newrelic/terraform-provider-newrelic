variable "account_id" {
  description = "New Relic Account ID"
  type        = string
}

variable "description" {
  description = "Description of the drop rule"
  type        = string
}

variable "action" {
  description = "Action to perform (drop_data, drop_attributes, drop_attributes_from_metric_aggregates)"
  type        = string
  validation {
    condition = contains([
      "drop_data",
      "drop_attributes",
      "drop_attributes_from_metric_aggregates"
    ], var.action)
    error_message = "Action must be one of: drop_data, drop_attributes, drop_attributes_from_metric_aggregates."
  }
}

variable "nrql" {
  description = "NRQL query for the drop rule"
  type        = string
}