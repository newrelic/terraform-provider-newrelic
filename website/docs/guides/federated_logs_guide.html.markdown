---
layout: "newrelic"
page_title: "Getting Started with New Relic Federated Logs"
sidebar_current: "docs-newrelic-provider-federated-logs-guide"
description: |-
  Use this guide to set up New Relic Federated Logs end-to-end through Terraform.
---

## Getting Started with New Relic Federated Logs

[New Relic Federated Logs](https://docs-preview.newrelic.com/docs/federated-logs/) lets your log data stay in your own AWS environment while you query it directly from the New Relic platform alongside the rest of your telemetry. The data plane (S3, SQS, Glue, AWS Managed Flink, and the Pipeline Control Gateway) runs in your AWS account; the control plane (the setup, partitions, and routing configuration) lives in New Relic.

This guide describes how to provision the full Federated Logs setup both the AWS data plane and the New Relic control-plane entities fully automated through Terraform.

Use the links below to go to the relevant sections:

- [Architecture overview](#architecture-overview)
- [Resources covered](#resources-covered)
- [Prerequisites](#prerequisites)
- [End-to-end onboarding with the `terraform-aws-federated-logs` modules](#onboarding)
  - [Stage 1: Data Processing (per fleet)](#stage-1-data-processing)
  - [Stage 2: Federated Logs Setup (per setup)](#stage-2-federated-logs-setup)
- [Activate the data path](#activate)
- [Adding partitions to an existing setup](#adding-partitions)
- [Querying federated logs](#query-federated-logs)
- [Cleaning up a setup while preserving stored logs](#cleanup)

**NOTE:** Federated Logs is currently provided as a limited preview. Your organization must be enrolled in the preview before these resources will function. See [Enroll in the limited preview program](https://docs-preview.newrelic.com/docs/federated-logs/#enroll) in the product docs.

### Architecture overview

Federated Logs has two distinct provisioning stages, and the New Relic Terraform provider exposes one control-plane resource for each artifact created along the way:

1. **Data Processing** — applied **once per Pipeline Control Gateway (PCG) fleet**. This stage creates the fleet-level cross-account IAM base role, the ABAC inline policy that lets it assume per-setup writer roles, the SQS queue and AWS Managed Flink job that update Iceberg metadata, and the [`newrelic_aws_connection`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/aws_connection) entity that stores the base role ARN against the fleet.

2. **Federated Logs Setup** — applied **once per log setup**. This stage creates the S3 bucket, the Glue catalog database, the `pcg-writer` IAM role (trusted by the fleet base role via ABAC tag matching), the New Relic reader IAM role for cross-account query access, the Iceberg tables, and the [`newrelic_federated_logs_setup`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/federated_logs_setup) entity plus any [`newrelic_federated_logs_partition`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/federated_logs_partition) entities you declare.

The `fleet_entity_guid` produced by your PCG installation is the key that links the two stages together. A single fleet can host multiple setups; you only need to run Stage 1 once per fleet.

### Resources covered

| Resource | What it does |
|---|---|
| [`newrelic_aws_connection`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/aws_connection) | Wraps an AWS IAM role ARN so New Relic can authenticate into your AWS account. Federated Logs uses two of these — one for the fleet base role (data ingest path) and one for the per-setup reader role (query path). |
| [`newrelic_federated_logs_setup`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/federated_logs_setup) | Creates the setup entity in New Relic. Ties together the storage bucket, Glue database, ingest and query connections, the default partition, and the forwarder (Pipeline Control Gateway fleet) routing rule. |
| [`newrelic_federated_logs_partition`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/federated_logs_partition) | Creates an additional partition under an existing setup. Each partition maps to its own Iceberg table and has its own retention policy and OTTL routing expression. |

### Prerequisites

Before you onboard Federated Logs through Terraform, confirm the following:

* Your New Relic organization is **enrolled in the Federated Logs limited preview**.
* A **Pipeline Control Gateway fleet** is already deployed in your Kubernetes environment, with its DNS endpoint reachable from your log sources, and you know its **`fleet_entity_guid`**. If you don't yet have a fleet, follow [Getting Started with New Relic Fleet Control](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/fleet_getting_started) and the [Pipeline Control Gateway documentation](https://docs.newrelic.com/docs/new-relic-control/pipeline-control/overview/).
* Your **log sources are pointed at the gateway endpoint** rather than sending logs directly to New Relic. The per-source configuration (APM agents, infrastructure agent, Fluent Bit, Fluentd) is covered in [Configure your log source](https://docs-preview.newrelic.com/docs/federated-logs/#configure-source) in the product docs.
* Federated Logs currently supports **AWS only**, and the Query engine is deployed in **`us-west-2`**. Provision your AWS resources in that region for the preview.
* You have an **AWS account** with permissions to create S3 buckets, SQS queues, Glue catalogs, AWS Managed Flink applications, and IAM roles.
* You have **Terraform >= 1.6.0**, the **New Relic provider >= 3.91.0**, and the **AWS provider >= 6.36.0** configured.
* Export your New Relic credentials as environment variables before running Terraform. They are read directly from the environment:

  ```sh
  export NEW_RELIC_API_KEY="NRAK-..."          # required by both stages (NerdGraph)
  export NEW_RELIC_LICENSE_KEY="..."           # required by Stage 1 only (Flink → New Relic metrics)
  ```

**NOTE:** These guides assume you've already configured the New Relic and AWS providers with the correct credentials. If you haven't done so, see the [New Relic provider getting started guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started) and the [AWS provider authentication docs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration).

<a id="onboarding"></a>
### End-to-end onboarding with the `terraform-aws-federated-logs` modules

The [`terraform-aws-federated-logs`](https://github.com/newrelic/terraform-aws-federated-logs) repository ships two modules, one for each stage, that provision the AWS infrastructure and create the matching New Relic control-plane entities (`newrelic_aws_connection`, `newrelic_federated_logs_setup`, and the default `newrelic_federated_logs_partition`). These modules capture all of the cross-resource wiring (ABAC tags, IAM trust policies, SQS notifications, Glue table parameters, optimizer settings) in one place.

**IMPORTANT:** Run Stage 1 (`data_processing`) before Stage 2 (`federated_logs_setup`). The `fleet_entity_guid` produced by your PCG installation is the input that links the two.

<a id="stage-1-data-processing"></a>
#### Stage 1: Data Processing (per fleet)

The `data_processing` module is deployed **once per PCG fleet**. It creates:

* A fleet-level IAM base role authenticated via either **IRSA** or **EKS Pod Identity**.
* An ABAC inline policy that lets the base role assume any per-setup `pcg-writer` role tagged with the matching `fleet_entity_guid` value.
* SQS queue, AWS Managed Flink job (the `flink-iceberg-commit-worker`), and supporting CloudWatch resources.
* A `newrelic_aws_connection` entity in New Relic storing the base role ARN, plus a `HAS_FED_LOGS_BASE_ROLE` relationship from the fleet entity to that connection.

```hcl
module "data_processing" {
  source = "git::https://github.com/newrelic/terraform-aws-federated-logs.git//modules/data_processing"

  data_processing_module_name = "my-app-logs"
  newrelic_org_id             = "YOUR_NR_ORG_ID"
  fleet_entity_guid           = "YOUR_FLEET_ENTITY_GUID"
  newrelic_region             = "US"           # "US" (default)

  clusters = {
    "prod-cluster" = {
      k8s_namespace            = "federated-logs"
      auth_mode                = "irsa"        # "irsa" or "pod_identity"
      k8s_service_account_name = "pcg-writer-sa"
      oidc_provider_arn        = "arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLE"
      # cluster_name           = "my-cluster" # required if auth_mode = "pod_identity"
    }
  }
}
```

Key variables:

* `data_processing_module_name` — Name used in resource naming (3–26 lowercase alphanumeric characters, hyphens allowed but not at the start or end).
* `newrelic_org_id` — Your New Relic organization ID.
* `fleet_entity_guid` — The NGEP entity GUID of the PCG fleet that will forward logs to this storage stack.
* `newrelic_region` (Optional) — New Relic region context. One of `US` (default).
* `clusters` — Map of EKS cluster configs used to build the base role trust policy. Every entry must share the same `auth_mode` (`irsa` or `pod_identity`); mixing modes is rejected.
  * `auth_mode = "irsa"` requires `oidc_provider_arn`.
  * `auth_mode = "pod_identity"` requires `cluster_name`.
* `parallelism`, `parallelism_per_kpu`, `auto_scaling_enabled` (Optional) — Flink runtime tuning. Defaults are sized for typical ingestion volume.

The module exports `base_role_arn` and `base_role_name`. The newly created `newrelic_aws_connection` entity is the fleet ingest connection that Stage 2 references implicitly via `fleet_entity_guid`.

<a id="stage-2-federated-logs-setup"></a>
#### Stage 2: Federated Logs Setup (per setup)

The root `terraform-aws-federated-logs` module is deployed **once per setup**. It creates:

* An S3 bucket and a Glue catalog database for log storage.
* A `pcg-writer` IAM role trusted by the fleet base role via ABAC tag matching.
* A New Relic reader IAM role for cross-account query access.
* Iceberg tables (one default plus any additional `partition_tables` you declare) with configurable optimizer and retention settings.
* A `newrelic_federated_logs_setup` entity and a default `newrelic_federated_logs_partition`, plus a `newrelic_aws_connection` entity wrapping the reader role (used as `query_connection_id` on the setup).

```hcl
module "federated_logs" {
  source = "git::https://github.com/newrelic/terraform-aws-federated-logs.git"

  setup_name          = "my-app-logs"
  fleet_entity_guid   = "YOUR_FLEET_ENTITY_GUID"
  newrelic_org_id     = "YOUR_NR_ORG_ID"
  newrelic_account_id = 1234567
  newrelic_region     = "US"

  data_retention_enabled = true

  default_table_setting = {
    retention_in_days = 30
    table_parameters = {
      "write.parquet.compression-codec"            = "zstd"
      "write.target-file-size-bytes"               = "67108864" # 64 MB
      "write.metadata.delete-after-commit.enabled" = "true"
      "write.metadata.previous-versions-max"       = "10"
    }
    optimizer_configuration = {
      orphan_file_deletion = {
        orphan_file_retention_period_in_days = 3
        run_rate_in_hours                    = 24
      }
      snapshot_retention = {
        snapshot_retention_period_in_days = 1
        number_of_snapshots_to_retain     = 1
        clean_expired_files               = false
        run_rate_in_hours                 = 3
      }
      compaction = {
        strategy              = "binpack"
        min_input_files       = 5
        delete_file_threshold = 1
      }
    }
  }

  # Additional partitions. Each entry can override retention, table_parameters,
  # optimizer_configuration, routing_expression, and/or description.
  partition_tables = {
    "application_log" = {
      retention_in_days = 5
    }
    "security_log" = {
      retention_in_days = 10
      optimizer_configuration = {
        compaction = {
          strategy              = "binpack"
          min_input_files       = 10
          delete_file_threshold = 2
        }
      }
    }
  }
}
```

Key variables:

* `setup_name` — Name for this setup. Used in AWS resource naming; 3–26 lowercase alphanumeric characters (hyphens allowed but not at the start or end).
* `fleet_entity_guid` — Same GUID you passed to Stage 1. The module uses it to scope the IAM trust on the `pcg-writer` role.
* `newrelic_account_id` / `newrelic_org_id` — The New Relic account and org receiving the federated logs.
* `newrelic_region` (Optional) — `US` (default).
* `data_retention_enabled` (Optional) — When `true`, the module deploys a Glue job that deletes data older than the per-table `retention_in_days`. Default is `true` (`false` in the underlying module — set explicitly when you need pruning).
* `default_table_setting` (Optional) — Iceberg parameters and optimizer settings for the default partition's table.
* `partition_tables` (Optional) — Map of additional partition tables; each entry may override `retention_in_days`, `table_parameters`, `optimizer_configuration`, `routing_expression`, and `description`. Each entry becomes a `newrelic_federated_logs_partition` under the setup.
* `e2e_validation_config` (Optional) — When enabled, the module deploys an AWS Lambda that POSTs a synthetic log through your PCG endpoint, polls NRDB for the log, and reports `HEALTHY` / `UNHEALTHY` back to New Relic. Useful in CI/CD to gate `terraform apply` on a working data path.

Useful outputs:

* `newrelic_federated_logs_setup_id` — The setup entity GUID. Pass this as `setup_id` when you later add partitions.
* `newrelic_default_partition_id` — The entity GUID of the default partition created alongside the setup.
* `newrelic_query_connection_id` — The entity GUID of the per-setup `newrelic_aws_connection` wrapping the reader role.
* `s3_bucket_name`, `glue_database_name`, `pcg_writer_role_arn`, `nr_reader_role_arn`, `iceberg_tables` — Underlying AWS resources.

<a id="activate"></a>
### Activate the data path

`terraform apply` creates the setup and ties it to your fleet, but logs only start flowing after two steps:

1. **Set routing conditions and deploy the gateway config** — [docs](https://docs-preview.newrelic.com/docs/federated-logs/#set-up-routing-conditions).
2. **Verify with "Simulate test log received"** — [docs](https://docs-preview.newrelic.com/docs/federated-logs/#test-the-setup).

<a id="adding-partitions"></a>
### Adding partitions to an existing setup

Add new partitions by extending the `partition_tables` map on your existing `federated_logs` module — there is no need to declare standalone `newrelic_federated_logs_partition` resources. Each new entry creates the Glue table, the retention Glue job, and the matching `newrelic_federated_logs_partition` entity together as a single unit.

```hcl
module "federated_logs" {

  # ... your existing config ...
  setup_name = "my-app-logs"

  # ───────────────────────────────────────────────────────────────────────
  #      To add a new partition, simply add a new entry below in the
  #      partition_tables map in the module corresponding to the setup
  # ───────────────────────────────────────────────────────────────────────
  partition_tables = {
    # ... your existing partitions ...

    "Log_partition_name" = {
      retention_in_days = 30
      description       = "<DESCRIPTION>"
    }
  }
}
```

Each entry in `partition_tables` may optionally override `retention_in_days`, `description`, `routing_expression`, `table_parameters`, and `optimizer_configuration`. Apply the change with the usual `terraform plan` / `terraform apply`; the new partition becomes visible in the New Relic UI under **Logs > Data Partitions > Federated Logs**.

As with the initial setup, a new partition only starts receiving logs once you define its routing condition and roll the gateway configuration out from the wizard. The product docs cover the same steps for partitions in [Set up routing conditions](https://docs-preview.newrelic.com/docs/federated-logs/#set-up-routing-conditions-1) and [Update gateway configuration](https://docs-preview.newrelic.com/docs/federated-logs/#update-gateway-configuration-1).

<a id="query-federated-logs"></a>
### Querying federated logs

After `terraform apply` completes and the gateway deployment has rolled out, your logs are queryable from the New Relic Logs UI (select **Federated logs** in the **Partitions** dropdown) and from NRQL. Always include `log.partition` in NRQL queries so they route to your S3 bucket rather than New Relic's standard log storage:

```sql
SELECT * FROM Log WHERE log.partition = 'log_federated' SINCE 1 hour ago
```

See [Query federated logs](https://docs-preview.newrelic.com/docs/federated-logs/#query) in the product docs for more examples.

<a id="cleanup"></a>
### Cleaning up a setup while preserving stored logs

The storage resources in the `federated_logs` module — the S3 bucket, the Glue catalog database, the partition folder objects, and the Iceberg tables are declared with `lifecycle { prevent_destroy = true }`, so `terraform destroy` will refuse to remove them and abort the plan. This is intentional: it keeps your historical log data from being deleted by accident.

To tear down a setup while leaving the stored logs intact, use Terraform's [`removed`](https://developer.hashicorp.com/terraform/language/resources/syntax#the-removed-block) block to drop the protected storage resources from state without destroying them, so Terraform can deprovision everything else around them:

1. **Comment out the existing `module "federated_logs"` block.** While the module declaration is still in your configuration, Terraform considers all of its resources to be under active management, and `removed` blocks pointing at the same addresses are rejected as conflicting. Commenting (or deleting) the module block first makes those addresses removable.
2. **Add `removed` blocks for the protected storage resources** so Terraform drops them from state without destroying them, while it deprovisions everything else around them:
```hcl
removed {
  from = module.federated_logs.module.setup
  lifecycle { destroy = false }
}

removed {
  from = module.federated_logs.module.partition.aws_s3_object.folder
  lifecycle { destroy = false }
}

removed {
  from = module.federated_logs.module.partition.aws_glue_catalog_table.iceberg_table
  lifecycle { destroy = false }
}
```

Then run `terraform plan` and `terraform apply`. The `removed` blocks drop the protected resources from Terraform state without deleting them in AWS, allowing the rest of the module including the New Relic entities to be deprovisioned cleanly. Afterwards, the New Relic setup and its partitions remain visible in the UI but are no longer editable or queryable; the S3 bucket and Glue tables stay intact so your historical log data is preserved.