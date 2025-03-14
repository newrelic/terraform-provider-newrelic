---
layout: "newrelic"
page_title: "New Relic: newrelic_application_settings"
sidebar_current: "docs-newrelic-resource-application-settings"
description: |-
  Manage configuration for an existing application in New Relic.
---

# Resource: newrelic_application_settings

-> **NOTE:** Applications are not created by this resource, but are created by a reporting agent.

Use this resource to manage configuration for an application that already exists in New Relic.

~> **WARNING:** We encourage you to use this resource to manage all application settings together, not just a few, to avoid potential issues like incompatibility or unexpected behavior.

## Example Usage

```hcl
resource "newrelic_application_settings" "app" {
  # please note that the arguments 'guid' and 'name' are mutually exclusive
  # using the 'guid' argument is preferred over using the 'name' argument
  guid                   = "Mxxxxxxxxxxxxxxxxxxxxx"
  name                   = "Sample New Relic APM Application"
  app_apdex_threshold    = "0.7"
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
    transaction_threshold_type  = "VALUE"
    transaction_threshold_value = "0.5"
  }
  error_collector {
    expected_error_classes = []
    expected_error_codes   = []
    ignored_error_classes  = []
    ignored_error_codes    = []
  }
  enable_slow_sql        = true
  tracer_type            = "NONE"
  enable_thread_profiler = true
}
```
## Argument Reference

The following arguments are supported:

* `guid` - (Required) The GUID of the application in New Relic APM.

-> **NOTE:** While the attribute `guid` is not mandatory at a schema level, it is recommended to use `guid` over `name`, as support for using `name` with this resource shall eventually be discontinued. Please see the note under `name` for more details.

* `name` - (Optional) A custom name or alias you can give the application in New Relic APM.

-> **NOTE:** <b style="color:red;">Please refrain from using the deprecated attribute `name`</b>with the resource `newrelic_application_settings` and use `guid` instead. For more information on the usage of `guid` against `name` and associated implications if the resource is upgraded from an older version of the New Relic Terraform Provider, please see the note in [this section](#deprecated-attribute-name) below.

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

### Deprecated Attribute `name`
Starting v3.59.0 of the New Relic Terraform Provider, the attribute `name` in this resource is no longer recommended to specify the application to apply the settings to - this has been made ineffective, in favor of the attribute `guid`, which shall henceforth be the *recommended* attribute to select the application to apply settings to, in order to eliminate problems caused by applications bearing the same name. In light of the above, please refrain from using `name`, and use `guid` (with the GUID of the application) instead.

However, if you have been using this resource prior to v3.59.0 and have upgraded to this new version, it would be required that the attribute `name` (already present with a non-null value in your configuration) needs to continue to exist until the first `terraform plan` and `terraform apply` after the upgrade to v3.59.0. This allows changes made to handle backward compatibility take effect, so the state of the resource is updated to match the expected format, starting v3.59.0. You may switch to using `guid` instead of `name` after the first `terraform plan` and `terraform apply` following the upgrade to v3.59.0, in such a case.

Additionally, when upgrading this resource from older versions, you may observe drift associated with attributes not yet controlled by your configuration (showing such attributes being set to null). This is expected due to the inclusion of new attributes in the resource that are required by the resource state. You can safely ignore these drifts shown during/after the first `terraform plan` and `terraform apply` after upgrading from older versions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the application.

## Import

Applications can be imported using notation `application_guid`, e.g.

```
$ terraform import newrelic_application_settings.main Mzk1NzUyNHQVRJNTxBUE18QVBQTElDc4ODU1MzYx
```

## Notes

-> **NOTE:** The `newrelic_application_settings` resource cannot be deleted directly via Terraform. It can only reset application settings to their initial state.
