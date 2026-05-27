variable "newrelic_account_id" {
  type        = number
  description = "The New Relic account ID to link the GCP project to."
}

variable "name" {
  type        = string
  default     = "production-gcp-v2"
  description = "Display name for this linked GCP account in New Relic."
}

variable "gcp_project_id" {
  type        = string
  description = "The GCP project ID to link (e.g. 'my-gcp-project-123')."
}

variable "wif_credential" {
  type        = string
  sensitive   = true
  description = "The Workload Identity Federation credential JSON string exported from GCP."
}

variable "metrics_polling_interval" {
  type        = number
  default     = 300
  description = "Default polling interval in seconds for all enabled integrations."
}
