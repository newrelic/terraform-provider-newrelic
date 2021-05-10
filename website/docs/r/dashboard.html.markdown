---
layout: "newrelic"
page_title: "New Relic: newrelic_dashboard (Deprecated)"
sidebar_current: "docs-newrelic-resource-dashboard"
description: |-
  Create and manage dashboards in New Relic. (Deprecated)
---

# Resource: newrelic\_dashboard (Deprecated)

Use this resource to create and manage New Relic dashboards.

-> **IMPORTANT!** The `newrelic_dashboard` resource is being deprecated, please move your dashboards to the new `newrelic_one_dashboard` resource. For more information check out [this Github issue](https://github.com/newrelic/terraform-provider-newrelic/issues/1297).

## Example Usage: Create a New Relic Dashboard

```hcl
data "newrelic_entity" "my_application" {
  name = "My Application"
  type = "APPLICATION"
  domain = "APM"
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
      data.newrelic_entity.my_application.application_id,
    ]
    metric {
        name = "Apdex"
        values = [ "score" ]
    }
    facet = "host"
    limit = 5
    order_by = "score"
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
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the dashboard.
  * `icon` - (Optional) The icon for the dashboard.  Valid values are `adjust`, `archive`, `bar-chart`, `bell`, `bolt`, `bug`, `bullhorn`, `bullseye`, `clock-o`, `cloud`, `cog`, `comments-o`, `crosshairs`, `dashboard`, `envelope`, `fire`, `flag`, `flask`, `globe`, `heart`, `leaf`, `legal`, `life-ring`, `line-chart`, `magic`, `mobile`, `money`, `none`, `paper-plane`, `pie-chart`, `puzzle-piece`, `road`, `rocket`, `shopping-cart`, `sitemap`, `sliders`, `tablet`, `thumbs-down`, `thumbs-up`, `trophy`, `usd`, `user`, and `users`.  Defaults to `bar-chart`.
  * `visibility` - (Optional) Determines who can see the dashboard in an account. Valid values are `all` or `owner`.  Defaults to `all`.
  * `editable` - (Optional) Determines who can edit the dashboard in an account. Valid values are `all`,  `editable_by_all`, `editable_by_owner`, or `read_only`.  Defaults to `editable_by_all`.
  * `grid_column_count` - (Optional) The number of columns to use when organizing and displaying widgets. New Relic One supports a 3 column grid and a 12 column grid. New Relic Insights supports a 3 column grid.
  * `filter` - (Optional) A nested block that describes a dashboard filter.  Exactly one nested `filter` block is allowed. See [Nested filter block](#nested-filter-block) below for details.
  * `widget` - (Optional) A nested block that describes a visualization.  Up to 300 `widget` blocks are allowed in a dashboard definition. See [Nested widget blocks](#nested-widget-blocks) below for details.

  <a name="widget-configuration-recommendation"></a>

  -> **Widget configuration recommendation** While the `newrelic_dashboard` resource attempts to avoid configuration drift where possible, there is still potential for drift to occur due to underlying API limitations, usually involving [cross-account widgets](#account_id). Also, even though the ordering of your widgets can be arbitrary, we recommend ordering your widgets in a consistent manner and maintaining that order if possible. An example could be ordering the `widget` blocks in the order they appear your New Relic Dashboard UI. The first `widget` block would be the widget displayed at the top left of your dashboard and the last `widget` block would be the widget at the bottom right of the dashboard.

## Attribute Reference

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
  * `account_id` - (Optional) The account ID to use when querying data. If `account_id` is omitted, the widget will use the account ID associated with the API key used in your provider configuration. You can also use `account_id` to configure cross-account widgets or simply to be explicit about which account the widget will be pulling data from.

<a name="cross-account-widget-help"></a>

-> **Configuring cross-account widgets** To configure a cross-account widget with an account different from the account associated with your API key, you must set the widget's `account_id` attribute to the account ID you wish to pull data from. Also note, the provider must be configured with an API Key that is scoped to a user with proper permissions to access and perform operations in other accounts that fall within or under the account associated with your API key. To facilitate cross-account widgets, we recommend [configuring the provider with a User API Key](../guides/provider_configuration.html#configuration-via-the-provider-block) from a user with **admin permissions** and access to the subaccount you would like to display data for in the widget.

~> **Note** Due to API limitations, cross-account widgets can cause configuration drift due to the API response omitting data for widgets that that pull data from New Relic accounts outside the primary scope of the API key being used. If you need to configure cross-account widgets and also want to bypass the configuration drift for widgets, you can use Terraform's [`ignore_changes`](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) using Terraforms `lifecycle` block. <br><br> ``` lifecycle { ignore_changes = [widget] }```

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
    * `limit` - (Optional) The limit of distinct data series to display.  Requires `order_by` to be set.
    * `order_by` - (Optional) Set the order of the results.  Required when using `limit`.
  * `application_breakdown`:
    * `entity_ids` - (Required) A collection of entity IDs to display data. These are typically application IDs.


### Nested `filter` block

The optional filter block supports the following arguments:
  * `event_types` - (Optional) A list of event types to enable filtering for.
  * `attributes` - (Optional) A list of attributes belonging to the specified event types to enable filtering for.

## Additional Examples

###  Create cross-account widgets in your dashboard.

The example below shows how you can display data for an application from a primary account and an application from a subaccount. In order to create cross-account widgets, you must use an API key from a user with admin permissions in the primary account. Please see the [`widget` attribute documentation](#cross-account-widget-help) for more details.

```hcl
# IMPORTANT!
# The User API Key must be from a user with admin permissions in the main account.
provider "newrelic" {
  api_key = "NRAK-*****"
  # ... additional configuration
}

# Fetch data for an application in your primary account to reference in the dashboard
data "newrelic_entity" "primary_account_application" {
  # Must be a unique name, otherwise use the `tags` attribute to get more specific if needed
  name   = "Main Account Application Name"
  type   = "APPLICATION"
  domain = "APM"
}

# Fetch data for an application in a subaccount to reference in the dashboard
data "newrelic_entity" "subaccount_application" {
  # Must be a unique name, otherwise use the `tags` attribute to get more specific if needed
  name     = "Subaccount Application Name"
  type     = "APPLICATION"
  domain   = "APM"
}

resource "newrelic_dashboard" "cross_account_widget_example" {
  title = "tf-test-cross-account-widget-dashboard"

  filter {
    event_types = [
      "Transaction"
    ]
    attributes = [
      "appName",
      "envName"
    ]
  }

  grid_column_count = 12

  # Omitting `account_id` will make this widget pull data from the primary account.
  widget {
    title         = "Apdex (primary account)"
    row           = 1
    column        = 1
    width         = 6
    height        = 3
    visualization = "metric_line_chart"
    duration      = 1800000

    metric {
      name   = "Apdex"
      values = ["score"]
    }

    entity_ids    = [
      data.newrelic_entity.primary_account_application.application_id
    ]
  }

  # Setting `account_id` to a subaccount ID will make this widget pull data from the subaccount.
  widget {
    account_id    = var.subaccount_id
    title         = "Apdex (subaccount)"
    row           = 1
    column        = 7
    width         = 6
    height        = 3
    visualization = "metric_line_chart"
    duration      = 1800000

    metric {
      name   = "Apdex"
      values = ["score"]
    }

    entity_ids    = [
      data.newrelic_entity.subaccount_application.application_id
    ]
  }
}
```

## Import

New Relic dashboards can be imported using their ID, e.g.

```
$ terraform import newrelic_dashboard.my_dashboard 8675309
```

~> **NOTE** Due to API restrictions, importing a dashboard resource will set the `grid_column_count` attribute to `3`. If your dashboard is a New Relic One dashboard _and_ uses a 12 column grid, you will need to make sure `grid_column_count` is set to `12` in your configuration, then run `terraform apply` after importing to sync remote state with Terraform state. Also note, cross-account widgets cannot be imported due to API restrictions.
