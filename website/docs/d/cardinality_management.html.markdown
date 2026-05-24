---
layout: "newrelic"
page_title: "New Relic: newrelic_cardinality_management"
sidebar_current: "docs-newrelic-datasource-cardinality-management"
description: |-
  Reads the current cardinality limits configured for a New Relic account.
---

# Data Source: newrelic\_cardinality\_management

Use this data source to read the current cardinality limits configured for a New Relic account. The returned list covers limits across all categories (e.g. ingest) and is useful for inspecting current values before creating overrides with the [`newrelic_cardinality_management`](/providers/newrelic/newrelic/latest/docs/resources/cardinality_management) resource.

## Example Usage

```hcl
data "newrelic_cardinality_management" "current" {}

output "all_limits" {
  value = data.newrelic_cardinality_management.current.limits
}
```

### Look up the default metric cardinality limit

```hcl
data "newrelic_cardinality_management" "current" {}

output "default_cardinality_limit" {
  value = [
    for l in data.newrelic_cardinality_management.current.limits :
    l.value if l.name == "Dimensional Metric per-metric cardinality ingested per day"
  ]
}
```

### Use alongside a cardinality management resource

```hcl
data "newrelic_cardinality_management" "current" {}

resource "newrelic_cardinality_management" "account_default" {
  mode              = "DEFAULT"
  cardinality_limit = 150000

  depends_on = [data.newrelic_cardinality_management.current]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID to query. Defaults to the account ID set in the provider configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The account ID used for the query (set automatically by the provider).
* `limits` - A list of limit objects for the account. Each object has the following attributes:
  * `name` - The unique name of the limit (e.g. `"Dimensional Metric per-metric cardinality ingested per day"`).
  * `value` - The current limit value.
  * `unit` - The unit for the limit value (e.g. `"COUNT"`).
  * `category` - The category of the limit (e.g. `"INGEST"`).
