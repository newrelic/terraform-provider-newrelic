# Create the vault replica in the dynamic destination region
resource "oci_kms_vault_replication" "vault_replica" {
  depends_on     = [oci_identity_policy.nr_metrics_policy]

  vault_id       = oci_kms_vault.newrelic_vault.id
  replica_region = "us-phoenix-1" # Use the variable for the replica region
}
