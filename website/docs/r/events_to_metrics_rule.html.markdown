---
layout: "newrelic"
page_title: "New Relic: newrelic_events_to_metrics_rule"
sidebar_current: "docs-newrelic-resource-events-to-metrics-rule"
description: |-
  Create and manage tags for a New Relic Events to Metrics rule.
---

# Resource: newrelic\_events\_to\_metrics\_rule

Use this resource to create, update, and delete New Relic Events to Metrics rules.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/providers/newrelic/newrelic/latest/docs/guides/migration_guide_v2) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

```hcl
resource "newrelic_events_to_metrics_rule" "foo" {
  account_id = 12345
  name = "Example events to metrics rule"
  description = "Example description"
  nrql = "SELECT uniqueCount(account_id) AS ``Transaction.account_id`` FROM Transaction FACET appName, name"
}
```

## Argument Reference

The following arguments are supported:

  * `account_id` - (Optional) The account ID where the Events to Metrics rule will be created. If not specified, defaults to *the provider's configured account ID*.

-> **NOTE:** Setting `account_id` explicitly is recommended so the target account is unambiguous in your configuration.
  * `name` - (Required) The name of the rule. This must be unique within an account.
  * `nrql` - (Required) Explains how to create metrics from events.
  * `description` - (Optional) Provides additional information about the rule.
  * `enabled` - (Optional) True means this rule is enabled. False means the rule is currently not creating metrics.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `account_id` - The account ID the rule belongs to. Set to the provider account ID if not specified in the configuration.
  * `rule_id` - The id, uniquely identifying the rule.

## Import

New Relic Events to Metrics rules can be imported using a concatenated string of the format
 `<account_id>:<rule_id>`, e.g.

```bash
$ terraform import newrelic_events_to_metrics_rule.foo 12345:34567
```
