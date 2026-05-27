locals {
  on = var.enabled_services
}

# ─────────────────────────────────────────────────────────────────────────────
# GCP: Workload Identity Federation infrastructure
#
# Creates a WIF pool + OIDC provider backed by New Relic's OIDC endpoint,
# a service account for New Relic to impersonate, and the IAM bindings that
# connect them.  The pool/provider/SA IDs are variables so multiple environments
# can coexist in the same GCP project.
# ─────────────────────────────────────────────────────────────────────────────

# New Relic OIDC issuer URI — includes /r/gcp-cmp path; region-specific.
locals {
  # GCP WIF requires the full issuer URI including path.  NR's OIDC discovery
  # document lives at {issuer_uri}/.well-known/openid-configuration.
  newrelic_oidc_issuer = (
    var.newrelic_region == "EU"      ? "https://oidc.eu.newrelic.com/r/gcp-cmp" :
    var.newrelic_region == "Staging" ? "https://oidc-staging.newrelic.com/r/gcp-cmp" :
    "https://oidc.newrelic.com/r/gcp-cmp"
  )
}

resource "google_iam_workload_identity_pool" "newrelic" {
  workload_identity_pool_id = var.wif_pool_id
  display_name              = "New Relic"
  description               = "WIF pool for the New Relic GCP v2 cloud integration"
}

resource "google_iam_workload_identity_pool_provider" "newrelic" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.newrelic.workload_identity_pool_id
  workload_identity_pool_provider_id = var.wif_provider_id
  display_name                       = "New Relic OIDC provider"

  oidc {
    # New Relic issues OIDC tokens from this full path.
    # GCP fetches {issuer_uri}/.well-known/openid-configuration to validate tokens.
    issuer_uri = local.newrelic_oidc_issuer

    # NR's OIDC tokens carry aud = "newrelic-gcp-integrations" — must match exactly.
    allowed_audiences = ["newrelic-gcp-integrations"]
  }

  # Map standard subject + the NR-specific account ID claim.
  attribute_mapping = {
    "google.subject"          = "assertion.sub"
    "attribute.nr_account_id" = "assertion.nr_account_id"
  }

  # Only allow tokens issued for this specific New Relic account.
  attribute_condition = "assertion.nr_account_id == \"${var.newrelic_account_id}\""
}

resource "google_service_account" "newrelic" {
  account_id   = var.newrelic_sa_name
  display_name = "New Relic Integration"
  description  = "Impersonated by New Relic via WIF to collect GCP metrics"
}

# Allow New Relic to collect metrics and enumerate services.
resource "google_project_iam_member" "newrelic_viewer" {
  project = var.gcp_project_id
  role    = "roles/monitoring.viewer"
  member  = "serviceAccount:${google_service_account.newrelic.email}"
}

resource "google_project_iam_member" "newrelic_service_usage" {
  project = var.gcp_project_id
  role    = "roles/serviceusage.serviceUsageConsumer"
  member  = "serviceAccount:${google_service_account.newrelic.email}"
}

# Allow any principal that authenticates through the WIF pool to impersonate
# the New Relic service account.
resource "google_service_account_iam_member" "newrelic_wif" {
  service_account_id = google_service_account.newrelic.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.newrelic.name}/*"
}

# ─────────────────────────────────────────────────────────────────────────────
# Step 1: Link the GCP project to New Relic.
#
# The provider builds the WIF credential JSON internally from the audience and
# service account email — no credential file needed.
#
#   audience              = the WIF pool provider resource name prefixed with
#                           //iam.googleapis.com/ — uniquely identifies the
#                           provider that New Relic must present a token for
#   service_account_email = the SA New Relic impersonates to call GCP APIs
# ─────────────────────────────────────────────────────────────────────────────
resource "newrelic_cloud_gcp_v2_link_account" "this" {
  account_id = var.newrelic_account_id
  name       = var.linked_account_name
  project_id = var.gcp_project_id

  audience              = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"
  service_account_email = google_service_account.newrelic.email

  depends_on = [
    google_project_iam_member.newrelic_viewer,
    google_project_iam_member.newrelic_service_usage,
    google_service_account_iam_member.newrelic_wif,
  ]
}

# ─────────────────────────────────────────────────────────────────────────────
# Step 2: Configure which GCP services New Relic polls for metrics.
# ─────────────────────────────────────────────────────────────────────────────
resource "newrelic_cloud_gcp_v2_integrations" "this" {
  account_id        = newrelic_cloud_gcp_v2_link_account.this.account_id
  linked_account_id = newrelic_cloud_gcp_v2_link_account.this.id

  # ── Existing GCP services ──────────────────────────────────────────────────

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
      fetch_table_metrics      = var.enable_fetch_tags
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

  dynamic "mem_cache" {
    for_each = contains(local.on, "mem_cache") ? [1] : []
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

  # ── GCP v2-only services ───────────────────────────────────────────────────

  dynamic "api_gateway" {
    for_each = contains(local.on, "api_gateway") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "firebase_auth" {
    for_each = contains(local.on, "firebase_auth") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "firebase_vertex_ai" {
    for_each = contains(local.on, "firebase_vertex_ai") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "istio" {
    for_each = contains(local.on, "istio") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "managed_kafka" {
    for_each = contains(local.on, "managed_kafka") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "memory_store" {
    for_each = contains(local.on, "memory_store") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }

  dynamic "firebase_app_hosting" {
    for_each = contains(local.on, "firebase_app_hosting") ? [1] : []
    content { metrics_polling_interval = var.metrics_polling_interval }
  }
}
