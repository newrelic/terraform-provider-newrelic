output "linked_account_id" {
  description = "The New Relic linked account ID for this GCP v2 project."
  value       = newrelic_cloud_gcp_v2_link_account.this.id
}

output "linked_account_name" {
  description = "The display name of the linked GCP account in New Relic."
  value       = newrelic_cloud_gcp_v2_link_account.this.name
}

output "newrelic_account_id" {
  description = "The New Relic account ID the GCP project is linked to."
  value       = newrelic_cloud_gcp_v2_link_account.this.account_id
}

output "gcp_project_id" {
  description = "The GCP project ID that is linked."
  value       = newrelic_cloud_gcp_v2_link_account.this.project_id
}

output "integrations_id" {
  description = "The resource ID of the integrations resource (equals linked_account_id)."
  value       = newrelic_cloud_gcp_v2_integrations.this.id
}

output "wif_pool_name" {
  description = "The full resource name of the GCP Workload Identity Pool."
  value       = google_iam_workload_identity_pool.newrelic.name
}

output "wif_provider_name" {
  description = "The full resource name of the GCP WIF pool provider (audience base)."
  value       = google_iam_workload_identity_pool_provider.newrelic.name
}

output "newrelic_service_account_email" {
  description = "The GCP service account email New Relic uses to collect metrics."
  value       = google_service_account.newrelic.email
}
