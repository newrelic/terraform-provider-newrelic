# Scenario: GCP Dimensional Metrics — multi-project with project-level IAM
#
# Creates a single WIF pool, OIDC provider, and service account in a
# designated SA project. Four IAM roles are bound directly on each project
# listed in gcp_projects — no GCP folder access required.
# One New Relic linked account is created for each entry in gcp_projects.
#
# Use this when: you do not have folder-level IAM access, or you are
# monitoring a small number of GCP projects without a shared folder.

locals {
  on = toset(var.enabled_services)

  # Derive the OIDC issuer URI based on the New Relic region.
  oidc_issuer_uri = {
    "US" = "https://oidc.newrelic.com/r/gcp-cmp"
    "EU" = "https://oidc.eu.newrelic.com/r/gcp-cmp"
    "JP" = "https://oidc.jp.newrelic.com/r/gcp-cmp"
  }[var.newrelic_region]
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

# ── Workload Identity Federation: Pool ────────────────────────────────────────
resource "google_iam_workload_identity_pool" "newrelic" {
  workload_identity_pool_id = var.wif_pool_id
  display_name              = "New Relic"
  description               = "WIF pool for the New Relic GCP Dimensional Metrics integration"
}

# ── Workload Identity Federation: OIDC Provider ───────────────────────────────
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
    issuer_uri        = local.oidc_issuer_uri
    allowed_audiences = ["newrelic-gcp-integrations"]
  }
}

# ── GCP Service Account ───────────────────────────────────────────────────────
resource "google_service_account" "newrelic" {
  account_id   = var.newrelic_sa_name
  display_name = "New Relic Integration"
  description  = "Impersonated by New Relic via WIF to collect GCP Dimensional Metrics"
}

# ── IAM: Project-level bindings — one set per project in gcp_projects ─────────
# Use this module when you do not have GCP folder-level IAM access.
# Each project in gcp_projects receives 4 IAM bindings directly.
resource "google_project_iam_member" "newrelic_monitoring_viewer" {
  for_each = var.gcp_projects
  project  = each.value
  role     = "roles/viewer"
  member   = google_service_account.newrelic.member
}

resource "google_project_iam_member" "newrelic_service_usage" {
  for_each = var.gcp_projects
  project  = each.value
  role     = "roles/serviceusage.serviceUsageConsumer"
  member   = google_service_account.newrelic.member
}

resource "google_project_iam_member" "newrelic_cloud_asset_viewer" {
  for_each = var.gcp_projects
  project  = each.value
  role     = "roles/cloudasset.viewer"
  member   = google_service_account.newrelic.member
}

resource "google_project_iam_member" "newrelic_browser_viewer" {
  for_each = var.gcp_projects
  project  = each.value
  role     = "roles/browser"
  member   = google_service_account.newrelic.member
}

# ── IAM: Allow WIF pool to impersonate the service account ────────────────────
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
    google_project_iam_member.newrelic_monitoring_viewer,
    google_project_iam_member.newrelic_service_usage,
    google_project_iam_member.newrelic_cloud_asset_viewer,
    google_project_iam_member.newrelic_browser_viewer,
  ]
}

# ── New Relic: Link one account per GCP project ───────────────────────────────
resource "newrelic_cloud_gcp_link_account" "this" {
  for_each = var.gcp_projects

  account_id                       = var.newrelic_account_id
  name                             = each.key
  project_id                       = each.value
  use_workload_identity_federation = true
  service_account_email            = google_service_account.newrelic.email
  audience                         = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [time_sleep.iam_propagation]
}

# ── New Relic: Enable Integrations ────────────────────────────────────────────
resource "newrelic_cloud_gcp_dm_integrations" "this" {
  for_each = newrelic_cloud_gcp_link_account.this

  account_id        = var.newrelic_account_id
  linked_account_id = tonumber(each.value.id)

  # All services use metrics_polling_interval (default 300 s / 5 min).
  # 1-minute polling is in Limited Preview (LP) for the following services:
  #   alloy_db, big_query, data_flow, data_proc, load_balancing,
  #   managed_kafka, pub_sub, spanner
  # To enable 1-minute polling for those services, set:
  #   metrics_polling_interval = 60
  dynamic "ai_platform" {
    for_each = contains(local.on, "ai_platform") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "alloy_db" {
    for_each = contains(local.on, "alloy_db") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "api_gateway" {
    for_each = contains(local.on, "api_gateway") ? [1] : []
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
  dynamic "istio" {
    for_each = contains(local.on, "istio") ? [1] : []
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
  description = "Map of display-name => New Relic linked account ID for each linked GCP project."
  value       = { for k, v in newrelic_cloud_gcp_link_account.this : k => v.id }
}

output "wif_pool_name" {
  description = "The full resource name of the WIF pool."
  value       = google_iam_workload_identity_pool.newrelic.name
}

output "wif_provider_name" {
  description = "The full resource name of the WIF OIDC provider."
  value       = google_iam_workload_identity_pool_provider.newrelic.name
}

output "newrelic_service_account_email" {
  description = "Email of the GCP service account impersonated by New Relic."
  value       = google_service_account.newrelic.email
}
