---
layout: "newrelic"
page_title: "New Relic: newrelic_service_level_alert_helper"
sidebar_current: "docs-newrelic-datasource-service-level-alert-helper"
description: |-
  Helper to set up alerts on Service Levels.
---

# Data Source: newrelic\_service\_level\_alert\_helper

Use this data source to obtain the necessary fields to set up alerts on your service levels. It can be used for a `custom` alert_type in order to set up an alert with custom tolerated budget consumption and custom evaluation period or for recommended ones like `fast_burn`. For more information check [the documentation](https://docs.newrelic.com/docs/service-level-management/alerts-slm/).

## Example Usage

Firstly set up your service level objective, we recommend to use local variables for the `target` and `time_window.rolling.count`, as they are also necessary for the helper.

```hcl
locals {
  foo_target = 99.9
  foo_period = 28
}

resource "newrelic_service_level" "foo" {
    guid = "MXxBUE18QVBQTElDQVRJT058MQ"
    name = "Latency"
    description = "Proportion of requests that are served faster than a threshold."

    events {
        account_id = 12345678
        valid_events {
            from = "Transaction"
            where = "appName = 'Example application' AND (transactionType='Web')"
        }
        bad_events {
            from = "Transaction"
            where = "appName = 'Example application' AND (transactionType= 'Web') AND duration > 0.1"
        }
    }

    objective {
        target = local.foo_target
        time_window {
            rolling {
                count = local.foo_period
                unit = "DAY"
            }
        }
    }
}
```
Then use the helper to obtain the necessary fields to set up an alert on that Service Level.

```hcl

data "newrelic_service_level_alert_helper" "foo_custom" {
    alert_type = "custom"
    sli_guid = newrelic_service_level.foo.sli_guid
    slo_target = local.foo_target
    slo_period = local.foo_period
    custom_tolerated_budget_consumption = 5
    custom_evaluation_period = 90
    is_bad_events = true
}

resource "newrelic_nrql_alert_condition" "your_condition" {
  account_id = 12345678
  policy_id = 67890
  type = "static"
  name = "Successs (custom)"
  enabled = true
  violation_time_limit_seconds = 259200

  nrql {
    query = data.newrelic_service_level_alert_helper.foo_custom.nrql
  }

  critical {
    operator = "above_or_equals"
    threshold = data.newrelic_service_level_alert_helper.foo_custom.threshold
    threshold_duration = data.newrelic_service_level_alert_helper.foo_custom.evaluation_period
    threshold_occurrences = "at_least_once"
  }
  fill_option = "none"
  aggregation_window = 3600
  aggregation_method = "event_flow"
  aggregation_delay = 120
  slide_by = 60
}
```


## Argument Reference

The following arguments are supported:

  * `alert_type` - (Required) The type of alert we want to set. Valid values are:
    * `custom` - Tolerated budget consumption and evaluation period have to be specified.
    * `fast_burn` - Tolerated budget consumption is 2% and evaluation period is 60min.
  * `sli_guid` - (Required) The guid of the sli we want to set the alert on.
  * `slo_target` - (Required) The target of the Service Level Objective, valid values between `0` and `100`.
  * `slo_period` - (Required) The time window of the Service Level Objective in days. Valid values are `1`, `7` and `28`.
  * `custom_tolerated_budget_consumption` - (Optional) How much budget you tolerate to consume during the custom evaluation period, valid values between `0` and `100`. Mandatory if `alert_type` is `custom`.
  * `custom_evaluation_period` - (Optional) Aggregation window taken into consideration in minutes. Mandatory if `alert_type` is `custom`.
  * `is_bad_events` - (Optional) If the SLI is defined using bad events. Defaults to `false`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `threshold` - (Computed) The computed threshold given the provided arguments.
  * `tolerated_budget_consumption` - (Computed) For non `custom` alert_type, this is the recommended for that type of alert. For `custom` alert_type it has the same value as `custom_tolerated_budget_consumption`.
  * `evaluation_period` - (Computed) For non `custom` alert_type, this is the recommended for that type of alert. For `custom` alert_type it has the same value as `custom_evaluation_period`.
  * `nrql` - (Computed) The nrql query for the selected type of alert.
