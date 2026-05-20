---
layout: "newrelic"
page_title: "New Relic: newrelic_account_cardinality_limit"
sidebar_current: "docs-newrelic-resource-account-cardinality-limit"
description: |-
  Create and manage New Relic account cardinality limit overrides.
---

# Resource: newrelic\_account\_cardinality\_limit

Use this resource to create and manage cardinality limit overrides for a New Relic account.

Dimensional metrics in New Relic are subject to a per-metric cardinality limit — the maximum number of unique attribute-value combinations a single metric name may produce per day. This resource allows you to override that limit either account-wide (for all metrics) or for a specific metric name.

Two modes are available, controlled by the required `mode` argument.

---

## DEFAULT Mode

In `DEFAULT` mode, the resource sets the account-wide default cardinality limit. This value applies to every dimensional metric in the account that does not have its own per-metric override. Terraform tracks the live value and will surface any drift on the next `terraform plan`. Running `terraform destroy` on this resource resets the account-wide limit back to the New Relic platform default of **100,000** — it does not remove the setting entirely.

### Example

```hcl
resource "newrelic_account_cardinality_limit" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}
```

-> **Note:** Destroying a `DEFAULT` mode resource resets the account-wide cardinality limit to the New Relic platform default of 100,000. Ensure this is intentional before running `terraform destroy`.

---

## PER\_METRIC Mode

In `PER_METRIC` mode, the resource overrides the cardinality limit for a single named metric, independent of the account-wide default. The updated limit takes effect as metric data is received and processed by the platform, and will be visible in the New Relic UI once the metric ingestion cycle completes. Terraform maintains the last applied value in state, as the enforced limit is tied to metric activity on the platform rather than being independently queryable — so `terraform plan` will not flag changes made outside of Terraform to this specific limit. Running `terraform destroy` resets the metric's limit to the current account-wide default rather than removing it entirely.

### Example

```hcl
resource "newrelic_account_cardinality_limit" "per_metric" {
  mode              = "PER_METRIC"
  metric_name       = "otelcol_nrreceiver_incoming_request_proxy"
  cardinality_limit = 200000
}
```

-> **Note:** Destroying a `PER_METRIC` resource resets the metric's limit to the current account-wide default. If no `DEFAULT` override exists, this will be the New Relic platform default of 100,000.

---

## Using Both Modes Together

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

---

## Argument Reference

The following arguments are supported:

* `mode` - (Required) The override mode. Must be `"DEFAULT"` or `"PER_METRIC"`. Forces re-creation when changed.
  * `DEFAULT` — sets the account-wide default limit for all metrics not individually overridden. `metric_name` must not be set.
  * `PER_METRIC` — overrides the limit for a single named metric. `metric_name` is required.
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

-> **Note:** When importing a `PER_METRIC` resource, `mode` and `metric_name` are restored from the resource ID. The `cardinality_limit` in state will be populated once you run `terraform apply` to establish the intended limit — the enforced value on the platform becomes fully active as metric data flows through.
