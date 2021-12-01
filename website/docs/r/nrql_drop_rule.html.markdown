---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_drop_rule"
sidebar_current: "docs-newrelic-resource-nrql-drop-rule"
description: |-
  Create and manage NRQL Drop Rules.
---

# Resource: newrelic\_nrql\_drop\_rule

Use this resource to create, and delete New Relic NRQL Drop Rules.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/docs/providers/newrelic/index.html) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

-> **Drop attributes on dimensional metric rollups is in preview**
Please contact your account team or fill in this [form](https://forms.gle/FGPZpy2y1zdSN9dn9) if you'd like to enroll your account. The feature is subject to change, and more information can be found in [the docs](https://docs.newrelic.com/docs/data-apis/manage-data/drop-data-using-nerdgraph/#drop-attributes-on-dimensional-metric-rollups).

## Example Usage

```hcl
resource "newrelic_nrql_drop_rule" "foo" {
  account_id  = 12345
  description = "Drops all data for MyCustomEvent that comes from the LoadGeneratingApp in the dev environment, because there is too much and we don’t look at it."
  action      = "drop_data"
  nrql        = "SELECT * FROM MyCustomEvent WHERE appName='LoadGeneratingApp' AND environment='development'"
}

resource "newrelic_nrql_drop_rule" "bar" {
  account_id  = 12345
  description = "Removes the user name and email fields from MyCustomEvent"
  action      = "drop_attributes"
  nrql        = "SELECT userEmail, userName FROM MyCustomEvent"
}

resource "newrelic_nrql_drop_rule" "baz" {
  account_id  = 12345
  description = "Removes containerId from metric aggregates to reduce metric cardinality."
  action      = "drop_attributes_from_metric_aggregates"
  nrql        = "SELECT containerId FROM Metric"
}
```

## Argument Reference

The following arguments are supported:

  * `account_id` - (Optional) Account where the drop rule will be put. Defaults to the account associated with the API key used.
  * `description` - (Optional) The description of the drop rule.
  * `nrql` - (Required) A NRQL string that specifies what data types to drop.
  * `action` - (Required) An action type specifying how to apply the NRQL string (either `drop_data`, `drop_attributes`, or ` drop_attributes_from_metric_aggregates`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `rule_id` - The id, uniquely identifying the rule.

## Import

New Relic NRQL drop rules can be imported using a concatenated string of the format
 `<account_id>:<rule_id>`, e.g.

```bash
$ terraform import newrelic_nrql_drop_rule.foo 12345:34567
```
