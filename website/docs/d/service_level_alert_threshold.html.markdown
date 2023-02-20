---
layout: "newrelic"
page_title: "New Relic: newrelic_service_leel_alert_threshold"
sidebar_current: "docs-newrelic-datasource-service-level-alert-threshold"
description: |-
  Calculates alert thresholds.
---

# Data Source: newrelic\_service\_level\_alert\_threshold

Use this data source to calculate the alert threshold of your Service
Level Objective.

## Example Usage

```hcl
data "newrelic_service_level_alert_threshold" "fast" {
    slo_target                   = 99.9
    slo_period                   = 28
    tolerated_budget_consumption = 2
    evaluation_period            = 1
}

resource "newrelic_alert_condition" "foo" {
  policy_id = 123456

  name            = "foo"
  type            = "apm_app_metric"
  entities        = 56789
  metric          = "apdex"
  runbook_url     = "https://www.example.com"
  condition_scope = "application"

  term {
    duration      = 60
    operator      = "below"
    priority      = "critical"
    threshold     = data.newrelic_service_level_alert_threshold.fast.alert_threshold
    time_function = "all"
  }
}
```


## Argument Reference

The following arguments are supported:

  * `slo_target` - (Required) The target of the objective, valid values between `0` and `100`.
  * `slo_period` - (Required) Time window is the period of the objective in days. Valid values are `1`, `7` and `28`.
  * `tolerated_budget_consumption` - (Required) How much budget you tolerate to consume in the evaluation period, valid values between `0` and `100`.
  * `evaluation_period` - (Required) Aggregation window taken into consideration in hours.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alert_threshold` - (Computed) The computed threshold given the provided arguments.
