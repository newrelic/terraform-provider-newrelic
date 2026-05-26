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

  role_arn = "arn:aws:iam::123456789012:role/newrelic-fed-logs-ingest"
  region   = "us-east-1"
  enabled  = true

  scope_type = "ORGANIZATION"
  scope_id   = "YOUR_ORG_ID_HERE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the AWS connection.
* `role_arn` - (Required) The ARN of the IAM role.
* `description` - (Optional) A description of AWS Connection.
* `enabled` - (Optional) Flag to indicate whether the connection is enabled. Defaults to `true`.
* `region` - (Optional) AWS region for this connection (e.g. `us-east-1`).
* `external_id` - (Optional) Consumer-managed identifier — useful for caller-side idempotent tracking.
* `scope_type` - (Optional) Scope type for the connection. Valid values: `ACCOUNT`, `ORGANIZATION`. 
* `scope_id` - (Optional) Scope ID matching `scope_type` — a New Relic account ID for `ACCOUNT` scope, or an organization ID for `ORGANIZATION` scope.
* `account_id` - (Optional) New Relic account ID where the connection will be created. Used when `scope_type = ACCOUNT`.
* `settings` - (Optional) Optional list of connection settings. Each entry takes:
  * `key` - (Required) The setting key.
  * `value` - (Required) The setting value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the AWS Connection.

## Import

AWS Connections can be imported using the entity GUID:

```bash
$ terraform import newrelic_aws_connection.foo <entity-guid>
```

The entity GUID is the same value `terraform output` would produce after a successful apply.
