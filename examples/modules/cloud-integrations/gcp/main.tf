resource "google_project_iam_member" "project" {
  project = var.gcp_project_id
  role    = "roles/viewer"
  member  = "serviceAccount:${var.gcp_service_account_id}"
}

resource "google_project_iam_binding" "project" {
  project = var.gcp_project_id
  role    = "roles/serviceusage.serviceUsageConsumer"

  members = [
    "serviceAccount:${var.gcp_service_account_id}",
  ]
}

resource "newrelic_cloud_gcp_link_account" "gcp_account" {
  account_id = var.newrelic_account_id
  project_id = var.gcp_project_id
  name       = var.name
}

resource "newrelic_cloud_gcp_integrations" "gcp_integrations" {
  account_id        = var.newrelic_account_id
  linked_account_id = newrelic_cloud_gcp_link_account.gcp_account.id
  app_engine {}
  big_query {}
  big_table {}
  composer {}
  data_flow {}
  data_proc {}
  data_store {}
  fire_base_database {}
  fire_base_hosting {}
  fire_base_storage {}
  fire_store {}
  functions {}
  interconnect {}
  kubernetes {}
  load_balancing {}
  mem_cache {}
  pub_sub {}
  redis {}
  router {}
  run {}
  spanner {}
  sql {}
  storage {}
  virtual_machines {}
  vpc_access {}
}
