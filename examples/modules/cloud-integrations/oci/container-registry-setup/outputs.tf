output "container_repository_id" {
  description = "OCID of the created Container Registry repository. Pass this to logs-integration's container_repository_id variable."
  value       = oci_artifacts_container_repository.log_forwarder_repo.id
}

output "container_repository_namespace" {
  description = "OCIR namespace the repository was created under."
  value       = oci_artifacts_container_repository.log_forwarder_repo.namespace
}
