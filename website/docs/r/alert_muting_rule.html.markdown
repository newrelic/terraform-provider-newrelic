---
layout: 'newrelic'
page_title: 'New Relic: newrelic_alert_muting_rule'
sidebar_current: 'docs-newrelic-resource-alert-muting-rule'
description: |-
  Create a muting rule for New Relic Alerts violations.
---

# Resource: newrelic_alert_muting_rule

Use this resource to create a muting rule for New Relic Alerts violations.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/docs/providers/newrelic/index.html) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

```hcl
resource "newrelic_alert_muting_rule" "foo" {
	name = "Example Muting Rule"
	enabled = true
	description = "muting rule test."
	condition {
		conditions {
			attribute 	= "product"
			operator 	= "EQUALS"
			values 		= ["APM"]
		}
		conditions {
			attribute 	= "targetId"
			operator 	= "EQUALS"
			values 		= ["Muted"]
		}
		operator = "AND"
	}
    schedule {
      start_time = "2021-01-28T15:30:00"
      end_time = "2021-01-28T16:30:00"
      time_zone = "America/Los_Angeles"
      repeat = "WEEKLY"
      weekly_repeat_days = ["MONDAY", "WEDNESDAY", "FRIDAY"]
      repeat_count = 42
    }
}
```

## Argument Reference

The following arguments are supported:
  * `account_id` - (Optional) The account id of the MutingRule.
  * `condition`  - (Required) The condition that defines which violations to target. See [Nested condition blocks](#nested-condition-blocks) below for details.
  * `enabled` - (Required) Whether the MutingRule is enabled.
  * `name` - The name of the MutingRule.
  * `description` - The description of the MutingRule.
  * `schedule` - (Optional) Specify a schedule for enabling the MutingRule. See [Schedule](#schedule) below for details


### Nested `condition` blocks

All nested `condition` blocks support the following arguments:
  * `conditions` - (Optional) The individual MutingRuleConditions within the group. See [Nested conditions blocks](#nested-conditions-blocks) below for details.
  * `operator` - (Required) The operator used to combine all the MutingRuleConditions within the group.


### Nested `conditions` blocks
* `attribute` - (Required) The attribute on a violation.
* `operator` - (Required) The operator used to compare the attribute's value with the supplied value(s). Valid values are `ANY`, `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IN`, `IS_BLANK`, `IS_NOT_BLANK`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_IN`, `NOT_STARTS_WITH`, `STARTS_WITH`
* `values` - (Required) The value(s) to compare against the attribute's value.

### Schedule
* `start_time` (Optional) The datetime stamp that represents when the muting rule starts. This is in local ISO 8601 format without an offset. Example: '2020-07-08T14:30:00'
* `end_time` (Optional) The datetime stamp that represents when the muting rule ends. This is in local ISO 8601 format without an offset. Example: '2020-07-15T14:30:00'
* `timeZone` (Required) The time zone that applies to the muting rule schedule. Example: 'America/Los_Angeles'. See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
* `repeat` (Optional) The frequency the muting rule schedule repeats. If it does not repeat, omit this field. Options are DAILY, WEEKLY, MONTHLY
* `end_repeat` (Optional) The datetime stamp when the muting rule schedule stops repeating. This is in local ISO 8601 format without an offset. Example: '2020-07-10T15:00:00'. Conflicts with `repeat_count`
* `repeat_count` (Optional) The number of times the muting rule schedule repeats. This includes the original schedule. For example, a repeatCount of 2 will recur one time. Conflicts with `end_repeat`
* `weekly_repeat_days` (Optional) The day(s) of the week that a muting rule should repeat when the repeat field is set to 'WEEKLY'. Example: ['MONDAY', 'WEDNESDAY']

## Import
Alert conditions can be imported using a composite ID of `<account_id>:<muting_rule_id>`, e.g.

```
$ terraform import newrelic_alert_muting_rule.foo 538291:6789035

```
