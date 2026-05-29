---
layout: "newrelic"
page_title: "New Relic: newrelic_aws_connection"
sidebar_current: "docs-newrelic-resource-aws-connection"
description: |-
  Create and manage an AWS Connection entity in New Relic.
---

# Resource: newrelic\_aws\_connection

Use this resource to create and manage an AWS Connection entity in New Relic.

## Example Usage

```hcl
resource "newrelic_aws_connection" "foo" {
  name        = "test-aws-connection"
  description = "AWS Connection wrapping the role"
  region      = "us-east-1"
  enabled     = true

  scope_type = "ORGANIZATION"
  scope_id   = "YOUR_ORG_ID_HERE"

  credential {
    assume_role {
      role_arn    = "arn:aws:iam::123456789012:role/newrelic-fed-logs-ingest"
      external_id = "external id"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AWS connection.
* `credential` - (Required) Credentials block describing how New Relic should authenticate into the AWS account. See [Nested credential block](#nested-credential-block) below.
* `description` - (Optional) A description of the AWS Connection.
* `enabled` - (Optional) Flag to indicate whether the connection is enabled. Defaults to `true`.
* `region` - (Optional) AWS region for this connection (e.g. `us-east-1`).
* `external_id` - (Optional) Consumer-managed identifier — useful for caller-side idempotent tracking. Distinct from `credential.assume_role.external_id` (the IAM cross-account External ID).
* `scope_type` - (Optional) Scope type for the connection. Valid values: `ACCOUNT`, `ORGANIZATION`.
* `scope_id` - (Optional) Scope ID matching `scope_type` — a New Relic account ID for `ACCOUNT` scope, or an organization ID for `ORGANIZATION` scope.
* `account_id` - (Optional) New Relic account ID where the connection will be created. Used when `scope_type = ACCOUNT`.
* `settings` - (Optional) Optional list of connection settings. Each entry takes:
  * `key` - (Required) The setting key.
  * `value` - (Required) The setting value.

### Nested `credential` block

The `credential` block describes how New Relic authenticates into the AWS account.

Each `credential` block supports:

* `assume_role` - (Required) STS:AssumeRole configuration. See below.

Each `assume_role` block supports:

* `role_arn` - (Required) ARN of the IAM role New Relic should assume.
* `external_id` - (Optional) External ID supplied by New Relic during STS:AssumeRole.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the AWS Connection.

## Import

AWS Connections can be imported using the entity GUID:

```bash
$ terraform import newrelic_aws_connection.foo <entity-guid>
```
