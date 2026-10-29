---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_macro"
sidebar_current: "docs-newrelic-nrql-macro"
description: |-
  Manage a New Relic NRQL Macro.
---

# Resource: newrelic\_nrql\_macro

Manages a [NRQL Macro](https://docs.newrelic.com/docs/nrql/nrql-macros/nrql-macros-overview/) — a reusable expression fragment that you reference in NRQL queries using the `$macroName` syntax.

Macros are scoped to your organization and available across every account within it. Any query, dashboard widget, or alert condition that references a macro will pick up changes to that macro's expression on its next execution.

-> **Note:** Macros are always created in your organization. The `organization_id` shown in the state is resolved automatically from your provider credentials and cannot be set in configuration.

## Example Usage

### Simple filter macro

```hcl
resource "newrelic_nrql_macro" "production_filter" {
  name       = "filterProductionTraffic"
  expression = "environment = 'production' AND appName NOT LIKE '%test%'"
}
```

Reference in a query:

```sql
SELECT count(*) FROM Transaction WHERE $filterProductionTraffic TIMESERIES
```

### Macro with positional parameters

Positional parameters use bare `$1`, `$2`, ... tokens (without quotes). The `parameter_count` must equal the number of distinct `$N` tokens in the expression.

```hcl
resource "newrelic_nrql_macro" "status_region_filter" {
  name            = "filterStatusRegion"
  expression      = "status = $1 AND region = $2"
  parameter_count = 2
}
```

### Using a macro in a dashboard widget

```hcl
resource "newrelic_nrql_macro" "prod_filter" {
  name       = "filterProduction"
  expression = "environment = 'production'"
}

resource "newrelic_one_dashboard" "example" {
  name = "Production Overview"

  page {
    name = "Traffic"

    widget_line {
      title  = "Requests (Production)"
      row    = 1
      column = 1
      width  = 12
      height = 3

      nrql_query {
        query = "SELECT count(*) FROM Transaction WHERE $filterProduction TIMESERIES"
      }
    }
  }
}
```

## Argument Reference

* `name` - (Required, Forces new resource) The macro name. Used to reference it as `$name` in NRQL queries. Cannot be changed after the resource is created.

* `expression` - (Required) The NRQL fragment this macro expands to. Must not be empty. If the expression contains positional parameters (`$1`, `$2`, ...), you must also set `parameter_count` to match the number of distinct tokens.

  -> **Warning:** Changes to `expression` take effect immediately. Every query, dashboard, and alert condition that references this macro will reflect the new expression on next execution.

* `parameter_count` - (Optional, Computed) The number of positional parameters the expression accepts. Required when the expression contains `$1`, `$2`, ... tokens; defaults to `0` for expressions without parameters. The count must equal the number of distinct `$N` tokens or the API will return an error.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `id` - The entity GUID of the macro.

* `organization_id` - The GUID of the organization that owns this macro.

## Import

Import using the macro's entity GUID:

```bash
$ terraform import newrelic_nrql_macro.production_filter <entity-guid>
```
