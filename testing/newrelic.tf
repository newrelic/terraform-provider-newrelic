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

resource "newrelic_one_dashboard" "exampledash" {
  name        = "New Relic Terraform Example Random abcdef"
  permissions = "public_read_only"

  page {
    name = "New Relic Terraform Example"

    widget_line {
      title  = "Average transaction duration and the request per minute, by application"
      row    = 4
      column = 7
      width  = 6
      height = 3

      nrql_query {
        account_id = 3806526
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

      units {
        unit = "ms"
        series_overrides {
          unit = "ms"
          series_name = "max duration"
        }
      }


    }

    widget_markdown {
      title  = "Dashboard Note xyz"
      row    = 7
      column = 1
      width  = 12
      height = 3

      text = "### Helpful Links\n\n* [New Relic One](https://one.newrelic.com)\n* [Developer Portal](https://developer.newrelic.com)"
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
      query       = "FROM Transaction SELECT average(duration) FACET appName"
    }
    replacement_strategy = "default"
    title                = "title"
    type                 = "nrql"
#    options {
#      ignore_time_range = true
#    }
  }
}