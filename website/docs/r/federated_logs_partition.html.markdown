---
layout: "newrelic"
page_title: "New Relic: newrelic_federated_logs_partition"
sidebar_current: "docs-newrelic-resource-federated-logs-partition"
description: |-
  Create and manage an additional partition under an existing Federated Logs setup.
---

# Resource: newrelic\_federated\_logs\_partition

Use this resource to create and manage an additional partition under an existing `newrelic_federated_logs_setup`. 

## Example Usage

```hcl
resource "newrelic_federated_logs_partition" "foo" {
  setup_id    = newrelic_federated_logs_setup.setup.id
  name        = "test-partition-logs"
  description = "test partition for logs"

  storage {
    table             = "my_app_partition_logs"
    data_location_uri = "s3://my-app-fed-logs/my_app_partition_logs"
  }

  data_retention_policy {
    duration = 365
    unit     = "DAYS"
  }

  forwarder_configuration {
    type = "PIPELINE_CONTROL"
    pipeline_control {
      partition_rule {
        expression = "attributes[\"log.type\"] == \"partition\""
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `setup_id` - (Required) Entity GUID of the setup.
* `name` - (Required) The name of the partition.
* `description` - (Optional) A description for partition.
* `active` - (Optional) Whether the partition is active.
* `storage` - (Required) Storage details for this partition. See [Nested storage block](#nested-storage-block) below.
* `data_retention_policy` - (Optional) Retention policy for logs in this partition. See [Nested data_retention_policy block](#nested-data_retention_policy-block) below.
* `forwarder_configuration` - (Optional) Forwarder configuration for routing specific logs to this partition. See [Nested forwarder_configuration block](#nested-forwarder_configuration-block) below.

### Nested `storage` block

Each `storage` block supports:

* `table` - (Required) Glue table name for the partition.
* `data_location_uri` - (Required) S3 URI of the partition's data location.

### Nested `data_retention_policy` block

Each `data_retention_policy` block supports:

* `duration` - (Required) Retention duration value.
* `unit` - (Required) Time unit. One of `DAYS`, `WEEKS`, or `MONTHS`.

### Nested `forwarder_configuration` block

Each `forwarder_configuration` block supports:

* `type` - (Required) Forwarder type. Must match the parent setup's forwarder type. Currently only `PIPELINE_CONTROL` is supported.
* `pipeline_control` - (Optional) Pipeline control configuration. See below.

Each `pipeline_control` block supports:

* `partition_rule` - (Optional) Rule that determines which logs are routed to this partition. See below.

Each `partition_rule` block supports:

* `expression` - (Required) [OTTL](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl) expression for routing logs to this partition. Example: `attributes["log.type"] == "partition"`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the partition.
* `is_default` - Whether this is the default partition for the parent setup.
* `lifecycle_status` - Current lifecycle status of the partition.
* `health_check` - Aggregate health check status for the partition.
* `created_at` - Creation timestamp.
* `updated_at` - Last-updated timestamp.

## Import

Federated Logs partitions can be imported using the entity GUID:

```bash
$ terraform import newrelic_federated_logs_partition.foo <entity-guid>
```
