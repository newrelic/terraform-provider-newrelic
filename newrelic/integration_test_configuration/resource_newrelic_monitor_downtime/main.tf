terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region     = var.NEW_RELIC_REGION
  account_id = var.NEW_RELIC_ACCOUNT_ID
  api_key    = var.NEW_RELIC_API_KEY
}

resource "newrelic_monitor_downtime" "sample_one_time_newrelic_monitor_downtime" {
  name          = var.name
  monitor_guids = var.monitor_guids
  mode          = var.mode
  start_time    = var.start_time
  end_time      = var.end_time
  time_zone     = var.time_zone
  dynamic "end_repeat" {
    for_each = var.include_end_repeat ? [1] : []
    content {
      on_date   = var.end_repeat_on_date != "" ? var.end_repeat_on_date : null
      on_repeat = var.end_repeat_on_repeat != -1 ? var.end_repeat_on_repeat : null
    }
  }
  maintenance_days = var.maintenance_days != [] ? var.maintenance_days : null
  dynamic "frequency" {
    for_each = var.include_frequency ? [1] : []
    content {
      days_of_month = length(var.frequency_days_of_month) != 0 ? var.frequency_days_of_month : null
      dynamic "days_of_week" {
        for_each = var.days_of_week_week_day != "" ? [1] : []
        content {
          ordinal_day_of_month = var.days_of_week_ordinal_day_of_month
          week_day             = var.days_of_week_week_day
        }
      }
    }
  }
}