---
layout: "newrelic"
page_title: "New Relic: newrelic_federated_log_setup"
sidebar_current: "docs-newrelic-resource-federated-log-setup"
description: |-
  Create and manage Federated Log Setup entities for storing logs in customer-owned cloud storage.
---

# Resource: newrelic\_federated\_log\_setup

Use this resource to create and manage Federated Log Setup entities in New Relic.

## Example Usage

### Basic Federated Log Setup

```hcl
resource "newrelic_federated_log_processor" "processor" {
  name           = "production-pcg"
  cloud_provider = "AWS"
  region         = "us-east-1"
  auth_mode      = "IRSA"

  auth_connection {
    role_arn       = "arn:aws:iam::123456789012:role/newrelic-fed-logs-base-role"
    aws_account_id = "123456789012"
  }
}

resource "newrelic_federated_log_setup" "example" {
  name                  = "production-logs-setup"
  description           = "Federated logs setup for production environment"
  cloud_provider        = "AWS"
  cloud_provider_region = "us-east-1"

  # Storage configuration
  data_location_bucket = "my-company-federated-logs-bucket"
  partition_database   = "federated_logs_db"

  # Link to data processing component
  data_processing_component_id = newrelic_federated_log_processor.processor.id

  # IAM roles for data access
  data_ingest_connection {
    role_arn       = "arn:aws:iam::123456789012:role/newrelic-fed-logs-prod-pcg-writer"
    aws_account_id = "123456789012"
  }

  query_connection {
    role_arn       = "arn:aws:iam::123456789012:role/newrelic-fed-logs-prod-query-role"
    aws_account_id = "123456789012"
    external_id    = "newrelic-query-access"
  }
}
```

### Federated Log Setup with Routing Rules

```hcl
resource "newrelic_federated_log_setup" "with_routing" {
  name                  = "application-logs-setup"
  description           = "Setup for application logs with routing rules"
  cloud_provider        = "AWS"
  cloud_provider_region = "us-west-2"

  data_location_bucket         = "app-logs-bucket"
  partition_database           = "app_logs_glue_db"
  data_processing_component_id = newrelic_federated_log_processor.processor.id

  data_ingest_connection {
    role_arn       = "arn:aws:iam::123456789012:role/app-logs-pcg-writer"
    aws_account_id = "123456789012"
  }

  query_connection {
    role_arn       = "arn:aws:iam::123456789012:role/app-logs-query-role"
    aws_account_id = "123456789012"
  }

  # Routing rule using OTTL expression
  routing_rule {
    expression = "attributes[\"service.name\"] == \"my-application\""
  }
}
```

### Multiple Setups with Single Processor

