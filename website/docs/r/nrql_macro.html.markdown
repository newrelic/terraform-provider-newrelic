---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_macro"
sidebar_current: "docs-newrelic-nrql-macro"
description: |-
  Use this resource to create and manage a New Relic NRQL Macro.
---

# Resource: newrelic\_nrql\_macro

Use this resource to create and manage a [NRQL Macro](https://docs.newrelic.com/docs/nrql/nrql-macros/nrql-macros-overview/).

A NRQL Macro is a named expression fragment you can reuse across NRQL queries using the `$macroName` syntax. Macros belong to your organization and are available in every account within it.

-> **NOTE:** Changes to a macro's `expression` take effect immediately. Every dashboard widget, alert condition, and query that references this macro will use the new expression on its next execution.

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

Positional parameters use bare `$1`, `$2`, ... tokens (not quoted). `parameter_count` must equal the number of distinct `$N` tokens in the expression.

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

The following arguments are supported:

* `name` - (Required) A unique name for the macro within your organization. Used to reference it as `$name` in NRQL queries. **Changing this value forces recreation of the resource.**

* `expression` - (Required) The NRQL expression fragment this macro expands to. When using positional parameters (`$1`, `$2`, ...), `parameter_count` must also be set to the number of distinct tokens used.

* `parameter_count` - (Optional) The number of positional parameters (`$1`, `$2`, ...) the expression uses. Required when the expression contains positional parameters; defaults to `0` for expressions without them. Must equal the exact count of distinct `$N` tokens in the expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the macro.

* `organization_id` - The GUID of the organization that owns this macro, resolved automatically from your provider credentials.

## Import

NRQL Macros can be imported using the entity GUID:

```bash
$ terraform import newrelic_nrql_macro.production_filter <entity-guid>
```
