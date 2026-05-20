---
layout: "newrelic"
page_title: "New Relic: newrelic_metric_pruning_rule"
sidebar_current: "docs-newrelic-resource-metric-pruning-rule"
description: |-
  Create and manage New Relic metric pruning rules.
---

# Resource: newrelic\_metric\_pruning\_rule

Use this resource to create and manage **metric pruning rules** for a New Relic account.

A metric pruning rule strips specific high-cardinality attributes from dimensional metric aggregates before they are stored. Unlike a full data drop, pruning keeps the metric itself but removes the nominated attributes, reducing cardinality without losing the metric signal entirely.

Internally, this resource uses the `nrqlDropRulesCreate` NerdGraph mutation with the `DROP_ATTRIBUTES_FROM_METRIC_AGGREGATES` action.

---

## Example Usage

```hcl
resource "newrelic_metric_pruning_rule" "example" {
  nrql        = "SELECT * FROM Metric WHERE metricName = 'scooter.speed.kmph'"
  description = "Drop high-cardinality rider_id from scooter speed metric"
}
```

### With an explicit account ID

```hcl
resource "newrelic_metric_pruning_rule" "example" {
  account_id  = 12345678
  nrql        = "SELECT * FROM Metric WHERE metricName = 'scooter.speed.kmph'"
  description = "Drop high-cardinality rider_id from scooter speed metric"
}
```

---

## Argument Reference

The following arguments are supported:

* `nrql` - (Required, Forces new resource) The NRQL query that identifies the metric and attributes to prune. The query must target a specific metric name. For example: `SELECT * FROM Metric WHERE metricName = 'my.metric.name'`.
* `description` - (Optional, Forces new resource) A human-readable description of the pruning rule.
* `account_id` - (Optional) The account ID in which to create the pruning rule. Defaults to the account ID configured on the provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:<rule_id>`.
* `rule_id` - The unique identifier of the pruning rule assigned by New Relic.

## Import

Metric pruning rules can be imported using the composite ID format `<account_id>:<rule_id>`:

```bash
$ terraform import newrelic_metric_pruning_rule.example 12345678:1234
```

-> **Note:** All attributes (`nrql`, `description`) are restored from the New Relic API on import.
