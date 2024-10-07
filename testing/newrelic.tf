terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}


resource "newrelic_one_dashboard" "shashankDashboard" {
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
        account_id = 3806526
        query      = "FROM Transaction SELECT average(duration) FACET appName"
      }

    }

    widget_bar {
      title  = "Average transaction duration, by application"
      row    = 4
      column = 1
      width  = 6
      height = 3

      refresh_rate = 300000 // 5 minutes

      nrql_query {
        account_id = 3806526
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
        account_id = 3806526
        query      = "FROM Transaction select max(duration) as 'max duration' where httpResponseCode = '504' timeseries since 5 minutes ago"
      }

      nrql_query {
        query = "FROM Transaction SELECT rate(count(*), 1 minute)"
      }
      legend_enabled = true
      y_axis_left_zero = true
      y_axis_left_min = 0
      y_axis_left_max = 1

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
      account_ids = [3806526]
      query       = "FROM Metric select count(*) since 1 day ago where dataType = 'Log API'"
    }
    replacement_strategy = "default"
    title                = "title"
    type                 = "nrql"
#     options{
#       ignore_time_range = false
#     }
  }
}
