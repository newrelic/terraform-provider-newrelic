variable "newrelic_account_id" {
  type        = number
  description = "The New Relic account ID to link the GCP project to."
}

variable "newrelic_api_key" {
  type        = string
  sensitive   = true
  description = "A New Relic User API key with the NerdGraph scope."
}

variable "newrelic_region" {
  type        = string
  default     = "US"
  description = "New Relic data-center region: US, EU, or JP. Determines which OIDC issuer URI is used."
}

variable "gcp_project_id" {
  type        = string
  description = "The GCP project ID to monitor (e.g. my-project-123)."
}

variable "gcp_folder_id" {
  type        = string
  description = "Numeric GCP folder ID (without the 'folders/' prefix) where the project lives. Required to grant roles/resourcemanager.folderViewer at the folder level."
}

variable "linked_account_name" {
  type        = string
  default     = "production-gcp-dm"
  description = "Display name shown in the New Relic UI for the linked account."
}

variable "wif_pool_id" {
  type        = string
  default     = "newrelic-pool"
  description = "ID for the Workload Identity Pool created in GCP."
}

variable "wif_provider_id" {
  type        = string
  default     = "newrelic-provider"
  description = "ID for the OIDC provider inside the WIF pool."
}

variable "newrelic_sa_name" {
  type        = string
  default     = "newrelic-integration"
  description = "Name for the GCP service account New Relic impersonates."
}

variable "metrics_polling_interval" {
  type        = number
  default     = 300
  description = <<-EOT
    How often (in seconds) New Relic polls each service for metrics. Default: 300 (5 minutes).
    1-minute polling is in Limited Preview (LP) and available only for:
      alloy_db, big_query, data_flow, data_proc, load_balancing, managed_kafka, pub_sub, spanner
    Set to 60 to enable 1-minute polling for those services.
  EOT
}
