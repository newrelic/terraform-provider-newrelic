variable "newrelic_account_id" {
  type        = number
  description = "The New Relic account ID to link the GCP project to."
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

variable "gcp_project_id" {
  type        = string
  description = "The GCP project ID to integrate with New Relic."
}

variable "linked_account_name" {
  type        = string
  description = "Display name for this linked GCP account in New Relic."
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
  type        = number
  default     = 300
  description = "Polling interval in seconds for services that do not support 1-minute polling."
}

variable "enable_fetch_tags" {
  type        = bool
  default     = false
  description = "Whether to fetch GCP resource tags/labels for supported services."
}

variable "enabled_services" {
  type        = list(string)
  default     = ["big_query", "pub_sub", "storage"]
  description = <<-EOT
    List of GCP services to enable. Supported values:
      alloy_db, ai_platform, api_gateway, app_engine, big_query, big_table,
      composer, data_flow, data_proc, data_store, firebase_app_hosting,
      firebase_auth, firebase_hosting, firebase_storage, firebase_vertex_ai,
      firestore, functions, interconnect, istio, kubernetes, load_balancing,
      managed_kafka, mem_cache, memory_store, pub_sub, redis, router, run,
      spanner, sql, storage, virtual_machines, vpc_access
  EOT
}
