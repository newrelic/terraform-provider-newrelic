---
layout: "newrelic"
page_title: "New Relic: newrelic_federated_log_processor"
sidebar_current: "docs-newrelic-resource-federated-log-processor"
description: |-
  Create and manage Federated Log Processor entities for data processing in Federated Logs setups.
---

# Resource: newrelic\_federated\_log\_processor

Use this resource to create and manage Federated Log Processor entities in New Relic.

## Example Usage

### Basic Federated Log Processor with IRSA

```hcl
resource "newrelic_federated_log_processor" "example" {
  name           = "production-pcg-processor"
  description    = "Data processor for production federated logs"
  cloud_provider = "AWS"
  region         = "us-east-1"
  auth_mode      = "IRSA"

  auth_connection {
    role_arn         = "arn:aws:iam::123456789012:role/newrelic-fed-logs-base-role"
    aws_account_id   = "123456789012"
    external_id      = "newrelic-federated-logs"
  }
}
```

### Federated Log Processor with EKS Pod Identity

```hcl
resource "newrelic_federated_log_processor" "pod_identity" {
  name           = "staging-pcg-processor"
  description    = "Data processor using EKS Pod Identity"
  cloud_provider = "AWS"
  region         = "us-west-2"
  auth_mode      = "POD_IDENTITY"

  auth_connection {
    role_arn         = "arn:aws:iam::123456789012:role/newrelic-fed-logs-base-role"
    aws_account_id   = "123456789012"
  }
}
```

### Federated Log Processor Linked to Fleet

```hcl
resource "newrelic_fleet" "pcg_fleet" {
  name                = "PCG Kubernetes Fleet"
  managed_entity_type = "KUBERNETESCLUSTER"
  description         = "Fleet for Pipeline Control Gateway"
}

resource "newrelic_federated_log_processor" "with_fleet" {
  name           = "production-pcg-processor"
  description    = "Data processor linked to Fleet Control"
  cloud_provider = "AWS"
  region         = "us-east-1"
  auth_mode      = "IRSA"
  fleet_id       = newrelic_fleet.pcg_fleet.id

  auth_connection {
    role_arn         = "arn:aws:iam::123456789012:role/newrelic-fed-logs-base-role"
    aws_account_id   = "123456789012"
    external_id      = "newrelic-federated-logs"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Federated Log Processor. This is used as the PCG instance identifier and becomes the `PCG_Instance` tag value for ABAC.
* `description` - (Optional) A description of the Federated Log Processor.
* `cloud_provider` - (Required) The cloud provider where this processor is deployed. Valid values are `AWS`, `AZURE`, `GCP`, or `OCI`. **Note**: Currently only `AWS` is fully supported.
* `region` - (Required) The cloud provider region where this processor is deployed (e.g., `us-east-1` for AWS).
* `auth_mode` - (Required) The authentication mode used by the processor runtime on Kubernetes. Valid values are:
  * `IRSA` - IAM Roles for Service Accounts
  * `POD_IDENTITY` - EKS Pod Identity
* `auth_connection` - (Required) The AWS connection configuration for the base IAM role. See [Auth Connection](#auth-connection) below.
* `fleet_id` - (Optional) The ID of the Fleet entity to link with this processor. This enables Fleet Control integration.
* `account_id` - (Optional) The New Relic account ID. Defaults to the account ID configured in the provider.
* `organization_id` - (Optional) The New Relic organization ID. If not provided, it will be automatically determined from the account.

### Auth Connection

The `auth_connection` block supports:

* `role_arn` - (Required) The ARN of the base IAM role that PCG pods will assume. This role should have an ABAC policy allowing it to assume target setup roles.
* `aws_account_id` - (Required) The AWS account ID where the IAM role is created.
* `external_id` - (Optional) The external ID for the IAM role trust policy. Recommended for security.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (NGEP entity GUID) of the Federated Log Processor.
* `status` - The current status of the processor. Possible values are:
  * `CREATING` - The processor is being created
  * `ACTIVE` - The processor is active and ready for use
  * `INACTIVE` - The processor is inactive
  * `ERROR` - The processor is in an error state
* `auth_connection_id` - The ID of the AWS Connection entity created for the base role.

## Import

Federated Log Processors can be imported using the entity GUID:

```bash
$ terraform import newrelic_federated_log_processor.example <entity_guid>
```

### Tags

The base role must be tagged with:
* `PCG_Instance` = `<processor_name>` (matches the `name` argument)

