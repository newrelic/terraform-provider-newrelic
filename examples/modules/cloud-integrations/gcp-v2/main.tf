# Link GCP project to New Relic using Workload Identity Federation (WIF).
# The wif_credential is a write-only field — it is used during Create and
# is never returned by the API. Changing it forces replacement of this resource.
resource "newrelic_cloud_gcp_dm_link_account" "main" {
  account_id     = var.newrelic_account_id
  name           = var.name
  project_id     = var.gcp_project_id
  wif_credential = var.wif_credential
}

# Configure which GCP services New Relic polls for metrics.
resource "newrelic_cloud_gcp_dm_integrations" "main" {
  account_id        = newrelic_cloud_gcp_dm_link_account.main.account_id
  linked_account_id = newrelic_cloud_gcp_dm_link_account.main.id

  # ── Existing GCP services ─────────────────────────────────────────────────

  ai_platform {
    metrics_polling_interval = var.metrics_polling_interval
  }

  alloy_db {
    metrics_polling_interval = var.metrics_polling_interval
  }

  app_engine {
    metrics_polling_interval = var.metrics_polling_interval
  }

  big_query {
    metrics_polling_interval = var.metrics_polling_interval
    fetch_tags               = true
    fetch_table_metrics      = true
  }

  big_table {
    metrics_polling_interval = var.metrics_polling_interval
  }

  composer {
    metrics_polling_interval = var.metrics_polling_interval
  }

  data_flow {
    metrics_polling_interval = var.metrics_polling_interval
  }

  data_proc {
    metrics_polling_interval = var.metrics_polling_interval
  }

  data_store {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_database {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_hosting {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_storage {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firestore {
    metrics_polling_interval = var.metrics_polling_interval
  }

  functions {
    metrics_polling_interval = var.metrics_polling_interval
  }

  interconnect {
    metrics_polling_interval = var.metrics_polling_interval
  }

  kubernetes {
    metrics_polling_interval = var.metrics_polling_interval
  }

  load_balancing {
    metrics_polling_interval = var.metrics_polling_interval
  }

  mem_cache {
    metrics_polling_interval = var.metrics_polling_interval
  }

  pub_sub {
    metrics_polling_interval = var.metrics_polling_interval
    fetch_tags               = true
  }

  redis {
    metrics_polling_interval = var.metrics_polling_interval
  }

  router {
    metrics_polling_interval = var.metrics_polling_interval
  }

  run {
    metrics_polling_interval = var.metrics_polling_interval
  }

  spanner {
    metrics_polling_interval = var.metrics_polling_interval
    fetch_tags               = true
  }

  sql {
    metrics_polling_interval = var.metrics_polling_interval
  }

  storage {
    metrics_polling_interval = var.metrics_polling_interval
    fetch_tags               = true
  }

  virtual_machines {
    metrics_polling_interval = var.metrics_polling_interval
  }

  vpc_access {
    metrics_polling_interval = var.metrics_polling_interval
  }

  # ── GCP v2-only services (Workload Identity Federation required) ───────────

  api_gateway {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_auth {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_vertex_ai {
    metrics_polling_interval = var.metrics_polling_interval
  }

  istio {
    metrics_polling_interval = var.metrics_polling_interval
  }

  managed_kafka {
    metrics_polling_interval = var.metrics_polling_interval
  }

  memory_store {
    metrics_polling_interval = var.metrics_polling_interval
  }

  firebase_app_hosting {
    metrics_polling_interval = var.metrics_polling_interval
  }
}
