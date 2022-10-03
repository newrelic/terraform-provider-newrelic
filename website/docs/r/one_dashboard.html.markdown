---
layout: "newrelic"
page_title: "New Relic: newrelic_one_dashboard"
sidebar_current: "docs-newrelic-resource-one-dashboard"
description: |-
  Create and manage dashboards in New Relic One.
---

# Resource: newrelic\_one\_dashboard

## Example Usage: Create a New Relic One Dashboard

```hcl
resource "newrelic_one_dashboard" "exampledash" {
  name = "New Relic Terraform Example"
  permissions = "public_read_only"

  page {
    name = "New Relic Terraform Example"

    widget_billboard {
      title = "Requests per minute"
      row = 1
      column = 1
      width = 6
      height = 3

      nrql_query {
        query       = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
    }

    widget_bar {
      title = "Average transaction duration, by application"
      row = 1
      column = 7
      width = 6
      height = 3

      nrql_query {
        account_id  = <Another Account ID>
        query       = "FROM Transaction SELECT average(duration) FACET appName"
      }

      # Must be another dashboard GUID
      linked_entity_guids = ["abc123"]
    }

    widget_bar {
      title = "Average transaction duration, by application"
      row = 4
      column = 1
      width = 6
      height = 3

      nrql_query {
        account_id  = <Another Account ID>
        query       = "FROM Transaction SELECT average(duration) FACET appName"
      }

      # Must be another dashboard GUID
      filter_current_dashboard = true
    }

    widget_line {
      title = "Average transaction duration and the request per minute, by application"
      row = 4
      column = 7
      width = 6
      height = 3

      nrql_query {
        account_id  = <Another Account ID>
        query       = "FROM Transaction SELECT average(duration) FACET appName"
      }

      nrql_query {
        query       = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
    }

    widget_markdown {
      title = "Dashboard Note"
      row    = 7
      column = 1
      width = 12
      height = 3

      text = "### Helpful Links\n\n* [New Relic One](https://one.newrelic.com)\n* [Developer Portal](https://developer.newrelic.com)"
    }
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The title of the dashboard.
  * `page` - (Required) A nested block that describes a page. See [Nested page blocks](#nested-page-blocks) below for details.
  * `account_id` - (Optional) Determines the New Relic account where the dashboard will be created. Defaults to the account associated with the API key used.
  * `description` - (Optional) Brief text describing the dashboard.
  * `permissions` - (Optional) Determines who can see the dashboard in an account. Valid values are `private`, `public_read_only`, or `public_read_write`.  Defaults to `public_read_only`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

  * `guid` - The unique entity identifier of the dashboard in New Relic.
  * `permalink` - The URL for viewing the dashboard.

### Nested `page` blocks

A New Relic One Dashboard is made up of one or more Pages. Each page contains
various widgets for displaying data.

The following arguments are supported:

  * `name` - (Required) The name of the page. **Note:** If there is only one page, this name will be the name of the Dashboard.
  * `description` - (Optional) Brief text describing the page.
  * `widget_area` - (Optional) A nested block that describes an Area widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_bar` - (Optional) A nested block that describes a Bar widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_billboard` - (Optional) A nested block that describes a Billboard widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_bullet` - (Optional) A nested block that describes a Bullet widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_funnel` - (Optional) A nested block that describes a Funnel widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_json` - (Optional) A nested block that describes a JSON widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_heatmap` - (Optional) A nested block that describes a Heatmap widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_histogram` - (Optional) A nested block that describes a Histogram widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_line` - (Optional) A nested block that describes a Line widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_markdown` - (Optional) A nested block that describes a Markdown widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_stacked_bar` - (Optional) A nested block that describes a Stacked Bar widget. See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_pie` - (Optional) A nested block that describes a Pie widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_log_table` - (Optional) A nested block that describes a Log Table widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_table` - (Optional) A nested block that describes a Table widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.


In addition to all arguments above, the following attributes are exported:

  * `guid` - The unique entity identifier of the dashboard page in New Relic.

