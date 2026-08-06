# GCP Dimensional Metrics — multi-service (per-group service selection)

Demonstrates linking **two groups of GCP projects** to New Relic for the **GCP Dimensional Metrics** (keyless / Workload Identity Federation) integration, where each group enables a **different set of GCP services**. All groups share one WIF pool and service account, with IAM granted at the **folder** level.

**Use this module when** you want different service coverage for different classes of projects — e.g. analytics projects monitored for BigQuery/PubSub/Spanner/Storage/Dataflow/Dataproc, and compute projects monitored for VMs/SQL/Cloud Run/Load Balancing/Functions/Kubernetes — from a single shared WIF setup.

## What it creates

In GCP:
- A shared **Workload Identity Federation pool** and **OIDC provider** (in `gcp_sa_project_id`) whose issuer is New Relic's region-specific OIDC endpoint, restricted to your New Relic account via `attribute_condition`.
- A shared **service account** impersonated by New Relic via WIF (no key file).
- **Folder-level IAM bindings** granting the four integration roles across every project in both groups.
- A `roles/iam.workloadIdentityUser` binding scoped to the pool via `attribute.nr_account_id`.
- A 90-second `time_sleep` to let IAM bindings propagate before New Relic authenticates.

In New Relic:
- One `newrelic_cloud_gcp_link_account` (WIF mode) per project in `analytics_projects` and `compute_projects`.
- One `newrelic_cloud_gcp_dm_integrations` per linked account — the **analytics** group enables the analytics service set; the **compute** group enables the compute service set.

## Requirements

- Terraform providers: `newrelic ~> 3.0`, `hashicorp/google ~> 7.0`, `hashicorp/time ~> 0.10` (see `providers.tf`).
- A New Relic User API key with the NerdGraph scope.
- GCP credentials with permission to create the WIF pool/SA in `gcp_sa_project_id` and to bind IAM at `gcp_folder_id`.

## Usage

```hcl
module "newrelic_gcp_dm" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/gcp-dimensional-metrics-multi-service"

  newrelic_account_id = "1234567"
  newrelic_api_key    = var.newrelic_api_key
  newrelic_region     = "US" # or "EU"

  gcp_sa_project_id = "my-shared-sa-project"
  gcp_folder_id     = "123456789012" # numeric, no "folders/" prefix

  analytics_projects = {
    "analytics-prod" = "my-analytics-project-123"
  }
  compute_projects = {
    "compute-prod" = "my-compute-project-456"
  }

  wif_pool_id      = "newrelic-wif-pool"
  wif_provider_id  = "newrelic-oidc-provider"
  newrelic_sa_name = "newrelic-integration"

  metrics_polling_interval = 300
}
```

## Inputs

| Variable | Type | Default | Description |
|---|---|---|---|
| `newrelic_account_id` | string | — | New Relic account ID to link the GCP projects to. |
| `newrelic_api_key` | string (sensitive) | — | New Relic User API key (`NRAK-…`). |
| `newrelic_region` | string | `"US"` | New Relic region: `US`, `EU`, or `JP`. Selects the OIDC issuer URI. |
| `gcp_sa_project_id` | string | — | GCP project where the service account and WIF pool are created. |
| `gcp_folder_id` | string | — | Numeric GCP folder ID (no `folders/` prefix). Folder-level IAM covers all projects in both groups. |
| `analytics_projects` | map(string) | — | `display-name => project ID` for analytics projects (monitored for BigQuery, PubSub, Spanner, Storage, Dataflow, Dataproc). |
| `compute_projects` | map(string) | — | `display-name => project ID` for compute projects (monitored for VMs, SQL, Cloud Run, Load Balancing, Functions, Kubernetes — Kubernetes is metrics only, no entity support; also supports 1-minute polling in Limited Preview). |
| `wif_pool_id` | string | — | ID for the Workload Identity Federation pool. |
| `wif_provider_id` | string | — | ID for the OIDC provider inside the pool. |
| `newrelic_sa_name` | string | — | Name for the impersonated GCP service account. |
| `metrics_polling_interval` | number | `300` | Polling interval in seconds for all enabled services. See LP note below. |

## Outputs

| Output | Description |
|---|---|
| `analytics_linked_account_ids` | Map of display-name => New Relic linked account ID for the analytics group. |
| `compute_linked_account_ids` | Map of display-name => New Relic linked account ID for the compute group. |
| `wif_pool_name` | Full resource name of the WIF pool. |
| `wif_provider_name` | Full resource name of the WIF OIDC provider. |
| `newrelic_service_account_email` | Email of the impersonated GCP service account. |

## Notes

- **1-minute polling (Limited Preview):** `metrics_polling_interval` defaults to `300` s. A 60-second floor is in Limited Preview for `big_query`, `data_flow`, `data_proc`, `kubernetes`, `load_balancing`, `pub_sub`, and `spanner`; set `metrics_polling_interval = 60` to enable it. All other services require `300`.
- **Write-only credential fields:** `audience` and `service_account_email` on `newrelic_cloud_gcp_link_account` are `ForceNew` and never returned by the API; changing them replaces the linked account.
- To change which services each group monitors, edit the `dm_integrations` blocks in `main.tf` for the analytics and compute resources.
