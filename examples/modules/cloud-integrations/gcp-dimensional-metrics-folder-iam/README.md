# GCP Dimensional Metrics — folder-level IAM

Links one or more GCP projects to New Relic for the **GCP Dimensional Metrics** (keyless / Workload Identity Federation) integration, granting the required IAM roles **once at the GCP folder level** so a single service account covers every project under that folder.

**Use this module when** you manage multiple GCP projects under a common folder and have permission to bind IAM at the folder level. If you can only bind IAM per project, use [`gcp-dimensional-metrics-project-iam`](../gcp-dimensional-metrics-project-iam) instead.

## What it creates

In GCP:
- A **Workload Identity Federation pool** and **OIDC provider** (in `gcp_sa_project_id`) whose issuer is New Relic's region-specific OIDC endpoint, restricted to your New Relic account via `attribute_condition`.
- A **service account** that New Relic impersonates via WIF (no key file).
- **Folder-level IAM bindings** granting the four roles the integration needs:
  `roles/viewer`, `roles/serviceusage.serviceUsageConsumer`, `roles/cloudasset.viewer`, and `roles/resourcemanager.folderViewer`.
- A `roles/iam.workloadIdentityUser` binding scoped to the pool via `attribute.nr_account_id`.
- A 90-second `time_sleep` to let IAM bindings propagate before New Relic authenticates.

In New Relic:
- One `newrelic_cloud_gcp_link_account` (WIF mode) per entry in `gcp_projects`.
- One `newrelic_cloud_gcp_dm_integrations` per linked account, enabling the services in `enabled_services`.

## Requirements

- Terraform providers: `newrelic ~> 3.0`, `hashicorp/google ~> 7.0`, `hashicorp/time ~> 0.10` (see `providers.tf`).
- A New Relic User API key with the NerdGraph scope.
- GCP credentials with permission to create the WIF pool/SA in `gcp_sa_project_id` and to bind IAM at `gcp_folder_id`.

## Usage

```hcl
module "newrelic_gcp_dm" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/gcp-dimensional-metrics-folder-iam"

  newrelic_account_id = 1234567
  newrelic_api_key    = var.newrelic_api_key
  newrelic_region     = "US" # or "EU"

  gcp_sa_project_id = "my-shared-sa-project"
  gcp_folder_id     = "123456789012" # numeric, no "folders/" prefix

  gcp_projects = {
    "prod-payments"  = "my-payments-project-123"
    "prod-analytics" = "my-analytics-project-456"
  }

  wif_pool_id      = "newrelic-wif-pool"
  wif_provider_id  = "newrelic-oidc-provider"
  newrelic_sa_name = "newrelic-integration"

  enabled_services         = ["big_query", "pub_sub", "storage"]
  metrics_polling_interval = 300
}
```

## Inputs

| Variable | Type | Default | Description |
|---|---|---|---|
| `newrelic_account_id` | number | — | New Relic account ID to link the GCP projects to. |
| `newrelic_api_key` | string (sensitive) | — | New Relic User API key (`NRAK-…`). |
| `newrelic_region` | string | `"US"` | New Relic region: `US`, `EU`, or `JP`. Selects the OIDC issuer URI. |
| `gcp_sa_project_id` | string | — | GCP project where the service account and WIF pool are created. |
| `gcp_folder_id` | string | — | Numeric GCP folder ID (no `folders/` prefix). All four roles are granted here. |
| `gcp_projects` | map(string) | — | `display-name => project ID` for each project to link. The display-name becomes the linked-account name in New Relic. |
| `wif_pool_id` | string | — | ID for the Workload Identity Federation pool. |
| `wif_provider_id` | string | — | ID for the OIDC provider inside the pool. |
| `newrelic_sa_name` | string | — | Name for the impersonated GCP service account. |
| `enabled_services` | list(string) | `["big_query","pub_sub","storage"]` | GCP services to enable (see `variables.tf` for the full list). |
| `metrics_polling_interval` | number | `300` | Polling interval in seconds for all enabled services. See LP note below. |

## Outputs

| Output | Description |
|---|---|
| `linked_account_ids` | Map of display-name => New Relic linked account ID. |
| `wif_pool_name` | Full resource name of the WIF pool. |
| `wif_provider_name` | Full resource name of the WIF OIDC provider. |
| `newrelic_service_account_email` | Email of the impersonated GCP service account. |

## Notes

- **1-minute polling (Limited Preview):** `metrics_polling_interval` defaults to `300` s. A 60-second floor is in Limited Preview for `alloy_db`, `big_query`, `data_flow`, `data_proc`, `load_balancing`, `managed_kafka`, `pub_sub`, and `spanner`; set `metrics_polling_interval = 60` to enable it. All other services require `300`.
- **Write-only credential fields:** `audience` and `service_account_email` on `newrelic_cloud_gcp_link_account` are `ForceNew` and never returned by the API; changing them replaces the linked account.
