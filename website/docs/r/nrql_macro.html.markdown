---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_macro"
sidebar_current: "docs-newrelic-nrql-macro"
description: |-
  Use this resource to create and manage a New Relic NRQL Macro.
---

# Resource: newrelic\_nrql\_macro

Use this resource to create and manage a New Relic NRQL Macro.

NRQL Macros are reusable NRQL expression fragments that can be referenced in queries using the `$macroName` syntax. They are scoped to an organization and can be used across dashboards, alert conditions, and other NRQL queries.

-> **NOTE:** NRQL Macros are scoped to your **organization**, not a specific account. The `organization_id` is resolved automatically from your provider credentials if not explicitly set.

-> **WARNING:** Updating the `expression` attribute of a NRQL macro takes effect **immediately** across all dashboards, queries, and alert conditions that reference this macro. Plan expression changes carefully to avoid unintended query behavior.

## Example Usage

```hcl
resource "newrelic_nrql_macro" "production_filter" {
  name       = "filterProductionTraffic"
  expression = "environment = 'production' AND appName NOT LIKE '%test%'"
}
```

### With positional parameters

NRQL Macros support positional parameters using bare `$1`, `$2`, ... tokens (not quoted).
The `num_parameters` attribute must match the number of distinct `$N` tokens in the expression.

```hcl
resource "newrelic_nrql_macro" "status_region_filter" {
  name           = "filterStatusRegion"
  expression     = "status = $1 AND region = $2"
  num_parameters = 2
}
```

### With explicit organization ID

```hcl
resource "newrelic_nrql_macro" "apdex_calc" {
  name            = "calculateApdex"
  expression      = "average(duration) / 1000 as 'avgDurationSeconds'"
  organization_id = "MXxPUkd8MjM0NTY3OA"
}
```

### Using a macro in a dashboard widget

```hcl
resource "newrelic_nrql_macro" "production_filter" {
  name       = "filterProductionTraffic"
  expression = "environment = 'production' AND appName NOT LIKE '%test%'"
}

resource "newrelic_one_dashboard" "example" {
  name = "Example Dashboard"

  page {
    name = "Overview"

    widget_line {
      title  = "Response Time (Production)"
      row    = 1
      column = 1
      width  = 12
      height = 3

      nrql_query {
        query = "SELECT average(duration) FROM Transaction WHERE $filterProductionTraffic TIMESERIES"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required, Forces new resource) The unique name of the NRQL macro. Used to reference it as `$name` in queries. **This field is immutable after creation.** Changing this value will destroy and recreate the resource.
*   `expression` - (Required) The NRQL expression fragment that this macro expands to when referenced in a query. **Updating this value will immediately affect all queries, dashboards, and alert conditions that reference this macro.**
*   `organization_id` - (Optional, Computed, Forces new resource) The organization GUID to which this macro belongs. If omitted, the organization ID is resolved from the provider credentials. Changing this value will destroy and recreate the resource.
*   `num_parameters` - (Optional, Computed) The number of positional parameters (`$1`, `$2`, ...) used in the expression. Must match the exact count of distinct bare `$N` tokens in the expression (e.g. `expression = "status = $1 AND region = $2"` requires `num_parameters = 2`). Defaults to `0` for expressions without parameters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   `id` - The entity GUID of the NRQL macro.

## Import

NRQL Macros can be imported using the entity GUID. For example:

```bash
$ terraform import newrelic_nrql_macro.production_filter <entity-guid>
```
