variable "newrelic_account_id" {
  type        = string
  description = "The New Relic account ID to link the GCP projects to."
}

variable "newrelic_api_key" {
  type        = string
  sensitive   = true
  description = "New Relic User API key (starts with NRAK-)."
}

variable "newrelic_region" {
  type        = string
  default     = "US"
  description = "New Relic data-center region: US or EU."
}

variable "gcp_sa_project_id" {
  type        = string
  description = "GCP project in which the service account and WIF pool are created."
}

variable "gcp_folder_id" {
  type        = string
  description = "Numeric GCP folder ID (without the 'folders/' prefix). Folder-level IAM covers all projects in both analytics_projects and compute_projects."
}

variable "wif_pool_id" {
  type        = string
  description = "ID for the Workload Identity Federation pool (e.g. 'newrelic-wif-pool')."
}

variable "wif_provider_id" {
  type        = string
  description = "ID for the WIF OIDC provider inside the pool (e.g. 'newrelic-oidc-provider')."
}

variable "newrelic_sa_name" {
  type        = string
  description = "Name for the GCP service account that New Relic will impersonate."
}

variable "analytics_projects" {
  type        = map(string)
  description = <<-EOT
    Map of display-name => GCP project ID for analytics projects.
    New Relic will monitor BigQuery, PubSub, Spanner, Storage, DataFlow, DataProc.
    The display-name becomes the linked account name in the New Relic UI.
    Example:
      {
        "analytics-prod" = "my-analytics-project-123"
      }
  EOT
}

variable "compute_projects" {
  type        = map(string)
  description = <<-EOT
    Map of display-name => GCP project ID for compute projects.
    New Relic will monitor VMs, SQL, Kubernetes, Cloud Run, Load Balancing, Functions.
    Example:
      {
        "compute-prod" = "my-compute-project-456"
      }
  EOT
}

variable "enable_fetch_tags" {
  type        = bool
  default     = false
  description = "When true, New Relic fetches resource tags for services that support it (BigQuery, PubSub, Spanner, Storage)."
}

variable "metrics_polling_interval" {
  type    = number
  default = 300
  description = <<-EOT
    Polling interval in seconds applied to all enabled services. Default: 300 (5 minutes).
    Limited Preview (LP) note: the following services support a minimum interval of 60 s (1 minute):
      big_query, data_flow, data_proc, load_balancing, pub_sub, spanner
    Set this variable to 60 to enable 1-minute polling.
  EOT
}
