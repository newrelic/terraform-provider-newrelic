# GCP Dimensional Metrics — single-project quickstart
#
# Creates the full Workload Identity Federation setup (pool, OIDC provider,
# service account, and IAM bindings) in GCP and links one project to New Relic
# for the GCP Dimensional Metrics integration. Copy it, set the variables, and
# run `terraform init && terraform apply`.
#
# Requires both the newrelic and google providers.

provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

provider "google" {
  project = var.gcp_project_id
}

locals {
  # Derive the OIDC issuer URI based on the New Relic region.
  newrelic_oidc_issuer = {
    "US" = "https://oidc.newrelic.com/r/gcp-cmp"
    "EU" = "https://oidc.eu.newrelic.com/r/gcp-cmp"
    "JP" = "https://oidc.jp.newrelic.com/r/gcp-cmp"
  }[var.newrelic_region]
}

# ── GCP: Workload Identity Federation ──────────────────────────────────────────

resource "google_iam_workload_identity_pool" "newrelic" {
  workload_identity_pool_id = var.wif_pool_id
  display_name              = "New Relic"
  description               = "WIF pool for New Relic GCP Dimensional Metrics integration"
}

resource "google_iam_workload_identity_pool_provider" "newrelic" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.newrelic.workload_identity_pool_id
  workload_identity_pool_provider_id = var.wif_provider_id
  display_name                       = "New Relic OIDC provider"

  oidc {
    # GCP fetches {issuer_uri}/.well-known/openid-configuration to validate tokens.
    issuer_uri        = local.newrelic_oidc_issuer
    allowed_audiences = ["newrelic-gcp-integrations"]
  }

  attribute_mapping = {
    "google.subject"          = "assertion.sub"
    "attribute.nr_account_id" = "assertion.nr_account_id"
  }

  # Restrict impersonation to tokens issued for this specific New Relic account.
  attribute_condition = "assertion.nr_account_id == \"${var.newrelic_account_id}\""
}

resource "google_service_account" "newrelic" {
  account_id   = var.newrelic_sa_name
  display_name = "New Relic Integration"
  description  = "Impersonated by New Relic via WIF to collect GCP metrics"
}

# Read access (metrics collection)
resource "google_project_iam_member" "newrelic_viewer" {
  project = var.gcp_project_id
  role    = "roles/viewer"
  member  = "serviceAccount:${google_service_account.newrelic.email}"
}

# API quota
resource "google_project_iam_member" "newrelic_service_usage" {
  project = var.gcp_project_id
  role    = "roles/serviceusage.serviceUsageConsumer"
  member  = "serviceAccount:${google_service_account.newrelic.email}"
}

# Resource discovery (Cloud Asset search/list APIs)
resource "google_project_iam_member" "newrelic_cloud_asset_viewer" {
  project = var.gcp_project_id
  role    = "roles/cloudasset.viewer"
  member  = "serviceAccount:${google_service_account.newrelic.email}"
}

# Folder-level resource discovery — must be granted at the folder level, not project level
resource "google_folder_iam_member" "newrelic_folder_viewer" {
  folder = "folders/${var.gcp_folder_id}"
  role   = "roles/resourcemanager.folderViewer"
  member = "serviceAccount:${google_service_account.newrelic.email}"
}

# Allow New Relic's WIF pool (scoped to this account) to impersonate the SA.
# Must use the attribute-scoped principalSet — the wildcard form does NOT grant
# iam.serviceAccounts.getAccessToken on the service account.
resource "google_service_account_iam_member" "newrelic_wif" {
  service_account_id = google_service_account.newrelic.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.newrelic.name}/attribute.nr_account_id/${var.newrelic_account_id}"
}

# ── New Relic: Step 1 — link GCP project ──────────────────────────────────────

