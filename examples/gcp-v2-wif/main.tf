terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.10"
    }
  }
}

variable "newrelic_account_id" { type = string }
variable "newrelic_api_key" {
  type      = string
  sensitive = true
}
variable "newrelic_region" {
  type    = string
  default = "US"
}
variable "gcp_sa_project_id" {
  type        = string
  description = "GCP project in which the service account and WIF pool are created (e.g. cmp-dev-proj-1)."
}
variable "gcp_folder_id" {
  type        = string
  description = "Numeric folder ID (without 'folders/' prefix). All 4 IAM roles are granted at this folder level, covering every project under it."
}
variable "gcp_projects" {
  type        = map(string)
  description = "Map of display-name => GCP project ID for each project to link. e.g. { \"proj-1\" = \"cmp-dev-proj-1\", \"proj-2\" = \"cmp-dev-proj-2\" }"
}
variable "wif_pool_id"      { type = string }
variable "wif_provider_id"  { type = string }
variable "newrelic_sa_name" { type = string }
variable "metrics_polling_interval" {
  type    = number
  default = 300
}
variable "enable_fetch_tags" {
  type    = bool
  default = false
}
variable "enabled_services" {
  type    = list(string)
  default = []
}

locals {
  on = toset(var.enabled_services)
}

provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

# SA and WIF pool are always created in the designated SA project.
provider "google" {
  project = var.gcp_sa_project_id
}

# ── WIF Pool ──────────────────────────────────────────────────────────────────
resource "google_iam_workload_identity_pool" "newrelic" {
  workload_identity_pool_id = var.wif_pool_id
  display_name              = "New Relic"
  description               = "WIF pool for New Relic GCP Dimensional Metrics integration"
}

# ── WIF Provider ──────────────────────────────────────────────────────────────
resource "google_iam_workload_identity_pool_provider" "newrelic" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.newrelic.workload_identity_pool_id
  workload_identity_pool_provider_id = var.wif_provider_id
  display_name                       = "New Relic OIDC provider"

  attribute_mapping = {
    "google.subject"          = "assertion.sub"
    "attribute.nr_account_id" = "assertion.nr_account_id"
  }
  attribute_condition = "assertion.nr_account_id == \"${var.newrelic_account_id}\""

  oidc {
    issuer_uri        = "https://oidc-staging.newrelic.com/r/gcp-cmp"
    allowed_audiences = ["newrelic-gcp-integrations"]
  }
}

# ── Service Account ───────────────────────────────────────────────────────────
resource "google_service_account" "newrelic" {
  account_id   = var.newrelic_sa_name
  display_name = "New Relic Integration"
  description  = "Impersonated by New Relic via WIF to collect GCP Dimensional Metrics"
}

# ── Folder-level IAM: all 4 roles cover every project under the folder ────────
resource "google_folder_iam_member" "newrelic_monitoring_viewer" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/monitoring.viewer"
  member = google_service_account.newrelic.member
}

resource "google_folder_iam_member" "newrelic_service_usage" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/serviceusage.serviceUsageConsumer"
  member = google_service_account.newrelic.member
}

resource "google_folder_iam_member" "newrelic_cloud_asset_viewer" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/cloudasset.viewer"
  member = google_service_account.newrelic.member
}

resource "google_folder_iam_member" "newrelic_folder_viewer" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/resourcemanager.folderViewer"
  member = google_service_account.newrelic.member
}

# ── WIF impersonation binding ─────────────────────────────────────────────────
resource "google_service_account_iam_member" "newrelic_wif" {
  service_account_id = google_service_account.newrelic.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.newrelic.name}/attribute.nr_account_id/${var.newrelic_account_id}"
}

# ── IAM propagation delay ─────────────────────────────────────────────────────
# GCP IAM bindings take ~60-90 s to propagate globally. Without this wait,
# cloudAuthenticateIntegration returns 403 on iam.serviceAccounts.getAccessToken.
resource "time_sleep" "iam_propagation" {
  create_duration = "90s"
  depends_on = [
    google_service_account_iam_member.newrelic_wif,
    google_folder_iam_member.newrelic_monitoring_viewer,
    google_folder_iam_member.newrelic_service_usage,
    google_folder_iam_member.newrelic_cloud_asset_viewer,
    google_folder_iam_member.newrelic_folder_viewer,
  ]
}

