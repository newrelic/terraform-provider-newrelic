---
layout: "newrelic"
page_title: "New Relic: newrelic_federated_logs_setup"
sidebar_current: "docs-newrelic-resource-federated-logs-setup"
description: |-
  Create and manage a Federated Logs setup in New Relic.
---

# Resource: newrelic\_federated\_logs\_setup

Use this resource to create and manage a Federated Logs setup. 

## Example Usage

```hcl
resource "newrelic_aws_connection" "ingest" {
  name       = "fed-logs-ingest"
  role_arn   = "arn:aws:iam::123456789012:role/newrelic-fed-logs-ingest"
  region     = "us-east-1"
  scope_type = "ORGANIZATION"
  scope_id   = "YOUR_ORG_ID_HERE"
}

resource "newrelic_aws_connection" "query" {
  name       = "fed-logs-query"
  role_arn   = "arn:aws:iam::123456789012:role/newrelic-fed-logs-query"
  region     = "us-east-1"
  scope_type = "ORGANIZATION"
  scope_id   = "YOUR_ORG_ID_HERE"
}

resource "newrelic_federated_logs_setup" "foo" {
  name        = "my-app-logs"
  description = "Federated logs setup for my-app"

  storage {
    data_location_bucket      = "my-app-fed-logs"
    database                  = "my_app_fed_logs_db"
    data_ingest_connection_id = newrelic_aws_connection.ingest.id
    query_connection_id       = newrelic_aws_connection.query.id

    cloud_provider_configuration {
      provider = "AWS"
      region   = "us-east-1"
    }
  }

  default_partition {
    storage {
      table             = "my_app_default_partition"
      data_location_uri = "s3://my-app-fed-logs/my_app_default_partition"
    }

    data_retention_policy {
      duration = 30
      unit     = "DAYS"
    }
  }

  forwarder {
    type = "PIPELINE_CONTROL"
    pipeline_control {
      fleet_id = "<fleet-entity-guid>"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the federated log setup.
* `description` - (Optional) A description for federated logs setup.
* `active` - (Optional) Whether the setup is active.
* `storage` - (Required) Storage configuration. Cannot be changed after creation. See [Nested storage block](#nested-storage-block) below.
* `default_partition` - (Required) Default partition created alongside the setup. See [Nested default_partition block](#nested-default_partition-block) below.
* `forwarder` - (Optional) Forwarder configuration that wires a fleet to this setup. See [Nested forwarder block](#nested-forwarder-block) below.

### Nested `storage` block

Each `storage` block supports:

* `data_location_bucket` - (Required) S3 bucket where log data is stored.
* `database` - (Required) Glue catalog database name associated with the setup.
* `data_ingest_connection_id` - (Required) Entity GUID of the `newrelic_aws_connection` used for writing data (the fleet ingest role).
* `query_connection_id` - (Required) Entity GUID of the `newrelic_aws_connection` used for reading data.
* `cloud_provider_configuration` - (Required) Cloud provider configuration. See below.

Each `cloud_provider_configuration` block supports:

* `provider` - (Required) The cloud provider. Currently only `AWS` is supported.
* `region` - (Required) The cloud provider region (e.g. `us-east-1`).

### Nested `default_partition` block

Each `default_partition` block supports:

* `storage` - (Required) Storage details for the default partition. See below.
* `data_retention_policy` - (Optional) Retention policy for logs in the default partition. See below.

Each `default_partition.storage` block supports:

* `table` - (Required) Glue table name for the default partition.
* `data_location_uri` - (Required) S3 URI of the default partition's data location.

Each `default_partition.data_retention_policy` block supports:

* `duration` - (Required) Retention duration value.
* `unit` - (Required) Time unit for the duration. One of `DAYS`, `WEEKS`, or `MONTHS`.

### Nested `forwarder` block

Each `forwarder` block supports:

* `type` - (Required) The forwarder type. Currently only `PIPELINE_CONTROL` is supported.
* `pipeline_control` - (Optional) Pipeline control configuration. Required when `type` is `PIPELINE_CONTROL`. See below.

Each `pipeline_control` block supports:

* `fleet_id` - (Required) The fleet entity GUID.
* `routing_rule` - (Optional) Routing rule that determines how incoming logs are routed to this setup. See below.

Each `routing_rule` block supports:

* `expression` - (Required) [OTTL](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl) expression for routing logs to this setup. Example: `attributes["service.name"] == "python-apm"`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the federated logs setup. Used as `setup_id` on `newrelic_federated_logs_partition`.
* `default_partition_id` - The entity GUID of the default partition created alongside this setup.
* `lifecycle_status` - Current lifecycle status of the setup.
* `health_check` - Aggregate health check status for the setup.
* `created_at` - Creation timestamp.
* `updated_at` - Last-updated timestamp.

## Import

Federated Logs setups can be imported using the entity GUID:

```bash
$ terraform import newrelic_federated_logs_setup.foo <entity-guid>
```