resource "newrelic_cloud_gcp_link_account" "main" {
  account_id = var.newrelic_account_id
  name       = var.linked_account_name
  project_id = var.gcp_project_id

  # Opt into keyless GCP Dimensional Metrics linking; the provider builds the WIF
  # credential JSON internally from audience + service_account_email.
  use_workload_identity_federation = true
  audience                         = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"
  service_account_email            = google_service_account.newrelic.email

  depends_on = [
    google_project_iam_member.newrelic_viewer,
    google_project_iam_member.newrelic_service_usage,
    google_project_iam_member.newrelic_cloud_asset_viewer,
    google_folder_iam_member.newrelic_folder_viewer,
    google_service_account_iam_member.newrelic_wif,
  ]
}

# ── New Relic: Step 2 — configure which services to poll ──────────────────────

resource "newrelic_cloud_gcp_dm_integrations" "main" {
  account_id        = newrelic_cloud_gcp_link_account.main.account_id
  linked_account_id = newrelic_cloud_gcp_link_account.main.id

  # All GCP services default to 300 s polling. 1-minute polling is in Limited Preview (LP)
  # and available only for: alloy_db, big_query, data_flow, data_proc, load_balancing,
  # managed_kafka, pub_sub, spanner — set metrics_polling_interval = 60 to enable.
  ai_platform { metrics_polling_interval = var.metrics_polling_interval }
  alloy_db { metrics_polling_interval = var.metrics_polling_interval }
  api_gateway { metrics_polling_interval = var.metrics_polling_interval } # DM only
  app_engine { metrics_polling_interval = var.metrics_polling_interval }
  big_query { metrics_polling_interval = var.metrics_polling_interval }
  big_table { metrics_polling_interval = var.metrics_polling_interval }
  composer { metrics_polling_interval = var.metrics_polling_interval }
  data_flow { metrics_polling_interval = var.metrics_polling_interval }
  data_proc { metrics_polling_interval = var.metrics_polling_interval }
  data_store { metrics_polling_interval = var.metrics_polling_interval }
  firebase_database { metrics_polling_interval = var.metrics_polling_interval }
  firebase_hosting { metrics_polling_interval = var.metrics_polling_interval }
  firebase_storage { metrics_polling_interval = var.metrics_polling_interval }
  firestore { metrics_polling_interval = var.metrics_polling_interval }
  functions { metrics_polling_interval = var.metrics_polling_interval }
  interconnect { metrics_polling_interval = var.metrics_polling_interval }
  istio { metrics_polling_interval = var.metrics_polling_interval }      # DM only, metrics only
  kubernetes { metrics_polling_interval = var.metrics_polling_interval } # metrics only, no entity support
  load_balancing { metrics_polling_interval = var.metrics_polling_interval }
  mem_cache { metrics_polling_interval = var.metrics_polling_interval }
  pub_sub { metrics_polling_interval = var.metrics_polling_interval }
  redis { metrics_polling_interval = var.metrics_polling_interval }
  router { metrics_polling_interval = var.metrics_polling_interval }
  run { metrics_polling_interval = var.metrics_polling_interval }
  spanner { metrics_polling_interval = var.metrics_polling_interval }
  sql { metrics_polling_interval = var.metrics_polling_interval }
  storage { metrics_polling_interval = var.metrics_polling_interval }
  virtual_machines { metrics_polling_interval = var.metrics_polling_interval }
  vpc_access { metrics_polling_interval = var.metrics_polling_interval }

  # GCP Dimensional Metrics-only services
  firebase_app_hosting { metrics_polling_interval = var.metrics_polling_interval } # metrics only
  firebase_auth { metrics_polling_interval = var.metrics_polling_interval }
  firebase_vertex_ai { metrics_polling_interval = var.metrics_polling_interval } # metrics only
  managed_kafka { metrics_polling_interval = var.metrics_polling_interval }
  memory_store { metrics_polling_interval = var.metrics_polling_interval }
}

# ── Outputs ────────────────────────────────────────────────────────────────────

output "linked_account_id" {
  description = "New Relic linked account ID for the GCP DM integration."
  value       = newrelic_cloud_gcp_link_account.main.id
}
