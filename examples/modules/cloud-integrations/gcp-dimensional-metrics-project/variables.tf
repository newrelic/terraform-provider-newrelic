variable "newrelic_account_id" {
  type        = number
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
  description = "GCP project in which the service account and WIF pool are created. The SA will be used to collect metrics from all projects in gcp_projects."
}

variable "gcp_projects" {
  type        = map(string)
  description = <<-EOT
    Map of display-name => GCP project ID for each project to link to New Relic.
    The display-name becomes the linked account name shown in the New Relic UI.
    IAM bindings are created directly on each project — no folder access required.
    Example:
      {
        "prod-payments" = "my-payments-project-123"
        "prod-analytics" = "my-analytics-project-456"
      }
  EOT
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
  description = "Name for the GCP service account that New Relic will impersonate (e.g. 'newrelic-integration')."
}

variable "metrics_polling_interval" {
  type    = number
  default = 300
  description = <<-EOT
    Polling interval in seconds applied to all enabled services. Default: 300 (5 minutes).
    Low-polling (LP) note: the following services support a minimum interval of 60 s (1 minute):
      alloy_db, big_query, data_flow, data_proc, load_balancing, managed_kafka, pub_sub, spanner
    Set this variable to 60 to enable 1-minute polling across all services.
  EOT
}

variable "enable_fetch_tags" {
  type        = bool
  default     = false
  description = "Whether to fetch GCP resource tags/labels for supported services (big_query, pub_sub, spanner, storage)."
}

variable "enabled_services" {
  type        = list(string)
  default     = ["big_query", "pub_sub", "storage"]
  description = <<-EOT
    List of GCP services to enable. Supported values:
      ai_platform, alloy_db, api_gateway, app_engine, big_query, big_table,
      composer, data_flow, data_proc, data_store, firebase_app_hosting,
      firebase_auth, firebase_database, firebase_hosting, firebase_storage,
      firebase_vertex_ai, firestore, functions, interconnect, istio,
      kubernetes, load_balancing, managed_kafka, mem_cache, memory_store,
      pub_sub, redis, router, run, spanner, sql, storage,
      virtual_machines, vpc_access

    Services marked LP (low-polling) support metrics_polling_interval = 60:
      alloy_db, big_query, data_flow, data_proc, load_balancing,
      managed_kafka, pub_sub, spanner
  EOT
}
