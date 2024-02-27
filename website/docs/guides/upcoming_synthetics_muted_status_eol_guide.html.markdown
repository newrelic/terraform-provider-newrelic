---
layout: "newrelic"
page_title: "New Relic Terraform Provider v4: Major Release Details and More"
sidebar_current: "docs-newrelic-provider-synthetic-monitors-muted-status-eol-guide"
description: |-
  Use this guide to find details on the end-of-life of the MUTED status of Synthetic Monitors, as well as alternatives which can replicate the same behavior.
---
## Synthetic Monitors' MUTED Status EOL: Implications, Alternatives and More 📢

Starting with version **3.33.0** of the New Relic Terraform Provider that was released on February 29, 2024, the `MUTED` status of Synthetic Monitors is no longer supported. The reason for the release is discontinued support for the `MUTED` status for Synthetic Monitors, to match the API change made to remove support for the `MUTED` status of Synthetic Monitors, effective on February 29, 2024.

The following is a comprehensive guide that lists the implications of this end-of-life, provides additional details, and offers alternatives to the `MUTED` status of Synthetic Monitors in the New Relic Terraform Provider.

## About: Synthetic Monitors' `MUTED` Status EOL and an Associated Release of the New Relic Terraform Provider

As mentioned in the initial section of this guide, **v3.33.0** of the New Relic Terraform Provider was released to discontinue support for `MUTED` as a valid `status` value for Synthetic Monitors within the New Relic Terraform Provider. At the EOL, the status of all `MUTED` monitors were changed to `ENABLED`.

For more detailed information regarding the end-of-life of the `MUTED` status of Synthetic Monitors, as well as the associated implications and alternatives, please refer to the subsequent sections of this guide.

### EOL of Synthetic Monitors' 'MUTED' Status
On February 29, 2024, Synthetics discontinued support for the `MUTED` status, as a result of which `MUTED` is no longer a valid `status` value and any operations performed using the `MUTED` status in any version of the New Relic Terraform Provider will fail, since NerdGraph mutations/queries will no longer recognize `MUTED` as a valid status value for Synthetic Monitors.

You would be affected by this end-of-life only if your Terraform configuration continues to comprise `MUTED` as the value of the argument `status` post the end-of-life too, in any of the following resources:

* [`newrelic_synthetics_monitor`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/synthetics_monitor)
* [`newrelic_synthetics_broken_links_monitor`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/synthetics_broken_links_monitor)
* [`newrelic_synthetics_cert_check_monitor`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/synthetics_cert_check_monitor)
* [`newrelic_synthetics_script_monitor`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/synthetics_script_monitor)
* [`newrelic_synthetics_step_monitor`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/synthetics_step_monitor)

As stated earlier, as of **v3.33.0** of the New Relic Terraform Provider, the `status` argument in all types of Synthetic Monitors can no longer be `MUTED`. The following sections of this guide provide alternatives to replicate the behavior of the `MUTED` status for Synthetic Monitors.

### Implications of the EOL, Solutions
The following breaking changes/implications would be experienced only by customers who have a Terraform configuration that fits into the criteria described above. 

##### Usage of the New Relic Terraform Provider < v3.33.0 Comprising `MUTED` Synthetic Monitors after the EOL 
* Since **v3.33.0** of the New Relic Terraform Provider invalidates the `MUTED` status of Synthetic Monitors, a change that came into effect starting with this version of the provider (which cannot retrospectively be applied to older versions of the provider), if versions prior to **v3.33.0** are used with muted Synthetic Monitors after the end-of-life, validation checks would not obstruct Terraform operations as `MUTED` would be deemed a valid status of Synthetic Monitors in older versions; however, upon performing `terraform apply` to create/update Synthetic Monitors with status `MUTED`, the operation would **fail**, as the API would throw an error, specifying that `MUTED` is not a valid status value anymore.
* In addition to the above, as communicated previously by New Relic, since all monitors with the status `MUTED` will have their status changed to `ENABLED` on the date of the end-of-life, any monitors in your Terraform configuration with the status `MUTED` (whose status has not been moved out of `MUTED`, prior to the end-of-life) will result in a drift being displayed when attempting to plan or apply the configuration containing monitors with the status `MUTED`. This is because Synthetics has changed the status of these monitors to `ENABLED`. 

##### Usage of the New Relic Terraform Provider >= v3.33.0 Comprising `MUTED` Synthetic Monitors after the EOL
* The change to invalidate the `MUTED` status of Synthetic Monitors works with **v3.33.0** of the New Relic Terraform Provider and above. As a consequence, running `terraform plan` on Terraform configuration comprising Synthetic Monitors with status `MUTED` would throw an error, as the `MUTED` status is no longer valid, leading to validation failure. This, in turn, would not allow planning or applying your configuration.

The **_solution_** to both of the implications listed above would be to replace all instances of `MUTED` in all Synthetic Monitors across your Terraform configuration to either `ENABLED` or `DISABLED`. This would allow performing a successful `terraform plan` and `terraform apply`. Additionally, please choose an appropriate alternative from the options described in the following section to enforce monitor muting through existing resources provided by the New Relic Terraform Provider. 

