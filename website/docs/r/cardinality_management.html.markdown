---
layout: "newrelic"
page_title: "New Relic: newrelic_cardinality_management"
sidebar_current: "docs-newrelic-resource-cardinality-management"
description: |-
  Manage New Relic account cardinality limit overrides.
---

# Resource: newrelic\_cardinality\_management

Use this resource to manage cardinality limit overrides for a New Relic account.

Dimensional metrics in New Relic are subject to a per-metric cardinality limit — the maximum number of unique attribute-value combinations a single metric name may produce per day. This resource lets you raise or lower that limit, either account-wide (for all metrics at once) or for individual named metrics.

Two modes are available, selected via the required `mode` argument.

-> **Note:** The New Relic API does not expose a delete operation for cardinality limit overrides. Destroying this resource resets the affected limit(s) back to the platform default of **100,000** rather than removing them. The destroy behaviour for each mode is described in the sections below.

---

## DEFAULT Mode

Sets a single account-wide limit that applies to every dimensional metric in the account that does not have its own per-metric override.

### Example

```hcl
resource "newrelic_cardinality_management" "account_default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}
```

### Behaviour

- **Create / Update** — submits the new default value. The change takes effect in the enforcement layer straight away.
- **Read** — reads the current account-wide default from the New Relic data management API and reconciles Terraform state. Drift is detected on the next `terraform plan`.
- **Destroy** — resets the account-wide default to the New Relic platform default of **100,000** and displays a confirmation warning.

-> **Note:** Changes may take a few minutes to be visible in the New Relic UI, particularly if affected metrics have not sent data recently.

---

## PER\_METRIC Mode

Sets individual cardinality limits for one or more named metrics. Each metric is configured in its own `metric` block and can have a different limit.

### Example — single metric

```hcl
resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"

  metric {
    name  = "http.server.duration"
    limit = 200000
  }
}
```

### Example — multiple metrics in one resource

```hcl
resource "newrelic_cardinality_management" "high_cardinality_metrics" {
  mode = "PER_METRIC"

  metric {
    name  = "http.server.duration"
    limit = 200000
  }

  metric {
    name  = "otelcol_nrreceiver_incoming_request_proxy"
    limit = 300000
  }

  metric {
    name  = "k8s.pod.cpu.usage"
    limit = 150000
  }
}
```

### Behaviour

- **Create / Update** — submits one override per `metric` block. A warning is displayed after apply as a reminder that updates may take a few minutes to be reflected in the UI.
- **Read** — the New Relic API does not return per-metric override values, so the `limit` values in state are preserved from the last successful `apply`. A warning is displayed on each `plan` and `apply` to reflect this limitation.
- **Destroy** — resets each managed metric's limit to the platform default of **100,000**. A warning is displayed to confirm this.

-> **Note:** Because the API does not return per-metric limit values, the `limit` attributes in state always reflect the last values applied by Terraform. If a limit was changed outside of Terraform, run `terraform apply` to re-apply the desired values.

---

## Using Both Modes Together

You can manage both the account-wide default and individual metric overrides at the same time:

```hcl
resource "newrelic_cardinality_management" "account_default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}

resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"

  metric {
    name  = "http.server.duration"
    limit = 250000
  }

  metric {
    name  = "otelcol_nrreceiver_incoming_request_proxy"
    limit = 300000
  }
}
```

---

## Argument Reference

The following arguments are supported:

* `mode` - (Required) The override mode. Accepted values: `"DEFAULT"` or `"PER_METRIC"`. Forces re-creation when changed.
  * `DEFAULT` — sets an account-wide limit for all metrics. `cardinality_limit` is required; `metric` blocks must not be set.
  * `PER_METRIC` — sets individual limits for named metrics. At least one `metric` block is required; `cardinality_limit` must not be set at the top level.

* `cardinality_limit` - (Optional) The account-wide cardinality limit value — the maximum unique dimension-value combinations allowed per metric per day. Required when `mode` is `"DEFAULT"`; must not be set when `mode` is `"PER_METRIC"`.

* `metric` - (Optional) One or more metric override blocks. Required when `mode` is `"PER_METRIC"`; must not be set when `mode` is `"DEFAULT"`. Each block supports:
  * `name` - (Required) The full name of the metric (e.g. `"http.server.duration"`).
  * `limit` - (Required) The maximum number of unique dimension-value combinations allowed per day for this metric.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:<mode>` (e.g. `12345678:DEFAULT` or `12345678:PER_METRIC`).

## Import

Cardinality management resources can be imported using the composite ID format `<account_id>:<mode>`.

For a **DEFAULT** override:

```bash
$ terraform import newrelic_cardinality_management.account_default 12345678:DEFAULT
```

For a **PER_METRIC** override:

```bash
$ terraform import newrelic_cardinality_management.per_metric 12345678:PER_METRIC
```

-> **Note:** When importing a `PER_METRIC` resource, the `metric` blocks cannot be populated from the API (per-metric override values are not readable). After import, add the correct `metric` blocks to your configuration and run `terraform apply` to re-apply the desired limits.
