---
layout: "newrelic"
page_title: "New Relic: newrelic_federated_log_partition"
sidebar_current: "docs-newrelic-resource-federated-log-partition"
description: |-
  Create and manage Federated Log Partition entities within a Federated Log Setup.
---

# Resource: newrelic\_federated\_log\_partition

Use this resource to create and manage Federated Log Partition entities in New Relic.

## Example Usage

### Basic Federated Log Partition

```hcl
resource "newrelic_federated_log_partition" "application_logs" {
  name        = "application_logs"
  description = "Partition for application log data"
  setup_id    = newrelic_federated_log_setup.example.id

  partition_table   = "application_logs_table"
  data_location_uri = "s3://my-company-federated-logs-bucket/data/application_logs/"
}
```

### Federated Log Partition with Retention Policy

```hcl
resource "newrelic_federated_log_partition" "security_logs" {
  name        = "security_logs"
  description = "Partition for security and audit logs with extended retention"
  setup_id    = newrelic_federated_log_setup.example.id

  partition_table   = "security_logs_table"
  data_location_uri = "s3://my-company-federated-logs-bucket/data/security_logs/"

  retention_policy {
    duration = 365
    unit     = "DAYS"
  }
}
```

### Federated Log Partition with Partition Rules

```hcl
resource "newrelic_federated_log_partition" "error_logs" {
  name        = "error_logs"
  description = "Partition for error-level logs only"
  setup_id    = newrelic_federated_log_setup.example.id

  partition_table   = "error_logs_table"
  data_location_uri = "s3://my-company-federated-logs-bucket/data/error_logs/"

  # Route only ERROR and FATAL logs to this partition
  partition_rule {
    expression = "severity_number >= 17"
  }

  retention_policy {
    duration = 90
    unit     = "DAYS"
  }
}
```

### Multiple Partitions for a Setup

