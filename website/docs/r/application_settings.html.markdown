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

~> **WARNING:** We encourage you to  use this resource to manage all application settings together not just a few to avoid potential issues like incompatibility or unexpected behavior.

## Example Usage

```hcl
resource "newrelic_application_settings" "app" {
  guid = "rhbwkguhfjkewqre4r9"
  name = "my-app"
  app_apdex_threshold = "0.7"
  use_server_side_config = true
  transaction_tracer {
    explain_query_plans {
      query_plan_threshold_type  = "VALUE"
      query_plan_threshold_value = "0.5"
    }
    sql {
      record_sql = "RAW"
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
  enable_slow_sql = true
  tracer_type = "NONE"
  enable_thread_profiler = true
}
```
## Argument Reference

The following arguments are supported:

* `guid` - (Required) The GUID of the application in New Relic APM.
* `name` - (Optional) A custom name or alias you can give the application in New Relic APM.
* `app_apdex_threshold` - (Optional) The acceptable response time limit (Apdex threshold) for the application.
* `use_server_side_config` - (Optional) Enable or disable server side monitoring for the New Relic application.
* `transaction_tracer` - (Optional) Configuration block for transaction tracer. Providing this block enables transaction tracing. The following arguments are supported:
  * `stack_trace_threshold_value` - (Optional) The response time threshold for collecting stack traces.
  * `transaction_threshold_type` - (Optional) The type of threshold for transactions. Valid values are `VALUE`,`APDEX_F`(4 times your apdex target)
  * `transaction_threshold_value` - (Optional) The threshold value for transactions(in seconds).
  * `explain_query_plans` - (Optional) Configuration block for query plans. Including this block enables the capture of query plans. The following arguments are supported:
    * `query_plan_threshold_value` - (Optional) The response time threshold for capturing query plans(in seconds).
    * `query_plan_threshold_type` - (Optional) The type of threshold for query plans. Valid values are `VALUE`,`APDEX_F`(4 times your apdex target)
  * `sql` - (Optional) Configuration block for SQL logging.  Including this block enables SQL logging. The following arguments are supported:
    * `record_sql` - (Required) The level of SQL recording. Valid values ar `OBFUSCATED`,`OFF`,`RAW` (Mandatory attribute when `sql` block is provided).
* `error_collector` - (Optional) Configuration block for error collection. Including this block enables the error collector. The following arguments are supported:
  * `expected_error_classes` - (Optional) A list of expected error classes.
  * `expected_error_codes` - (Optional) A list of expected error codes(any status code between 100-900).
  * `ignored_error_classes` - (Optional) A list of ignored error classes.
  * `ignored_error_codes` - (Optional) A list of ignored error codes(any status code between 100-900).
* `enable_slow_sql` - (Optional) Enable or disable the collection of slowest database queries in your traces.
* `tracer_type` - (Optional) Configures the type of tracer used. Valid values are `CROSS_APPLICATION_TRACER`, `DISTRIBUTED_TRACING`, `NONE`, `OPT_OUT`.
* `enable_thread_profiler` - (Optional) Enable or disable the collection of thread profiling data.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the application.

## Import

Applications can be imported using notation `application_guid`, e.g.

```
$ terraform import newrelic_application_settings.main Mzk1NzUyNHQVRJNTxBUE18QVBQTElDc4ODU1MzYx
```

## Notes

-> **NOTE:** The `newrelic_application_settings` resource cannot be deleted directly via Terraform. It can only reset application settings to their initial state.
