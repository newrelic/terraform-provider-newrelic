---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_condition"
sidebar_current: "docs-newrelic-resource-alert-condition"
description: |-
  Create and manage alert conditions for APM, Browser, and Mobile in New Relic.
---

# Resource: newrelic\_alert\_condition

Use this resource to create and manage alert conditions for APM, Browser, and Mobile in New Relic.

## Example Usage

```hcl
data "newrelic_application" "app" {
  name = "my-app"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo"
  type        = "apm_app_metric"
  entities    = [data.newrelic_application.app.id]
  metric      = "apdex"
  runbook_url = "https://www.example.com"
  condition_scope = "application"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the policy where this condition should be used.
  * `name` - (Required) The title of the condition. Must be between 1 and 64 characters, inclusive.
  * `type` - (Required) The type of condition. One of: `apm_app_metric`, `apm_kt_metric`, `browser_metric`, `mobile_metric`
  * `entities` - (Required) The instance IDs associated with this condition.
  * `metric` - (Required) The metric field accepts parameters based on the `type` set. One of these metrics based on `type`:
    * `apm_app_metric`
      * `apdex`
      * `error_percentage`
      * `response_time_background`
      * `response_time_web`
      * `throughput_background`
      * `throughput_web`
      * `user_defined`
    * `apm_kt_metric`
      * `apdex`
      * `error_count`
      * `error_percentage`
      * `response_time`
      * `throughput`
    * `browser_metric`
      * `ajax_response_time`
      * `ajax_throughput`
      * `dom_processing`
      * `end_user_apdex`
      * `network`
      * `page_rendering`
      * `page_view_throughput`
      * `page_views_with_js_errors`
      * `request_queuing`
      * `total_page_load`
      * `user_defined`
      * `web_application`
    * `mobile_metric`
      * `database`
      * `images`
      * `json`
      * `mobile_crash_rate`
      * `network_error_percentage`
      * `network`
      * `status_error_percentage`
      * `user_defined`
      * `view_loading`
  * `condition_scope` - (Required for some types) `application` or `instance`.  Choose `application` for most scenarios.  If you are using the JVM plugin in New Relic, the `instance` setting allows your condition to trigger [for specific app instances](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/scope-alert-thresholds-specific-instances).
  * `enabled` - (Optional) Whether the condition is enabled or not. Defaults to true.
  * `gc_metric` - (Optional) A valid Garbage Collection metric e.g. `GC/G1 Young Generation`.
  * `violation_close_timer` - (Optional) Automatically close instance-based violations, including JVM health metric violations, after the number of hours specified. Must be: `1`, `2`, `4`, `8`, `12` or `24`.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `term` - (Required) A list of terms for this condition. See [Terms](#terms) below for details.
  * `user_defined_metric` - (Optional) A custom metric to be evaluated.
  * `user_defined_value_function` - (Optional) One of: `average`, `min`, `max`, `total`, or `sample_size`.

## Terms

The `term` mapping supports the following arguments:

  * `duration` - (Required) In minutes, must be in the range of `5` to `120`, inclusive.
  * `operator` - (Optional) `above`, `below`, or `equal`.  Defaults to `equal`.
  * `priority` - (Optional) `critical` or `warning`.  Defaults to `critical`. Terms must include at least one `critical` priority term
  * `threshold` - (Required) Must be 0 or greater.
  * `time_function` - (Required) `all` or `any`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the alert condition.

## Import

Alert conditions can be imported using notation `alert_policy_id:alert_condition_id`, e.g.

```
$ terraform import newrelic_alert_condition.main 123456:6789012345
```
