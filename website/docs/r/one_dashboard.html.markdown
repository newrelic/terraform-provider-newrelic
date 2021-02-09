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

  page {
    name = "New Relic Terraform Example"

    widget_billboard {
      title = "Requests per minute"
      row = 1
      column = 1

      nrql_query {
        query       = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
    }

    widget_bar {
      title = "Average transaction duration, by application"
      row = 1
      column = 5

      nrql_query {
        account_id = <Another Account ID>
        query       = "FROM Transaction SELECT average(duration) FACET appName"
      }

      # Must be another dashboard GUID
      linked_entity_guids = ["abc123"]
    }

    widget_markdown {
      title = "Dashboard Note"
      row    = 1
      column = 9

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
  * `widget_line` - (Optional) A nested block that describes a Line widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_markdown` - (Optional) A nested block that describes a Markdown widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
  * `widget_pie` - (Optional) A nested block that describes a Pie widget.  See [Nested widget blocks](#nested-widget-blocks) below for details.
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

Each widget type supports an additional set of arguments:

  * `widget_bar`, `widget_line`, `widget_pie`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
  * `widget_table`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
  * `widget_billboard`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql-query-blocks) below for details.
    * `critical` - (Optional) Threshold above which the displayed value will be styled with a red color.
    * `warning` - (Optional) Threshold above which the displayed value will be styled with a yellow color.
  * `widget_markdown`:
    * `text` - (Required) The markdown source to be rendered in the widget.

### Nested `nrql_query` blocks

Nested `nrql_query` blocks allow you to make one or more NRQL queries within a widget, against a specified account.

The following arguments are supported:

  * `account_id` - (Optional) The New Relic account ID to issue the query against. Defaults to the Account ID where the dashboard was created.
  * `query` - (Required) Valid NRQL query string. See [Writing NRQL Queries](https://docs.newrelic.com/docs/insights/nrql-new-relic-query-language/using-nrql/introduction-nrql) for help.

## Additional Examples

###  Create a two page dashboard

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
