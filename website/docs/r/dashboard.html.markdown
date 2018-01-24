---
layout: "newrelic"
page_title: "New Relic: newrelic_dashboard"
sidebar_current: "docs-newrelic-resource-dashboard"
description: |-
  Create and manage dashboards in New Relic.
---

# newrelic\_dashboard

## Example Usage

```hcl
resource "newrelic_dashboard" "exampledash" {
  title = "New Relic Terraform Example"

  widget {
    title         = "Average Transaction Duration"
    row           = 1
    column        = 1
    width         = 2
    visualization = "faceted_line_chart"
    nrql          = "SELECT AVERAGE(duration) from Transaction FACET appName TIMESERIES auto"
  }

  widget {
    title         = "Page Views"
    row           = 1
    column        = 3
    visualization = "billboard"
    nrql          = "SELECT count(*) FROM PageView SINCE 1 week ago"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the dashboard.
  * `icon` - (Optional) The icon for the dashboard.  Defaults to `bar-chart`.
  * `visibility` - (Optional) Who can see the dashboard in an account. Must be `owner` or `all`. Defaults to `all`.
  * `widget` - (Optional) A widget that describes a visualization. See [Widgets](#widgets) below for details.
  * `editable` - (Optional) Who can edit the dashboard in an account. Must be `read_only`, `editable_by_owner`, `editable_by_all`, or `all`. Defaults to `editable_by_all`.

## Widgets

The `widget` mapping supports the following arguments:

  * `title` - (Required) A title for the widget.
  * `visualization` - (Required) How the widget visualizes data.
  * `row` - (Optional) Row position of widget from top left, starting at `1`. Defaults to `1`.
  * `column` - (Optional) Column position of widget from top left, starting at `1`.  Defaults to `1`.
  * `notes` - (Optional) Description of the widget.
  * `nrql` - (Optional) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the dashboard.
