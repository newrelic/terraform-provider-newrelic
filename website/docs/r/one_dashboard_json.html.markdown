---
layout: "newrelic"
page_title: "New Relic: newrelic_one_dashboard_json"
sidebar_current: "docs-newrelic-resource-one-dashboard-json"
description: |-
  Create and manage dashboards from a JSON file.
---

# Resource: newrelic_one_dashboard_json

## Example Usage: Create a New Relic One Dashboard from a JSON file

```hcl
resource "newrelic_one_dashboard_json" "foo" {
   json = file("dashboard.json")
}
```

## Argument Reference

The following arguments are supported:

- `json` - (Required) The JSON export of a dashboard. [The JSON can be exported from the UI](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/#dashboards)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `guid` - The unique entity identifier of the dashboard in New Relic.
- `permalink` - The URL for viewing the dashboard.
- `updated_at` - The date and time when the dashboard was last updated.
