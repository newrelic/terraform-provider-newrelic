---
layout: "newrelic"
page_title: "New Relic: newrelic_metric_pruning_rule"
sidebar_current: "docs-newrelic-resource-metric-pruning-rule"
description: |-
  Create and manage New Relic metric pruning rules.
---

# Resource: newrelic\_metric\_pruning\_rule

Use this resource to create and manage metric pruning rules for a New Relic account.

A metric pruning rule strips specific high-cardinality attributes from dimensional metric aggregates before they are stored. Unlike dropping a metric entirely, pruning keeps the metric signal intact while removing the nominated attributes — reducing cardinality without any loss of the metric itself.

---

## Example Usage

```hcl
resource "newrelic_metric_pruning_rule" "example" {
  nrql        = "SELECT collector.name FROM Metric WHERE metricName = 'http.server.duration'"
  description = "Remove collector.name attribute from http.server.duration to reduce cardinality"
}
```

### With an explicit account ID

```hcl
resource "newrelic_metric_pruning_rule" "example" {
  account_id  = 12345678
  nrql        = "SELECT collector.name FROM Metric WHERE metricName = 'http.server.duration'"
  description = "Remove collector.name attribute from http.server.duration to reduce cardinality"
}
```

---

## Argument Reference

The following arguments are supported:

* `nrql` - (Required, Forces new resource) The NRQL query that identifies the metric and the specific attributes to prune. The `SELECT` clause must name the attributes to remove (not `SELECT *`), and the `FROM` clause must target `Metric`. Example: `SELECT collector.name FROM Metric WHERE metricName = 'my.metric.name'`.
* `description` - (Optional, Forces new resource) A human-readable description of what this pruning rule does.
* `account_id` - (Optional, Forces new resource) The New Relic account ID in which to create the pruning rule. Defaults to the account ID configured on the provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:<rule_id>`.
* `rule_id` - The unique identifier of the pruning rule assigned by New Relic.

## Import

Metric pruning rules can be imported using the composite ID format `<account_id>:<rule_id>`:

```bash
$ terraform import newrelic_metric_pruning_rule.example 12345678:1234
```

-> **Note:** All attributes (`nrql`, `description`) are restored from New Relic on import.
