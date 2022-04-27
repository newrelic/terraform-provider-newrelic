/*

    Complete example to enable AWS integration with New Relic

*/


variable "NEW_RELIC_SERVICE_ACCOUNT_ID" {
  type = string
}

variable "GCP_PROJECT_ID" {
  type = string
}

variable "NR_ACCOUNT_ID" {
  type = string
}

resource "google_project_iam_member" "project" {
  project = var.GCP_PROJECT_ID
  role    = "roles/editor"
  member  = "serviceAccount:${var.NEW_RELIC_SERVICE_ACCOUNT_ID}"
}

resource "google_project_iam_binding" "project" {
  project = var.GCP_PROJECT_ID
  role    = "roles/editor"

  members = [
    "serviceAccount:${var.NEW_RELIC_SERVICE_ACCOUNT_ID}",
  ]
}

resource "newrelic_cloud_gcp_link_account" "gcp_account" {
  account_id=var.NR_ACCOUNT_ID
  project_id =var.GCP_PROJECT_ID
  name       = "GCP linked account name"
}

resource "newrelic_cloud_gcp_integrations" "gcp_integrations" {
  account_id=var.NR_ACCOUNT_ID
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
