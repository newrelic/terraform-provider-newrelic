# ── New Relic credentials ─────────────────────────────────────────────────────

variable "newrelic_account_id" {
  type        = number
  description = "New Relic account ID to link the GCP project to."
}

variable "newrelic_api_key" {
  type        = string
  sensitive   = true
  description = "New Relic User API key (starts with NRAK-)."
}

variable "newrelic_region" {
  type        = string
  default     = "US"
  description = "New Relic data center region: 'US', 'EU', or 'Staging'."

  validation {
    condition     = contains(["US", "EU", "JP", "Staging"], var.newrelic_region)
    error_message = "newrelic_region must be one of: 'US', 'EU', 'JP', 'Staging'."
  }
}

# ── GCP project ───────────────────────────────────────────────────────────────

variable "gcp_project_id" {
  type        = string
  description = "GCP project ID (e.g. 'my-project-123')."
}

variable "linked_account_name" {
  type        = string
  default     = "production-gcp-v2"
  description = "Display name for this linked GCP account in New Relic."
}

variable "gcp_org_id" {
  type        = string
  description = "GCP organization ID (numeric, e.g. '123456789012'). Required to grant roles/resourcemanager.folderViewer at the org level for folder-level resource discovery."
}

# ── WIF infrastructure IDs ────────────────────────────────────────────────────
# These control the resource IDs created inside GCP for the WIF pool/provider
# and the New Relic service account. Change them if you need to avoid collisions.

variable "wif_pool_id" {
  type        = string
  default     = "newrelic-pool"
  description = "ID for the GCP Workload Identity Pool (alphanumeric and hyphens, 4-32 chars)."
}

variable "wif_provider_id" {
  type        = string
  default     = "newrelic-provider"
  description = "ID for the GCP Workload Identity Pool Provider (alphanumeric and hyphens, 4-32 chars)."
}

variable "newrelic_sa_name" {
  type        = string
  default     = "newrelic-integration"
  description = "account_id portion of the GCP service account email created for New Relic (max 30 chars)."
}

# ── Integration tuning ────────────────────────────────────────────────────────

variable "metrics_polling_interval" {
  type        = number
  default     = 300
  description = "Polling interval in seconds for all enabled integrations (minimum: 60)."

  validation {
    condition     = var.metrics_polling_interval >= 60
    error_message = "metrics_polling_interval must be at least 60 seconds."
  }
}

variable "enable_fetch_tags" {
  type        = bool
  default     = false
  description = "Whether to enable fetch_tags (and fetch_table_metrics for BigQuery). Disable on staging."
}

variable "enabled_services" {
  type        = set(string)
  description = "Set of GCP service names to enable. Omit a service to disable it."
  default = [
    # Existing GCP services
    "ai_platform",
    "alloy_db",
    "app_engine",
    "big_query",
    "big_table",
    "composer",
    "data_flow",
    "data_proc",
    "data_store",
    "firebase_database",
    "firebase_hosting",
    "firebase_storage",
    "firestore",
    "functions",
    "interconnect",
    "kubernetes",
    "load_balancing",
    "mem_cache",
    "pub_sub",
    "redis",
    "router",
    "run",
    "spanner",
    "sql",
    "storage",
    "virtual_machines",
    "vpc_access",
    # GCP v2-only services
    "api_gateway",
    "firebase_auth",
    "firebase_vertex_ai",
    "istio",
    "managed_kafka",
    "memory_store",
    "firebase_app_hosting",
  ]
}
