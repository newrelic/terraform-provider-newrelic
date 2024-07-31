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

resource "newrelic_one_dashboard" "threashold-to-from-float-type-test-tf1" {
  name        = "random-threshold-to-from-float-type-test-tf1"
  permissions = "public_read_only"

  page {
    name = "threashold-to-from-float-type-test-tf1"

    widget_table {
      title  = "Table Widget 1"
      row    = 1
      column = 5
      width  = 4
      height = 4

      nrql_query {
        query = "FROM DistributedTraceSummary SELECT *"
      }

      threshold {
        to =  1.1
        from= 6768678.787
        column_name = "Table Column 11"
        severity    = "severe"
      }
    }
    

     widget_line {
      title  = "Average transaction duration and the request per minute, by application"
      row    = 4
      column = 7
      width  = 6
      height = 3

      nrql_query {
        query = "FROM Transaction select max(duration) as 'max duration' where httpResponseCode = '504' timeseries since 5 minutes ago"
      }
      legend_enabled = true
      ignore_time_range = false
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
        from     = 1.33
        to       = 2.56
        severity = "critical"
      }


      units {
        unit = "ms"
        series_overrides {
          unit = "ms"
          series_name = "max duration"
        }
      }
    }
  }
}
