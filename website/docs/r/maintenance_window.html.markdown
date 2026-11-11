---
layout: "newrelic"
page_title: "New Relic: newrelic_maintenance_window"
sidebar_current: "docs-newrelic-resource-maintenance-window"
description: |-
  Create and manage a New Relic Maintenance Window.
---

# Resource: newrelic\_maintenance\_window

Use this resource to create, update, and delete New Relic Maintenance Windows. Maintenance windows let you schedule planned downtime for one or more affected entities (Service Levels) so that activity occurring during the window is excluded from calculations such as SLI/SLO compliance.

A New Relic User API key is required to provision this resource. Set the `api_key`
attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment
variable with your User API key.

Important:
- Maintenance Windows are currently behind a feature flag and not yet generally available to all accounts.

## Example Usage

```hcl
resource "newrelic_maintenance_window" "foo" {
  name        = "One-time patching window"
  description = "Excludes a single planned maintenance event from SLI calculations."

  scope_id   = 12345678
  scope_type = "ACCOUNT"

  start_time = "2026-12-01T22:00:00"
  duration   = "PT2H"
  timezone   = "America/Los_Angeles"
  rrule = "FREQ=WEEKLY;BYDAY=MO,WE,FR"

  affected_entity_type = "SERVICE_LEVEL"
  affected_entities    = ["MXxBUE18QVBQTElDQVRJT058MQ"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the maintenance window.
* `description` - (Optional) A free-text description of the maintenance window, for your own context (e.g. why it's scheduled, what work it covers).
* `scope_id` - (Required) The unique identifier of the object the maintenance window is scoped to. In practice, this is your New Relic account ID. Changing this forces a new resource.
* `scope_type` - (Required) The kind of scope `scope_id` refers to. Valid values are `ACCOUNT` and `ORGANIZATION`. For Service Level maintenance windows, this should be `ACCOUNT`. Changing this forces a new resource.
* `start_time` - (Required) The start date and time of the maintenance window (or, for a recurring window, of its first occurrence), as an ISO 8601 date-time **without** a UTC offset or time zone suffix (e.g. `2026-12-01T22:00:00`). The time zone to interpret this value in is supplied separately via the `timezone` argument.
* `duration` - (Required) How long each occurrence of the maintenance window lasts, expressed in ISO 8601 duration format (e.g. `PT2H` for two hours, `PT30M` for thirty minutes, `P1D` for one day). The window's end time is computed as `start_time + duration`.
* `timezone` - (Required) The IANA time zone name that `start_time` (and the recurrence pattern, if `rrule` is set) should be interpreted in (e.g. `America/Los_Angeles`, `America/New_York`, `UTC`).
* `rrule` - (Optional) The recurrence rule for the maintenance window, expressed in iCalendar/RFC 5545 `RRULE` format, for a recurring maintenance window. This corresponds to the **Frequency** dropdown in the UI (`Once`, `Daily`, `Weekly`, `Monthly`). Omit this argument entirely for `Once` — a maintenance window that occurs only a single time, at `start_time`. This is the default behavior when `rrule` isn't set.
  * The `FREQ` component controls how often the window recurs: `FREQ=DAILY`, `FREQ=WEEKLY`, or `FREQ=MONTHLY`.
  * `BYDAY` is required when `FREQ=WEEKLY` — e.g. `FREQ=WEEKLY;BYDAY=MO,WE,FR` to repeat every Monday, Wednesday, and Friday.
  * Either `BYDAY` or `BYMONTHDAY` is required when `FREQ=MONTHLY` — e.g. `FREQ=MONTHLY;BYMONTHDAY=15` for a specific date each month, or `FREQ=MONTHLY;BYDAY=1SU` for a relative day (e.g. "the first Sunday").
  * When you pick a recurring frequency in the UI, it also asks for **Active from** and **Active until** dates, in addition to the **Start date and time**/**End date and time** of a single occurrence. These map onto existing arguments rather than needing separate ones: **Active from** is just `start_time` (the first occurrence also marks when the recurrence pattern becomes active), and **Active until** is expressed by appending `;UNTIL=20271231T235959Z` to the `rrule` string. Selecting **Never ends** in the UI simply means omitting `UNTIL` (and `COUNT`) from `rrule`, so the recurrence continues indefinitely.
* `affected_entity_type` - (Required) The type of the entities affected by this maintenance window. The only currently supported value is `SERVICE_LEVEL`.
* `affected_entities` - (Required) The list of entity GUIDs affected by this maintenance window (e.g. the `guid` of one or more `newrelic_service_level` resources, or other entities of the type specified in `affected_entity_type`). At least one entity GUID must be provided.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the maintenance window, assigned after creation.

## Additional Examples

Monthly recurring maintenance window, on the 15th of each month through the end of 2027:

```hcl
resource "newrelic_maintenance_window" "monthly" {
  name       = "Monthly maintenance"
  scope_id   = 12345678
  scope_type = "ACCOUNT"

  start_time = "2026-12-15T22:00:00"
  duration   = "PT2H"
  timezone   = "America/Los_Angeles"
  rrule      = "FREQ=MONTHLY;BYMONTHDAY=15;UNTIL=20271231T235959Z"

  affected_entity_type = "SERVICE_LEVEL"
  affected_entities    = ["MXxBUE18QVBQTElDQVRJT058MQ"]
}
```

One-time maintenance window covering multiple affected entities:

```hcl
resource "newrelic_maintenance_window" "multi_entity" {
  name       = "Multi-entity maintenance"
  scope_id   = 12345678
  scope_type = "ACCOUNT"

  start_time = "2026-12-01T22:00:00"
  duration   = "PT2H"
  timezone   = "America/Los_Angeles"

  affected_entity_type = "SERVICE_LEVEL"
  affected_entities    = ["MXxBUE18QVBQTElDQVRJT058MQ", "MXxBUE18QVBQTElDQVRJT058Mg", "MXxBUE18QVBQTElDQVRJT058Mw"]
}
```

## Import

New Relic Maintenance Windows can be imported using their ID, e.g.

```bash
$ terraform import newrelic_maintenance_window.foo <id>
```
