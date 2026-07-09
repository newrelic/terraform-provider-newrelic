# Example: GCP Dimensional Metrics — all services enabled
#
# Creates WIF infrastructure (pool, OIDC provider, service account) in a
# designated SA project and binds IAM at the folder level so every project
# under that folder is covered without per-project changes.
#
# All 34 GCP services are enabled.  The following Limited Preview (LP) services support
# 1-minute polling and are set to 60 s; the rest run at 300 s:
#   alloy_db, big_query, data_flow, data_proc, load_balancing,
#   managed_kafka, pub_sub, spanner

locals {
  oidc_issuer_uri = (var.newrelic_region == "EU"
    ? "https://oidc.eu.newrelic.com/r/gcp-cmp"
    : var.newrelic_region == "Staging"
    ? "https://oidc-staging.newrelic.com/r/gcp-cmp"
    : "https://oidc.newrelic.com/r/gcp-cmp")
}

provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

provider "google" { project = var.gcp_sa_project_id }

# ── Workload Identity Federation ──────────────────────────────────────────────

resource "google_iam_workload_identity_pool" "newrelic" {
  workload_identity_pool_id = var.wif_pool_id
  display_name              = "New Relic"
  description               = "WIF pool for New Relic GCP Dimensional Metrics integration"
}

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

# ── IAM: Folder-level bindings ────────────────────────────────────────────────

resource "google_folder_iam_member" "newrelic_monitoring_viewer" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/viewer"
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

resource "google_service_account_iam_member" "newrelic_wif" {
  service_account_id = google_service_account.newrelic.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.newrelic.name}/attribute.nr_account_id/${var.newrelic_account_id}"
}

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

# ── New Relic: Link GCP projects ──────────────────────────────────────────────

resource "newrelic_cloud_gcp_dm_link_account" "this" {
  for_each = var.gcp_projects

  account_id            = var.newrelic_account_id
  name                  = each.key
  project_id            = each.value
  service_account_email = google_service_account.newrelic.email
  audience              = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [time_sleep.iam_propagation]
}

# ── New Relic: Enable all GCP integrations ────────────────────────────────────
# LP = Limited Preview services (60 s min); all other services use 300 s (5 min).

resource "newrelic_cloud_gcp_dm_integrations" "this" {
  for_each = newrelic_cloud_gcp_dm_link_account.this

  account_id        = var.newrelic_account_id
  linked_account_id = tonumber(each.value.id)

  # ── 300 s services ───────────────────────────────────────────────────────────
  ai_platform {
    metrics_polling_interval = 300
  }
  api_gateway {
    metrics_polling_interval = 300
  }
  app_engine {
    metrics_polling_interval = 300
  }
  big_table {
    metrics_polling_interval = 300
  }
  composer {
    metrics_polling_interval = 300
  }
  data_store {
    metrics_polling_interval = 300
  }
  firebase_app_hosting {
    metrics_polling_interval = 300
  }
  firebase_auth {
    metrics_polling_interval = 300
  }
  firebase_database {
    metrics_polling_interval = 300
  }
  firebase_hosting {
    metrics_polling_interval = 300
  }
  firebase_storage {
    metrics_polling_interval = 300
  }
  firebase_vertex_ai {
    metrics_polling_interval = 300
  }
  firestore {
    metrics_polling_interval = 300
  }
  functions {
    metrics_polling_interval = 300
  }
  interconnect {
    metrics_polling_interval = 300
  }
  istio {
    metrics_polling_interval = 300
  }
  kubernetes {
    metrics_polling_interval = 300
  }
  mem_cache {
    metrics_polling_interval = 300
  }
  memory_store {
    metrics_polling_interval = 300
  }
  redis {
    metrics_polling_interval = 300
  }
  router {
    metrics_polling_interval = 300
  }
  run {
    metrics_polling_interval = 300
  }
  sql {
    metrics_polling_interval = 300
  }
  storage {
    metrics_polling_interval = 300
    fetch_tags               = var.enable_fetch_tags
  }
  virtual_machines {
    metrics_polling_interval = 300
  }
  vpc_access {
    metrics_polling_interval = 300
  }

  # ── 60 s services (LP = Limited Preview) ─────────────────────────────────────
  alloy_db {
    metrics_polling_interval = 60
  }
  big_query {
    metrics_polling_interval = 60
    fetch_tags               = var.enable_fetch_tags
  }
  data_flow {
    metrics_polling_interval = 60
  }
  data_proc {
    metrics_polling_interval = 60
  }
  load_balancing {
    metrics_polling_interval = 60
  }
  managed_kafka {
    metrics_polling_interval = 60
  }
  pub_sub {
    metrics_polling_interval = 60
    fetch_tags               = var.enable_fetch_tags
  }
  spanner {
    metrics_polling_interval = 60
    fetch_tags               = var.enable_fetch_tags
  }
}

# ── Outputs ───────────────────────────────────────────────────────────────────

output "linked_account_ids" {
  description = "Map of display-name => New Relic linked account ID for each linked GCP project."
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