### Nested `widget` blocks

All nested `widget` blocks support the following common arguments:

  * `title` - (Required) A title for the widget.
  * `row` - (Required) Row position of widget from top left, starting at `1`.
  * `column` - (Required) Column position of widget from top left, starting at `1`.
  * `width` - (Optional) Width of the widget.  Valid values are `1` to `12` inclusive.  Defaults to `4`.
  * `height` - (Optional) Height of the widget.  Valid values are `1` to `12` inclusive.  Defaults to `3`.
  * `ignore_time_range` - (Optional) With this turned on, the time range in this query will override the time picker on dashboards and other pages. Defaults to `false`.

Each widget type supports an additional set of arguments:

  * `widget_area`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_bar`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_billboard`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `critical` - (Optional) Threshold above which the displayed value will be styled with a red color.
    * `warning` - (Optional) Threshold above which the displayed value will be styled with a yellow color.
  * `widget_bullet`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `limit` - (Required) Visualization limit for the widget.
  * `widget_funnel`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_json`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_heatmap`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_histogram`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_line`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_markdown`:
    * `text` - (Required) The markdown source to be rendered in the widget.
  * `widget_stacked_bar`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_pie`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_log_table`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_table`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.

### Nested `nrql_query` blocks

Nested `nrql_query` blocks allow you to make one or more NRQL queries within a widget, against a specified account.

The following arguments are supported:

  * `account_id` - (Optional) The New Relic account ID to issue the query against. Defaults to the Account ID where the dashboard was created. When using an account ID you don't have permissions for the widget will be replaced with a widget showing the data is inaccessible. Terraform will not throw an error, so this widget will only be visible in the UI.
  * `query` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.

```hcl
widget_line {
  title = "Average transaction duration and the request per minute, by application"
  row = 4
  column = 7
  width = 6
  height = 3

  nrql_query {
    account_id  = <Another Account ID>
    query       = "FROM Transaction SELECT average(duration) FACET appName"
  }

  nrql_query {
    query       = "FROM Transaction SELECT rate(count(*), 1 minute)"
  }
}
```

## Additional Examples

### Create a two page dashboard

The example below shows how you can display data for an application from a primary account and an application from a subaccount. In order to create cross-account widgets, you must use an API key from a user with admin permissions in the primary account. Please see the [`widget` attribute documentation](#cross-account-widget-help) for more details.

```hcl
resource "newrelic_one_dashboard" "multi_page_dashboard" {
  name = "My Multi-page dashboard"

  # Only I can see this dashboard
  permissions = "private"

  page {
    name = "My Multi-page dashboard"

    widget_bar {
      title = "foo"
      row    = 1
      column = 1

      nrql_query {
        query      = "FROM Transaction SELECT count(*) FACET name"
      }

      # Must be another dashboard GUID
      linked_entity_guids = ["abc123"]
    }
  }

  page {
    name = "Multi-query Page"

    widget_line {
      title = "Comparing throughput cross-account"
      row    = 1
      column = 1
      width  = 12
      nrql_query {
        account_id = <First Account ID>
        query      = "FROM Metric SELECT rate(count(apm.service.transaction.duration), 1 minute) as 'First Account Throughput' TIMESERIES"
      }
      nrql_query {
        account_id = <Second Account ID>
        query      = "FROM Metric SELECT rate(count(apm.service.transaction.duration), 1 minute) as 'Second Account Throughput' TIMESERIES"
      }
    }
  }
}
```

## Import

New Relic dashboards can be imported using their GUID, e.g.

```
$ terraform import newrelic_one_dashboard.my_dashboard <Dashboard GUID>
```

In addition you can use the [New Relic CLI](https://github.com/newrelic/newrelic-cli#readme) to convert existing dashboards to HCL. [Copy your dashboards as JSON using the UI](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/), save it as a file (for example `terraform.json`), and use the following command to convert it to HCL: `cat terraform.json | newrelic utils terraform dashboard --label my_dashboard_resource`.
