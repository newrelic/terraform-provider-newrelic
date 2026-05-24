---
layout: "newrelic"
page_title: "New Relic: newrelic_cardinality_management"
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

In `DEFAULT` mode, the resource sets the account-wide default cardinality limit. This value is applied to every dimensional metric in the account that does not have its own per-metric override.

### Behaviour

- **Create / Update**: Submits the new default value via the `dataManagementCreateAccountLimit` mutation. The change takes effect immediately in the enforcement layer.
- **Read**: Reads the current account-wide default from the New Relic data management API and reconciles Terraform state. Drift is detected and surfaced on the next `plan`.
- **Destroy**: Because the New Relic API does not expose a delete operation for cardinality limit overrides, destroying this resource resets the account-wide default back to the New Relic platform default of **100,000**. A warning is displayed to confirm the reset value.

### Example

```hcl
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}
```

-> **Note:** Destroying a `DEFAULT` mode resource resets the account-wide cardinality limit to the New Relic platform default of 100,000. Ensure this is intentional before running `terraform destroy`.

---

## PER\_METRIC Mode

In `PER_METRIC` mode, the resource overrides the cardinality limit for a single named metric, independent of the account-wide default.

### Behaviour

- **Create / Update**: Submits the per-metric override via the `dataManagementCreateAccountLimit` mutation with the metric name as the qualifier. A warning is displayed after apply indicating that the change may take a few minutes to be reflected in the New Relic UI.
- **Read**: The New Relic data management API does not expose per-metric qualifier values in its response, and the NRDB event stream for limit enforcement lags behind the mutation API. As a result, the `cardinality_limit` value in state is preserved from the last successful `apply` and is not synchronised from the API on each `plan`. A warning is displayed to indicate this limitation.
- **Destroy**: Because the New Relic API does not expose a delete operation, destroying this resource resets the metric's cardinality limit to the **current account-wide default** (fetched live from the API at destroy time). A warning is displayed confirming the value the limit was reset to.

### Example

```hcl
resource "newrelic_cardinality_management" "per_metric" {
  mode              = "PER_METRIC"
  metric_name       = "otelcol_nrreceiver_incoming_request_proxy"
  cardinality_limit = 200000
}
```

-> **Note:** Destroying a `PER_METRIC` resource does not remove the override from New Relic. It resets the metric's limit to the current account-wide default. If no `DEFAULT` override exists, this will be the New Relic platform default of 100,000.

-> **Note:** Due to API limitations, the `cardinality_limit` value for a `PER_METRIC` resource reflects the last value applied by Terraform and may not represent the currently enforced limit if the override was modified outside of Terraform.

---

## Using Both Modes Together

```hcl
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}

resource "newrelic_cardinality_management" "per_metric" {
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
$ terraform import newrelic_cardinality_management.default 12345678:
```

For a **PER_METRIC** override, append the metric name after the colon:

```bash
$ terraform import newrelic_cardinality_management.per_metric 12345678:otelcol_nrreceiver_incoming_request_proxy
```

-> **Note:** When importing a `PER_METRIC` resource, `mode` and `metric_name` are restored from the resource ID. However, because the API does not return the current per-metric override value, `cardinality_limit` will reflect the account-wide default read at import time. Run `terraform apply` after import to re-apply the intended limit.
