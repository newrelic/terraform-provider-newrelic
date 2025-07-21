---
layout: "newrelic"
page_title: "New Relic: newrelic_one_dashboard"
sidebar_current: "docs-newrelic-resource-one-dashboard"
description: |-
  Create and manage dashboards in New Relic One.
---

# Resource: newrelic\_one\_dashboard

-> **IMPORTANT!**
When configuring the `newrelic_one_dashboard` resource, it is important to understand that widgets should ideally be sorted by row and column order to maintain the stability and accuracy of your dashboard setup. If this specified order is not adhered to, it can lead to resource drift, which might result in discrepancies between the intended setup and the actual deployed dashboard.

## Example Usage: Create a New Relic One Dashboard

```hcl
resource "newrelic_one_dashboard" "exampledash" {
  name        = "New Relic Terraform Example"
  permissions = "public_read_only"

  page {
    name = "New Relic Terraform Example"

    widget_table {
      title  = "List of Transactions"
      row    = 1
      column = 4
      width  = 6
      height = 3

      refresh_rate = 60000 // data refreshes every 60 seconds

      nrql_query {
        query = "FROM Transaction SELECT *"
      }

      initial_sorting {
        direction = "desc"
        name      = "timestamp"
      }

      data_format {
        name = "duration"
        type = "decimal"
      }
    }

    widget_billboard {
      title  = "Requests per minute"
      row    = 1
      column = 1
      width  = 6
      height = 3

      refresh_rate = 60000 // 60 seconds

      data_format {
        name = "rate"
        type = "recent-relative"
      }
      
      nrql_query {
        query = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
    }

    widget_bar {
      title  = "Average transaction duration, by application"
      row    = 1
      column = 7
      width  = 6
      height = 3

      nrql_query {
        account_id = 12345
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }

      # Must be another dashboard GUID
      linked_entity_guids = ["abc123"]
    }

    widget_bar {
      title  = "Average transaction duration, by application"
      row    = 4
      column = 1
      width  = 6
      height = 3

      refresh_rate = 300000 // 5 minutes

      nrql_query {
        account_id = 12345
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }

      # Must be another dashboard GUID
      filter_current_dashboard = true

      # color customization
      colors {
        color = "#722727"
        series_overrides {
          color = "#722322"
          series_name = "Node"
        }
        series_overrides {
          color = "#236f70"
          series_name = "Java"
        }
      }
    }

    widget_line {
      title  = "Average transaction duration and the request per minute, by application"
      row    = 4
      column = 7
      width  = 6
      height = 3

      refresh_rate = 30000 // 30 seconds

      nrql_query {
        account_id = 12345
        query      = "FROM Transaction select max(duration) as 'max duration' where httpResponseCode = '504' timeseries since 5 minutes ago"
      }

      nrql_query {
        query = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
      legend_enabled = true
      ignore_time_range = false
      y_axis_left_zero = true
      y_axis_left_min = 0
      y_axis_left_max = 1
      
      tooltip {
        mode = "single"
      }
      
      y_axis_right {
        y_axis_right_zero   = true
        y_axis_right_min    = 0
        y_axis_right_max    = 300
        y_axis_right_series = ["A", "B"]
      }
      
      is_label_visible = true
      
      threshold {
        name     = "Duration Threshold"
        from     = 1 
        to       = 2
        severity = "critical"
      }

      threshold {
        name     = "Duration Threshold Two"
        from     = 2.1
        to       = 3.3
        severity = "warning"
      }
      
      units {
        unit = "ms"
        series_overrides {
          unit = "ms"
          series_name = "max duration"
        }
      }


    }

    widget_markdown {
      title  = "Dashboard Note"
      row    = 7
      column = 1
      width  = 12
      height = 3

      text = "### Helpful Links\n\n* [New Relic One](https://one.newrelic.com)\n* [Developer Portal](https://developer.newrelic.com)"
    }

    widget_line {
      title = "Overall CPU % Statistics"
      row = 1
      column = 5
      height = 3
      width = 4

      nrql_query {
        query = <<EOT
SELECT average(cpuSystemPercent), average(cpuUserPercent), average(cpuIdlePercent), average(cpuIOWaitPercent) FROM SystemSample  SINCE 1 hour ago TIMESERIES
EOT
      }
      facet_show_other_series = false
      legend_enabled = true
      ignore_time_range = false
      y_axis_left_zero = true
      y_axis_left_min = 0
      y_axis_left_max = 0
      null_values {
        null_value = "default"

        series_overrides {
          null_value = "remove"
          series_name = "Avg Cpu User Percent"
        }

        series_overrides {
          null_value = "zero"
          series_name = "Avg Cpu Idle Percent"
        }

        series_overrides {
          null_value = "default"
          series_name = "Avg Cpu IO Wait Percent"
        }

        series_overrides {
          null_value = "preserve"
          series_name = "Avg Cpu System Percent"
        }
      }

    }

  }

  variable {
      default_values     = ["value"]
      is_multi_selection = true
      item {
        title = "item"
        value = "ITEM"
      }
      name = "variable"
      nrql_query {
        account_ids = [12345]
        query       = "FROM Transaction SELECT average(duration) FACET appName"
      }
      replacement_strategy = "default"
      title                = "title"
      type                 = "nrql"
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
  * `variable` - (Optional) A nested block that describes a dashboard-local variable. See [Nested variable blocks](#nested-variable-blocks) below for details.

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
  * `facet_show_other_series` - (Optional) Enable or disable the Other group in visualisations. The other group is used if a facet on a query returns more than 2000 items for bar charts, pie charts, and tables. The other group aggregates the rest of the facets. Defaults to `false`
  * `y_axis_left_min`, `y_axis_left_max` - (Optional) Adjust the Y axis to display the data within certain values by setting a minimum and maximum value for the axis for line charts and area charts. If no customization option is selected, dashboards automatically displays the full Y axis from 0 to the top value plus a margin.
  * `legend_enabled` - (Optional) With this turned on, the legend will be displayed. Defaults to `true`.
  * `null_values` - (Optional) A nested block that describes a Null Values. See [Nested Null Values blocks](#nested-null-values-blocks) below for details.
  * `units` - (Optional) A nested block that describes units on your Y axis. See [Nested Units blocks](#nested-units-blocks) below for details.
  * `colors` - (Optional) A nested block that describes colors of your charts per series. See [Nested Colors blocks](#nested-colors-blocks) below for details.
  *  `refresh_rate` - (Optional) This attribute determines the frequency for data refresh specified in milliseconds. Accepted values are `auto` for default value, `0` for no refresh, `5000` for 5 seconds, `30000` for 30 seconds, `60000` for 60 seconds, `300000` for 5 minutes, `1800000` for 30 minutes, `3600000` for 60 minute, `10800000` for 3 hours, `43200000` for 12 hours and `86400000` for 24 hours.
  * `tooltip` - (Optional) A nested block that describes tooltip configuration for area, line, and stacked bar widgets. See [Nested tooltip blocks](#nested-tooltip-blocks) below for details. 

Each widget type supports an additional set of arguments:

  * `widget_area`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_bar`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_billboard`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `critical` - (Optional) Threshold above which the displayed value will be styled with a red color.
    * `warning` - (Optional) Threshold above which the displayed value will be styled with a yellow color.
    * `data_format` - (Optional) A nested block that describes data format. See [Nested data_format blocks](#nested-data_format-blocks) below for details.
  * `widget_bullet`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `limit` - (Required) Visualization limit for the widget.
  * `widget_funnel`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_json`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_heatmap`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_histogram`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_line`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `y_axis_left_zero` - (Optional) An attribute that specifies if the values on the graph to be rendered need to be fit to scale, or printed within the specified range from `y_axis_left_min` (or 0 if it is not defined) to `y_axis_left_max`. Use `y_axis_left_zero = true` with a combination of `y_axis_left_min` and `y_axis_left_max` to render values from 0 or the specified minimum to the maximum, and `y_axis_left_zero = false` to fit the graph to scale.
    * `y_axis_right` - (Optional) An attribute which helps specify the configuration of the Y-Axis displayed on the right side of the line widget. This is a nested block, which includes the following attributes:
      * `y_axis_right_zero` - (Optional) An attribute that specifies if the values on the graph to be rendered need to be fit to scale, or printed within the specified range from `y_axis_right_min` (or 0 if it is not defined) to `y_axis_right_max`. Use `y_axis_right_zero = true` with a combination of `y_axis_right_min` and `y_axis_right_max` to render values from 0 or the specified minimum to the maximum, and `y_axis_right_zero = false` to fit the graph to scale.
      * `y_axis_right_min`, `y_axis_right_max` - (Optional) Attributes which help specify a range of minimum and maximum values, which adjust the right Y axis to display the data within the specified minimum and maximum value for the axis. 
      * `y_axis_right_series` - (Optional) An attribute which takes a list of strings, specifying a selection of series' displayed in the line chart to be adjusted against the values of the right Y-axis.
    * `threshold` - (Optional) An attribute that helps specify multiple thresholds, each inclusive of a range of values between which the threshold would need to function, the name of the threshold and its severity. Multiple thresholds can be defined in a line widget. The `threshold` attribute requires specifying the following attributes in a nested block - 
      * `name` - The name of the threshold.
      * `from` - The value 'from' which the threshold would need to be applied.
      * `to` - The value until which the threshold would need to be applied.
      * `severity` - The severity of the threshold, which would affect the visual appearance of the threshold (such as its color) accordingly. The value of this attribute would need to be one of the following - `warning`, `severe`, `critical`, `success`, `unavailable` which correspond to the severity labels _Warning_, _Approaching critical_, _Critical_, _Good_, _Neutral_ in the dropdown that helps specify the severity of thresholds in line widgets in the UI, respectively.
    * `is_label_visible` - (Optional) A boolean value, which when true, sets the label to be visibly displayed within thresholds. In other words, if this attribute is set to true, the _label always visible_ toggle in the _Thresholds_ section in the settings of the widget is enabled.
  * `widget_markdown`:
    * `text` - (Required) The markdown source to be rendered in the widget.
  * `widget_stacked_bar`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_pie`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
  * `widget_log_table`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
  * `widget_table`
    * `nrql_query` - (Required) A nested block that describes a NRQL Query. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) below for details.
    * `linked_entity_guids`: (Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.
    * `filter_current_dashboard`: (Optional) Use this item to filter the current dashboard.
    * `threshold` - (Optional) An attribute that helps specify multiple thresholds, each inclusive of a range of values between which the threshold would need to function, the name of the threshold and its severity. Multiple thresholds can be defined in a table widget. The `threshold` attribute requires specifying the following attributes in a nested block - 
        * `column_name` - The name of the column in the table, to which the threshold would need to be applied.
        * `from` - The value 'from' which the threshold would need to be applied.
        * `to` - The value until which the threshold would need to be applied.
        * `severity` - The severity of the threshold, which would affect the visual appearance of the threshold (such as its color) accordingly. The value of this attribute would need to be one of the following - `warning`, `severe`, `critical`, `success`, `unavailable` which correspond to the severity labels _Warning_, _Approaching critical_, _Critical_, _Good_, _Neutral_ in the dropdown that helps specify the severity of thresholds in table widgets in the UI, respectively.
    * `initial_sorting` - (Optional) An attribute that describes the sorting mechanism for the table. This attribute requires specifying the following attributes in a nested block -
        * `name` - (Required) The name of column to be sorted. Examples of few valid values are `timestamp`, `appId`, `appName`, etc.
        * `direction` - (Required) Defines the sort order. Accepted values are `asc` for ascending or `desc` for descending.
    * `data_format` - (Optional) A nested block that describes data format. See [Nested data_format blocks](#nested-data_format-blocks) below for details.
        

### Nested `data_format` blocks

Nested `data_format` blocks help specify the format of data displayed by a widget per attribute in the data returned by the NRQL query rendering the widget; thereby defining how the data fetched is best interpreted. This is supported for **billboards** and **tables**, as these are the only widgets in dashboards which return single or multi-faceted data that may be formatted based on the type of data returned.
This attribute requires specifying the following attributes in a nested block -

  * `name` - (Required) This attribute mandates the specification of the column name to be formatted. It identifies which column the data format should be applied to.
  * `type` - (Required) This attribute sets the format category for your data. Accepted values include - 
    - `decimal` for numeric values
    - `date` for date/time values
    - `duration` for length of time
    - `recent-relative` for values referencing a relative point in time
    - `custom` to be used with date/time values, in order to select a specific format the date/time value would need to be rendered as
    - `humanized` to be used with numeric values, in order to enable Autoformat
  * `format` - (Optional) This attribute is provided when the `name` is that of a column comprising date/time values and the `type` attribute is set to `custom` defining the specific date format to be applied to your data.

      |     Accepted value  |        Format           |                     
      |---------------------|-------------------------| 
      | `%b %d, %Y`         | `MMM DD,YYYY`           | 
      | `%d/%m/%Y`          | `DD/MM/YYYY(EU)`        | 
      | `%x`                | `DD/MM/YYYY(USA)`       | 
      | `%I:%M%p`           | `12-hour format`        | 
      | `%H:%Mh`            | `24-hour format`        | 
      | `%H:%Mh UTC (%Z)`   | `24-hour with timezone` | 
      | `%Y-%m-%dT%X.%L%Z`  | `ISO with timezone`     | 
      | `%b %d, %Y, %X`     | `MMM DD, YYYY,hh:mm:ss` | 
      | `%X`                | `hh:mm:ss`              | 


  * `precision` - (Optional) This attribute is utilized when the `type` attribute is set to `decimal`, stipulating the precise number of digits after the decimal point for your data.

-> **IMPORTANT!**
  As specified in the description of arguments of `data_format` above, using certain arguments requires using a specific `type` with the arguments, on a case-to-case basis. Please see the examples below for more details on such argument combinations.
  
* The following example illustrates using `data_format{}` with values of type `duration`, `recent-relative` and `timestamp` with no additional arguments specified.
```hcl
 widget_table {
  title  = "List of Transactions"
  row    = 1
  column = 4
  width  = 6
  height = 3

  nrql_query {
    account_id = Account_ID
    query = "SELECT average(duration), max(duration), min(duration) FROM Transaction FACET name SINCE 1 day ago"
  }

  data_format {
    name = "Max duration"
    Type = "duration"
  }

  data_format {
    name = "Max duration"
    type = "recent-relative"
  }

  initial_sorting {
    direction = "desc"
    name      = "timestamp"
  }
}
```
* In order to add a `data_format` block for date/time values, the `type` would need to be set to `date`. However, if you would also like to specify a format of the date/time value (with the `format` argument), the type would need to be set to `custom`.
```hcl
  data_format {
    name = "timestamp"
    Type = "date"
  }

  data_format {
    name = "timestamp"
    Type = "custom"
    Format = "%Y-%m-%dT%X.%L%Z"
  }
```
* Similarly, in order to use `data_format{}` with numeric values, the `type` would be need to set to `decimal`. The `precision` of the value may also be specified with type `decimal`. However, in order to have "Autoformat" enabled on the numeric value, specify the type as `humanized`.
```hcl
  data_format {
    name = "count"
    type = "decimal"
    precision = 4
  }
  data_format {
    name = "count"
    type = "humanized"
  }
```

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
    account_id  = Another_Account_ID
    query       = "FROM Transaction SELECT average(duration) FACET appName"
  }

  nrql_query {
    query       = "FROM Transaction SELECT rate(count(*), 1 minute)"
  }
}
```

### Nested `variable` blocks

The following arguments are supported:

  * `default_values` - (Optional) A list of default values for this variable. To select **all** default values, the appropriate value to be used with this argument would be `["*"]`.
  * `is_multi_selection` - (Optional) Indicates whether this variable supports multiple selection or not. Only applies to variables of type `nrql` or `enum`.
  * `item` - (Optional) List of possible values for variables of type `enum`. See [Nested item blocks](#nested-item-blocks) below for details.
  * `name` - (Required) The variable identifier.
  * `nrql_query` - (Optional) Configuration for variables of type `nrql`. See [Nested nrql\_query blocks](#nested-nrql_query-blocks) for details.
  * `replacement_strategy` - (Optional) Indicates the strategy to apply when replacing a variable in a NRQL query. One of `default`, `identifier`, `number` or `string`.
  * `title` - (Optional) Human-friendly display string for this variable.
  * `type` - (Required) Specifies the data type of the variable and where its possible values may come from. One of `enum`, `nrql` or `string`
  * `options` - (Optional) Specifies additional options to be added to dashboard variables. Supports the following nested attribute(s) -
    * `ignore_time_range` - (Optional) An argument with a boolean value that is supported only by variables of `type` _nrql_ - when true, the time range specified in the query will override the time picker on dashboards and other pages.
    * `excluded` - (Optional) An argument with a boolean value. With this turned on, the query condition defined with the variable will not be included in the query. Defaults to `false`.
### Nested `item` blocks

The following arguments are supported:

  * `title` - (Optional) A human-friendly display string for this value.
  * `value` - (Required) A possible variable value

### Nested `Null Values` blocks

The following arguments are supported:

* `null_value` -  Choose an option in displaying null values. Accepted values are `default`, `remove`, `preserve`, or `zero`.
* `series_overrides` - (Optional) A Nested block which will take two string attributes `null_value` and `series_name`. This nested block is used to customize null values of individual.

### Nested `Units` blocks

The following arguments are supported:

* `unit` - (Optional) Choose a unit to customize the unit on your Y axis and in each of your series.
* `series_overrides` - (Optional) A Nested block which will take two string attributes `unit` and `series_name`. This nested block is used to customize null values of individual.

### Nested `Colors` blocks

The following arguments are supported:

* `color` - (Optional) Choose a color to customize the color of your charts per series in area, bar, line, pie, and stacked bar charts. Accepted values are RGB, HEX, or HSL code.
* `series_overrides` - (Optional) A Nested block which will take two string attributes `color` and `series_name`. This nested block is used to customize colors of individual.

### Nested `tooltip` blocks

The following arguments are supported:

* `mode` - (Required) The tooltip display mode. Valid values are:
  * `all` - Show tooltip for all data points.
  * `single` - Show tooltip for a single data point.
  * `hidden` - Hide tooltips completely.

## Additional Examples

### Use the New Relic CLI to convert an existing dashboard

You can use the New Relic CLI to convert an existing dashboard into HCL code for use in Terraform.

1. [Download and install the New Relic CLI](https://github.com/newrelic/newrelic-cli#installation--upgrades)
2. [Export the dashboard you want to add to Terraform from the UI](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/#dashboards). Copy the JSON from the UI and paste it into a `.json` file.
3. Convert the `.json` file to HCL using the CLI: `cat dashboard.json | newrelic utils terraform dashboard --label my_dashboard_resource`

If you encounter any issues converting your dashboard, [please create a ticket on the New Relic CLI Github repository](https://github.com/newrelic/newrelic-cli/issues/new/choose).

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
        account_id = First_Account_ID
        query      = "FROM Metric SELECT rate(count(apm.service.transaction.duration), 1 minute) as 'First Account Throughput' TIMESERIES"
      }
      nrql_query {
        account_id = Second_Account_ID
        query      = "FROM Metric SELECT rate(count(apm.service.transaction.duration), 1 minute) as 'Second Account Throughput' TIMESERIES"
      }
      y_axis_left_zero = false
    }
  }
}
```

## Import

New Relic dashboards can be imported using their GUID, e.g.

```bash
$ terraform import newrelic_one_dashboard.my_dashboard <dashboard GUID>
```