# ── New Relic: one linked account per GCP project ─────────────────────────────
resource "newrelic_cloud_gcp_dm_link_account" "this" {
  for_each = var.gcp_projects

  account_id            = tonumber(var.newrelic_account_id)
  name                  = each.key
  project_id            = each.value
  service_account_email = google_service_account.newrelic.email
  audience              = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [time_sleep.iam_propagation]
}

# ── New Relic: integrations per linked account ────────────────────────────────
resource "newrelic_cloud_gcp_dm_integrations" "this" {
  for_each = newrelic_cloud_gcp_dm_link_account.this

  account_id        = tonumber(var.newrelic_account_id)
  linked_account_id = tonumber(each.value.id)

  # All services use metrics_polling_interval (default 300 s / 5 min).
  # LP (low-polling) services support down to 60 s:
  #   alloy_db, big_query, data_flow, data_proc, load_balancing,
  #   managed_kafka, pub_sub, spanner
  dynamic "ai_platform" {
    for_each = contains(local.on, "ai_platform") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "alloy_db" {
    for_each = contains(local.on, "alloy_db") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "app_engine" {
    for_each = contains(local.on, "app_engine") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "big_query" {
    for_each = contains(local.on, "big_query") ? [1] : []
    content {
      metrics_polling_interval = var.metrics_polling_interval
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "big_table" {
    for_each = contains(local.on, "big_table") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "composer" {
    for_each = contains(local.on, "composer") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "data_flow" {
    for_each = contains(local.on, "data_flow") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "data_proc" {
    for_each = contains(local.on, "data_proc") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "data_store" {
    for_each = contains(local.on, "data_store") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_app_hosting" {
    for_each = contains(local.on, "firebase_app_hosting") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_auth" {
    for_each = contains(local.on, "firebase_auth") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_database" {
    for_each = contains(local.on, "firebase_database") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_hosting" {
    for_each = contains(local.on, "firebase_hosting") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_storage" {
    for_each = contains(local.on, "firebase_storage") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firebase_vertex_ai" {
    for_each = contains(local.on, "firebase_vertex_ai") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "firestore" {
    for_each = contains(local.on, "firestore") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "functions" {
    for_each = contains(local.on, "functions") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "interconnect" {
    for_each = contains(local.on, "interconnect") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "kubernetes" {
    for_each = contains(local.on, "kubernetes") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "load_balancing" {
    for_each = contains(local.on, "load_balancing") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "managed_kafka" {
    for_each = contains(local.on, "managed_kafka") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "mem_cache" {
    for_each = contains(local.on, "mem_cache") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "memory_store" {
    for_each = contains(local.on, "memory_store") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "pub_sub" {
    for_each = contains(local.on, "pub_sub") ? [1] : []
    content {
      metrics_polling_interval = var.metrics_polling_interval
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "redis" {
    for_each = contains(local.on, "redis") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "router" {
    for_each = contains(local.on, "router") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "run" {
    for_each = contains(local.on, "run") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "spanner" {
    for_each = contains(local.on, "spanner") ? [1] : []
    content {
      metrics_polling_interval = var.metrics_polling_interval
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "sql" {
    for_each = contains(local.on, "sql") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "storage" {
    for_each = contains(local.on, "storage") ? [1] : []
    content {
      metrics_polling_interval = var.metrics_polling_interval
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "virtual_machines" {
    for_each = contains(local.on, "virtual_machines") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "vpc_access" {
    for_each = contains(local.on, "vpc_access") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
}

# ── Outputs ───────────────────────────────────────────────────────────────────
output "linked_account_ids" {
  description = "Map of display-name => New Relic linked account ID."
  value       = { for k, v in newrelic_cloud_gcp_dm_link_account.this : k => v.id }
}

output "wif_pool_name" {
  value = google_iam_workload_identity_pool.newrelic.name
}

output "wif_provider_name" {
  value = google_iam_workload_identity_pool_provider.newrelic.name
}

output "newrelic_service_account_email" {
  value = google_service_account.newrelic.email
}
