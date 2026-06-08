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

### Pruning the same attribute from many metrics in bulk

When the same attribute needs to be stripped from a set of metrics, keep the metric list in `locals` and use `for_each` with a templated `nrql` and `description`. Only the metric name varies between rules — the rest of the configuration is shared.

```hcl
locals {
  # The attribute to prune from every metric in the list below.
  pruned_attribute = "collector.name"

  # Metrics to apply the pruning rule to.
  # Add or remove entries here to manage rules in bulk.
  metrics_to_prune = toset([
    "http.server.duration",
    "http.client.duration",
    "rpc.server.duration",
    "k8s.pod.cpu.usage",
    "k8s.pod.memory.usage",
    # ...add more entries here
  ])
}

resource "newrelic_metric_pruning_rule" "bulk" {
  for_each    = local.metrics_to_prune
  nrql        = "SELECT ${local.pruned_attribute} FROM Metric WHERE metricName = '${each.value}'"
  description = "Remove ${local.pruned_attribute} from ${each.value} to reduce cardinality"
}
```

Each entry in `metrics_to_prune` produces an independent `newrelic_metric_pruning_rule` resource (e.g. `newrelic_metric_pruning_rule.bulk["http.server.duration"]`) that can be inspected, imported, or destroyed individually.

---

## Behaviour

- **`terraform apply`** — creates the pruning rule in New Relic. The rule begins stripping the nominated attributes from matching metric aggregates immediately after creation.
- **`terraform plan` / `terraform refresh` on an existing resource** — reads the current state of the pruning rule from New Relic and surfaces any drift (e.g. if the rule was deleted outside of Terraform).
- **`terraform destroy`** — permanently deletes the pruning rule. Once removed, the nominated attributes will no longer be stripped from incoming metric data. There is no reset to a default state; the rule is deleted outright.

-> **Note:** Because all arguments are immutable, any in-place change (e.g. updating the NRQL or description) will trigger a destroy-and-recreate. The old rule is deleted before the new one is created, so there will be a brief window during which no pruning is active for the affected metric.

---

## Argument Reference

The following arguments are supported:

* `nrql` - (Required) The NRQL query that identifies the metric and the specific attributes to prune. The `SELECT` clause must name the attributes to remove (not `SELECT *`), and the `FROM` clause must target `Metric`. Example: `SELECT collector.name FROM Metric WHERE metricName = 'my.metric.name'`.
* `description` - (Optional) A human-readable description of what this pruning rule does.
* `account_id` - (Optional) The New Relic account ID in which to create the pruning rule. Defaults to the account ID configured on the provider.

-> **Note:** All arguments on this resource are immutable. Any change to an existing `newrelic_metric_pruning_rule` will cause the resource to be destroyed and recreated with the updated configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:<rule_id>`.
* `rule_id` - The unique identifier of the pruning rule assigned by New Relic.

## Import

Metric pruning rules can be imported using the composite ID format `<account_id>:<rule_id>`:

```bash
$ terraform import newrelic_metric_pruning_rule.example 12345678:1234
```
