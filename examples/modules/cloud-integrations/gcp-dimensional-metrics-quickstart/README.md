# GCP Dimensional Metrics — single-project quickstart

A complete, ready-to-run example that creates the full **Workload Identity Federation (WIF)** setup in GCP and links **one project** to New Relic for the **GCP Dimensional Metrics** (keyless) integration — no service-account key files. This is the flat, single-file starting point referenced from the [cloud integrations guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp-dimensional-metrics); copy it, replace the placeholder values, and run `terraform init && terraform apply`.

For reusable, multi-project setups see the sibling modules: [`gcp-dimensional-metrics-folder-iam`](../gcp-dimensional-metrics-folder-iam), [`gcp-dimensional-metrics-project-iam`](../gcp-dimensional-metrics-project-iam), and [`gcp-dimensional-metrics-multi-service`](../gcp-dimensional-metrics-multi-service).

## What it creates

In GCP:
- A **Workload Identity Federation pool** and **OIDC provider**, with the issuer set to New Relic's region-specific endpoint and `attribute_condition` restricting tokens to your New Relic account.
- A **service account** impersonated by New Relic via WIF (no key file).
- IAM bindings granting the four roles the integration needs: `roles/viewer`, `roles/serviceusage.serviceUsageConsumer`, and `roles/cloudasset.viewer` at the **project** level, and `roles/resourcemanager.folderViewer` at the **folder** level.
- A `roles/iam.workloadIdentityUser` binding scoped to the pool via `attribute.nr_account_id`.

In New Relic:
- A `newrelic_cloud_gcp_link_account` in keyless (WIF) mode.
- A `newrelic_cloud_gcp_dm_integrations` enabling all supported GCP services.

## Requirements

- Terraform providers: `newrelic ~> 3.0`, `hashicorp/google ~> 7.0` (see `providers.tf`).
- A New Relic User API key with the NerdGraph scope.
- GCP credentials with permission to create the WIF pool/SA and bind IAM at both the project and folder level.

## Usage

```hcl
# terraform.tfvars
newrelic_account_id = 1234567
newrelic_api_key    = "NRAK-…"
newrelic_region     = "US" # US, EU, or JP
gcp_project_id      = "my-project-123"
gcp_folder_id       = "123456789012" # numeric, no "folders/" prefix
```

```bash
terraform init
terraform apply
```

## Inputs

| Variable | Type | Default | Description |
|---|---|---|---|
| `newrelic_account_id` | number | — | New Relic account ID to link the GCP project to. |
| `newrelic_api_key` | string (sensitive) | — | New Relic User API key with the NerdGraph scope. |
| `newrelic_region` | string | `"US"` | New Relic region: `US`, `EU`, or `JP`. Selects the OIDC issuer URI. |
| `gcp_project_id` | string | — | GCP project ID to monitor. |
| `gcp_folder_id` | string | — | Numeric GCP folder ID (no `folders/` prefix) for the folder-level role. |
| `linked_account_name` | string | `"production-gcp-dm"` | Display name shown in the New Relic UI. |
| `wif_pool_id` | string | `"newrelic-pool"` | ID for the Workload Identity Pool. |
| `wif_provider_id` | string | `"newrelic-provider"` | ID for the OIDC provider inside the pool. |
| `newrelic_sa_name` | string | `"newrelic-integration"` | Name for the impersonated GCP service account. |
| `metrics_polling_interval` | number | `300` | Polling interval in seconds for all enabled services. See LP note below. |

## Outputs

| Output | Description |
|---|---|
| `linked_account_id` | New Relic linked account ID for the GCP DM integration. |

## Notes

- **1-minute polling (Limited Preview):** `metrics_polling_interval` defaults to `300` s. A 60-second floor is in Limited Preview for `alloy_db`, `big_query`, `data_flow`, `data_proc`, `load_balancing`, `managed_kafka`, `pub_sub`, and `spanner`; set `metrics_polling_interval = 60` to enable it. All other services require `300`.
- **Write-only credential fields:** `audience` and `service_account_email` on `newrelic_cloud_gcp_link_account` are write-only, `ForceNew` fields — they're used to construct the WIF credential internally and are never returned by the API. To import an existing linked account, run `terraform import newrelic_cloud_gcp_link_account.main <linked_account_id>` and then `terraform apply` to reconcile those fields (the resource is replaced).