It may be noted that Synthetic Monitors need to be in the state `ENABLED` in order to use alert muting rules and monitor downtime with Synthetic Monitors.

## Alternatives To `MUTED` Status

There are two key alternatives one can opt for, to replicate the behavior of the `MUTED` status of Synthetic Monitors.
* [**Alert Muting Rules**](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-notifications/muting-rules-suppress-notifications/)
  * Alert Muting Rules, by definition, are similar to the `MUTED` status of Synthetic Monitors in terms of behavior, though these cater to a wider scope; all kinds of alerts. These help mute alerts on the basis of pre-defined schedules and attribute matching (when alerts match the condition(s) prescribed by the user in terms of incident event attributes, operators and values, they can be muted). See [this page](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/advanced-alerts/understand-technical-concepts/incident-event-attributes/) for a comprehensive list of attributes supported by Muting Rules.
  * This feature may be availed from the New Relic One UI, NerdGraph, and also via the resource [`newrelic_alert_muting_rule`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/alert_muting_rule) in the New Relic Terraform Provider. Check out the example below to find how this resource can _exactly_ be used to substitute the `MUTED` status of a Synthetic Monitor.
* [**Monitor Downtime**](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/using-monitors/monitor-downtimes-disable-monitoring-during-scheduled-maintenance-times/)
  * A 'Monitor Downtime', as the name suggests, helps set up a 'downtime' or a maintenance window for Synthetic Monitors, in which period they do not run, as a result of which alerts are not raised, and no notifications are received.
  * This feature may be availed from the New Relic One UI, NerdGraph and also via a _new_ resource that's been built to facilitate managing Monitor Downtimes via Terraform, [`newrelic_monitor_downtime`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/monitor_downtime). Check out the example below for an understanding of how this resource helps mute monitors.

It is important to note the difference between the functioning of Alert Muting Rules and Monitor Downtimes when these are considered to be used as substitutes to the `MUTED` status of Synthetic Monitors.
* Since Alert Muting Rules are designed to mute alerts based on conditions they match, alerts generated by checks performed by monitors are muted; however, this does not affect the checks performed by these monitors, which would continue to function as usual.
* Since Monitor Downtimes are dedicated to scheduling "downtime"s of monitors, no alerts would be generated by monitors in this case, as they would stop running checks for the period defined in the downtime.
  Users may need to choose the right alternative, based on the expected behavior they desire, when monitors are muted.

### Substituting Synthetic Monitors `MUTED` Status With Alert Muting Rules

The simplest method to mute a Synthetic Monitor via the [`newrelic_alert_muting_rule`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/alert_muting_rule) resource is to use the GUID of the Synthetic Monitor in the condition used to define the muting rule, and match it against the `entity.guid` attribute in the condition. The following code snippet gives an example of the same.

```hcl
resource "newrelic_synthetics_monitor" "sample_synthetics_monitor" {
  status           = "ENABLED"
  name             = "Sample Monitor"
  period           = "EVERY_MINUTE"
  uri              = "https://www.one.newrelic.com"
  type             = "BROWSER"
  locations_public = ["AP_EAST_1"]
  custom_header {
    name  = "some_name"
    value = "some_value"
  }
  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
resource "newrelic_alert_muting_rule" "sample_alert_muting_rule" {
  name        = "Sample Muting Rule"
  enabled     = true
  description = "A muting rule, deployed to test a muting rule with a monitor."
  condition {
    conditions {
      attribute = "entity.guid"
      operator  = "EQUALS"
      values    = [newrelic_synthetics_monitor.sample_synthetics_monitor.id]
    }
    operator = "AND"
  }
  schedule {
    start_time = "2023-10-31T06:30:00"
    end_time   = "2023-11-01T16:30:00"
    time_zone  = "America/Los_Angeles"
  }
}
```

The configuration of the muting rule may be customized, based on any apt approach identified - one of which could be to use the attribute `conditionId` in the condition, to match it against the ID of the alert condition that checks for failures of tests run by the preferred Synthetic Monitor (for instance, a `SyntheticCheck` based NRQL Alert Condition), so the muting rule is applied to any alerts originating out of the condition that evaluates monitor failures/successes. Head over to the [documentation of the resource](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/alert_muting_rule) for more details on attributes supported by `condition` blocks.

### Substituting Synthetic Monitors `MUTED` Status With a Monitor Downtime

Setting up a Monitor Downtime resource with GUIDs of the right monitors would disable Synthetic checks in the window specified, thereby, muting the monitor (please read the differences between Muting Rules and Monitor Downtimes, explained above).

```hcl
resource "newrelic_monitor_downtime" "foo" {
  name = "Sample Monitor Downtime"
  monitor_guids = [
    "<GUID-1>",
    "<GUID-2>"
  ]
  mode       = "WEEKLY"
  start_time = "2023-11-30T10:30:00"
  end_time   = "2023-12-10T10:30:00"
  time_zone  = "Asia/Kolkata"
  end_repeat {
    on_date = "2023-12-20"
  }
  maintenance_days = [
    "MONDAY",
    "TUESDAY",
  ]
}
```

For more examples and details of arguments of the `newrelic_monitor_downtime` resource, head over to [this page](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/monitor_downtime). 