```hcl
resource "newrelic_federated_log_processor" "shared_processor" {
  name           = "shared-pcg-processor"
  cloud_provider = "AWS"
  region         = "us-east-1"
  auth_mode      = "IRSA"

  auth_connection {
    role_arn       = "arn:aws:iam::123456789012:role/newrelic-fed-logs-base-role"
    aws_account_id = "123456789012"
  }
}

# Setup 1: Production logs
resource "newrelic_federated_log_setup" "production" {
  name                         = "production-logs"
  cloud_provider               = "AWS"
  cloud_provider_region        = "us-east-1"
  data_location_bucket         = "prod-logs-bucket"
  partition_database           = "prod_logs_db"
  data_processing_component_id = newrelic_federated_log_processor.shared_processor.id

  data_ingest_connection {
    role_arn       = "arn:aws:iam::123456789012:role/prod-logs-pcg-writer"
    aws_account_id = "123456789012"
  }

  query_connection {
    role_arn       = "arn:aws:iam::123456789012:role/prod-logs-query-role"
    aws_account_id = "123456789012"
  }

  routing_rule {
    expression = "attributes[\"environment\"] == \"production\""
  }
}

# Setup 2: Staging logs (same processor, different storage)
resource "newrelic_federated_log_setup" "staging" {
  name                         = "staging-logs"
  cloud_provider               = "AWS"
  cloud_provider_region        = "us-east-1"
  data_location_bucket         = "staging-logs-bucket"
  partition_database           = "staging_logs_db"
  data_processing_component_id = newrelic_federated_log_processor.shared_processor.id

  data_ingest_connection {
    role_arn       = "arn:aws:iam::123456789012:role/staging-logs-pcg-writer"
    aws_account_id = "123456789012"
  }

  query_connection {
    role_arn       = "arn:aws:iam::123456789012:role/staging-logs-query-role"
    aws_account_id = "123456789012"
  }

  routing_rule {
    expression = "attributes[\"environment\"] == \"staging\""
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Federated Log Setup. Must be unique within the account.
* `description` - (Optional) A description of the Federated Log Setup.
* `cloud_provider` - (Required) The cloud provider where storage is deployed. Valid values are `AWS`, `AZURE`, `GCP`, or `OCI`. **Note**: Currently only `AWS` is fully supported.
* `cloud_provider_region` - (Required) The cloud provider region where storage is deployed (e.g., `us-east-1` for AWS).
* `data_location_bucket` - (Required) The name of the S3 bucket (or equivalent object storage) where log data is stored.
* `partition_database` - (Required) The name of the AWS Glue Catalog database (or equivalent) containing the Iceberg tables.
* `data_processing_component_id` - (Required) The ID of the `newrelic_federated_log_processor` resource that will write to this setup. This enables multiple setup support.
* `data_ingest_connection` - (Required) The AWS connection configuration for the PCG writer role. See [Data Ingest Connection](#data-ingest-connection) below.
* `query_connection` - (Required) The AWS connection configuration for the query/reader role. See [Query Connection](#query-connection) below.
* `routing_rule` - (Optional) The routing rule that determines which logs are sent to this setup. See [Routing Rule](#routing-rule) below.
* `account_id` - (Optional) The New Relic account ID. Defaults to the account ID configured in the provider.

### Data Ingest Connection

The `data_ingest_connection` block configures the IAM role used by PCG to write logs:

* `role_arn` - (Required) The ARN of the IAM role for PCG to write to S3 and Glue. This role must:
  * Have a trust policy allowing the base role (from `newrelic_federated_log_processor`) to assume it
  * Be tagged with `PCG_Instance` matching the processor's name
  * Have permissions to write to the S3 bucket and update Glue tables
* `aws_account_id` - (Required) The AWS account ID where the IAM role exists.
* `external_id` - (Optional) The external ID for the IAM role trust policy.

### Query Connection

The `query_connection` block configures the IAM role used by New Relic to query logs:

* `role_arn` - (Required) The ARN of the IAM role for New Relic query workers to read from S3 and Glue. This role must have read-only permissions.
* `aws_account_id` - (Required) The AWS account ID where the IAM role exists.
* `external_id` - (Optional) The external ID for the IAM role trust policy. Recommended for security.

### Routing Rule

The `routing_rule` block defines which logs are routed to this setup:

* `expression` - (Required) An OTTL (OpenTelemetry Transformation Language) expression that evaluates to true for logs that should be routed to this setup.

Common OTTL expression examples:
* `attributes["service.name"] == "my-service"` - Route logs from a specific service
* `attributes["environment"] == "production"` - Route logs from a specific environment
* `attributes["cloud.provider"] == "aws"` - Route logs from AWS resources
* `severity_number >= 17` - Route ERROR and above logs (severity 17+)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (NGEP entity GUID) of the Federated Log Setup.
* `setup_id` - A unique identifier used internally for linking setups across UI, Terraform, and data processing components.
* `status` - The current onboarding status of the setup. Possible values are:
  * `RESOURCE_CREATION_COMPLETE` - AWS resources have been created
  * `DATA_PROCESSING_COMPONENT_LINKED` - The processor has been linked
  * `DATA_PROCESSING_COMPONENT_DEPLOYED` - PCG has been deployed with this configuration
  * `ACTIVE` - The setup is fully active and processing logs
  * `INACTIVE` - The setup is inactive (data can still be queried)
  * `ERROR` - The setup is in an error state
* `default_partition_id` - The ID of the default partition created with this setup.
* `data_ingest_connection_id` - The ID of the AWS Connection entity for data ingestion.
* `query_connection_id` - The ID of the AWS Connection entity for querying.

## Import

Federated Log Setups can be imported using the entity GUID:

```bash
$ terraform import newrelic_federated_log_setup.example <entity_guid>
```

**Tags:**
* `PCG_Instance` = `<processor_name>` (must match the linked processor)

### Query Role (query_connection)

The query role must have read-only access to S3 and Glue. New Relic's query workers will assume this role to execute queries.
