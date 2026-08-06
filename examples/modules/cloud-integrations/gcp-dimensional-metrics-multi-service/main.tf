# Scenario: GCP Dimensional Metrics — different services per project group
#
# Creates shared GCP infrastructure (WIF pool, OIDC provider, service account,
# and folder-level IAM) used across two distinct project groups that each have
# different monitoring needs:
#   - analytics_projects: BigQuery, PubSub, Spanner, Storage, DataFlow, DataProc
#   - compute_projects:   VMs, SQL, Cloud Run, Load Balancing, Functions, Kubernetes (metrics only, no
#                         entity support; also supports 1-minute polling in Limited Preview)
#
# Use this when: your GCP projects serve different workloads and you want to
# enable only the relevant GCP services in New Relic per group, while sharing
# a single WIF identity setup to avoid duplicating GCP infrastructure.

# ── Locals ────────────────────────────────────────────────────────────────────
locals {
  oidc_issuer_uri = {
    "US" = "https://oidc.newrelic.com/r/gcp-cmp"
    "EU" = "https://oidc.eu.newrelic.com/r/gcp-cmp"
    "JP" = "https://oidc.jp.newrelic.com/r/gcp-cmp"
  }[var.newrelic_region]
}

# ── Providers ─────────────────────────────────────────────────────────────────
provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

provider "google" {
  project = var.gcp_sa_project_id
}

# ── Shared GCP Infrastructure ─────────────────────────────────────────────────
# A single WIF pool, OIDC provider, and service account serve both project
# groups. IAM is granted at the folder level so all projects are covered.

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

resource "google_service_account" "newrelic" {
  account_id   = var.newrelic_sa_name
  display_name = "New Relic Integration"
  description  = "Impersonated by New Relic via WIF to collect GCP Dimensional Metrics"
}

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
  create_duration = "120s"
  depends_on = [
    google_service_account_iam_member.newrelic_wif,
    google_folder_iam_member.newrelic_monitoring_viewer,
    google_folder_iam_member.newrelic_service_usage,
    google_folder_iam_member.newrelic_cloud_asset_viewer,
    google_folder_iam_member.newrelic_folder_viewer,
  ]
}

# ── Group 1: Analytics Projects ───────────────────────────────────────────────
# Enabled: BigQuery, PubSub, Spanner, Storage, DataFlow, DataProc

resource "newrelic_cloud_gcp_link_account" "analytics" {
  for_each = var.analytics_projects

  account_id                       = tonumber(var.newrelic_account_id)
  name                             = each.key
  project_id                       = each.value
  use_workload_identity_federation = true
  service_account_email            = google_service_account.newrelic.email
  audience                         = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [time_sleep.iam_propagation]
}

resource "newrelic_cloud_gcp_dm_integrations" "analytics" {
  for_each = newrelic_cloud_gcp_link_account.analytics

  account_id        = tonumber(var.newrelic_account_id)
  linked_account_id = tonumber(each.value.id)

  big_query {
    metrics_polling_interval = var.metrics_polling_interval
  }
  data_flow {
    metrics_polling_interval = var.metrics_polling_interval
  }
  data_proc {
    metrics_polling_interval = var.metrics_polling_interval
  }
  pub_sub {
    metrics_polling_interval = var.metrics_polling_interval
  }
  spanner {
    metrics_polling_interval = var.metrics_polling_interval
  }
  storage {
    metrics_polling_interval = var.metrics_polling_interval
  }
}

# ── Group 2: Compute Projects ─────────────────────────────────────────────────
# Enabled: VMs, SQL, Cloud Run, Load Balancing, Cloud Functions, Kubernetes (metrics only, no entity
# support; also supports 1-minute polling in Limited Preview)

resource "newrelic_cloud_gcp_link_account" "compute" {
  for_each = var.compute_projects

  account_id                       = tonumber(var.newrelic_account_id)
  name                             = each.key
  project_id                       = each.value
  use_workload_identity_federation = true
  service_account_email            = google_service_account.newrelic.email
  audience                         = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [time_sleep.iam_propagation]
}

resource "newrelic_cloud_gcp_dm_integrations" "compute" {
  for_each = newrelic_cloud_gcp_link_account.compute

  account_id        = tonumber(var.newrelic_account_id)
  linked_account_id = tonumber(each.value.id)

  functions {
    metrics_polling_interval = var.metrics_polling_interval
  }
  kubernetes {
    metrics_polling_interval = var.metrics_polling_interval
  }
  load_balancing {
    metrics_polling_interval = var.metrics_polling_interval
  }
  run {
    metrics_polling_interval = var.metrics_polling_interval
  }
  sql {
    metrics_polling_interval = var.metrics_polling_interval
  }
  virtual_machines {
    metrics_polling_interval = var.metrics_polling_interval
  }
}

# ── Outputs ───────────────────────────────────────────────────────────────────
output "analytics_linked_account_ids" {
  description = "Map of display-name => New Relic linked account ID for analytics projects."
  value       = { for k, v in newrelic_cloud_gcp_link_account.analytics : k => v.id }
}

output "compute_linked_account_ids" {
  description = "Map of display-name => New Relic linked account ID for compute projects."
  value       = { for k, v in newrelic_cloud_gcp_link_account.compute : k => v.id }
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
