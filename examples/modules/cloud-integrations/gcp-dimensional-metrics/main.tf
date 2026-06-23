locals {
  on = toset(var.enabled_services)

  # Derive the OIDC issuer URI based on the New Relic region.
  oidc_issuer_uri = var.newrelic_region == "EU" ? "https://oidc.eu.newrelic.com/r/gcp-cmp" : "https://oidc.newrelic.com/r/gcp-cmp"
}

provider "newrelic" {
  account_id = var.newrelic_account_id
  api_key    = var.newrelic_api_key
  region     = var.newrelic_region
}

provider "google" {
  project = var.gcp_project_id
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

# ── IAM: Grant service account monitoring and asset-viewer access ─────────────
resource "google_project_iam_member" "newrelic_viewer" {
  project = var.gcp_project_id
  role    = "roles/monitoring.viewer"
  member  = google_service_account.newrelic.member
}

resource "google_project_iam_member" "newrelic_service_usage" {
  project = var.gcp_project_id
  role    = "roles/serviceusage.serviceUsageConsumer"
  member  = google_service_account.newrelic.member
}

resource "google_project_iam_member" "newrelic_cloud_asset_viewer" {
  project = var.gcp_project_id
  role    = "roles/cloudasset.viewer"
  member  = google_service_account.newrelic.member
}

# ── IAM: Folder-level resource discovery — must be at org level, not project ──
resource "google_organization_iam_member" "newrelic_folder_viewer" {
  org_id = var.gcp_org_id
  role   = "roles/resourcemanager.folderViewer"
  member = google_service_account.newrelic.member
}

# ── IAM: Allow WIF pool to impersonate the service account ────────────────────
resource "google_service_account_iam_member" "newrelic_wif" {
  service_account_id = google_service_account.newrelic.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.newrelic.name}/attribute.nr_account_id/${var.newrelic_account_id}"
}

# ── New Relic: Link GCP Project ───────────────────────────────────────────────
resource "newrelic_cloud_gcp_dm_link_account" "this" {
  account_id            = var.newrelic_account_id
  name                  = var.linked_account_name
  project_id            = var.gcp_project_id
  service_account_email = google_service_account.newrelic.email
  audience              = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"

  depends_on = [
    google_project_iam_member.newrelic_viewer,
    google_project_iam_member.newrelic_service_usage,
    google_project_iam_member.newrelic_cloud_asset_viewer,
    google_organization_iam_member.newrelic_folder_viewer,
    google_service_account_iam_member.newrelic_wif,
  ]
}

# ── New Relic: Enable Integrations ────────────────────────────────────────────
resource "newrelic_cloud_gcp_dm_integrations" "this" {
  account_id        = var.newrelic_account_id
  linked_account_id = tonumber(newrelic_cloud_gcp_dm_link_account.this.id)

  # Services with confirmed 1-minute polling support
  dynamic "data_flow" {
    for_each = contains(local.on, "data_flow") ? [1] : []
    content { metrics_polling_interval = 60 }
  }
  dynamic "data_proc" {
    for_each = contains(local.on, "data_proc") ? [1] : []
    content { metrics_polling_interval = 60 }
  }
  dynamic "spanner" {
    for_each = contains(local.on, "spanner") ? [1] : []
    content {
      metrics_polling_interval = 60
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "managed_kafka" {
    for_each = contains(local.on, "managed_kafka") ? [1] : []
    content { metrics_polling_interval = 60 }
  }
  dynamic "alloy_db" {
    for_each = contains(local.on, "alloy_db") ? [1] : []
    content { metrics_polling_interval = 60 }
  }
  dynamic "big_query" {
    for_each = contains(local.on, "big_query") ? [1] : []
    content {
      metrics_polling_interval = 60
      fetch_tags               = var.enable_fetch_tags
    }
  }
  dynamic "load_balancing" {
    for_each = contains(local.on, "load_balancing") ? [1] : []
    content { metrics_polling_interval = 60 }
  }
  dynamic "pub_sub" {
    for_each = contains(local.on, "pub_sub") ? [1] : []
    content {
      metrics_polling_interval = 60
      fetch_tags               = var.enable_fetch_tags
    }
  }

  # Standard-polling services
  dynamic "ai_platform" {
    for_each = contains(local.on, "ai_platform") ? [1] : []
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
  dynamic "big_table" {
    for_each = contains(local.on, "big_table") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "composer" {
    for_each = contains(local.on, "composer") ? [1] : []
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
  dynamic "mem_cache" {
    for_each = contains(local.on, "mem_cache") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
  dynamic "memory_store" {
    for_each = contains(local.on, "memory_store") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
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
output "linked_account_id" {
  description = "The New Relic linked account ID for this GCP project."
  value       = newrelic_cloud_gcp_dm_link_account.this.id
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
