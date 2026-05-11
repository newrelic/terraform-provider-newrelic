---
layout: "newrelic"
page_title: "New Relic: newrelic_account_cardinality_limit"
sidebar_current: "docs-newrelic-resource-account-cardinality-limit"
description: |-
  Create and manage New Relic account cardinality limit overrides.
---

# Resource: newrelic\_account\_cardinality\_limit

Use this resource to create and manage cardinality limit overrides for a New Relic account. Two modes are supported, controlled by the required `mode` argument:

- **`DEFAULT`**: Override the account-wide default cardinality limit applied to all dimensional metrics that do not have a per-metric override.
- **`PER_METRIC`**: Override the cardinality limit for a single named metric (requires `metric_name`).

-> **NOTE:** The New Relic API does not provide a delete operation for cardinality limit overrides. Destroying this resource removes it from Terraform state only; the override remains active in New Relic until it is changed by a subsequent `apply` or modified externally.

## Example Usage

### DEFAULT Mode — Account-Wide Default Override

```hcl
resource "newrelic_account_cardinality_limit" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}
```

### PER_METRIC Mode — Per-Metric Override

```hcl
resource "newrelic_account_cardinality_limit" "per_metric" {
  mode              = "PER_METRIC"
  metric_name       = "otelcol_nrreceiver_incoming_request_proxy"
  cardinality_limit = 200000
}
```

### Both Together

```hcl
resource "newrelic_account_cardinality_limit" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}

resource "newrelic_account_cardinality_limit" "per_metric" {
  mode              = "PER_METRIC"
  metric_name       = "otelcol_nrreceiver_incoming_request_proxy"
  cardinality_limit = 200000
}
```

## Argument Reference

The following arguments are supported:

* `mode` - (Required) The override mode. Must be `"DEFAULT"` or `"PER_METRIC"`. Forces re-creation when changed.
  * `DEFAULT` — sets the account-wide default limit for all metrics not individually overridden. `metric_name` must not be set.
  * `PER_METRIC` — overrides the limit for a single metric. `metric_name` is required.
* `metric_name` - (Optional) The name of the metric to override. Required when `mode` is `"PER_METRIC"`; must not be set when `mode` is `"DEFAULT"`. Forces re-creation when changed.
* `cardinality_limit` - (Required) The cardinality limit value — the maximum number of unique dimension-value combinations allowed per day.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:<metric_name>`. For `DEFAULT` mode the ID ends with `:` (e.g. `12345678:`). For `PER_METRIC` mode the metric name follows the colon (e.g. `12345678:my_metric`).

## Import

Cardinality limit overrides can be imported using the composite ID format `<account_id>:<metric_name>`.

For a **DEFAULT** override (no metric name), use a trailing colon:

```bash
$ terraform import newrelic_account_cardinality_limit.default 12345678:
```

For a **PER_METRIC** override, append the metric name after the colon:

```bash
$ terraform import newrelic_account_cardinality_limit.per_metric 12345678:otelcol_nrreceiver_incoming_request_proxy
```

-> **NOTE:** Because the New Relic API does not return per-metric qualifier information via the management query, importing a `PER_METRIC` override restores `mode` and `metric_name` from the resource ID, but `cardinality_limit` will reflect the value read via NRQL for that qualifier. Run `terraform apply` after import if the state drifts.
