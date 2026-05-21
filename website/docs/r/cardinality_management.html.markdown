---
layout: "newrelic"
page_title: "New Relic: newrelic_cardinality_management"
sidebar_current: "docs-newrelic-resource-cardinality-management"
description: |-
  Manage cardinality limit overrides for a New Relic account.
---

# Resource: newrelic\_cardinality\_management

Use this resource to manage cardinality limit overrides for a New Relic account.

Dimensional metrics in New Relic are subject to a per-metric cardinality limit — the maximum number of unique attribute-value combinations a single metric name may produce per day. This resource lets you override that limit account-wide (for all metrics at once) or for specific metric names individually.

Two modes are supported, controlled by the required `mode` argument.

---

## DEFAULT Mode

In `DEFAULT` mode, the resource sets an account-wide cardinality limit that applies to every metric in the account unless a per-metric override is in place. Terraform reads the live value from the platform on each plan and will flag any drift. Running `terraform destroy` on this resource resets the account-wide limit back to the New Relic platform default of **100,000**.

### Example

```hcl
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}
```

-> **Note:** Destroying a `DEFAULT` mode resource resets the account-wide cardinality limit to the New Relic platform default of 100,000. Make sure this is intentional before running `terraform destroy`.

---

## PER\_METRIC Mode

In `PER_METRIC` mode, the resource overrides the cardinality limit for one or more named metrics. Each `metric` block specifies a metric name and its target limit. You can manage as many metrics as needed within a single resource by adding multiple `metric` blocks.

The updated limit takes effect as the platform receives and processes metric data. Terraform maintains the last applied values in state — so `terraform plan` will not detect changes made to these limits outside of Terraform. Running `terraform destroy` resets all managed metrics back to the New Relic platform default of **100,000**.

### Example — single metric

```hcl
resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"
  metric {
    name              = "otelcol_nrreceiver_incoming_request_proxy"
    cardinality_limit = 200000
  }
}
```

### Example — multiple metrics

```hcl
resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"
  metric {
    name              = "otelcol_nrreceiver_incoming_request_proxy"
    cardinality_limit = 200000
  }
  metric {
    name              = "custom.app.requests"
    cardinality_limit = 150000
  }
}
```

-> **Note:** Destroying a `PER_METRIC` resource resets each managed metric's limit back to the New Relic platform default of 100,000.

---

## Using Both Modes Together

```hcl
resource "newrelic_cardinality_management" "default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000
}

resource "newrelic_cardinality_management" "per_metric" {
  mode = "PER_METRIC"
  metric {
    name              = "otelcol_nrreceiver_incoming_request_proxy"
    cardinality_limit = 200000
  }
}
```

---

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID. Defaults to the account ID set in the provider configuration.
* `mode` - (Required) The override mode. Must be `"DEFAULT"` or `"PER_METRIC"`. Forces re-creation when changed.
  * `DEFAULT` — sets the account-wide default limit for all metrics. `cardinality_limit` is required; `metric` blocks must not be used.
  * `PER_METRIC` — overrides the limit for specific named metrics. At least one `metric` block is required; `cardinality_limit` at the top level must not be set.
* `cardinality_limit` - (Optional) The account-wide cardinality limit. Required when `mode` is `"DEFAULT"`.
* `metric` - (Optional) One or more per-metric limit overrides. Required when `mode` is `"PER_METRIC"`. Each block supports:
  * `name` - (Required) The name of the metric to override.
  * `cardinality_limit` - (Required) The cardinality limit for this metric.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform resource ID in the format `<account_id>:DEFAULT` or `<account_id>:PER_METRIC`.

## Import

Cardinality management resources can be imported using the composite ID format `<account_id>:<mode>`.

For a **DEFAULT** override:

```bash
$ terraform import newrelic_cardinality_management.default 12345678:DEFAULT
```

For a **PER_METRIC** override:

```bash
$ terraform import newrelic_cardinality_management.per_metric 12345678:PER_METRIC
```

-> **Note:** When importing a `PER_METRIC` resource, `mode` is restored from the resource ID but the `metric` blocks are not — since per-metric override values cannot be read back from the platform API. Run `terraform apply` after import to re-establish the intended limits.
