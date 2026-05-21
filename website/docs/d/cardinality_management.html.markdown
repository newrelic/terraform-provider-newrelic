---
layout: "newrelic"
page_title: "New Relic: newrelic_cardinality_management"
sidebar_current: "docs-newrelic-datasource-cardinality-management"
description: |-
  Reads all current cardinality limits for a New Relic account.
---

# Data Source: newrelic\_cardinality\_management

Use this data source to read the current cardinality limits configured for a New Relic account. The returned list covers all limit categories (e.g. ingest, API) and can be used to inspect current values before creating overrides with [`newrelic_cardinality_management`](/providers/newrelic/newrelic/latest/docs/resources/cardinality_management).

## Example Usage

```hcl
data "newrelic_cardinality_management" "current" {
  account_id = 12345678
}

output "all_limits" {
  value = data.newrelic_cardinality_management.current.limits
}
```

### Look Up the Default Metric Cardinality Limit

```hcl
data "newrelic_cardinality_management" "current" {
  account_id = 12345678
}

output "default_cardinality_limit" {
  value = [
    for l in data.newrelic_cardinality_management.current.limits :
    l.value if l.name == "Dimensional Metric per-metric cardinality ingested per day"
  ]
}
```

### Use Alongside a Cardinality Limit Override

```hcl
data "newrelic_cardinality_management" "current" {
  account_id = 12345678
}

resource "newrelic_cardinality_management" "default" {
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
* `limits` - A list of limit objects. Each object has the following attributes:
  * `name` - The unique name of the limit (e.g. `"Dimensional Metric per-metric cardinality ingested per day"`).
  * `value` - The current limit value.
  * `unit` - The unit for the limit value (e.g. `"COUNT"`).
  * `category` - The category of the limit (e.g. `"INGEST"`).
