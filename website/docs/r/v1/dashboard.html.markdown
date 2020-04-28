---
layout: "newrelic"
page_title: "New Relic: newrelic_dashboard"
sidebar_current: "docs-newrelic-resource-dashboard"
description: |-
  Create and manage dashboards in New Relic.
---

# Resource: newrelic\_dashboard

-> **NOTE:** This page refers to version **1.x** of the New Relic Terraform provider. For the latest documentation, please view the [latest docs for newrelic_dashboard](/docs/providers/newrelic/r/dashboard.html).

Use this resource to create and manage New Relic dashboards.

## Example Usage: Create a New Relic Dashboard

```hcl
data "newrelic_application" "my_application" {
  name = "My Application"
}

resource "newrelic_dashboard" "exampledash" {
  title = "New Relic Terraform Example"

  filter {
    event_types = [
        "Transaction"
    ]
    attributes = [
        "appName",
        "name"
    ]
  }

  widget {
    title = "Requests per minute"
    visualization = "billboard"
    nrql = "SELECT rate(count(*), 1 minute) FROM Transaction"
    row = 1
    column = 1
  }

  widget {
    title = "Error rate"
    visualization = "gauge"
    nrql = "SELECT percentage(count(*), WHERE error IS True) FROM Transaction"
    threshold_red = 2.5
    row = 1
    column = 2
  }

  widget {
    title = "Average transaction duration, by application"
    visualization = "facet_bar_chart"
    nrql = "SELECT average(duration) FROM Transaction FACET appName"
    row = 1
    column = 3
  }

  widget {
    title = "Apdex, top 5 by host"
    duration = 1800000
    visualization = "metric_line_chart"
    entity_ids = [
      data.newrelic_application.my_application.id,
    ]
    metric {
        name = "Apdex"
        values = [ "score" ]
    }
    facet = "host"
    limit = 5
    row = 2
    column = 1
  }

  widget {
    title = "Requests per minute, by transaction"
    visualization = "facet_table"
    nrql = "SELECT rate(count(*), 1 minute) FROM Transaction FACET name"
    row = 2
    column = 2
  }

  widget {
    title = "Dashboard Note"
    visualization = "markdown"
    source = "### Helpful Links\n\n* [New Relic One](https://one.newrelic.com)\n* [Developer Portal](https://developer.newrelic.com)"
    row = 2
    column = 3
  }
}
```

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the dashboard.
  * `icon` - (Optional) The icon for the dashboard.  Valid values are `adjust`, `archive`, `bar-chart`, `bell`, `bolt`, `bug`, `bullhorn`, `bullseye`, `clock-o`, `cloud`, `cog`, `comments-o`, `crosshairs`, `dashboard`, `envelope`, `fire`, `flag`, `flask`, `globe`, `heart`, `leaf`, `legal`, `life-ring`, `line-chart`, `magic`, `mobile`, `money`, `none`, `paper-plane`, `pie-chart`, `puzzle-piece`, `road`, `rocket`, `shopping-cart`, `sitemap`, `sliders`, `tablet`, `thumbs-down`, `thumbs-up`, `trophy`, `usd`, `user`, and `users`.  Defaults to `bar-chart`.
  * `visibility` - (Optional) Determines who can see the dashboard in an account. Valid values are `all` or `owner`.  Defaults to `all`.
  * `editable` - (Optional) Determines who can edit the dashboard in an account. Valid values are `all`,  `editable_by_all`, `editable_by_owner`, or `read_only`.  Defaults to `editable_by_all`.
  * `grid_column_count` - (Optional) The number of columns to use when organizing and displaying widgets. New Relic One supports a 3 column grid and a 12 column grid. New Relic Insights supports a 3 column grid.
  * `widget` - (Optional) A nested block that describes a visualization.  Up to 300 `widget` blocks are allowed in a dashboard definition.  See [Nested widget blocks](#nested-`widget`-blocks) below for details.
  * `filter` - (Optional) A nested block that describes a dashboard filter.  Exactly one nested `filter` block is allowed. See [Nested filter block](#nested-`filter`-block) below for details.

## Attribute Refence

In addition to all arguments above, the following attributes are exported:

  * `dashboard_url` - The URL for viewing the dashboard.

### Nested `widget` blocks

All nested `widget` blocks support the following common arguments:

  * `title` - (Required) A title for the widget.
  * `visualization` - (Required) How the widget visualizes data.  Valid values are `billboard`, `gauge`, `billboard_comparison`, `facet_bar_chart`, `faceted_line_chart`, `facet_pie_chart`, `facet_table`, `faceted_area_chart`, `heatmap`, `attribute_sheet`, `single_event`, `histogram`, `funnel`, `raw_json`, `event_feed`, `event_table`, `uniques_list`, `line_chart`, `comparison_line_chart`, `markdown`, and `metric_line_chart`.
  * `row` - (Required) Row position of widget from top left, starting at `1`.
  * `column` - (Required) Column position of widget from top left, starting at `1`.
  * `width` - (Optional) Width of the widget.  Valid values are `1` to `3` inclusive.  Defaults to `1`.
  * `height` - (Optional) Height of the widget.  Valid values are `1` to `3` inclusive.  Defaults to `1`.
  * `notes` - (Optional) Description of the widget.

Each `visualization` type supports an additional set of arguments:

  * `billboard`, `billboard_comparison`:
    * `nrql` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.
    * `threshold_red` - (Optional) Threshold above which the displayed value will be styled with a red color.
    * `threshold_yellow` - (Optional) Threshold above which the displayed value will be styled with a yellow color.
  * `gauge`:
    * `nrql` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.
    * `threshold_red` - (Required) Threshold above which the displayed value will be styled with a red color.
    * `threshold_yellow` - (Optional) Threshold above which the displayed value will be styled with a yellow color.
  * `facet_bar_chart`, `facet_pie_chart`, `facet_table`, `faceted_area_chart`, `faceted_line_chart`, or `heatmap`:
    * `nrql` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.
    * `drilldown_dashboard_id` - (Optional) The ID of a dashboard to link to from the widget's facets.
  * `attribute_sheet`, `comparison_line_chart`, `event_feed`, `event_table`, `funnel`, `histogram`, `line_chart`, `raw_json`, `single_event`, or `uniques_list`:
    * `nrql` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.
  * `markdown`:
    * `source` - (Required) The markdown source to be rendered in the widget.
  * `metric_line_chart`:
    * `entity_ids` - (Required) A collection of entity ids to display data for.  These are typically application IDs.
    * `metric` - (Required) A nested block that describes a metric.  Nested `metric` blocks support the following arguments:
      * `name` - (Required) The metric name to display.
      * `values` - (Required) The metric values to display.
    * `duration` - (Required) The duration, in ms, of the time window represented in the chart.
    * `end_time` - (Optional) The end time of the time window represented in the chart in epoch time.  When not set, the time window will end at the current time.
    * `facet` - (Optional) Can be set to "host" to facet the metric data by host.
    * `limit` - (Optional) The limit of distinct data series to display.
  * `application_breakdown`:
    * `entity_ids` - (Required) A collection of entity IDs to display data. These are typically application IDs.


### Nested `filter` block

The optional filter block supports the following arguments:
  * `event_types` - (Optional) A list of event types to enable filtering for.
  * `attributes` - (Optional) A list of attributes belonging to the specified event types to enable filtering for.

## Import

New Relic dashboards can be imported using their ID, e.g.

```
$ terraform import newrelic_dashboard.my_dashboard 8675309
```

~> **NOTE:** Due to API restrictions, importing a dashboard resource will set the `grid_column_count` attribute to `3`. If your dashboard is a New Relic One dashboard _and_ uses a 12 column grid, you will need to make sure `grid_column_count` is set to `12` in your configuration, then run `terraform apply` after importing to sync remote state with Terraform state.
