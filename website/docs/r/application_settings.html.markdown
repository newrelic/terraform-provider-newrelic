---
layout: "newrelic"
page_title: "New Relic: newrelic_application_settings"
sidebar_current: "docs-newrelic-resource-application-settings"
description: |-
  Manage configuration for an existing application in New Relic.
---

# Resource: newrelic_application_settings

-> **NOTE:** Applications are not created by this resource, but are created by
a reporting agent.

Use this resource to manage configuration for an application that already
exists in New Relic.

## Example Usage

```hcl
resource "newrelic_application_settings" "app" {
  guid = "rhbwkguhfjkewqre4r9"
  name = "my-app"
  app_apdex_threshold = "0.7"
  end_user_apdex_threshold = "0.8"
  enable_real_user_monitoring = false
  transaction_tracing {
    explain_query_plans {
      query_plan_threshold_type  = "VALUE"
      query_plan_threshold_value = "0.5"
    }
    sql {
      record_sql = "OFF"
    }
    stack_trace_threshold_value = "0.5"
    transaction_threshold_type = "VALUE"
    transaction_threshold_value = "0.5"
  }
  error_collector {
    expected_error_classes = []
    expected_error_codes = []
    ignored_error_classes = []
    ignored_error_codes = []
  }
  tracer_type = "RAW"
  enable_thread_profiler = true
}
```

## Argument Reference

The following arguments are supported:

* `guid` - (Required) The name of the application in New Relic APM.
* `name` - (Optional) The name of the application in New Relic APM.
* `app_apdex_threshold` - (Optional) The apdex threshold for the New Relic application.
* `end_user_apdex_threshold` - (Optional) The user's apdex threshold for the New Relic application.
* `enable_real_user_monitoring` - (Optional) Enable or disable real user monitoring for the New Relic application.
* `transaction_tracing` - (Optional) Configuration block for transaction tracing. If provided, it enables transaction tracing; otherwise, it disables transaction tracing. The following arguments are supported:
  * `stack_trace_threshold_value` - (Optional) The threshold value for stack traces.
  * `transaction_threshold_type` - (Optional) The type of threshold for transactions (e.g., VALUE).
  * `transaction_threshold_value` - (Optional) The threshold value for transactions.
  * `explain_query_plans` - (Optional) Configuration block for query plans. If provided, it enables explain query plans; otherwise, it disables explain query plans. The following arguments are supported:
    * `query_plan_threshold_value` - (Optional) The threshold value for query plans.
    * `query_plan_threshold_type` - (Optional) The type of threshold for query plans (e.g., VALUE).
  * `sql` - (Optional) Configuration block for SQL logging. If provided, it enables sql logging; otherwise, it disables sql logging. The following arguments are supported:
    * `record_sql` - (Optional) The level of SQL recording (e.g., OFF).
* `error_collector` - (Optional) Configuration block for error collection. The following arguments are supported:
  * `expected_error_classes` - (Optional) A list of expected error classes.
  * `expected_error_codes` - (Optional) A list of expected error codes.
  * `ignored_error_classes` - (Optional) A list of ignored error classes.
  * `ignored_error_codes` - (Optional) A list of ignored error codes.
* `tracer_type` - (Optional) The type of tracer (e.g., RAW).
* `enable_thread_profiler` - (Optional) Enable or disable the thread profiler.
```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the application.

## Import

Applications can be imported using notation `application_guid`, e.g.

```
$ terraform import newrelic_application_settings.main 6789012345
```

## Notes

-> **NOTE:** Applications that have reported data in the last twelve hours
cannot be deleted.
