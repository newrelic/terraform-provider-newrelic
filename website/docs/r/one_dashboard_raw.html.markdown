---
layout: "newrelic"
page_title: "New Relic: newrelic_one_dashboard_raw"
sidebar_current: "docs-newrelic-resource-one-dashboard-raw"
description: |-
  Create and manage dashboards with custom visualizations and/or RawConfiguration in New Relic One.
---

# Resource: newrelic_one_dashboard_raw

## Example Usage: Create a New Relic One Dashboard with RawConfiguration

```hcl
resource "newrelic_one_dashboard_raw" "exampledash" {
  name = "New Relic Terraform Example"

    page {
    name = "Page Name"
    widget {
      title = "Custom widget"
      row = 1
      column = 1
      width = 1
      height = 1
      visualization_id = "viz.custom"
      configuration = <<EOT
      {
        "legend": {
          "enabled": false
        },
        "nrqlQueries": [
          {
            "accountId": ` + accountID + `,
            "query": "SELECT average(loadAverageOneMinute), average(loadAverageFiveMinute), average(loadAverageFifteenMinute) from SystemSample SINCE 60 minutes ago    TIMESERIES"
          }
        ],
        "yAxisLeft": {
          "max": 100,
          "min": 50,
          "zero": false
        }
      }
      EOT
    }
    widget {
      title = "Server CPU"
      row = 1
      column = 2
      width = 1
      height = 1
      visualization_id = "viz.testing"
      configuration = <<EOT
      {
        "nrqlQueries": [
          {
            "accountId": ` + accountID + `,
            "query": "SELECT average(cpuPercent) FROM SystemSample since 3 hours ago facet hostname limit 400"
          }
        ]
      }
      EOT
    }
    widget {
      title  = "Docker Server CPU"
      row    = 1
      column = 3
      height = 1
      width  = 1
      visualization_id = "viz.bar"
      configuration = jsonencode(
      {
        "facet": {
          "showOtherSeries": false
        },
        "nrqlQueries": [
          {
            "accountId": local.accountID,
            "query": "SELECT average(cpuPercent) FROM SystemSample since 3 hours ago facet hostname limit 400"
          }
        ]
      }
      )
      linked_entity_guids = ["MzI5ODAxNnxWSVp8REFTSEJPQVJEfDI2MTcxNDc"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The title of the dashboard.
- `page` - (Required) A nested block that describes a page. See [Nested page blocks](#nested-page-blocks) below for details.
- `account_id` - (Optional) Determines the New Relic account where the dashboard will be created. Defaults to the account associated with the API key used.
- `description` - (Optional) Brief text describing the dashboard.
- `permissions` - (Optional) Determines who can see the dashboard in an account. Valid values are `private`, `public_read_only`, or `public_read_write`. Defaults to `public_read_only`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `guid` - The unique entity identifier of the dashboard in New Relic.
- `permalink` - The URL for viewing the dashboard.

### Nested `page` blocks

A New Relic One Dashboard is made up of one or more Pages. Each page contains
various widgets for displaying data.

The following arguments are supported:

- `name` - (Required) The name of the page. **Note:** If there is only one page, this name will be the name of the Dashboard.
- `description` - (Optional) Brief text describing the page.
- `widget` - (Optional) A nested block that describes a widget. See [Nested widget blocks](#nested-widget-blocks) below for details.

In addition to all arguments above, the following attributes are exported:

- `guid` - The unique entity identifier of the dashboard page in New Relic.

### Nested `widget` blocks

Nested `widget` blocks support the following common arguments:

- `title` - (Required) A title for the widget.
- `row` - (Required) Row position of widget from top left, starting at `1`.
- `column` - (Required) Column position of widget from top left, starting at `1`.
- `width` - (Optional) Width of the widget. Valid values are `1` to `12` inclusive. Defaults to `4`.
- `height` - (Optional) Height of the widget. Valid values are `1` to `12` inclusive. Defaults to `3`.
- `visualization_id` - (Required) The visualization ID of the widget
- `configuration` - (Required) The configuration of the widget.
- `linked_entity_guids` - (Optional) Related entity GUIDs. 