```hcl
resource "newrelic_federated_log_setup" "main_setup" {
  name                         = "main-logs-setup"
  cloud_provider               = "AWS"
  cloud_provider_region        = "us-east-1"
  data_location_bucket         = "company-logs-bucket"
  partition_database           = "logs_glue_db"
  data_processing_component_id = newrelic_federated_log_processor.processor.id

  data_ingest_connection {
    role_arn       = "arn:aws:iam::123456789012:role/logs-pcg-writer"
    aws_account_id = "123456789012"
  }

  query_connection {
    role_arn       = "arn:aws:iam::123456789012:role/logs-query-role"
    aws_account_id = "123456789012"
  }
}

# Default partition (created automatically, but can be managed explicitly)
resource "newrelic_federated_log_partition" "default" {
  name        = "Log"
  description = "Default partition for unrouted logs"
  setup_id    = newrelic_federated_log_setup.main_setup.id
  is_default  = true

  partition_table   = "default_logs_table"
  data_location_uri = "s3://company-logs-bucket/data/default/"

  retention_policy {
    duration = 30
    unit     = "DAYS"
  }
}

# Application logs partition
resource "newrelic_federated_log_partition" "app_logs" {
  name        = "Log_Application"
  description = "Application service logs"
  setup_id    = newrelic_federated_log_setup.main_setup.id

  partition_table   = "app_logs_table"
  data_location_uri = "s3://company-logs-bucket/data/application/"

  partition_rule {
    expression = "attributes[\"log.type\"] == \"application\""
  }

  retention_policy {
    duration = 60
    unit     = "DAYS"
  }
}

# Infrastructure logs partition
resource "newrelic_federated_log_partition" "infra_logs" {
  name        = "Log_Infrastructure"
  description = "Infrastructure and system logs"
  setup_id    = newrelic_federated_log_setup.main_setup.id

  partition_table   = "infra_logs_table"
  data_location_uri = "s3://company-logs-bucket/data/infrastructure/"

  partition_rule {
    expression = "attributes[\"log.type\"] == \"infrastructure\""
  }

  retention_policy {
    duration = 14
    unit     = "DAYS"
  }
}

# Security logs partition with longer retention
resource "newrelic_federated_log_partition" "security_logs" {
  name        = "Log_Security"
  description = "Security and compliance logs"
  setup_id    = newrelic_federated_log_setup.main_setup.id

  partition_table   = "security_logs_table"
  data_location_uri = "s3://company-logs-bucket/data/security/"

  partition_rule {
    expression = "attributes[\"log.type\"] == \"security\" or attributes[\"log.type\"] == \"audit\""
  }

  retention_policy {
    duration = 1
    unit     = "YEARS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the partition. This corresponds to the event type used in NRQL queries (e.g., `Log_Application`). Must be unique within the setup.
* `description` - (Optional) A description of the partition.
* `setup_id` - (Required) The ID of the `newrelic_federated_log_setup` this partition belongs to.
* `partition_table` - (Required) The name of the Iceberg table in the Glue Catalog database that stores this partition's data.
* `data_location_uri` - (Required) The S3 URI where this partition's data is stored (e.g., `s3://bucket-name/path/to/partition/`).
* `is_default` - (Optional) Whether this is the default partition for the setup. Logs that don't match any partition rules are routed to the default partition. Defaults to `false`. **Note**: Only one partition per setup can be marked as default.
* `partition_rule` - (Optional) The rule that determines which logs are routed to this partition. See [Partition Rule](#partition-rule) below. If not specified, this partition will only receive logs explicitly routed to it or (if `is_default = true`) logs that don't match other partitions.
* `retention_policy` - (Optional) The retention policy for logs in this partition. See [Retention Policy](#retention-policy) below.
* `account_id` - (Optional) The New Relic account ID. Defaults to the account ID configured in the provider.

### Partition Rule

The `partition_rule` block defines which logs are routed to this partition:

* `expression` - (Required) An OTTL (OpenTelemetry Transformation Language) expression that evaluates to true for logs that should be routed to this partition.

OTTL expression examples:
* `attributes["service.name"] == "payment-service"` - Route logs from a specific service
* `attributes["log.type"] == "application"` - Route based on custom log type attribute
* `severity_number >= 17` - Route ERROR level and above (17 = ERROR, 21 = FATAL)
* `attributes["k8s.namespace.name"] == "production"` - Route Kubernetes logs by namespace
* `IsMatch(attributes["message"], ".*ERROR.*")` - Route logs containing "ERROR" in message

### Retention Policy

The `retention_policy` block defines how long logs are retained:

* `duration` - (Required) The numeric value for retention duration.
* `unit` - (Required) The time unit for retention. Valid values are:
  * `DAYS` - Duration in days
  * `WEEKS` - Duration in weeks
  * `MONTHS` - Duration in months

**Note:** Retention is enforced by a Glue ETL job running in the customer's AWS account. The job reads retention settings from NGEP and deletes expired data from Iceberg tables.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (NGEP entity GUID) of the Federated Log Partition.
* `partition_id` - A unique identifier used internally for linking partitions across UI, Terraform, and data processing components.
* `status` - The current onboarding status of the partition. Possible values are:
  * `RESOURCE_CREATION_COMPLETE` - AWS resources have been created
  * `DATA_PROCESSING_COMPONENT_LINKED` - The processor has been linked via the parent setup
  * `DATA_PROCESSING_COMPONENT_DEPLOYED` - PCG has been deployed with this partition configuration
  * `ACTIVE` - The partition is fully active and receiving logs
  * `INACTIVE` - The partition is inactive (existing data can still be queried, but no new logs are written)
  * `ERROR` - The partition is in an error state
* `retention_policy_status` - The status of retention policy enforcement:
  * `ACTIVE` - Retention policy is being enforced
  * `ERROR` - Retention policy enforcement encountered an error

## Import

Federated Log Partitions can be imported using the entity GUID:

```bash
$ terraform import newrelic_federated_log_partition.example <entity_guid>
```

### Deleting a Partition

**Warning:** Deleting a partition via Terraform removes the NGEP entity but does **not** delete the underlying AWS resources (S3 data, Iceberg table). You must clean up AWS resources separately if desired.


