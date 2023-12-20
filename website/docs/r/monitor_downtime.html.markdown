---
layout: "newrelic"
page_title: "New Relic: newrelic_monitor_downtime"
sidebar_current: "docs-newrelic-resource-monitor_downtime"
description: |-
    Create and manage Monitor Downtimes in New Relic.
---

# Resource: newrelic\_monitor\_downtime

Use this resource to create, update, and delete [Monitor Downtimes](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/using-monitors/monitor-downtimes-disable-monitoring-during-scheduled-maintenance-times/) in New Relic.

## Example Usage

```hcl
resource "newrelic_monitor_downtime" "foo" {
  name = "Sample Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>"
  ]
  mode       = "WEEKLY"
  start_time = "2023-11-30T10:30:00"
  end_time   = "2023-12-10T02:45:30"
  time_zone  = "Asia/Kolkata"
  end_repeat {
    on_date = "2023-12-20"
  }
  maintenance_days = [
    "FRIDAY",
    "SATURDAY",
  ]
}
```
Monitor Downtimes are of four types; **one-time**, **daily**, **weekly** and **monthly**. For more details on each type and the right arguments that go with them, check out the [argument reference](#argument-reference) and [examples](#examples) sections below.

## Argument Reference

### Arguments Common To All Four Types of Monitor Downtimes

* `account_id`- (Optional) The account in which the monitor downtime would be created. Defaults to the value of the environment variable `NEW_RELIC_ACCOUNT_ID` (or the `account_id` specified in the `provider{}`), if not specified.
* `name` - (Required) Name of the monitor downtime to be created.
* `mode` - (Required) One of the four modes of operation of monitor downtimes - `ONE_TIME`, `DAILY`, `MONTHLY` or `WEEKLY`.
* `monitor_guids` - (Optional) A list of GUIDs of synthetic monitors the monitor downtime would need to be applied to.
* `start_time` - (Required) The time at which the monitor downtime would begin to operate, a timestamp specified in the ISO 8601 format without the offset/timezone - for instance, `2023-12-20T10:48:53`.
* `end_time` - (Required) The time at which the monitor downtime would end operating, a timestamp specified in the ISO 8601 format without the offset/timezone - for instance, `2024-01-05T14:27:07`.
* `timezone` - (Required) The timezone in which timestamps `start_time` and `end_time` have been specified. Valid timezones which may be specified with this argument can be found in this [list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List), under the column '**TZ identifier**'. 

### Arguments Specific Only to Certain Types of Monitor Downtimes

The following arguments go only with certain types of monitor downtimes, and are hence, optional, at a resource schema level. However, _some of these_ are **required** to be specified with certain types of monitor downtimes - please see notes adjoining arguments in the list below, and [examples](#examples) to obtain an understanding of the usage of apt arguments with each type of monitor downtimes.

* `end_repeat` - Options which may be used to specify when the repeat cycle of the monitor should end. This argument comprises the following nested arguments -
  * `on_date` - The date on which the monitor downtime's repeat cycle would need to come to an end, a string in `DDDD-MM-YY` format.
  * `on_repeat` - An integer that specifies the number of occurrences, after which the monitor downtime's repeat cycle would need to come to an end.

-> **NOTE:** `end_repeat` **can only be used with the modes** `DAILY`, `MONTHLY` and `WEEKLY` and **is an optional argument** when monitor downtimes of these modes are created. Additionally, **either** `on_date` or `on_repeat` **are required to be specified with** `end_repeat`, but not both, as `on_date` and `on_repeat` are mutually exclusive.

* `maintenance_days` - A list of days on which weekly monitor downtimes would function. Valid values which go into this list would be `"SUNDAY"`, `"MONDAY"`, `"TUESDAY"`, `"WEDNESDAY"`, `"THURSDAY"`, `"FRIDAY"` and/or `"SATURDAY"`.

-> **NOTE:** `maintenance_days` **can only be used with the mode** `WEEKLY`, and **is a required argument** with weekly monitor downtimes (i.e. if the `mode` is `WEEKLY`). 

* `frequency` - Options which may be used to specify the configuration of a monthly monitor downtime. This argument comprises the following nested arguments -
    * `days_of_month` - A list of integers, specifying the days of a month on which the monthly monitor downtime would function, e.g. [3, 6, 14, 23].
    * `days_of_week` - An argument that specifies a day of a week and its occurrence in a month, on which the monthly monitor downtime would function. This argument, further, comprises the following nested arguments - 
      * `week_day` - A day of the week (one of `"SUNDAY"`, `"MONDAY"`, `"TUESDAY"`, `"WEDNESDAY"`, `"THURSDAY"`, `"FRIDAY"` or `"SATURDAY"`).
      * `ordinal_day_of_month` - The occurrence of `week_day` in a month (one of `"FIRST"`, `"SECOND"`, `"THIRD"`, `"FOURTH"`, `"LAST"`).

-> **NOTE:** `frequency` **can only be used with the mode** `MONTHLY`, and **is a required argument** with monthly monitor downtimes (if the `mode` is `MONTHLY`). Additionally, **either** `days_of_month` or `days_of_week` **are required to be specified with** `frequency`, but not both, as `days_of_month` and `days_of_week` are mutually exclusive. If `days_of_week` is specified, values of **both** of its nested arguments, `week_day` and `ordinal_day_of_month` **would need to be specified** too.

## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the monitor downtime.

## Examples

### One-Time Monitor Downtime

The below example illustrates creating a **one-time** monitor downtime. 

```hcl
resource "newrelic_monitor_downtime" "sample_one_time_newrelic_monitor_downtime" {
  name = "Sample One Time Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>",
  ]
  mode       = "ONE_TIME"
  start_time = "2023-12-04T10:15:00"
  end_time   = "2024-01-04T16:24:30"
  time_zone  = "America/Los_Angeles"
}
```

### Daily Monitor Downtime

The below example illustrates creating a **daily** monitor downtime. 

Note that `end_repeat` has been specified in the configuration; however, this is optional, in accordance with the rules of `end_repeat` specified in the [argument reference](#argument-reference) section above. This example uses the `on_date` nested argument of `end_repeat`, however, the other nested argument, `on_repeat` may also be used _instead_, as you may see in some of the other examples below; though both `on_date` and `on_repeat` cannot be specified together, as they are mutually exclusive.

```hcl
resource "newrelic_monitor_downtime" "sample_daily_newrelic_monitor_downtime" {
  name = "Sample Daily Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>",
  ]
  mode       = "DAILY"
  start_time = "2023-12-04T18:15:00"
  end_time   = "2024-01-04T07:15:00"
  end_repeat {
    on_date = "2023-12-25"
  }
  time_zone = "Asia/Kolkata"
}
```

### Weekly Monitor Downtime

The below example illustrates creating a **weekly** monitor downtime. 

Note that `maintenance_days` has been specified in the configuration as it is required with weekly monitor downtimes; and `end_repeat` has not been specified as it is optional, all in accordance with the rules of these arguments specified in the [argument reference](#argument-reference) section above.

```hcl
resource "newrelic_monitor_downtime" "sample_weekly_newrelic_monitor_downtime" {
  name = "Sample Weekly Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>",
  ]
  mode       = "WEEKLY"
  start_time = "2023-12-04T14:15:00"
  end_time   = "2024-01-04T23:55:00"
  time_zone  = "US/Hawaii"
  maintenance_days = [
    "SATURDAY",
    "SUNDAY"
  ]
} 
```

### Monthly Monitor Downtime

The below example illustrates creating a **monthly** monitor downtime.

Note that `frequency` has been specified in the configuration as it is required with monthly monitor downtimes, and `end_repeat` has been specified too, though it is optional. `frequency` has been specified with `days_of_week` comprising both of its nested arguments, `ordinal_day_of_month` and `week_day`; all in accordance with the rules of these arguments specified in the [argument reference](#argument-reference) section above. 

```hcl
resource "newrelic_monitor_downtime" "sample_monthly_newrelic_monitor_downtime" {
  name = "Sample Monthly Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>",
  ]
  mode       = "MONTHLY"
  start_time = "2023-12-04T07:15:00"
  end_time   = "2024-01-04T19:15:00"
  end_repeat {
    on_repeat = 6
  }
  time_zone = "Europe/Dublin"
  frequency {
    days_of_week {
      ordinal_day_of_month = "SECOND"
      week_day             = "SATURDAY"
    }
  }
} 
```
However, the `frequency` block in monthly monitor downtimes may also be specified with its other nested argument, `days_of_month`, as shown in the example below - though both `days_of_month` and `days_of_week` cannot be specified together, as they are mutually exclusive.
```hcl
resource "newrelic_monitor_downtime" "sample_monthly_newrelic_monitor_downtime" {
  name = "Sample Monthly Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>",
  ]
  mode       = "MONTHLY"
  start_time = "2023-12-04T07:15:00"
  end_time   = "2024-01-04T19:15:00"
  end_repeat {
    on_repeat = 6
  }
  time_zone = "Europe/Dublin"
  frequency {
    days_of_month = [3, 6, 14, 23]
  }
} 
```
## Import

A monitor downtime can be imported into Terraform configuration using its `guid`, i.e.

```bash
$ terraform import newrelic_monitor_downtime.monitor <guid>
```